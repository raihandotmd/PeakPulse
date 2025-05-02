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

// ListExercises godoc
// @Summary      List all exercises
// @Description  Returns all exercises in the catalog along with the total count.
// @Tags         exercises
// @Produce      json
// @Success      200  {object}  models.ExerciseListResponse  "Successful response"
// @Failure      500  {object}  map[string]string             "Internal server error"
// @Router       /exercises [get]
func (h *ExerciseHandler) ListExercises(c *gin.Context) {
	DB := db.GetGormClient()

	var exercises []models.Exercise
	if err := DB.Find(&exercises).Error; err != nil {
		c.JSON(http.StatusInternalServerError, map[string]string{"error": "Failed to fetch exercises"})
		return
	}

	var count int64
	if err := DB.Model(&models.Exercise{}).Count(&count).Error; err != nil {
		c.JSON(http.StatusInternalServerError, map[string]string{"error": "Failed to fetch exercises count"})
		return
	}

	// Respond using the named response struct
	c.JSON(http.StatusOK, models.ExerciseListResponse{
		Data:  exercises,
		Count: count,
	})
}
