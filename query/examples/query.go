package examples

import (
	"gorm.io/gorm"
)

// User is a basic GORM model for query examples.
type User struct {
	gorm.Model
	Name string
	Age  int
}

// BasicQuery demonstrates basic query operations.
func BasicQuery(db *gorm.DB) {
	var user User

	// Get the first record
	db.First(&user)

	// Get one record, no specified order
	db.Take(&user)

	// Get last record
	db.Last(&user)

	// Get all records
	var users []User
	db.Find(&users)
}

// QueryWithConditions demonstrates querying with conditions.
func QueryWithConditions(db *gorm.DB) {
	var user User
	// Get first matched record
	db.Where("name = ?", "jinzhu").First(&user)

	// Get all matched records
	var users []User
	db.Where("name <> ?", "jinzhu").Find(&users)

	// IN
	db.Where("name IN ?", []string{"jinzhu", "jinzhu 2"}).Find(&users)

	// LIKE
	db.Where("name LIKE ?", "%jin%").Find(&users)

	// AND
	db.Where("name = ? AND age >= ?", "jinzhu", "22").Find(&users)
}
