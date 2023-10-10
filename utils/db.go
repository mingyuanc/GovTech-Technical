package utils

import (
	"os"

	"github.com/mingyuanc/GovTech-Technical/models"
	"gorm.io/driver/postgres"
	"gorm.io/gorm"
)

// Connects to the db and return a connection variable
func Connect() *gorm.DB {
	// TODO use env
	dsn := os.Getenv("DATABASE_URL")
	db, err := gorm.Open(postgres.Open(dsn), &gorm.Config{})
	if err != nil {
		panic(err)
	}
	db.AutoMigrate(&models.Teacher{}, &models.Student{})
	return db
}
