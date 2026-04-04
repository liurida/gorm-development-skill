package examples

import (
	"fmt"

	"gorm.io/gorm"
	"gorm.io/gorm/clause"
)

// --- Models for Preload Examples ---

// User has many Orders and one Profile
type User struct {
	gorm.Model
	Username  string
	Profile   Profile
	Orders    []Order
	CompanyID uint
	Company   Company
}

// Profile belongs to User
type Profile struct {
	gorm.Model
	UserID uint
	Bio    string
	Avatar string
}

// Order belongs to User, has many OrderItems
type Order struct {
	gorm.Model
	UserID     uint
	Price      float64
	State      string
	OrderItems []OrderItem
}

// OrderItem belongs to Order and Product
type OrderItem struct {
	gorm.Model
	OrderID   uint
	ProductID uint
	Quantity  int
	Product   Product
}

// Product model
type Product struct {
	gorm.Model
	Name  string
	Price float64
}

// Company model
type Company struct {
	gorm.Model
	Name  string
	Alive bool
}

// --- Basic Preload ---

// PreloadOrders demonstrates basic preloading of a has-many association.
func PreloadOrders(db *gorm.DB) ([]User, error) {
	var users []User
	// Executes two queries:
	// SELECT * FROM users;
	// SELECT * FROM orders WHERE user_id IN (1,2,3,4);
	err := db.Preload("Orders").Find(&users).Error
	return users, err
}

// PreloadMultipleAssociations demonstrates preloading multiple associations.
func PreloadMultipleAssociations(db *gorm.DB) ([]User, error) {
	var users []User
	// Preloads all three associations
	err := db.Preload("Orders").Preload("Profile").Preload("Company").Find(&users).Error
	return users, err
}

// --- Preload All ---

// PreloadAll demonstrates using clause.Associations to preload all associations.
func PreloadAll(db *gorm.DB) ([]User, error) {
	var users []User
	// Preloads Orders, Profile, and Company
	// Note: Does NOT preload nested associations like Orders.OrderItems
	err := db.Preload(clause.Associations).Find(&users).Error
	return users, err
}

// --- Preload with Conditions ---

// PreloadWithConditions demonstrates filtering preloaded records.
func PreloadWithConditions(db *gorm.DB) ([]User, error) {
	var users []User
	// Only preload paid orders
	err := db.Preload("Orders", "state = ?", "paid").Find(&users).Error
	return users, err
}

// PreloadWithCustomQuery demonstrates using a function for custom preload logic.
func PreloadWithCustomQuery(db *gorm.DB) ([]User, error) {
	var users []User
	// Order the preloaded records
	err := db.Preload("Orders", func(db *gorm.DB) *gorm.DB {
		return db.Order("orders.price DESC").Limit(5) // Only top 5 expensive orders
	}).Find(&users).Error
	return users, err
}

// --- Nested Preloading ---

// PreloadNested demonstrates preloading nested associations using dot notation.
func PreloadNested(db *gorm.DB) ([]User, error) {
	var users []User
	// Preloads Orders -> OrderItems -> Product
	err := db.Preload("Orders.OrderItems.Product").Find(&users).Error
	return users, err
}

// PreloadNestedWithConditions demonstrates conditional nested preloading.
func PreloadNestedWithConditions(db *gorm.DB) ([]User, error) {
	var users []User
	// Only preload paid orders and their items
	// GORM won't preload OrderItems for unpaid orders
	err := db.Preload("Orders", "state = ?", "paid").
		Preload("Orders.OrderItems").
		Find(&users).Error
	return users, err
}

// PreloadAllWithNested demonstrates combining clause.Associations with nested preloading.
func PreloadAllWithNested(db *gorm.DB) ([]User, error) {
	var users []User
	// First explicitly preload nested, then all direct associations
	err := db.Preload("Orders.OrderItems.Product").
		Preload(clause.Associations).
		Find(&users).Error
	return users, err
}

// --- Joins Preloading (for one-to-one/belongs-to) ---

// JoinsPreload demonstrates using Joins for belongs-to/has-one associations.
// This uses LEFT JOIN instead of separate queries.
func JoinsPreload(db *gorm.DB) ([]User, error) {
	var users []User
	// Uses LEFT JOIN - more efficient for belongs-to relationships
	// Works with: has_one, belongs_to
	// Does NOT work with: has_many, many_to_many
	err := db.Joins("Company").Joins("Profile").Find(&users).Error
	return users, err
}

// JoinsWithConditions demonstrates Joins with conditions.
func JoinsWithConditions(db *gorm.DB) ([]User, error) {
	var users []User
	// Only join companies that are alive
	err := db.Joins("Company", db.Where(&Company{Alive: true})).Find(&users).Error
	return users, err
}

// JoinsNestedModel demonstrates joining nested models.
func JoinsNestedModel(db *gorm.DB) ([]User, error) {
	var users []User
	// For models like User -> Manager (self-referential) -> Company
	// This example assumes User has a Manager field
	err := db.Joins("Company").Find(&users).Error
	return users, err
}

// --- Practical Examples ---

// GetUserOrderHistory returns a user with their complete order history.
func GetUserOrderHistory(db *gorm.DB, userID uint) (*User, error) {
	var user User
	err := db.Preload("Orders", func(db *gorm.DB) *gorm.DB {
		return db.Order("created_at DESC")
	}).Preload("Orders.OrderItems.Product").First(&user, userID).Error

	if err != nil {
		return nil, fmt.Errorf("failed to get user order history: %w", err)
	}

	return &user, nil
}

// GetAllUsersWithCompanyInfo efficiently loads users with company data.
func GetAllUsersWithCompanyInfo(db *gorm.DB) ([]User, error) {
	var users []User
	// Use Joins for single-record associations (more efficient)
	// Use Preload for collections
	err := db.Joins("Company").Joins("Profile").Preload("Orders").Find(&users).Error

	if err != nil {
		return nil, fmt.Errorf("failed to get users with company info: %w", err)
	}

	return users, nil
}

// GetRecentPaidOrders demonstrates combining query conditions with preloading.
func GetRecentPaidOrders(db *gorm.DB) ([]User, error) {
	var users []User
	err := db.Where("created_at > ?", "2024-01-01").
		Preload("Orders", "state = ? AND price > ?", "paid", 100).
		Preload("Orders.OrderItems").
		Find(&users).Error

	return users, err
}

// CountOrderItems demonstrates accessing preloaded nested data.
func CountOrderItems(user *User) int {
	total := 0
	for _, order := range user.Orders {
		total += len(order.OrderItems)
	}
	return total
}
