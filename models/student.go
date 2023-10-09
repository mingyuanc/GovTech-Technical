package models

import "gorm.io/gorm"

type Student struct {
	gorm.Model
	Email    string     `json:"email" gorm:"unique;not null"`
	Teachers []*Teacher `json:"teachers" gorm:"many2many:teacher_student;"`
}