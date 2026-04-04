
package examples

import (
	"context"
	"errors"
	"log"
	"os"
	"time"

	"gorm.io/driver/sqlite"
	"gorm.io/gorm"
	"gorm.io/gorm/clause"
	"gorm.io/gorm/logger"
)

// User model for performance examples
type User struct {
	gorm.Model
	Name     string
	Age      int
	Email    string `gorm:"uniqueIndex"`
	IsAdmin  bool
	CompanyID uint
}

// APIUser is a smaller struct for demonstrating smart select fields.
type APIUser struct {
	ID   uint
	Name string
}

// --- Disable Default Transaction ---

// DisableDefaultTx demonstrates initializing GORM without default transactions for writes.
// This can improve performance by about 30% for single write operations.
func DisableDefaultTx(dsn string) (*gorm.DB, error) {
	return gorm.Open(sqlite.Open(dsn), &gorm.Config{
		SkipDefaultTransaction: true,
	})
}

// --- Prepared Statement Caching ---

// EnablePreparedStmtCache demonstrates caching prepared statements for faster subsequent calls.
func EnablePreparedStmtCache(dsn string) (*gorm.DB, error) {
	return gorm.Open(sqlite.Open(dsn), &gorm.Config{
		PrepareStmt: true,
	})
}

// UsePreparedStmtInSession shows using prepared statements in a session.
func UsePreparedStmtInSession(db *gorm.DB) {
	tx := db.Session(&gorm.Session{PrepareStmt: true})

	var user User
	tx.First(&user, 1)

	var users []User
	tx.Find(&users, "age > ?", 30)
}

// --- Select Fields ---

// SelectSpecificFields demonstrates selecting only the required fields.
func SelectSpecificFields(db *gorm.DB) ([]User, error) {
	var users []User
	// Avoids fetching all fields of the User struct
	err := db.Select("name", "age").Find(&users).Error
	return users, err
}

// SmartSelectFields demonstrates using a smaller struct to select fields automatically.
func SmartSelectFields(db *gorm.DB) ([]APIUser, error) {
	var apiUsers []APIUser
	// GORM automatically selects only `id` and `name` from the `users` table
	err := db.Model(&User{}).Find(&apiUsers).Error
	return apiUsers, err
}

// --- Iteration & Batch Processing ---

// FindInBatches demonstrates processing records in batches.
func FindInBatches(db *gorm.DB) error {
	var results []User

	// Process 100 records at a time
	result := db.Where("age > ?", 20).FindInBatches(&results, 100, func(tx *gorm.DB, batch int) error {
		fmt.Printf("Processing batch %d with %d users\n", batch, len(results))

		// Do something with the batch of users
		for _, user := range results {
			tx.Model(&user).Update("age", user.Age+1)
		}

		return nil // continue processing
	})

	return result.Error
}

// RowsIteration demonstrates processing records one by one using Rows.
func RowsIteration(db *gorm.DB) error {
	rows, err := db.Model(&User{}).Where("is_admin = ?", true).Rows()
	if err != nil {
		return err
	}
	defer rows.Close()

	for rows.Next() {
		var user User
		// ScanRows is a GORM method to scan a row into a struct.
		db.ScanRows(rows, &user)

		// Process user
		fmt.Printf("Processing admin user: %s\n", user.Name)
	}

	return nil
}

// --- Index Hints ---

// UseIndexHint demonstrates forcing the query optimizer to use a specific index.
func UseIndexHint(db *gorm.DB, username string) (*User, error) {
	var user User
	// Assumes an index `idx_user_name` exists on the `name` column
	err := db.Clauses(hints.UseIndex("idx_user_name")).Where("name = ?", username).First(&user).Error
	return &user, err
}

// --- Read/Write Splitting (DBResolver) ---

// This is a conceptual example. For a full implementation, see the `dbresolver` skill.
// func UseReadWriteSplitting(dsn string) (*gorm.DB, error) {
// 	db, err := gorm.Open(sqlite.Open(dsn), &gorm.Config{})
// 	if err != nil {
// 		return nil, err
// 	}
//
// 	db.Use(dbresolver.Register(dbresolver.Config{
// 		Sources:  []gorm.Dialector{mysql.Open("write_dsn")},
// 		Replicas: []gorm.Dialector{mysql.Open("read_dsn_1"), mysql.Open("read_dsn_2")},
// 		Policy:   dbresolver.RandomPolicy{},
// 	}))
//
// 	return db, nil
// }

// --- DryRun Mode for Debugging ---

// DryRunQuery generates the SQL for a query without executing it.
func DryRunQuery(db *gorm.DB) string {
	var user User
	stmt := db.Session(&gorm.Session{DryRun: true}).First(&user, 1).Statement
	return stmt.SQL.String()
}

// --- Efficient Updates ---

// UpdateWithMap demonstrates updating with a map for specific fields.
func UpdateWithMap(db *gorm.DB, userID uint, newAge int) error {
	// Only updates `age` and `updated_at`
	// Doesn't trigger hooks/time tracking for other fields
	return db.Model(&User{}).Where("id = ?", userID).Update("age", newAge).Error
}

// UpdateWithExpr demonstrates using SQL expressions for updates.
func UpdateWithExpr(db *gorm.DB, userID uint) error {
	// Server-side calculation, avoids fetching the record first
	return db.Model(&User{}).Where("id = ?", userID).Update("age", gorm.Expr("age + ?", 1)).Error
}

// --- Performance Best Practices Example ---

// GetActiveAdminUsers shows a combination of performance techniques.
func GetActiveAdminUsers(db *gorm.DB) ([]APIUser, error) {
	var apiUsers []APIUser

	// - Uses prepared statements (if enabled globally)
	// - Uses smart select for `APIUser`
	// - Assumes an index on `is_admin` for efficient filtering
	err := db.Model(&User{}).Where(&User{IsAdmin: true}).Find(&apiUsers).Error
	if err != nil {
		return nil, err
	}

	return apiUsers, nil
}

// --- Avoid Leaking Goroutines in Transactions ---

// CorrectTxErrorHandling ensures context cancellation is handled in transactions.
func CorrectTxErrorHandling(ctx context.Context, db *gorm.DB, userID uint) error {
	return db.WithContext(ctx).Transaction(func(tx *gorm.DB) error {
		var user User

		// If context is cancelled, this will return an error
		if err := tx.First(&user, userID).Error; err != nil {
			return err // Rollback
		}

		// Check context error before next operation
		if ctx.Err() != nil {
			return ctx.Err() // Rollback
		}

		if err := tx.Model(&user).Update("age", user.Age+1).Error; err != nil {
			return err // Rollback
		}

		return nil // Commit
	})
}
