package v1

import (
	metav1 "k8s.io/apimachinery/pkg/apis/meta/v1"
)

// KubeArmorAlertRuleSpec defines the desired state of KubeArmorAlertRule
type KubeArmorAlertRuleSpec struct {
	// +kubebuilder:validation:optional
	Selector SelectorType `json:"selector,omitempty"`

	// +kubebuilder:validation:Enum=File;Network
	Operation string `json:"operation"`

	// +kubebuilder:validation:Minimum:=1
	// +kubebuilder:validation:Maximum:=10
	Severity int `json:"severity"`

	Message string `json:"message"`

	// +kubebuilder:validation:optional
	Condition ConditionType `json:"condition,omitempty"`
}

// KubeArmorAlertRuleStatus defines the observed state of KubeArmorAlertRule
type KubeArmorAlertRuleStatus struct {
}

//+kubebuilder:object:root=true
//+kubebuilder:subresource:status

// KubeArmorAlertRule is the Schema for the kubearmoralertrules API
// +kubebuilder:resource:scope=Cluster,shortName=ar
// +genclient
type KubeArmorAlertRule struct {
	metav1.TypeMeta   `json:",inline"`
	metav1.ObjectMeta `json:"metadata,omitempty"`

	Spec   KubeArmorAlertRuleSpec   `json:"spec,omitempty"`
	Status KubeArmorAlertRuleStatus `json:"status,omitempty"`
}

//+kubebuilder:object:root=true

// KubeArmorAlertRuleList contains a list of KubeArmorAlertRule
type KubeArmorAlertRuleList struct {
	metav1.TypeMeta `json:",inline"`
	metav1.ListMeta `json:"metadata,omitempty"`
	Items           []KubeArmorAlertRule `json:"items"`
}

func init() {
	SchemeBuilder.Register(&KubeArmorAlertRule{}, &KubeArmorAlertRuleList{})
}
