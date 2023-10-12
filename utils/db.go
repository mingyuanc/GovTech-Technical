package utils

import (
	"errors"
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

	// dsn := "postgres://pg:pg@localhost:5432/pg"
	// db, err := gorm.Open(postgres.Open(dsn), &gorm.Config{})
	if err != nil {
		panic(err)
	}
	db.AutoMigrate(&models.Teacher{}, &models.Student{})
	return db
}

// Checks if specified teacher is present in the database
func IsStudentPresent(db *gorm.DB, studentEmail string) bool {
	var student models.Student
	err := db.Where(&models.Student{Email: studentEmail}).First(&student).Error
	if errors.Is(err, gorm.ErrRecordNotFound) {
		return false
	}
	if err != nil {
		panic(err)
	}
	return true
}

// Checks if specified teacher is present in the database
func IsTeacherPresent(db *gorm.DB, teacherEmail string) bool {
	var teacher models.Teacher
	err := db.Where(&models.Teacher{Email: teacherEmail}).First(&teacher).Error
	if errors.Is(err, gorm.ErrRecordNotFound) {
		return false
	}
	if err != nil {
		panic(err)
	}
	return true
}

// Checks if specified student is suspended or not
func IsStudentSuspended(db *gorm.DB, studentEmail string) bool {
	var student models.Student
	err := db.Where(&models.Student{Email: studentEmail}).First(&student).Error
	if err != nil {
		panic(err)
	}
	return *student.IsSuspended
}

// Returns an array of unique students from a given list of teachers
func GetCommonStudentFromTeachersEmail(db *gorm.DB, teachers []string) ([]models.Student, error) {

	var students []models.Student

	err := db.
		Table("students").
		Joins("JOIN teacher_student ON students.id = teacher_student.student_id").
		Joins("JOIN teachers ON teachers.id = teacher_student.teacher_id").
		Where("teachers.email IN ?", teachers).
		Group("students.id").
		Having("COUNT(DISTINCT teachers.email) = ?", len(teachers)).
		Find(&students).Error
	return students, err
}

// Returns students from a specific teacher
func GetStudentFromTeacher(db *gorm.DB, teacherEmail string) ([]*models.Student, error) {
	var teacher models.Teacher
	err := db.Preload("Students").Where("email = ?", teacherEmail).First(&teacher).Error
	return teacher.Students, err
}

func SuspendStudent(db *gorm.DB, studentEmail string) error {
	return db.Model(&models.Student{}).Where("email = ?", studentEmail).Update("is_suspended", true).Error
}
