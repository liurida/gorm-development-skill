
package examples

import (
	"fmt"

	"gorm.io/gorm"
)

// ChainingUser model for method chaining examples
type ChainingUser struct {
	gorm.Model
	Name string
	Age  int
}

// BasicChaining demonstrates a simple chain of GORM methods.
func BasicChaining(db *gorm.DB) ([]ChainingUser, error) {
	var users []ChainingUser
	// Chain of methods: Model -> Where -> Order -> Limit -> Offset -> Find
	err := db.Model(&ChainingUser{}).
		Where("age > ?", 20).
		Order("name asc").
		Limit(10).
		Offset(5).
		Find(&users).Error

	return users, err
}

// UnsafeReuse demonstrates the incorrect way to reuse a GORM query object.
func UnsafeReuse(db *gorm.DB) {
	// Create a base query
	query := db.Model(&ChainingUser{}).Where("name LIKE ?", "user%")

	// First use of the query object
	var users1 []ChainingUser
	query.Where("age < ?", 30).Find(&users1)
	fmt.Printf("Unsafe 1 - Found %d users under 30\n", len(users1))

	// --- THIS IS INCORRECT ---
	// Second use of the same query object. The `age < 30` condition from the
	// previous query will leak into this one.
	var users2 []ChainingUser
	query.Where("age > ?", 25).Find(&users2)
	fmt.Printf("Unsafe 2 - Found %d users over 25\n", len(users2))

	// The actual SQL for the second query will be something like:
	// SELECT * FROM users WHERE name LIKE 'user%' AND age < 30 AND age > 25
}

// SafeReuse demonstrates the correct way to reuse a GORM query object using a new session.
func SafeReuse(db *gorm.DB) {
	// Create a base query and immediately create a new session from it.
	// This makes the baseQuery object safe for reuse.
	baseQuery := db.Model(&ChainingUser{}).Where("name LIKE ?", "user%").Session(&gorm.Session{})

	// First use of the base query
	var users1 []ChainingUser
	// Conditions added here are local to this chain
	baseQuery.Where("age < ?", 30).Find(&users1)
	fmt.Printf("Safe 1 - Found %d users under 30\n", len(users1))

	// Second use of the base query
	var users2 []ChainingUser
	// Conditions here are independent from the first query
	baseQuery.Where("age > ?", 25).Find(&users2)
	fmt.Printf("Safe 2 - Found %d users over 25\n", len(users2))

	// The SQL for the second query will be correct:
	// SELECT * FROM users WHERE name LIKE 'user%' AND age > 25
}

// DynamicQueryBuilding demonstrates building a query dynamically through chaining.
func DynamicQueryBuilding(db *gorm.DB, nameFilter string, minAge int) ([]ChainingUser, error) {
	query := db.Model(&ChainingUser{})

	if nameFilter != "" {
		// Add a WHERE clause if the filter is provided
		query = query.Where("name = ?", nameFilter)
	}

	if minAge > 0 {
		// Add another WHERE clause if the age is provided
		query = query.Where("age > ?", minAge)
	}

	var users []ChainingUser
	err := query.Find(&users).Error
	return users, err
}
