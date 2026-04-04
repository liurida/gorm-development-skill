package examples

import (
	"fmt"

	"gorm.io/driver/sqlite"
	"gorm.io/gorm"
)

// --- Check Constraint ---

// UserWithCheck demonstrates a CHECK constraint to enforce a rule on a column.
type UserWithCheck struct {
	gorm.Model
	Name string
	Age  int `gorm:"check:age_checker,age >= 18"` // Named CHECK constraint
	Role string `gorm:"check:role IN ('admin', 'user', 'guest')"` // Unnamed CHECK constraint
}

// DemonstrateCheckConstraint shows how CHECK constraints work.
func DemonstrateCheckConstraint(db *gorm.DB) error {
	db.AutoMigrate(&UserWithCheck{})

	// This will succeed
	user1 := UserWithCheck{Name: "Adult User", Age: 25, Role: "admin"}
	if err := db.Create(&user1).Error; err != nil {
		return fmt.Errorf("failed to create valid user: %w", err)
	}

	// This will fail because age is less than 18
	user2 := UserWithCheck{Name: "Young User", Age: 16, Role: "user"}
	if err := db.Create(&user2).Error; err == nil {
		return fmt.Errorf("expected error for user with age < 18 but got none")
	}
	fmt.Println("Successfully blocked user with age < 18.")

	// This will fail because the role is invalid
	user3 := UserWithCheck{Name: "Invalid Role User", Age: 30, Role: "super-admin"}
	if err := db.Create(&user3).Error; err == nil {
		return fmt.Errorf("expected error for user with invalid role but got none")
	}
	fmt.Println("Successfully blocked user with invalid role.")

	return nil
}

// --- Foreign Key Constraint ---

// Company model for the foreign key relationship.
type Company struct {
	ID   int `gorm:"primaryKey"`
	Name string
}

// UserWithFK demonstrates foreign key constraints on associations.
type UserWithFK struct {
	gorm.Model
	Name      string
	CompanyID int
	Company   Company `gorm:"constraint:OnUpdate:CASCADE,OnDelete:SET NULL;"`
}

// DemonstrateForeignKeyConstraint shows how foreign key constraints work.
func DemonstrateForeignKeyConstraint(db *gorm.DB) error {
	// Disable foreign key constraints for SQLite during migration for this example
	db.Exec("PRAGMA foreign_keys = OFF;")
	db.AutoMigrate(&Company{}, &UserWithFK{})
	db.Exec("PRAGMA foreign_keys = ON;")

	// Create a company
	company := Company{ID: 1, Name: "Tech Corp"}
	db.Create(&company)

	// Create a user associated with the company
	user := UserWithFK{Name: "John Doe", CompanyID: 1}
	db.Create(&user)

	// 1. OnUpdate: CASCADE
	// Update the company's ID. The user's CompanyID should be updated automatically.
	db.Model(&company).Update("ID", 2)

	var updatedUser UserWithFK
	db.First(&updatedUser, user.ID)
	if updatedUser.CompanyID != 2 {
		return fmt.Errorf("expected user CompanyID to be 2 after cascade update, but got %d", updatedUser.CompanyID)
	}
	fmt.Println("OnUpdate:CASCADE successfully updated the foreign key.")

	// 2. OnDelete: SET NULL
	// Delete the company. The user's CompanyID should be set to NULL.
	db.Delete(&Company{}, 2)

	db.First(&updatedUser, user.ID)
	// Note: In GORM, if the foreign key is a non-pointer int, it will be set to 0 (the zero value).
	// For it to be set to NULL in the database, CompanyID should be a pointer (*int).
	if updatedUser.CompanyID != 0 {
		return fmt.Errorf("expected user CompanyID to be 0 (NULL) after delete, but got %d", updatedUser.CompanyID)
	}
	fmt.Println("OnDelete:SET NULL successfully nulled the foreign key.")

	return nil
}

// --- Disabling Foreign Key Constraints during Migration ---

// NoForeignKeyDB demonstrates initializing GORM to skip creating foreign key constraints.
func NoForeignKeyDB() (*gorm.DB, error) {
	return gorm.Open(sqlite.Open("file::memory:"), &gorm.Config{
		DisableForeignKeyConstraintWhenMigrating: true,
	})
}

func main() {
	db, err := gorm.Open(sqlite.Open("file::memory:?cache=shared"), &gorm.Config{})
	if err != nil {
		panic("failed to connect database")
	}

	fmt.Println("--- Testing CHECK Constraints ---")
	if err := DemonstrateCheckConstraint(db); err != nil {
		panic(err)
	}

	fmt.Println("\n--- Testing Foreign Key Constraints ---")
	if err := DemonstrateForeignKeyConstraint(db); err != nil {
		panic(err)
	}
}
