package examples

import (
	"gorm.io/gorm"
	"gorm.io/gorm/clause"
)

// User is a basic GORM model for advanced query examples.
type User struct {
	gorm.Model
	Name    string
	Age     int
	Role    string
	Company Company
}

// Company represents a user's company.
type Company struct {
	gorm.Model
	Name   string
	UserID uint
}

// Order represents an order for scope examples.
type Order struct {
	gorm.Model
	Amount      float64
	PayModeSign string
	Status      string
}

// SmartSelectExample demonstrates automatic field selection.
func SmartSelectExample(db *gorm.DB) {
	type APIUser struct {
		ID   uint
		Name string
	}

	var apiUsers []APIUser
	db.Model(&User{}).Limit(10).Find(&apiUsers)
	// SQL: SELECT `id`, `name` FROM `users` LIMIT 10
}

// LockingExample demonstrates different locking strategies.
func LockingExample(db *gorm.DB) {
	var users []User

	// FOR UPDATE lock
	db.Clauses(clause.Locking{Strength: "UPDATE"}).Find(&users)

	// FOR SHARE lock
	db.Clauses(clause.Locking{Strength: "SHARE"}).Find(&users)

	// NOWAIT option
	db.Clauses(clause.Locking{Strength: "UPDATE", Options: "NOWAIT"}).Find(&users)

	// SKIP LOCKED for high concurrency
	db.Clauses(clause.Locking{Strength: "UPDATE", Options: "SKIP LOCKED"}).Find(&users)
}

// SubQueryExample demonstrates using subqueries.
func SubQueryExample(db *gorm.DB) {
	var orders []Order

	// Simple subquery - find orders above average
	db.Where("amount > (?)", db.Table("orders").Select("AVG(amount)")).Find(&orders)

	// FROM subquery
	var users []User
	db.Table("(?) as u", db.Model(&User{}).Select("name", "age")).Where("age = ?", 18).Find(&users)
}

// GroupConditionsExample demonstrates nested conditions.
func GroupConditionsExample(db *gorm.DB) {
	var users []User

	// Complex nested conditions
	db.Where(
		db.Where("role = ?", "admin").Where(db.Where("age >= ?", 18).Or("verified = ?", true)),
	).Or(
		db.Where("role = ?", "superuser"),
	).Find(&users)
}

// ScopeExample demonstrates reusable query scopes.
func AmountGreaterThan1000(db *gorm.DB) *gorm.DB {
	return db.Where("amount > ?", 1000)
}

func PaidWithCreditCard(db *gorm.DB) *gorm.DB {
	return db.Where("pay_mode_sign = ?", "C")
}

func OrderStatus(statuses []string) func(db *gorm.DB) *gorm.DB {
	return func(db *gorm.DB) *gorm.DB {
		return db.Where("status IN ?", statuses)
	}
}

func ScopesExample(db *gorm.DB) {
	var orders []Order

	// Combine multiple scopes
	db.Scopes(AmountGreaterThan1000, PaidWithCreditCard).Find(&orders)

	// Scope with parameters
	db.Scopes(AmountGreaterThan1000, OrderStatus([]string{"paid", "shipped"})).Find(&orders)
}

// FindInBatchesExample demonstrates batch processing.
func FindInBatchesExample(db *gorm.DB) {
	var results []User

	db.Where("age > ?", 18).FindInBatches(&results, 100, func(tx *gorm.DB, batch int) error {
		for i := range results {
			results[i].Role = "processed"
		}
		tx.Save(&results)
		return nil
	})
}

// CountExample demonstrates counting records.
func CountExample(db *gorm.DB) {
	var count int64

	// Basic count
	db.Model(&User{}).Where("name = ?", "jinzhu").Count(&count)

	// Count with distinct
	db.Model(&User{}).Distinct("name").Count(&count)

	// Count with group
	db.Model(&User{}).Group("role").Count(&count)
}

// PluckExample demonstrates extracting single columns.
func PluckExample(db *gorm.DB) {
	var ages []int64
	db.Model(&User{}).Pluck("age", &ages)

	var names []string
	db.Model(&User{}).Distinct().Pluck("name", &names)
}
