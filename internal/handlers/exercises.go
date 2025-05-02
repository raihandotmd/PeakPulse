package handlers

import (
	"net/http"

	"github.com/gin-gonic/gin"
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
	data, count, err := h.Supa.From("exercises").Select("*", "", false).Execute()
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}
	c.JSON(http.StatusOK, gin.H{"data": data, "count": count})
}
