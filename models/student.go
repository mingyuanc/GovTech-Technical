package models

import "gorm.io/gorm"

// Represents the stored student object
type Student struct {
	gorm.Model
	Email       string     `json:"email" gorm:"unique;not null"`
	IsSuspended *bool      `json:"isSuspended" gorm:"default:false"`
	Teachers    []*Teacher `json:"teachers" gorm:"many2many:teacher_student;"`
}
