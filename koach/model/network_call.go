package model

type NetworkCall struct {
	ID       string `gorm:"type:varchar(255)"`
	Protocol string `gorm:"type:varchar(255)"`
	Address  string `gorm:"type:varchar(255)"`
	BaseDate
}

func (NetworkCall) TableName() string {
	return "network_calls"
}

func (NetworkCall) IsOperationDetail() {}
