package utils

import (
	"github.com/mingyuanc/GovTech-Technical/models"
	"gorm.io/driver/postgres"
	"gorm.io/gorm"
)

func Connect() *gorm.DB {
	dsn := "host=localhost user=pg password=pg dbname=pg port=5432"
	db, err := gorm.Open(postgres.Open(dsn), &gorm.Config{})
	if err != nil {
		panic(err)
	}
	db.AutoMigrate(&models.Teacher{}, &models.Student{})
	return db
}
