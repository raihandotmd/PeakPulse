package main

import (
	"net/http"

	"github.com/gin-gonic/gin"
	"github.com/raihandotmd/peakPulse/internal/db"
	"github.com/raihandotmd/peakPulse/internal/routes"
)

// @title PeakPulse API
// @version 1.0
// @description This is the API documentation for PeakPulse.
// @host localhost:8585
// @BasePath /
// @schemes http
// @contact.name API Support
// @contact.url http://www.example.com/support
// @contact.email support@example.com

func setupRouter() *gin.Engine {
	r := gin.Default()

	// Ping test
	r.GET("/ping", func(c *gin.Context) {
		c.String(http.StatusOK, "pong")
	})

	return r
}

func main() {
	db.InitSupabase()
	db.InitGorm()

	r := setupRouter()

	routes.Setup(r)

	r.Run(":8585")

}
