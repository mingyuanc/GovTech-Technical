package models

// Represents the required json format for the register endpoint
type RetrieveNotificationBody struct {
	Teacher      string `json:"teacher" binding:"required,email"`
	Notification string `json:"notification"`
}
