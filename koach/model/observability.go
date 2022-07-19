package model

import (
	"regexp"
	"strings"
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
	NamespaceID      string
	PodID            string
	ContainerID      string
	OperationType    Operation
	Labels           string
	SinceTimeSeconds int
}

func FilterObservabilitiesByFilter(observabilities []Observability, labelsFilter LabelsFilter) []Observability {
	filteredObservabilities := []Observability{}

	for _, observability := range observabilities {
		observabilityLabel := map[string]string{}

		for _, label := range strings.Split(observability.Labels, ",") {
			labelSplit := strings.Split(label, "=")

			labelKey := labelSplit[0]
			labelValue := labelSplit[1]

			observabilityLabel[labelKey] = labelValue
		}

		isObservabilityValid := true

		for filterKey, filterValue := range labelsFilter.Filter {
			labelValue, isLabelKeyExist := observabilityLabel[filterKey]

			if !isLabelKeyExist {
				isObservabilityValid = false
				break
			}

			match, _ := regexp.MatchString(filterValue, labelValue)
			if match {
				continue
			}

			if filterValue != labelValue {
				isObservabilityValid = false
				break
			}
		}

		if isObservabilityValid {
			filteredObservabilities = append(filteredObservabilities, observability)
		}
	}

	return filteredObservabilities
}
