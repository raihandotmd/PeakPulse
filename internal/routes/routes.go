package routes

import (
	"github.com/gin-gonic/gin"

	"github.com/raihandotmd/peakPulse/internal/db"
	"github.com/raihandotmd/peakPulse/internal/handlers"
	"github.com/raihandotmd/peakPulse/internal/middleware"
)

// Setup godoc
// @Summary Setup all routes
// @Description Initialize all API routes and their handlers
// @Tags routes
// @Router / [get]
// Setup registers all routes on the provided router
func Setup(r *gin.Engine) {
	initSupabase := db.InitSupabase()
	initGorm := db.InitGorm()

	// auth
	authH := handlers.NewAuthHandler(initSupabase)
	auth := r.Group("/auth")
	{
		auth.POST("/signup", authH.Signup)
		auth.POST("/login", authH.Login)
	}

	// plan routes
	planH := handlers.NewPlanHandler(initGorm)

	// protected routes
	a := r.Group("/api/v1")
	a.Use(middleware.JWTMiddleware())
	{
		exH := handlers.NewExerciseHandler(initSupabase)
		a.GET("/exercises", exH.ListExercises)

		// TODO: add workout plans, schedules, sessions
		a.POST("/plans", planH.CreatePlan)
		a.GET("/plans", planH.ListPlans)
		a.GET("/plans/:id", planH.GetPlan)
		a.PUT("/plans/:id", planH.UpdatePlan)
		a.DELETE("/plans/:id", planH.DeletePlan)
	}
}
