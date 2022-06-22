package model

type Observability struct {
	ID            string          `gorm:"type:varchar(255)"`
	NamespaceID   string          `gorm:"type:varchar(255)"`
	DeploymentID  string          `gorm:"type:varchar(255)"`
	NodeID        string          `gorm:"type:varchar(255)"`
	PodID         string          `gorm:"type:varchar(255)"`
	ContainerID   string          `gorm:"type:varchar(255)"`
	OperationType Operation       `gorm:"type:varchar(255)"`
	DetailID      string          `gorm:"type:varchar(255)"`
	Detail        OperationDetail `gorm:"-"`
	BaseDate
}

func (Observability) TableName() string {
	return "observabilities"
}

type ObservabilityFilter struct {
	NamespaceID   string
	DeploymentID  string
	NodeID        string
	PodID         string
	ContainerID   string
	OperationType Operation
}
