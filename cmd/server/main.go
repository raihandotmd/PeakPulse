package main

import (
	"net/http"

	"github.com/gin-gonic/gin"
	"github.com/raihandotmd/peakPulse/internal/db"
	"github.com/raihandotmd/peakPulse/internal/routes"
)

func setupRouter() *gin.Engine {
	// Disable Console Color
	// gin.DisableConsoleColor()
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
