package server

import (
	"context"
	"net"
	"sync"
	"time"

	"github.com/google/uuid"
	"github.com/kubearmor/koach/koach/config"
	kg "github.com/kubearmor/koach/koach/log"
	"github.com/kubearmor/koach/koach/model"
	"github.com/kubearmor/koach/koach/repository"
	"google.golang.org/grpc"
	"google.golang.org/grpc/reflection"
	"gorm.io/gorm"
	"k8s.io/apimachinery/pkg/apis/meta/v1/unstructured"
	"k8s.io/apimachinery/pkg/runtime"

	kubearmor_proto "github.com/kubearmor/KubeArmor/protobuf"
	"github.com/kubearmor/koach/koach/service"
	proto "github.com/kubearmor/koach/protobuf"

	koach_api "github.com/kubearmor/koach/pkg/KoachController/api/v1"
)

// ================== //
// == Koach Server == //
// ================== //

// KoachServer Structure
type KoachServer struct {
	// port
	Port string

	// gRPC listener
	Listener net.Listener

	// observability service
	ObservabilityService service.IObservabilityService

	// log client
	LogClient *LogClient

	proto.UnimplementedObservabilityServiceServer

	// alert
	AlertRules   map[string]koach_api.KubeArmorAlertRuleSpec
	AlertStructs map[string]model.AlertStruct
	AlertLock    *sync.RWMutex
}

// NewKoachServer Function
func NewKoachServer(port string, db *gorm.DB) *KoachServer {
	ks := &KoachServer{}

	ks.Port = port

	// listen to gRPC port
	listener, err := net.Listen("tcp", ":"+ks.Port)
	if err != nil {
		kg.Errf("Failed to listen a port (%s)\n", ks.Port)
		return nil
	}
	ks.Listener = listener

	// initialize repositories
	observabilityRepository := repository.NewObservabilityRepository(db)

	// initialize services
	observabilityService := service.NewObservabilityService(
		observabilityRepository,
	)

	ks.ObservabilityService = observabilityService

	grpcServer := grpc.NewServer()
	proto.RegisterObservabilityServiceServer(grpcServer, ks)
	reflection.Register(grpcServer)

	go func() {
		if err := grpcServer.Serve(ks.Listener); err != nil {
			kg.Errf("Failed to serve gRPC server\n")
		}
	}()

	ks.AlertRules = map[string]koach_api.KubeArmorAlertRuleSpec{}
	ks.AlertStructs = map[string]model.AlertStruct{}
	ks.AlertLock = &sync.RWMutex{}

	return ks
}

func (ks *KoachServer) WatchLogs() {
	var err error

	for ks.LogClient.Running {
		var log *kubearmor_proto.Log

		if log, err = ks.LogClient.logStream.Recv(); err != nil {
			kg.Warnf("Failed to receive a log (%s)", ks.LogClient.server)
			break
		}

		kg.Printf("Succesfully receive an observability data")

		observability := model.Observability{
			ClusterName:       log.ClusterName,
			HostName:          log.HostName,
			NamespaceName:     log.NamespaceName,
			PodName:           log.PodName,
			Labels:            log.Labels,
			ContainerID:       log.ContainerID,
			ContainerName:     log.ContainerImage,
			ContainerImage:    log.ContainerImage,
			ParentProcessName: log.ParentProcessName,
			ProcessName:       log.ProcessName,
			HostPPID:          log.HostPPID,
			HostPID:           log.HostPID,
			PPID:              log.PPID,
			PID:               log.PID,
			UID:               log.UID,
			Type:              log.Type,
			Source:            log.Source,
			Operation:         model.Operation(log.Operation),
			Resource:          log.Resource,
			Data:              log.Data,
			Result:            log.Result,
		}

		err := ks.ObservabilityService.Save(&observability)
		if err != nil {
			kg.Errf("Failed to save observability data")
		}

		kg.Printf("Succesfully save an observability data to DB")

		ks.CheckObservabilityWithAlertRule(observability)
	}
}

func (ks *KoachServer) GetFeedsFromRelay(relayConfig config.RelayConfig) {
	err := ks.NewLogClient(relayConfig.RelayIP + ":" + relayConfig.RelayPort)
	if err != nil {
		return
	}

	go ks.WatchLogs()
}

func (ks *KoachServer) PeriodicDataDeletion(ageString string) {
	t := time.NewTicker(24 * time.Hour)
	defer t.Stop()

	for {
		select {
		case <-t.C:
			kg.Printf("Deleting data from database")

			age := model.Duration{}
			err := age.FromString(ageString)
			if err != nil {
				kg.Printf("Invalid format for data age")
				return
			}

			err = ks.ObservabilityService.DeleteByAgeSeconds(age.GetSeconds())
			if err != nil {
				kg.Printf("Error on deleting old data")
			}
		}
	}
}

func (ks *KoachServer) WatchAlertRule() {
	watcher := K8s.WatchKubeArmorAlertRules()

	if watcher != nil {
		for event := range watcher.ResultChan() {
			if event.Type != "ADDED" && event.Type != "MODIFIED" && event.Type != "DELETED" {
				continue
			}

			alertRuleUnstructured, ok := event.Object.(*unstructured.Unstructured)
			if !ok {
				continue
			}

			alertRule := koach_api.KubeArmorAlertRule{}

			err := runtime.DefaultUnstructuredConverter.FromUnstructured(alertRuleUnstructured.UnstructuredContent(), &alertRule)
			if err != nil {
				continue
			}

			switch event.Type {
			case "ADDED":
				kg.Printf("Succesfully added an alert rule")
				ks.AlertRules[string(alertRule.ObjectMeta.UID)] = alertRule.Spec
			case "MODIFIED":
				kg.Printf("Succesfully modified an alert rule")
				ks.AlertRules[string(alertRule.ObjectMeta.UID)] = alertRule.Spec
			case "DELETED":
				kg.Printf("Succesfully delete an alert rule")
				delete(ks.AlertRules, string(alertRule.ObjectMeta.UID))
			}
		}
	} else {
		kg.Errf("Failed to watch KubeArmorAlertRule")
	}
}

func (ks *KoachServer) CheckObservabilityWithAlertRule(observability model.Observability) {
	for _, alertRule := range ks.AlertRules {
		isViolatingRule := false

		filter := model.ObservabilityFilter{}
		filter.FromAlertRule(alertRule, observability)

		if observability.IsValid(filter) {
			result, err := ks.ObservabilityService.Get(filter)
			if err != nil {
				kg.Errf("Error on fetching data while checking with alert rule")
			}

			isViolatingRule = len(result) >= alertRule.Condition.Occurrence.Count
		}

		if isViolatingRule {
			alert := &model.Alert{
				Message:       alertRule.Message,
				Severity:      alertRule.Severity,
				Observability: observability,
			}

			ks.AlertLock.RLock()
			for _, alertStruct := range ks.AlertStructs {
				if alert.IsValid(alertStruct.Filter) {
					select {
					case alertStruct.Broadcast <- alert:
					default:
					}
				}
			}
			ks.AlertLock.RUnlock()
		}
	}
}

func (ks *KoachServer) Get(ctx context.Context, request *proto.GetRequest) (*proto.GetResponse, error) {
	observabilityFilter := model.ObservabilityFilter{
		NamespaceID:      request.GetNamespaceId(),
		PodID:            request.GetPodId(),
		ContainerID:      request.GetContainerId(),
		OperationType:    model.Operation(request.GetOperationType()),
		Labels:           model.LabelsFromString(request.GetLabels()),
		SinceTimeSeconds: 0,
	}

	if request.GetTime() != "" {
		duration := model.Duration{}

		if err := duration.FromString(request.GetTime()); err != nil {
			return nil, err
		}

		observabilityFilter.SinceTimeSeconds = duration.GetSeconds()
	}

	observabilities, err := ks.ObservabilityService.Get(observabilityFilter)
	if err != nil {
		return nil, err
	}

	observabilitiesResponse := []*proto.ObservabilityData{}

	for _, observability := range observabilities {
		observabilityResponse := proto.ObservabilityData{
			ClusterName:       observability.ClusterName,
			HostName:          observability.HostName,
			NamespaceName:     observability.NamespaceName,
			PodName:           observability.PodName,
			Labels:            observability.Labels,
			ContainerId:       observability.ContainerID,
			ContainerName:     observability.ContainerName,
			ContainerImage:    observability.ContainerImage,
			ParentProcessName: observability.ParentProcessName,
			ProcessName:       observability.ProcessName,
			HostPpid:          observability.HostPPID,
			HostPid:           observability.HostPID,
			Ppid:              observability.PPID,
			Pid:               observability.PID,
			Uid:               observability.UID,
			Type:              observability.Type,
			Source:            observability.Source,
			Operation:         string(observability.Operation),
			Resource:          observability.Resource,
			Data:              observability.Data,
			Result:            observability.Result,
			CreatedAt:         observability.CreatedAt.String(),
		}

		observabilitiesResponse = append(observabilitiesResponse, &observabilityResponse)
	}

	response := proto.GetResponse{
		Data: observabilitiesResponse,
	}

	return &response, nil
}

func (ks *KoachServer) addAlertChans(id string, alertChan chan *model.Alert, filter model.ListenAlertFilter) {
	ks.AlertLock.Lock()
	defer ks.AlertLock.Unlock()

	ks.AlertStructs[id] = model.AlertStruct{
		Filter:    filter,
		Broadcast: alertChan,
	}

	kg.Printf("Added a new client (%s) for ListenAlert", id)
}

func (ks *KoachServer) removeAlertChans(id string) {
	ks.AlertLock.Lock()
	defer ks.AlertLock.Unlock()

	delete(ks.AlertStructs, id)

	kg.Printf("Deleted client (%s) for ListenAlert", id)
}

func (ks *KoachServer) ListenAlert(req *proto.ListenAlertRequest, stream proto.ObservabilityService_ListenAlertServer) error {
	id := uuid.Must(uuid.NewRandom()).String()

	alertChan := make(chan *model.Alert, 8192)
	defer close(alertChan)

	alertFilter := model.ListenAlertFilter{
		NamespaceID: req.NamespaceId,
		PodID:       req.PodId,
		ContainerID: req.ContainerId,
	}

	ks.addAlertChans(id, alertChan, alertFilter)
	defer ks.removeAlertChans(id)

	for {
		time.Sleep(time.Second * 1)

		select {
		case <-stream.Context().Done():
			return nil
		case msg, valid := <-alertChan:
			if !valid {
				continue
			}

			alertResponse := proto.ListenAlertResponse{
				Message:  msg.Message,
				Severity: int32(msg.Severity),
				Observability: &proto.ObservabilityData{
					ClusterName:       msg.Observability.ClusterName,
					HostName:          msg.Observability.HostName,
					NamespaceName:     msg.Observability.NamespaceName,
					PodName:           msg.Observability.PodName,
					Labels:            msg.Observability.Labels,
					ContainerId:       msg.Observability.ContainerID,
					ContainerName:     msg.Observability.ContainerName,
					ContainerImage:    msg.Observability.ContainerImage,
					ParentProcessName: msg.Observability.ParentProcessName,
					ProcessName:       msg.Observability.ProcessName,
					HostPpid:          msg.Observability.PPID,
					HostPid:           msg.Observability.HostPID,
					Ppid:              msg.Observability.PPID,
					Pid:               msg.Observability.PID,
					Uid:               msg.Observability.UID,
					Type:              msg.Observability.Type,
					Source:            msg.Observability.Source,
					Operation:         string(msg.Observability.Operation),
					Resource:          msg.Observability.Resource,
					Data:              msg.Observability.Data,
					Result:            msg.Observability.Result,
					CreatedAt:         msg.Observability.CreatedAt.String(),
				},
			}

			if err := stream.Send(&alertResponse); err != nil {
				return err
			}
		default:
		}
	}
}

// DestroyKoachServer Function
func (ks *KoachServer) DestroyKoachServer() error {
	// wait for a while
	time.Sleep(time.Second * 1)

	// close listener
	if ks.Listener != nil {
		if err := ks.Listener.Close(); err != nil {
			kg.Err(err.Error())
		}
		ks.Listener = nil
	}

	return nil
}

// =============== //
// == Log Feeds == //
// =============== //

// LogClient Structure
type LogClient struct {
	// flags
	Running bool

	// server
	server string

	// connection
	conn *grpc.ClientConn

	// client
	client kubearmor_proto.LogServiceClient

	// logs
	logStream kubearmor_proto.LogService_WatchLogsClient
}

func (ks *KoachServer) NewLogClient(server string) error {
	var err error

	lc := &LogClient{}

	lc.Running = true

	lc.server = server

	lc.conn, err = grpc.Dial(lc.server, grpc.WithInsecure())
	if err != nil {
		kg.Warnf("Failed to connect to Relay's gRPC service (%s)", server)
		return err
	}
	defer func() {
		if err != nil {
			err = lc.DestroyClient()
			if err != nil {
				kg.Warnf("DestroyClient() failed err=%s", err.Error())
			}
		}
	}()

	lc.client = kubearmor_proto.NewLogServiceClient(lc.conn)

	logIn := kubearmor_proto.RequestMessage{}
	logIn.Filter = "system"

	lc.logStream, err = lc.client.WatchLogs(context.Background(), &logIn)
	if err != nil {
		kg.Warnf("Failed to call WatchLogs (%s)\n err=%s", server, err.Error())
		return err
	}

	ks.LogClient = lc
	return nil
}

func (lc *LogClient) DestroyClient() error {
	lc.Running = false

	if err := lc.conn.Close(); err != nil {
		return err
	}

	return nil
}
