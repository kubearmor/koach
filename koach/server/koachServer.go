package server

import (
	"context"
	"net"
	"time"

	kg "github.com/kubearmor/koach/koach/log"
	"github.com/kubearmor/koach/koach/model"
	"github.com/kubearmor/koach/koach/repository"
	"google.golang.org/grpc"
	"google.golang.org/grpc/reflection"
	"gorm.io/gorm"

	"github.com/kubearmor/koach/koach/service"
	proto "github.com/kubearmor/koach/protobuf"
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

	proto.UnimplementedObservabilityServiceServer
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
	fileAccessRepository := repository.NewFileAccessRepository(db)
	networkCallRepository := repository.NewNetworkCallRepository(db)
	processSpawnRepository := repository.NewProcessSpawnRepository(db)

	// initialize services
	observabilityService := service.NewObservabilityService(
		fileAccessRepository,
		networkCallRepository,
		processSpawnRepository,
	)

	ks.ObservabilityService = observabilityService

	grpcServer := grpc.NewServer()
	proto.RegisterObservabilityServiceServer(grpcServer, ks)
	reflection.Register(grpcServer)

	if err := grpcServer.Serve(ks.Listener); err != nil {
		kg.Errf("Failed to server gRPC server\n")
		return nil
	}

	return ks
}

func (ks *KoachServer) Get(ctx context.Context, request *proto.GetRequest) (*proto.GetResponse, error) {
	observabilityFilter := model.ObservabilityFilter{
		NamespaceID:   request.GetNamespaceId(),
		DeploymentID:  request.GetDeploymentId(),
		NodeID:        request.GetNodeId(),
		PodID:         request.GetPodId(),
		ContainerID:   request.GetContainerId(),
		OperationType: model.Operation(request.GetOperationType()),
	}

	observabilities, err := ks.ObservabilityService.Get(observabilityFilter)
	if err != nil {
		return nil, err
	}

	observabilitiesResponse := []*proto.ObservabilityData{}

	for _, observability := range observabilities {
		observabilityResponse := proto.ObservabilityData{
			NamespaceId:  observability.NamespaceID,
			DeploymentId: observability.DeploymentID,
			NodeId:       observability.NodeID,
			PodId:        observability.PodID,
			ContainerId:  observability.ContainerID,
			CreatedAt:    observability.CreatedAt.String(),
		}

		switch detail := observability.Detail.(type) {
		case *model.FileAccess:
			observabilityResponse.OperationData = &proto.ObservabilityData_File{
				File: &proto.FileAccessData{
					Path: detail.Path,
				},
			}
		case *model.NetworkCall:
			observabilityResponse.OperationData = &proto.ObservabilityData_Network{
				Network: &proto.NetworkCallData{
					Protocol: detail.Protocol,
					Address:  detail.Address,
				},
			}
		case *model.ProcessSpawn:
			observabilityResponse.OperationData = &proto.ObservabilityData_Process{
				Process: &proto.ProcessSpawnData{
					User:    detail.User,
					Command: detail.Command,
				},
			}
		}

		observabilitiesResponse = append(observabilitiesResponse, &observabilityResponse)
	}

	response := proto.GetResponse{
		Data: observabilitiesResponse,
	}

	return &response, nil
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
