package examples

import (
	"errors"

	"github.com/google/uuid"
	"gorm.io/gorm"
)

// User is a basic GORM model for hook examples.
type User struct {
	gorm.Model
	Name   string
	UUID   string
	Role   string
	Active bool
}

// BeforeCreate is a GORM hook that is called before a record is created.
func (u *User) BeforeCreate(tx *gorm.DB) (err error) {
	u.UUID = uuid.New().String()
	if u.Role == "admin" {
		return errors.New("invalid role")
	}
	return
}

// AfterCreate is a GORM hook that is called after a record is created.
func (u *User) AfterCreate(tx *gorm.DB) (err error) {
	if u.ID == 1 {
		tx.Model(u).Update("Role", "admin")
	}
	return
}
