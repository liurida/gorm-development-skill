package examples

import (
	"gorm.io/gorm"
)

// User is a basic GORM model for session examples.
type User struct {
	gorm.Model
	Name string
	Age  int
}

// DryRunSession demonstrates using DryRun mode to simulate queries without executing them.
func DryRunSession(db *gorm.DB) *gorm.DB {
	return db.Session(&gorm.Session{DryRun: true})
}

// PreparedStatementSession demonstrates using PrepareStmt mode for prepared statement caching.
func PreparedStatementSession(db *gorm.DB) *gorm.DB {
	return db.Session(&gorm.Session{PrepareStmt: true})
}

// NewDBSession demonstrates creating a new DB session without inherited conditions.
func NewDBSession(db *gorm.DB) *gorm.DB {
	return db.Session(&gorm.Session{NewDB: true})
}

// SkipHooksSession demonstrates creating a session that skips hooks.
func SkipHooksSession(db *gorm.DB) *gorm.DB {
	return db.Session(&gorm.Session{SkipHooks: true})
}

// AllowGlobalUpdateSession demonstrates creating a session that allows global updates.
func AllowGlobalUpdateSession(db *gorm.DB) *gorm.DB {
	return db.Session(&gorm.Session{AllowGlobalUpdate: true})
}

// ContinuousSessionExample demonstrates using a session for multiple operations.
func ContinuousSessionExample(db *gorm.DB) {
	tx := db.Session(&gorm.Session{SkipDefaultTransaction: true})

	var user User
	tx.First(&user, 1)

	var users []User
	tx.Find(&users)

	tx.Model(&user).Update("Age", 18)
}
