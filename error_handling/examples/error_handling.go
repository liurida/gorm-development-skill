
package examples

import (
	"errors"
	"fmt"

	"gorm.io/gorm"
)

// ErrorUser model for error handling examples
type ErrorUser struct {
	gorm.Model
	Name string `gorm:"unique"`
}

// HandleBasicError demonstrates checking the .Error field after an operation.
func HandleBasicError(db *gorm.DB) error {
	// This will cause an error because the table doesn't exist yet.
	if err := db.First(&ErrorUser{}).Error; err != nil {
		fmt.Printf("Caught an error: %v\n", err)
		return err
	}
	return nil
}

// HandleRecordNotFound demonstrates how to specifically handle gorm.ErrRecordNotFound.
func HandleRecordNotFound(db *gorm.DB) error {
	// Ensure the table exists first
	db.AutoMigrate(&ErrorUser{})

	var user ErrorUser
	// Try to find a user that doesn't exist
	err := db.Where("name = ?", "non_existent_user").First(&user).Error

	if errors.Is(err, gorm.ErrRecordNotFound) {
		fmt.Println("User not found, which is an expected outcome.")
		return nil // Not treating this as a fatal error
	}
	if err != nil {
		// Handle other potential database errors
		return err
	}

	return nil
}

// HandleTransactionError demonstrates automatic rollback on error in a transaction.
func HandleTransactionError(db *gorm.DB) error {
	db.AutoMigrate(&ErrorUser{})

	err := db.Transaction(func(tx *gorm.DB) error {
		// Create a user
		if err := tx.Create(&ErrorUser{Name: "trans_user"}).Error; err != nil {
			return err
		}

		// This second create will fail because the name must be unique.
		if err := tx.Create(&ErrorUser{Name: "trans_user"}).Error; err != nil {
			fmt.Println("Got an error, transaction will be rolled back.")
			// By returning the error, GORM knows to roll back the transaction.
			return err
		}

		return nil
	})

	// Verify that the first user was not created due to the rollback.
	var count int64
	db.Model(&ErrorUser{}).Where("name = ?", "trans_user").Count(&count)
	fmt.Printf("Count of 'trans_user': %d\n", count) // Should be 0

	return err // Return the error from the transaction
}

// UseTranslatedErrors demonstrates handling generic GORM errors.
// This requires the GORM DB instance to be configured with `TranslateError: true`.
func UseTranslatedErrors(db *gorm.DB) error {
	db.AutoMigrate(&ErrorUser{})

	// Create an initial user
	db.Create(&ErrorUser{Name: "unique_user"})

	// Attempt to create a user with the same name
	err := db.Create(&ErrorUser{Name: "unique_user"}).Error

	if errors.Is(err, gorm.ErrDuplicatedKey) {
		fmt.Println("Caught a duplicate key error, as expected.")
		return nil // Not a fatal error for this example
	}

	return err
}
