package apicontroller

import (
	"log"

	"gorm.io/gorm"
)

// Connection struct to contain the database connection
// Instead of using a global function
type Connection struct {
	db     *gorm.DB
	logger *log.Logger
}

// Creates and returns a new pointer to a connection
func NewConnection(db *gorm.DB) *Connection {
	return &Connection{db: db, logger: log.Default()}
}
