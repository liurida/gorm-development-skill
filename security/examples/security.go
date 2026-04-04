package examples

import (
	"fmt"
	"strconv"

	"gorm.io/gorm"
)

// User is a basic GORM model for security examples.
type User struct {
	gorm.Model
	Name  string
	Email string
}

// SafeQueryWithPlaceholder demonstrates the safe way to query with user input.
// Always use parameterized queries to prevent SQL injection.
func SafeQueryWithPlaceholder(db *gorm.DB, userInput string) (*User, error) {
	var user User
	// Safe: userInput is passed as an argument and will be escaped
	result := db.Where("name = ?", userInput).First(&user)
	if result.Error != nil {
		return nil, fmt.Errorf("failed to find user: %w", result.Error)
	}
	return &user, nil
}

// SafeInlineCondition demonstrates safe inline conditions.
func SafeInlineCondition(db *gorm.DB, userInput string) (*User, error) {
	var user User
	// Safe: parameterized inline condition
	result := db.First(&user, "name = ?", userInput)
	if result.Error != nil {
		return nil, fmt.Errorf("failed to find user: %w", result.Error)
	}
	return &user, nil
}

// SafeNumericID demonstrates how to safely handle numeric IDs from user input.
// Always validate and convert user input to the expected type.
func SafeNumericID(db *gorm.DB, userInputID string) (*User, error) {
	// Validate that the input is actually a number
	id, err := strconv.Atoi(userInputID)
	if err != nil {
		return nil, fmt.Errorf("invalid user ID format: %w", err)
	}

	var user User
	// Safe: id is now a validated integer
	result := db.First(&user, id)
	if result.Error != nil {
		return nil, fmt.Errorf("failed to find user: %w", result.Error)
	}
	return &user, nil
}

// SafeMultipleConditions demonstrates safe queries with multiple conditions.
func SafeMultipleConditions(db *gorm.DB, name, email string) ([]User, error) {
	var users []User
	// Safe: all user inputs are passed as arguments
	result := db.Where("name = ? AND email = ?", name, email).Find(&users)
	if result.Error != nil {
		return nil, fmt.Errorf("failed to find users: %w", result.Error)
	}
	return users, nil
}

// SafeMapConditions demonstrates using a map for conditions.
// Maps are automatically escaped by GORM.
func SafeMapConditions(db *gorm.DB, name string) ([]User, error) {
	var users []User
	// Safe: map values are escaped
	result := db.Where(map[string]interface{}{"name": name}).Find(&users)
	if result.Error != nil {
		return nil, fmt.Errorf("failed to find users: %w", result.Error)
	}
	return users, nil
}

// SafeStructConditions demonstrates using a struct for conditions.
// Struct fields are automatically used as conditions.
func SafeStructConditions(db *gorm.DB, name string) ([]User, error) {
	var users []User
	// Safe: struct fields are escaped
	result := db.Where(&User{Name: name}).Find(&users)
	if result.Error != nil {
		return nil, fmt.Errorf("failed to find users: %w", result.Error)
	}
	return users, nil
}

// ValidatedOrderBy demonstrates how to safely handle ORDER BY with user input.
// Use a whitelist approach for column names since Order() does not escape input.
func ValidatedOrderBy(db *gorm.DB, orderField string) ([]User, error) {
	// Whitelist of allowed order fields
	allowedFields := map[string]bool{
		"name":       true,
		"email":      true,
		"created_at": true,
		"updated_at": true,
	}

	if !allowedFields[orderField] {
		return nil, fmt.Errorf("invalid order field: %s", orderField)
	}

	var users []User
	// Safe: orderField has been validated against whitelist
	result := db.Order(orderField).Find(&users)
	if result.Error != nil {
		return nil, fmt.Errorf("failed to find users: %w", result.Error)
	}
	return users, nil
}

// ValidatedTableName demonstrates how to safely handle table names from user input.
// Use a whitelist approach since Table() does not escape input.
func ValidatedTableName(db *gorm.DB, tableName string) (int64, error) {
	// Whitelist of allowed table names
	allowedTables := map[string]bool{
		"users":    true,
		"products": true,
		"orders":   true,
	}

	if !allowedTables[tableName] {
		return 0, fmt.Errorf("invalid table name: %s", tableName)
	}

	var count int64
	// Safe: tableName has been validated against whitelist
	result := db.Table(tableName).Count(&count)
	if result.Error != nil {
		return 0, fmt.Errorf("failed to count records: %w", result.Error)
	}
	return count, nil
}
