package examples

import (
	"gorm.io/gorm"
)

// User is a basic GORM model for delete examples.
type User struct {
	gorm.Model
	Name string
	Age  int
}

// DeleteRecord demonstrates deleting a single record.
func DeleteRecord(db *gorm.DB, user *User) {
	db.Delete(user)
}

// DeleteByPrimaryKey demonstrates deleting by primary key.
func DeleteByPrimaryKey(db *gorm.DB) {
	db.Delete(&User{}, 10)
}

// BatchDelete demonstrates a batch delete operation.
func BatchDelete(db *gorm.DB) {
	db.Where("age < ?", 18).Delete(&User{})
}

// SoftDelete demonstrates soft deleting a record.
func SoftDelete(db *gorm.DB, user *User) {
	db.Delete(user)
}

// FindSoftDeleted demonstrates how to find soft-deleted records.
func FindSoftDeleted(db *gorm.DB) {
	var users []User
	db.Unscoped().Where("age = 20").Find(&users)
}

// PermanentDelete demonstrates a permanent delete.
func PermanentDelete(db *gorm.DB, user *User) {
	db.Unscoped().Delete(user)
}
