package handlers

import (
	"net/http"

	"github.com/gin-gonic/gin"
	"github.com/raihandotmd/peakPulse/internal/models"
	"gorm.io/gorm"
)

// PlanHandler handles CRUD operations for workout plans
type PlanHandler struct{ DB *gorm.DB }

// NewPlanHandler returns a new PlanHandler
func NewPlanHandler(db *gorm.DB) *PlanHandler {
	return &PlanHandler{DB: db}
}

// CreatePlan godoc
// @Summary      Create a new workout plan
// @Description  Creates a workout plan for the authenticated user, with a unique title and a set of exercises.
// @Tags         plans
// @Accept       json
// @Produce      json
// @Param        payload  body      models.CreatePlanRequest  true  "Workout plan payload"
// @Success      201      {object}  models.WorkoutPlan        "Created plan"
// @Failure      400      {object}  gin.H                     "Invalid request or duplicate title"
// @Failure      401      {object}  gin.H                     "Unauthorized or missing user"
// @Failure      500      {object}  gin.H                     "Internal server error"
// @Security     BearerAuth
// @Router       /plans [post]
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
	var req models.CreatePlanRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "invalid payload", "details": err.Error()})
		return
	}

	// Check for duplicate title per user
	var existingPlan models.WorkoutPlan
	if err := h.DB.Where("user_id = ? AND title = ?", user.ID, req.Title).First(&existingPlan).Error; err == nil {
		c.JSON(http.StatusConflict, gin.H{"error": "workout plan with this title already exists"})
		return
	}

	// Create the workout plan
	plan := models.WorkoutPlan{
		UserID:      user.ID,
		Title:       req.Title,
		Description: req.Description,
	}
	if err := h.DB.Create(&plan).Error; err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "failed to create workout plan", "details": err.Error()})
		return
	}

	// Bulk insert plan exercises
	var planExercises []models.PlanExercise
	for _, ex := range req.Exercises {
		planExercises = append(planExercises, models.PlanExercise{
			PlanID:     plan.ID,
			ExerciseID: ex.ExerciseID,
			Sets:       ex.Sets,
			Reps:       ex.Reps,
		})
	}
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

// ListPlans godoc
// @Summary      List all workout plans
// @Description  Retrieves all workout plans for the authenticated user.
// @Tags         plans
// @Produce      json
// @Success      200  {object}  []models.WorkoutPlansListView  "List of plans"
// @Failure      500  {object}  gin.H                         "Internal server error"
// @Security     BearerAuth
// @Router       /plans [get]
func (h *PlanHandler) ListPlans(c *gin.Context) {
	var plans []models.WorkoutPlansListView
	if err := h.DB.Find(&plans).Error; err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "failed to fetch workout plans", "details": err.Error()})
		return
	}
	c.JSON(http.StatusOK, gin.H{"plans": plans})
}

// GetPlan godoc
// @Summary      Get a workout plan
// @Description  Retrieves a single workout plan by ID, including its exercises.
// @Tags         plans
// @Produce      json
// @Param        id   path      string  true  "Plan ID"
// @Success      200  {object}  models.WorkoutPlansListView  "The workout plan"
// @Failure      404  {object}  gin.H                         "Plan not found"
// @Failure      500  {object}  gin.H                         "Internal server error"
// @Security     BearerAuth
// @Router       /plans/{id} [get]
func (h *PlanHandler) GetPlan(c *gin.Context) {
	planID := c.Param("id")
	var plan []models.WorkoutPlansListView
	if err := h.DB.Where("plan_id = ?", planID).Find(&plan).Error; err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "failed to fetch workout plan", "details": err.Error()})
		return
	}
	if len(plan) == 0 {
		c.JSON(http.StatusNotFound, gin.H{"error": "workout plan not found"})
		return
	}
	c.JSON(http.StatusOK, gin.H{"plan": plan})
}

// UpdatePlan godoc
// @Summary      Update a workout plan
// @Description  Updates the title, description, and exercises of an existing workout plan.
// @Tags         plans
// @Accept       json
// @Produce      json
// @Param        id       path      string                   true  "Plan ID"
// @Param        payload  body      models.UpdatePlanRequest true  "Update payload"
// @Success      200      {object}  gin.H                     "Updated plan"
// @Failure      400      {object}  gin.H                     "Bad request"
// @Failure      404      {object}  gin.H                     "Plan not found"
// @Failure      500      {object}  gin.H                     "Internal server error"
// @Security     BearerAuth
// @Router       /plans/{id} [put]
func (h *PlanHandler) UpdatePlan(c *gin.Context) {
	planID := c.Param("id")
	if planID == "" {
		c.JSON(http.StatusBadRequest, gin.H{"error": "plan_id is required"})
		return
	}

	var req models.UpdatePlanRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "invalid payload", "details": err.Error()})
		return
	}

	var plan models.WorkoutPlan
	if err := h.DB.First(&plan, "id = ?", planID).Error; err != nil {
		c.JSON(http.StatusNotFound, gin.H{"error": "workout plan not found", "details": err.Error()})
		return
	}

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

	if req.Exercises != nil {
		if err := h.DB.Where("plan_id = ?", plan.ID).Delete(&models.PlanExercise{}).Error; err != nil {
			c.JSON(http.StatusInternalServerError, gin.H{"error": "failed to delete existing exercises", "details": err.Error()})
			return
		}
		var planExercises []models.PlanExercise
		for _, ex := range *req.Exercises {
			planExercises = append(planExercises, models.PlanExercise{
				PlanID:     plan.ID,
				ExerciseID: ex.ExerciseID,
				Sets:       ex.Sets,
				Reps:       ex.Reps,
			})
		}
		if err := h.DB.Create(&planExercises).Error; err != nil {
			c.JSON(http.StatusInternalServerError, gin.H{"error": "failed to add new exercises to the plan", "details": err.Error()})
			return
		}
	}

	c.JSON(http.StatusOK, gin.H{
		"plan":      plan,
		"exercises": req.Exercises,
	})
}

// DeletePlan godoc
// @Summary      Delete a workout plan
// @Description  Deletes a workout plan and all its associated exercises.
// @Tags         plans
// @Produce      json
// @Param        id   path      string  true  "Plan ID"
// @Success      200  {object}  gin.H    "Deletion confirmation"
// @Failure      400  {object}  gin.H    "Bad request"
// @Failure      404  {object}  gin.H    "Plan not found"
// @Failure      500  {object}  gin.H    "Internal server error"
// @Security     BearerAuth
// @Router       /plans/{id} [delete]
func (h *PlanHandler) DeletePlan(c *gin.Context) {
	planID := c.Param("id")
	if planID == "" {
		c.JSON(http.StatusBadRequest, gin.H{"error": "plan_id is required"})
		return
	}

	var plan models.WorkoutPlan
	if err := h.DB.First(&plan, "id = ?", planID).Error; err != nil {
		c.JSON(http.StatusNotFound, gin.H{"error": "workout plan not found", "details": err.Error()})
		return
	}

	if err := h.DB.Transaction(func(tx *gorm.DB) error {
		if err := tx.Where("plan_id = ?", plan.ID).Delete(&models.PlanExercise{}).Error; err != nil {
			return err
		}
		if err := tx.Delete(&plan).Error; err != nil {
			return err
		}
		return nil
	}); err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "failed to delete workout plan", "details": err.Error()})
		return
	}

	c.JSON(http.StatusOK, gin.H{"message": "workout plan deleted successfully"})
}
