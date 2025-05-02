package handlers

import (
	"net/http"

	"github.com/gin-gonic/gin"
	"github.com/google/uuid"
	"github.com/raihandotmd/peakPulse/internal/models"
	"gorm.io/gorm"
)

type PlanHandler struct{ DB *gorm.DB }

// NewPlanHandler returns a new PlanHandler
func NewPlanHandler(db *gorm.DB) *PlanHandler {
	return &PlanHandler{DB: db}
}

func (h *PlanHandler) CreatePlan(c *gin.Context) {
	// Extract user ID from the context
	userID := c.GetString("userID")
	if userID == "" {
		c.JSON(http.StatusUnauthorized, gin.H{"error": "user ID not found in context"})
		return
	}

	// Validate the user ID
	var user models.User
	if err := h.DB.First(&user, "user_id = ?", userID).Error; err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "invalid user ID"})
		return
	}

	// Parse the request body
	var req struct {
		Title       string `json:"title" binding:"required"`
		Description string `json:"description"`
		Exercises   []struct {
			ExerciseID uuid.UUID `json:"exercise_id" binding:"required"`
			Sets       int       `json:"sets" binding:"required"`
			Reps       int       `json:"reps" binding:"required"`
		} `json:"exercises" binding:"required"`
	}
	if err := c.ShouldBindJSON(&req); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "invalid payload", "details": err.Error()})
		return
	}

	// title workout plan must be unique.
	// check the title if the plan that user is trying to create already exists
	var existingPlan models.WorkoutPlan
	if err := h.DB.Where("user_id = ? AND title = ?", user.ID, req.Title).First(&existingPlan).Error; err == nil {
		c.JSON(http.StatusConflict, gin.H{"error": "workout plan with this title already exists"})
		return
	}

	// Create a new workout plan
	plan := models.WorkoutPlan{
		UserID:      user.ID,
		Title:       req.Title,
		Description: req.Description,
	}
	if err := h.DB.Create(&plan).Error; err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "failed to create workout plan", "details": err.Error()})
		return
	}

	// Prepare plan exercises for bulk insertion
	var planExercises []models.PlanExercise
	for _, ex := range req.Exercises {
		planExercises = append(planExercises, models.PlanExercise{
			PlanID:     plan.ID,
			ExerciseID: ex.ExerciseID,
			Sets:       ex.Sets,
			Reps:       ex.Reps,
		})
	}

	// Bulk insert plan exercises
	if err := h.DB.Create(&planExercises).Error; err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "failed to add exercises to the plan", "details": err.Error()})
		return
	}

	// Respond with the created plan
	c.JSON(http.StatusCreated, gin.H{
		"plan":      plan,
		"exercises": planExercises,
	})
}

func (h *PlanHandler) ListPlans(c *gin.Context) {
	// Fetch all records from the workout_plans_list_view
	var plans []models.WorkoutPlansListView
	if err := h.DB.Find(&plans).Error; err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "failed to fetch workout plans", "details": err.Error()})
		return
	}

	// Respond with the list of plans
	c.JSON(http.StatusOK, gin.H{"plans": plans})
}

func (h *PlanHandler) GetPlan(c *gin.Context) {
	planID := c.Param("id")
	var plan []models.WorkoutPlansListView
	// Fetch all the exercises plan by ID
	if err := h.DB.Where("plan_id = ?", planID).Find(&plan).Error; err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "failed to fetch workout plan", "details": err.Error()})
		return
	}
	// Check if the plan exists
	if len(plan) == 0 {
		c.JSON(http.StatusNotFound, gin.H{"error": "workout plan not found"})
		return
	}

	c.JSON(http.StatusOK, gin.H{"plan": plan})
}

func (h *PlanHandler) UpdatePlan(c *gin.Context) {
	// Extract the plan_id from the URL parameters
	planID := c.Param("id")
	if planID == "" {
		c.JSON(http.StatusBadRequest, gin.H{"error": "plan_id is required"})
		return
	}

	// Parse the request body
	var req struct {
		Title       *string `json:"title"`       // Optional
		Description *string `json:"description"` // Optional
		Exercises   *[]struct {
			ExerciseID uuid.UUID `json:"exercise_id" binding:"required"`
			Sets       int       `json:"sets" binding:"required"`
			Reps       int       `json:"reps" binding:"required"`
			Weight     float64   `json:"weight"`
		} `json:"exercises"` // Optional
	}
	if err := c.ShouldBindJSON(&req); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "invalid payload", "details": err.Error()})
		return
	}

	// Fetch the existing plan
	var plan models.WorkoutPlan
	if err := h.DB.First(&plan, "id = ?", planID).Error; err != nil {
		c.JSON(http.StatusNotFound, gin.H{"error": "workout plan not found", "details": err.Error()})
		return
	}

	// Update the plan details if provided
	if req.Title != nil {
		plan.Title = *req.Title
	}
	if req.Description != nil {
		plan.Description = *req.Description
	}
	if err := h.DB.Save(&plan).Error; err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "failed to update workout plan", "details": err.Error()})
		return
	}

	// Update exercises if provided
	if req.Exercises != nil {
		// Delete existing exercises for the plan
		if err := h.DB.Where("plan_id = ?", plan.ID).Delete(&models.PlanExercise{}).Error; err != nil {
			c.JSON(http.StatusInternalServerError, gin.H{"error": "failed to delete existing exercises", "details": err.Error()})
			return
		}

		// Prepare new exercises for bulk insertion
		var planExercises []models.PlanExercise
		for _, ex := range *req.Exercises {
			planExercises = append(planExercises, models.PlanExercise{
				PlanID:     plan.ID,
				ExerciseID: ex.ExerciseID,
				Sets:       ex.Sets,
				Reps:       ex.Reps,
			})
		}

		// Bulk insert new exercises
		if err := h.DB.Create(&planExercises).Error; err != nil {
			c.JSON(http.StatusInternalServerError, gin.H{"error": "failed to add new exercises to the plan", "details": err.Error()})
			return
		}
	}

	// Respond with the updated plan
	c.JSON(http.StatusOK, gin.H{
		"plan":      plan,
		"exercises": req.Exercises, // Return the updated exercises if provided
	})
}

func (h *PlanHandler) DeletePlan(c *gin.Context) {
	// Extract the plan_id from the URL parameters
	planID := c.Param("id")
	if planID == "" {
		c.JSON(http.StatusBadRequest, gin.H{"error": "plan_id is required"})
		return
	}

	// Fetch the existing plan to ensure it exists
	var plan models.WorkoutPlan
	if err := h.DB.First(&plan, "id = ?", planID).Error; err != nil {
		c.JSON(http.StatusNotFound, gin.H{"error": "workout plan not found", "details": err.Error()})
		return
	}

	// Delete the plan and its associated exercises
	if err := h.DB.Transaction(func(tx *gorm.DB) error {
		// Delete associated exercises
		if err := tx.Where("plan_id = ?", plan.ID).Delete(&models.PlanExercise{}).Error; err != nil {
			return err
		}

		// Delete the plan
		if err := tx.Delete(&plan).Error; err != nil {
			return err
		}

		return nil
	}); err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "failed to delete workout plan", "details": err.Error()})
		return
	}

	// Respond with success
	c.JSON(http.StatusOK, gin.H{"message": "workout plan deleted successfully"})
}
