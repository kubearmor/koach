package model

import (
	"regexp"
	"strings"

	turnip "github.com/nyrahul/turnip/api"

	koach_api "github.com/kubearmor/koach/pkg/KoachController/api/v1"
)

type Observability struct {
	ID                string    `gorm:"type:varchar(255)"`
	ClusterName       string    `gorm:"type:varchar(255)"`
	HostName          string    `gorm:"type:varchar(255)"`
	NamespaceName     string    `gorm:"type:varchar(255)"`
	PodName           string    `gorm:"type:varchar(255)"`
	Labels            string    `gorm:"type:varchar(255)"`
	ContainerID       string    `gorm:"type:varchar(255)"`
	ContainerName     string    `gorm:"type:varchar(255)"`
	ContainerImage    string    `gorm:"type:varchar(255)"`
	ParentProcessName string    `gorm:"type:varchar(255)"`
	ProcessName       string    `gorm:"type:varchar(255)"`
	HostPPID          int32     `gorm:"type:int"`
	HostPID           int32     `gorm:"type:int"`
	PPID              int32     `gorm:"type:int"`
	PID               int32     `gorm:"type:int"`
	UID               int32     `gorm:"type:int"`
	Type              string    `gorm:"type:varchar(255)"`
	Source            string    `gorm:"type:varchar(255)"`
	Operation         Operation `gorm:"type:varchar(255)"`
	Resource          string    `gorm:"type:varchar(255)"`
	Data              string    `gorm:"type:varchar(255)"`
	Result            string    `gorm:"type:varchar(255)"`
	BaseDate
}

func (Observability) TableName() string {
	return "observabilities"
}

type ObservabilityFilter struct {
	NamespaceID                    string
	PodID                          string
	ContainerID                    string
	OperationType                  Operation
	Labels                         map[string]string
	PID                            int32
	SinceTimeSeconds               int
	ResourceRegex                  string
	Syscall                        SystemCall
	NetworkProtocol                string
	SourceAddressIP                string
	SourceAddressIsSuspicious      bool
	DestinationAddressIP           string
	DestinationAddressIsSuspicious bool
}

func (filter *ObservabilityFilter) FromAlertRule(alertRule koach_api.KubeArmorAlertRuleSpec, observability Observability) {
	filter.Labels = alertRule.Selector.MatchLabels

	if alertRule.Condition.Occurrence.Timeframe != "" {
		duration := Duration{}
		duration.FromString(alertRule.Condition.Occurrence.Timeframe)

		filter.SinceTimeSeconds = duration.GetSeconds()
	}

	if alertRule.Condition.IsSamePID {
		filter.PID = observability.PID
	}

	switch {
	case alertRule.Operation == string(OperationFileAccess) &&
		alertRule.Condition.File.Action == string(FileActionDelete):
		filter.OperationType = OperationSystemCall
		filter.ResourceRegex = alertRule.Condition.File.Path
		filter.Syscall = SystemCallUnlinkat

	case alertRule.Operation == string(OperationNetworkCall):
		filter.OperationType = OperationNetworkCall
		filter.NetworkProtocol = alertRule.Condition.Network.Protocol
		filter.SourceAddressIP = alertRule.Condition.Network.SourceAddress.IP
		filter.SourceAddressIsSuspicious = alertRule.Condition.Network.SourceAddress.IsSuspicious
		filter.DestinationAddressIP = alertRule.Condition.Network.DestinationAddress.IP
		filter.DestinationAddressIsSuspicious = alertRule.Condition.Network.DestinationAddress.IsSuspicious
	}
}

func (observability *Observability) IsValid(filter ObservabilityFilter) bool {
	if observability.Operation != filter.OperationType {
		return false
	}

	if filter.ResourceRegex != "" {
		if match, _ := regexp.MatchString(filter.ResourceRegex, observability.Resource); !match {
			return false
		}
	}

	if !observability.IsValidLabels(filter.Labels) {
		return false
	}

	rawDataString := []string{}
	rawDataString = append(rawDataString, strings.Fields(observability.Resource)...)
	rawDataString = append(rawDataString, strings.Fields(observability.Data)...)

	data := map[string]string{}
	for _, rawData := range rawDataString {
		if strings.Contains(rawData, "=") {
			dataSplit := strings.Split(rawData, "=")
			key := dataSplit[0]
			value := dataSplit[1]

			data[key] = value
		}
	}

	if filter.Syscall != "" {
		value, found := data["syscall"]
		if !found || value != string(filter.Syscall) {
			return false
		}
	}

	if filter.NetworkProtocol != "" {
		value, found := data["protocol"]
		if !found || !strings.EqualFold(value, filter.NetworkProtocol) {
			return false
		}
	}

	if filter.DestinationAddressIP != "" {
		possibleKeys := []string{"sin_addr", "remoteip"}
		foundOverall := false

		for _, key := range possibleKeys {
			value, found := data[key]

			if found {
				if value != filter.DestinationAddressIP {
					return false
				}

				foundOverall = true
				break
			}
		}

		if !foundOverall {
			return false
		}
	}

	if filter.DestinationAddressIsSuspicious {
		err := turnip.Setup(turnip.TurnipDefSrc)
		if err != nil {
			return false
		}

		destinationAddress, found := data["sin_addr"]
		if !found {
			return false
		}

		src, _ := turnip.AddressIsBlocked(destinationAddress)
		return src != nil
	}

	return true
}

func (observability Observability) IsValidLabels(labels map[string]string) bool {
	observabilityLabel := LabelsFromString(observability.Labels)

	for filterKey, filterValue := range labels {
		labelValue, isLabelKeyExist := observabilityLabel[filterKey]

		if !isLabelKeyExist {
			return false
		}

		match, _ := regexp.MatchString(filterValue, labelValue)
		if match {
			continue
		}

		if filterValue != labelValue {
			return false
		}
	}

	return true
}
