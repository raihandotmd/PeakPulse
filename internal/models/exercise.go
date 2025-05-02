package models

import (
	"time"

	"github.com/google/uuid"
)

// Exercise represents master exercise catalog
type Exercise struct {
	ID          uuid.UUID `gorm:"type:uuid;default:uuid_generate_v4();primaryKey"`
	Name        string    `gorm:"unique;not null"`
	Description string
	Category    string
	MuscleGroup string
	CreatedAt   time.Time `gorm:"autoCreateTime"`
	UpdatedAt   time.Time `gorm:"autoUpdateTime"`
}

func (Exercise) TableName() string {
	return "exercises"
}
