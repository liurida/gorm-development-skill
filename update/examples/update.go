package examples

import (
	"gorm.io/gorm"
)

// User is a basic GORM model for update examples.
type User struct {
	gorm.Model
	Name   string
	Age    int
	Active bool
}

// UpdateSingleColumn demonstrates updating a single column.
func UpdateSingleColumn(db *gorm.DB, user *User) {
	db.Model(user).Update("name", "hello")
}

// UpdateMultipleColumns demonstrates updating multiple columns with a struct.
func UpdateMultipleColumns(db *gorm.DB, user *User) {
	db.Model(user).Updates(User{Name: "hello", Age: 18})
}

// UpdateWithMap demonstrates updating with a map.
func UpdateWithMap(db *gorm.DB, user *User) {
	db.Model(user).Updates(map[string]interface{}{"name": "hello", "age": 18, "active": false})
}

// UpdateSelectedFields demonstrates updating only selected fields.
func UpdateSelectedFields(db *gorm.DB, user *User) {
	db.Model(user).Select("name").Updates(map[string]interface{}{"name": "hello", "age": 18, "active": false})
}
