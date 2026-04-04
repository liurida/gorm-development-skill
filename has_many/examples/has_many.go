package examples

import (
	"fmt"

	"gorm.io/gorm"
)

// User has many CreditCards - basic one-to-many relationship
type User struct {
	gorm.Model
	Name        string
	MemberNumber string
	CreditCards []CreditCard
}

// CreditCard belongs to User via UserID foreign key
type CreditCard struct {
	gorm.Model
	Number string
	UserID uint
}

// CreateUserWithCreditCards demonstrates creating a user with associated credit cards.
func CreateUserWithCreditCards(db *gorm.DB) error {
	user := User{
		Name: "John Doe",
		CreditCards: []CreditCard{
			{Number: "4111111111111111"},
			{Number: "5500000000000004"},
		},
	}

	result := db.Create(&user)
	if result.Error != nil {
		return fmt.Errorf("failed to create user with credit cards: %w", result.Error)
	}

	fmt.Printf("Created user %d with %d credit cards\n", user.ID, len(user.CreditCards))
	return nil
}

// GetUserWithCreditCards demonstrates preloading has many associations.
func GetUserWithCreditCards(db *gorm.DB, userID uint) (*User, error) {
	var user User
	err := db.Preload("CreditCards").First(&user, userID).Error
	if err != nil {
		return nil, fmt.Errorf("failed to get user with credit cards: %w", err)
	}

	return &user, nil
}

// AddCreditCardToUser demonstrates appending to a has many association.
func AddCreditCardToUser(db *gorm.DB, userID uint, cardNumber string) error {
	var user User
	if err := db.First(&user, userID).Error; err != nil {
		return fmt.Errorf("user not found: %w", err)
	}

	card := CreditCard{Number: cardNumber}
	err := db.Model(&user).Association("CreditCards").Append(&card)
	if err != nil {
		return fmt.Errorf("failed to add credit card: %w", err)
	}

	return nil
}

// RemoveCreditCardFromUser demonstrates removing from a has many association.
func RemoveCreditCardFromUser(db *gorm.DB, userID uint, cardID uint) error {
	var user User
	if err := db.First(&user, userID).Error; err != nil {
		return fmt.Errorf("user not found: %w", err)
	}

	var card CreditCard
	if err := db.First(&card, cardID).Error; err != nil {
		return fmt.Errorf("credit card not found: %w", err)
	}

	err := db.Model(&user).Association("CreditCards").Delete(&card)
	if err != nil {
		return fmt.Errorf("failed to remove credit card: %w", err)
	}

	return nil
}

// ClearAllCreditCards demonstrates clearing all associations.
func ClearAllCreditCards(db *gorm.DB, userID uint) error {
	var user User
	if err := db.First(&user, userID).Error; err != nil {
		return fmt.Errorf("user not found: %w", err)
	}

	err := db.Model(&user).Association("CreditCards").Clear()
	if err != nil {
		return fmt.Errorf("failed to clear credit cards: %w", err)
	}

	return nil
}

// CountCreditCards demonstrates counting associations.
func CountCreditCards(db *gorm.DB, userID uint) (int64, error) {
	var user User
	if err := db.First(&user, userID).Error; err != nil {
		return 0, fmt.Errorf("user not found: %w", err)
	}

	count := db.Model(&user).Association("CreditCards").Count()
	return count, nil
}

// --- Custom Foreign Key Example ---

// Employee has many Tasks with custom foreign key
type Employee struct {
	gorm.Model
	Name  string
	Tasks []Task `gorm:"foreignKey:AssigneeID"`
}

// Task belongs to Employee via custom AssigneeID foreign key
type Task struct {
	gorm.Model
	Title      string
	AssigneeID uint
}

// CreateEmployeeWithTasks demonstrates custom foreign key usage.
func CreateEmployeeWithTasks(db *gorm.DB) error {
	employee := Employee{
		Name: "Jane Smith",
		Tasks: []Task{
			{Title: "Review PR #123"},
			{Title: "Deploy to staging"},
		},
	}

	return db.Create(&employee).Error
}

// --- Custom References Example ---

// Company uses MemberNumber as the reference for employees
type Company struct {
	gorm.Model
	MemberNumber string
	Employees    []Staff `gorm:"foreignKey:CompanyNumber;references:MemberNumber"`
}

// Staff belongs to Company via custom reference
type Staff struct {
	gorm.Model
	Name          string
	CompanyNumber string
}

// --- Self-Referential Has Many ---

// Manager has many team members (self-referential)
type Manager struct {
	gorm.Model
	Name      string
	ManagerID *uint
	Team      []Manager `gorm:"foreignKey:ManagerID"`
}

// CreateTeamHierarchy demonstrates self-referential has many.
func CreateTeamHierarchy(db *gorm.DB) error {
	manager := Manager{
		Name: "Team Lead",
		Team: []Manager{
			{Name: "Developer 1"},
			{Name: "Developer 2"},
		},
	}

	return db.Create(&manager).Error
}

// --- Has Many with Constraints ---

// Customer has many Orders with ON DELETE CASCADE
type Customer struct {
	gorm.Model
	Name   string
	Orders []Order `gorm:"constraint:OnUpdate:CASCADE,OnDelete:CASCADE;"`
}

// Order belongs to Customer
type Order struct {
	gorm.Model
	Amount     float64
	CustomerID uint
}

// DeleteCustomerCascade demonstrates cascading delete behavior.
func DeleteCustomerCascade(db *gorm.DB, customerID uint) error {
	// When customer is deleted, all associated orders are also deleted
	return db.Delete(&Customer{}, customerID).Error
}
