package model

type ProcessSpawn struct {
	ID      string `gorm:"type:varchar(255)"`
	User    string `gorm:"type:varchar(255)"`
	Command string `gorm:"type:varchar(255)"`
	BaseDate
}

func (ProcessSpawn) TableName() string {
	return "process_spawns"
}

func (ProcessSpawn) IsOperationDetail() {}
