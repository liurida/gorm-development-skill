package examples

import (
	"time"

	"gorm.io/gorm"
)

// User is a basic GORM model for create examples.
type User struct {
	gorm.Model
	Name     string
	Age      int
	Birthday time.Time
}

// CreateRecord demonstrates creating a single record.
func CreateRecord(db *gorm.DB) {
	user := User{Name: "Jinzhu", Age: 18, Birthday: time.Now()}
	db.Create(&user)
}

// CreateWithSelectedFields demonstrates creating a record with selected fields.
func CreateWithSelectedFields(db *gorm.DB) {
	user := User{Name: "Jinzhu", Age: 18, Birthday: time.Now()}
	db.Select("Name", "Age", "CreatedAt").Create(&user)
}

// BatchInsert demonstrates inserting multiple records at once.
func BatchInsert(db *gorm.DB) {
	var users = []User{{Name: "jinzhu1"}, {Name: "jinzhu2"}, {Name: "jinzhu3"}}
	db.Create(&users)
}
