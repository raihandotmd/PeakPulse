package routes

import (
	"github.com/gin-gonic/gin"

	"github.com/raihandotmd/peakPulse/internal/db"
	"github.com/raihandotmd/peakPulse/internal/handlers"
	"github.com/raihandotmd/peakPulse/internal/middleware"
)

// Setup registers all routes on the provided router
func Setup(r *gin.Engine) {
	initSupabase := db.InitSupabase()

	// auth
	authH := handlers.NewAuthHandler(initSupabase)
	r.POST("/signup", authH.Signup)
	r.POST("/login", authH.Login)

	// protected routes
	a := r.Group("/api/v1")
	a.Use(middleware.JWTMiddleware())
	{
		exH := handlers.NewExerciseHandler(initSupabase)
		a.GET("/exercises", exH.ListExercises)
		// TODO: add workout plans, schedules, sessions
	}
}
