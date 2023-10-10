package models

import (
	"gorm.io/gorm"
)

type Teacher struct {
	gorm.Model
	Email    string     `json:"email" binding:"required,email" gorm:"unique;not null"`
	Students []*Student `json:"students" gorm:"many2many:teacher_student;"`
}
