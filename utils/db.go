package utils

import (
	"github.com/mingyuanc/GovTech-Technical/models"
	"gorm.io/driver/postgres"
	"gorm.io/gorm"
)

// Connects to the db and return a connection variable
func Connect() *gorm.DB {
	// TODO use env
	dsn := "host=localhost user=pg password=pg dbname=pg port=5432"
	db, err := gorm.Open(postgres.Open(dsn), &gorm.Config{})
	if err != nil {
		panic(err)
	}
	db.AutoMigrate(&models.Teacher{}, &models.Student{})
	return db
}

// Connects to the test db and return a connection variable
func ConnectTest() *gorm.DB {
	dsn := "host=localhost user=pg-test password=pg-test dbname=pg-test port=5433"
	db, err := gorm.Open(postgres.Open(dsn), &gorm.Config{})
	if err != nil {
		panic(err)
	}
	db.AutoMigrate(&models.Teacher{}, &models.Student{})
	return db
}
