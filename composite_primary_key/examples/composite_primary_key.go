
package examples

import (
	"fmt"

	"gorm.io/driver/sqlite"
	"gorm.io/gorm"
)

// --- Basic Composite Primary Key ---

// Product model with a composite primary key made of two strings.
// The combination of (ID, LanguageCode) must be unique.
type Product struct {
	ID           string `gorm:"primaryKey"`
	LanguageCode string `gorm:"primaryKey"`
	Name         string
	Description  string
}

// CreateAndFindProduct demonstrates creating and finding a record with a composite key.
func CreateAndFindProduct(db *gorm.DB) error {
	db.AutoMigrate(&Product{})

	// Create a product
	product := Product{
		ID:           "product-001",
		LanguageCode: "en-US",
		Name:         "GORM for Beginners",
		Description:  "A book about GORM",
	}
	if err := db.Create(&product).Error; err != nil {
		return fmt.Errorf("failed to create product: %w", err)
	}

	// Find the product using its composite key
	var foundProduct Product
	// When querying with a struct, GORM automatically uses the primary key fields.
	if err := db.First(&foundProduct, Product{ID: "product-001", LanguageCode: "en-US"}).Error; err != nil {
		return fmt.Errorf("failed to find product: %w", err)
	}

	fmt.Printf("Found product: %s\n", foundProduct.Name)
	return nil
}

// --- Composite Primary Key with Integers ---

// OrderItem model with a composite primary key made of two integers.
// `autoIncrement` is explicitly disabled to prevent the database from managing the keys.
type OrderItem struct {
	OrderID   uint `gorm:"primaryKey;autoIncrement:false"`
	ProductID uint `gorm:"primaryKey;autoIncrement:false"`
	Quantity  int
}

// CreateAndFindOrderItem demonstrates working with integer-based composite keys.
func CreateAndFindOrderItem(db *gorm.DB) error {
	db.AutoMigrate(&OrderItem{})

	orderItem := OrderItem{OrderID: 100, ProductID: 200, Quantity: 3}
	if err := db.Create(&orderItem).Error; err != nil {
		return fmt.Errorf("failed to create order item: %w", err)
	}

	var foundItem OrderItem
	if err := db.First(&foundItem, OrderItem{OrderID: 100, ProductID: 200}).Error; err != nil {
		return fmt.Errorf("failed to find order item: %w", err)
	}

	fmt.Printf("Found order item with quantity: %d\n", foundItem.Quantity)
	return nil
}

// --- Composite Primary Keys and Associations ---

// Tag model with composite key, used in a many-to-many relationship.
type Tag struct {
	ID     uint   `gorm:"primaryKey"`
	Locale string `gorm:"primaryKey"`
	Value  string
}

// Blog model with a many-to-many relationship to Tag.
// GORM automatically uses the composite key of Tag for the join table.
type Blog struct {
	gorm.Model
	Title string
	Tags  []Tag `gorm:"many2many:blog_tags;"`
}

// CreateBlogWithTags demonstrates using composite keys in associations.
func CreateBlogWithTags(db *gorm.DB) error {
	db.AutoMigrate(&Blog{}, &Tag{})

	blog := Blog{
		Title: "Advanced GORM",
		Tags: []Tag{
			{ID: 1, Locale: "en", Value: "golang"},
			{ID: 2, Locale: "en", Value: "orm"},
		},
	}
	if err := db.Create(&blog).Error; err != nil {
		return fmt.Errorf("failed to create blog with tags: %w", err)
	}

	// Preload the tags
	var foundBlog Blog
	db.Preload("Tags").First(&foundBlog, blog.ID)
	fmt.Printf("Blog '%s' has %d tags.\n", foundBlog.Title, len(foundBlog.Tags))

	return nil
}

func main() {
	// This main function is for demonstration purposes.
	db, err := gorm.Open(sqlite.Open("file::memory:?cache=shared"), &gorm.Config{})
	if err != nil {
		panic("failed to connect database")
	}

	fmt.Println("--- Testing Basic Composite Key ---")
	if err := CreateAndFindProduct(db); err != nil {
		panic(err)
	}

	fmt.Println("\n--- Testing Integer Composite Key ---")
	if err := CreateAndFindOrderItem(db); err != nil {
		panic(err)
	}

	fmt.Println("\n--- Testing Composite Key with Associations ---")
	if err := CreateBlogWithTags(db); err != nil {
		panic(err)
	}
}
