package examples

import (
	"database/sql"

	"gorm.io/gorm"
)

// Result is a struct to hold raw query results.
type Result struct {
	ID   int
	Name string
	Age  int
}

// RawQueryWithScan demonstrates using raw SQL with Scan.
func RawQueryWithScan(db *gorm.DB) {
	var result Result
	db.Raw("SELECT id, name, age FROM users WHERE id = ?", 3).Scan(&result)
}

// RawExec demonstrates executing raw SQL.
func RawExec(db *gorm.DB) {
	db.Exec("DROP TABLE users")
}

// NamedArguments demonstrates using named arguments in raw SQL.
func NamedArguments(db *gorm.DB) {
	db.Raw("SELECT * FROM users WHERE name = @name", sql.Named("name", "jinzhu")).Scan(&Result{})
}

// DryRun demonstrates the DryRun mode for generating SQL without execution.
func DryRun(db *gorm.DB) string {
	stmt := db.Session(&gorm.Session{DryRun: true}).First(&User{}, 1).Statement
	return stmt.SQL.String()
}

// User is a basic GORM model for raw_sql examples.
type User struct {
	gorm.Model
	Name string
}
