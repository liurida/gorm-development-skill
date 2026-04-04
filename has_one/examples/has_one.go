package examples

import "gorm.io/gorm"

// User has one CreditCard, UserID is the foreign key
type User struct {
	gorm.Model
	CreditCard CreditCard
}

// CreditCard has one User
type CreditCard struct {
	gorm.Model
	Number string
	UserID uint
}
