package models

import (
	"time"

	"github.com/google/uuid"
)

// AuthUser represents the auth.users table in the database
type AuthUser struct {
	ID        uuid.UUID `json:"id" gorm:"type:uuid;default:gen_random_uuid();primaryKey"`
	Email     string    `json:"email" gorm:"type:text;unique;not null"`
	CreatedAt time.Time `json:"created_at" gorm:"type:timestamp with time zone;default:CURRENT_TIMESTAMP"`
	UpdatedAt time.Time `json:"updated_at" gorm:"type:timestamp with time zone;default:CURRENT_TIMESTAMP"`
}

func (AuthUser) TableName() string {
	return "auth.users"
}
