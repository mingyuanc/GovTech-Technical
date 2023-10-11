package utils

import (
	"errors"

	"github.com/mingyuanc/GovTech-Technical/models"
	"gorm.io/driver/postgres"
	"gorm.io/gorm"
)

// Connects to the db and return a connection variable
func Connect() *gorm.DB {
	// TODO use env
	// dsn := os.Getenv("DATABASE_URL")
	// db, err := gorm.Open(postgres.Open(dsn), &gorm.Config{})

	dsn := "postgres://pg:pg@localhost:5432/pg"
	db, err := gorm.Open(postgres.Open(dsn), &gorm.Config{})
	if err != nil {
		panic(err)
	}
	db.AutoMigrate(&models.Teacher{}, &models.Student{})
	return db
}

// Return an array of unique students from a given list of teachers
func GetUniqueStudentFromTeachersEmail(db *gorm.DB, teachers []string) ([]models.Student, error) {

	uniqueStudent := make(map[string]models.Student)
	// Get unique student frm each teacher
	for _, teacher := range teachers {
		students, err := GetStudentFromTeacher(db, teacher)
		if err != nil {
			return nil, errors.New("Unable to find a teacher with email: " + teacher)
		}
		for _, student := range students {
			println(student.Email)
			uniqueStudent[student.Email] = student
		}
	}

	// Get back an array
	var students []models.Student
	for _, stuArr := range uniqueStudent {
		students = append(students, stuArr)
	}
	return students, nil
}

// return students from a specific teacher
func GetStudentFromTeacher(db *gorm.DB, teacherEmail string) ([]models.Student, error) {
	var teacher models.Teacher
	err := db.Preload("Students").Where("email = ?", teacherEmail).First(&teacher).Error
	var students []models.Student
	if err != nil {
		return nil, err
	}
	err = db.Model(&teacher).Association("Students").Find(&students)
	return students, err
}
