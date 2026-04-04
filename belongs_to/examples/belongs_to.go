package examples

import "gorm.io/gorm"

// User has one CreditCard, CreditCardID is the foreign key
type User struct {
	gorm.Model
	Name         string
	CreditCardID uint
	CreditCard   CreditCard
}

// CreditCard belongs to a User
type CreditCard struct {
	gorm.Model
	Number string
}
