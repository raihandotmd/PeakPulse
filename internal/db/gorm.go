package db

import (
	"log"
	"os"

	"github.com/raihandotmd/peakPulse/internal/config"
	"gorm.io/driver/postgres"
	"gorm.io/gorm"
)

var DB *gorm.DB

// InitGorm initializes GORM with the DSN from env
func InitGorm() {
	config.Load()
	dsn := os.Getenv("DB_DSN")

	var err error
	DB, err = gorm.Open(postgres.Open(dsn), &gorm.Config{})
	if err != nil {
		log.Fatalf("failed to connect database: %v", err)
	}
}

// GetDB returns the GORM DB instance
func GetGormDB() *gorm.DB {
	if DB == nil {
		log.Fatal("GORM DB is not initialized")
	}
	return DB
}
