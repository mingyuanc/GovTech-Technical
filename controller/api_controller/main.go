package apicontroller

import "gorm.io/gorm"

// Connection struct to contain the database connection
// Instead of using a global function
type Connection struct {
	db *gorm.DB
}

// Creates and returns a new pointer to a connection
func NewConnection(db *gorm.DB) *Connection {
	return &Connection{db: db}
}
