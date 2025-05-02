package models

import (
	"time"

	"github.com/google/uuid"
)

// User represents the public.users table in the database
type User struct {
	ID           uuid.UUID  `json:"id" gorm:"type:uuid;default:gen_random_uuid();primaryKey"`
	Email        string     `json:"email" gorm:"type:text;unique;not null"`
	PasswordHash string     `json:"password_hash" gorm:"type:text;not null"`
	Name         string     `json:"name" gorm:"type:text"`
	CreatedAt    time.Time  `json:"created_at" gorm:"type:timestamp with time zone;default:CURRENT_TIMESTAMP"`
	UpdatedAt    time.Time  `json:"updated_at" gorm:"type:timestamp with time zone;default:CURRENT_TIMESTAMP"`
	UserID       *uuid.UUID `json:"user_id" gorm:"type:uuid"`

	// Relationships
	AuthUser AuthUser `json:"auth_user" gorm:"foreignKey:UserID;references:ID;constraint:OnDelete:CASCADE"`
}

func (User) TableName() string {
	return "public.users"
}
