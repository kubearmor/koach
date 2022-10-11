package model

type Alert struct {
	Message       string
	Severity      int
	Observability Observability
}

type ListenAlertFilter struct {
	NamespaceID string
	PodID       string
	ContainerID string
}

type AlertStruct struct {
	Filter    ListenAlertFilter
	Broadcast chan *Alert
}

func (alert Alert) IsValid(filter ListenAlertFilter) bool {
	if filter.NamespaceID != "" && alert.Observability.NamespaceName != filter.NamespaceID {
		return false
	}

	if filter.PodID != "" && alert.Observability.PodName != filter.PodID {
		return false
	}

	if filter.ContainerID != "" && alert.Observability.ContainerID != filter.ContainerID {
		return false
	}

	return true
}
