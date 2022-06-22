package model

type FileAccess struct {
	ID   string `gorm:"type:varchar(255)"`
	Path string `gorm:"type:varchar(255)"`
	BaseDate
}

func (FileAccess) TableName() string {
	return "file_accesses"
}

func (FileAccess) IsOperationDetail() {}
