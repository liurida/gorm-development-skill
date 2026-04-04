package examples

import "gorm.io/gorm"

// User has one CreditCard
type User struct {
	gorm.Model
	Name       string
	CreditCard CreditCard
}

// CreditCard belongs to a User
type CreditCard struct {
	gorm.Model
	Number string
	UserID uint
}

// CreateWithAssociation demonstrates creating a record with associations.
func CreateWithAssociation(db *gorm.DB) {
	user := User{
		Name:       "jinzhu",
		CreditCard: CreditCard{Number: "411111111111"},
	}
	db.Create(&user)
}

// AssociationMode demonstrates using association mode to manage relationships.
func AssociationMode(db *gorm.DB, user *User) {
	// Append
	db.Model(user).Association("CreditCard").Append(&CreditCard{Number: "411111111111"})

	// Delete
	db.Model(user).Association("CreditCard").Delete(&CreditCard{Number: "411111111111"})

	// Clear
	db.Model(user).Association("CreditCard").Clear()
}
