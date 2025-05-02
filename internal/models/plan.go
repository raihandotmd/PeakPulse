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

// PlanExerciseInput represents the input for a single exercise in a plan
type PlanExerciseInput struct {
	ExerciseID uuid.UUID `json:"exercise_id" binding:"required"`
	Sets       int       `json:"sets" binding:"required"`
	Reps       int       `json:"reps" binding:"required"`
}

// CreatePlanRequest represents the payload for creating a workout plan
type CreatePlanRequest struct {
	Title       string              `json:"title" binding:"required"`
	Description string              `json:"description"`
	Exercises   []PlanExerciseInput `json:"exercises" binding:"required"`
}

// UpdatePlanRequest represents the payload for updating a workout plan
type UpdatePlanRequest struct {
	Title       *string              `json:"title"`       // Optional
	Description *string              `json:"description"` // Optional
	Exercises   *[]PlanExerciseInput `json:"exercises"`   // Optional
}
