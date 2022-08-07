package model

import metav1 "k8s.io/apimachinery/pkg/apis/meta/v1"

// K8sPolicyStatus Structure
type K8sPolicyStatus struct {
	Status string `json:"status,omitempty"`
}

// K8sKubeArmorAlertRuleEvent Structure
type K8sKubeArmorAlertRuleEvent struct {
	Type   string                `json:"type"`
	Object K8sKubeArmorAlertRule `json:"object"`
}

// K8sKubeArmorAlertRule Structure
type K8sKubeArmorAlertRule struct {
	Metadata metav1.ObjectMeta `json:"metadata"`
	Spec     AlertRuleSpec     `json:"spec"`
	Status   K8sPolicyStatus   `json:"status,omitempty"`
}

// AlertRuleSpec Structure
type AlertRuleSpec struct {
	Operation string             `json:"operation"`
	Severity  int                `json:"severity"`
	Message   string             `json:"message,omitempty"`
	Selector  AlertRuleSelector  `json:"selector"`
	Condition AlertRuleCondition `json:"condition"`
}

type AlertRuleSelector struct {
	Labels map[string]string `json:"matchLabels"`
}

type AlertRuleCondition struct {
	IsSamePID  bool                         `json:"isSamePID"`
	Occurrence AlertRuleConditionOccurrence `json:"occurrence"`
	File       AlertRuleConditionFile       `json:"file"`
	Network    AlertRuleConditionNetwork    `json:"network"`
}

type AlertRuleConditionOccurrence struct {
	Count     int    `json:"count"`
	Timeframe string `json:"timeframe"`
}

type AlertRuleConditionFile struct {
	Path   string `json:"path"`
	Action string `json:"action"`
}

type AlertRuleConditionNetwork struct {
	Protocol           string                           `json:"protocol"`
	SourceAddress      AlertRuleConditionNetworkAddress `json:"sourceAddress"`
	DestinationAddress AlertRuleConditionNetworkAddress `json:"destinationAddress"`
}

type AlertRuleConditionNetworkAddress struct {
	IP           string `json:"ip"`
	IsSuspicious bool   `json:"isSuspicious"`
}
