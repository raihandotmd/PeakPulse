package handlers

import (
	"net/http"

	"github.com/gin-gonic/gin"
	"github.com/raihandotmd/peakPulse/internal/db"
	"github.com/raihandotmd/peakPulse/internal/models"
	supa "github.com/supabase-community/supabase-go"
)

// ExerciseHandler handles exercise-related endpoints
type ExerciseHandler struct {
	Supa *supa.Client
}

// NewExerciseHandler creates an ExerciseHandler
func NewExerciseHandler(supaCli *supa.Client) *ExerciseHandler {
	return &ExerciseHandler{Supa: supaCli}
}

// ListExercises returns all exercises
func (h *ExerciseHandler) ListExercises(c *gin.Context) {
	DB := db.GetGormClient()

	// Use the Exercise model to fetch all exercises
	var exercises []models.Exercise
	if err := DB.Find(&exercises).Error; err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to fetch exercises"})
		return
	}

	// Get the count of exercises
	var count int64
	if err := DB.Model(&models.Exercise{}).Count(&count).Error; err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to fetch exercises count"})
		return
	}

	c.JSON(http.StatusOK, gin.H{"data": exercises, "count": count})
}
