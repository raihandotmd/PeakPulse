package models

import (
	"github.com/google/uuid"
)

// PlanExercise represents the plan_exercises table in the database
type PlanExercise struct {
	ID         uuid.UUID `json:"id" gorm:"type:uuid;default:gen_random_uuid();primaryKey"`
	PlanID     uuid.UUID `json:"plan_id" gorm:"type:uuid;not null"`
	ExerciseID uuid.UUID `json:"exercise_id" gorm:"type:uuid;not null"`
	Sets       int       `json:"sets" gorm:"not null"`
	Reps       int       `json:"reps" gorm:"not null"`
}

func (PlanExercise) TableName() string {
	return "plan_exercises"
}
