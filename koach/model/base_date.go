package model

import (
	"time"

	"gorm.io/gorm"
)

type BaseDate struct {
	CreatedAt time.Time
	UpdatedAt time.Time
	DeletedAt gorm.DeletedAt
}
