package models

import (
	"time"

	"github.com/google/uuid"
)

// WorkoutPlan represents the workout_plans table in the database
type WorkoutPlan struct {
	ID          uuid.UUID `json:"id" gorm:"type:uuid;default:gen_random_uuid();primaryKey"`
	UserID      uuid.UUID `json:"user_id" gorm:"type:uuid;not null"`
	Title       string    `json:"title" gorm:"type:text;not null"`
	Description string    `json:"description" gorm:"type:text"`
	CreatedAt   time.Time `json:"created_at" gorm:"type:timestamp with time zone;default:CURRENT_TIMESTAMP"`
	UpdatedAt   time.Time `json:"updated_at" gorm:"type:timestamp with time zone;default:CURRENT_TIMESTAMP"`
}

func (WorkoutPlan) TableName() string {
	return "workout_plans"
}

// WorkoutPlansListView represents the workout_plans_list_view in the database
type WorkoutPlansListView struct {
	PlanID              uuid.UUID `json:"plan_id" gorm:"column:plan_id"`
	PlanTitle           string    `json:"plan_title" gorm:"column:plan_title"`
	PlanDescription     string    `json:"plan_description" gorm:"column:plan_description"`
	ExerciseID          uuid.UUID `json:"exercise_id" gorm:"column:exercise_id"`
	ExerciseName        string    `json:"exercise_name" gorm:"column:exercise_name"`
	ExerciseDescription string    `json:"exercise_description" gorm:"column:exercise_description"`
	Sets                int       `json:"sets" gorm:"column:sets"`
	Reps                int       `json:"reps" gorm:"column:reps"`
}

func (WorkoutPlansListView) TableName() string {
	return "workout_plans_list_view"
}
