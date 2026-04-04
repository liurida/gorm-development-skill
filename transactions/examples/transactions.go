package examples

import (
	"errors"

	"gorm.io/gorm"
)

// User is a basic GORM model for transaction examples.
type User struct {
	gorm.Model
	Name string
}

// BasicTransaction demonstrates a basic transaction.
func BasicTransaction(db *gorm.DB) error {
	return db.Transaction(func(tx *gorm.DB) error {
		if err := tx.Create(&User{Name: "Giraffe"}).Error; err != nil {
			return err
		}

		if err := tx.Create(&User{Name: "Lion"}).Error; err != nil {
			return err
		}

		return nil
	})
}

// ManualTransaction demonstrates manual transaction control.
func ManualTransaction(db *gorm.DB) error {
	tx := db.Begin()
	defer func() {
		if r := recover(); r != nil {
			tx.Rollback()
		}
	}()

	if err := tx.Error; err != nil {
		return err
	}

	if err := tx.Create(&User{Name: "Giraffe"}).Error; err != nil {
		tx.Rollback()
		return err
	}

	if err := tx.Create(&User{Name: "Lion"}).Error; err != nil {
		tx.Rollback()
		return err
	}

	return tx.Commit().Error
}
