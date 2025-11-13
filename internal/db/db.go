package db

import "gorm.io/gorm"

// package db provides a tiny global holder for the *gorm.DB instance used by
// handlers. It's intentionally minimal so imports like "admin-api/internal/db"
// compile and allow the main program to set the DB connection.

var conn *gorm.DB

// Set stores the gorm DB connection for use by other packages.
func Set(d *gorm.DB) {
	conn = d
}

// Get returns the stored *gorm.DB connection (may be nil if not set).
func Get() *gorm.DB {
	return conn
}
