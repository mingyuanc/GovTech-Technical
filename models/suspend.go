package models

// Represents the required json format for the suspend
type SuspendBody struct {
	Student []*string `json:"student" binding:"required,email"`
}
