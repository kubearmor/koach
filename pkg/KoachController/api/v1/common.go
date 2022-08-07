package v1

type SelectorType struct {
	// +kubebuilder:validation:optional
	MatchLabels map[string]string `json:"matchLabels,omitempty"`
}

type ConditionType struct {
	// +kubebuilder:validation:optional
	IsSamePID bool `json:"isSamePID,omitempty"`

	// +kubebuilder:validation:optional
	Occurrence OccurenceType `json:"occurrence,omitempty"`

	// +kubebuilder:validation:optional
	File FileConditionType `json:"file,omitempty"`

	// +kubebuilder:validation:optional
	Network NetworkConditionType `json:"network,omitempty"`
}

type OccurenceType struct {
	// +kubebuilder:validation:Minimum:=1
	Count int `json:"count"`

	// +kubebuilder:validation:optional
	// +kubebuilder:validation:Pattern=([0-9]+(s|m|h|d|M|y))+
	Timeframe string `json:"timeframe,omitempty"`
}

type FileConditionType struct {
	// +kubebuilder:validation:optional
	Path string `json:"path,omitempty"`

	// +kubebuilder:validation:Enum=delete
	Action string `json:"action,omitempty"`
}

type NetworkConditionType struct {
	// +kubebuilder:validation:optional
	// +kubebuilder:validation:Pattern=(icmp|ICMP|tcp|TCP|udp|UDP|raw|RAW)$
	Protocol string `json:"protocol,omitempty"`

	// +kubebuilder:validation:optional
	SourceAddress NetworkAddressType `json:"sourceAddress,omitempty"`

	// +kubebuilder:validation:optional
	DestinationAddress NetworkAddressType `json:"destinationAddress,omitempty"`
}

type NetworkAddressType struct {
	// +kubebuilder:validation:optional
	IP string `json:"ip,omitempty"`

	// +kubebuilder:validation:optional
	IsSuspicious bool `json:"isSuspicious,omitempty"`
}
