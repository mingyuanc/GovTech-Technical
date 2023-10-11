package models

// Represents the required json format for the register endpoint
type RegisterBody struct {
	Teacher  *string   `json:"teacher" binding:"required,email"`
	Students []*string `json:"students" binding:"required,dive,email"`
}
