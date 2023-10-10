package models

type Register struct {
	Teacher  *string   `json:"teacher" binding:"required,email"`
	Students []*string `json:"students" binding:"required,dive,email"`
}
