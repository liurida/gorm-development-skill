package examples

import (
	"fmt"

	"gorm.io/gorm"
)

// --- Basic Polymorphic Association ---

// Dog has many toys with polymorphic association
type Dog struct {
	gorm.Model
	Name string
	Toys []Toy `gorm:"polymorphic:Owner;"`
}

// Cat has many toys with polymorphic association
type Cat struct {
	gorm.Model
	Name string
	Toys []Toy `gorm:"polymorphic:Owner;"`
}

// Toy can belong to either Dog or Cat (polymorphic)
type Toy struct {
	gorm.Model
	Name      string
	OwnerID   uint
	OwnerType string
}

// CreateDogWithToys demonstrates polymorphic association.
func CreateDogWithToys(db *gorm.DB) error {
	dog := Dog{
		Name: "Buddy",
		Toys: []Toy{
			{Name: "Ball"},
			{Name: "Rope"},
		},
	}

	result := db.Create(&dog)
	if result.Error != nil {
		return fmt.Errorf("failed to create dog with toys: %w", result.Error)
	}

	// OwnerType will be "dogs" (pluralized table name)
	// OwnerID will be the dog's ID
	fmt.Printf("Created dog %s with %d toys\n", dog.Name, len(dog.Toys))
	return nil
}

// CreateCatWithToys demonstrates the same polymorphic association for a different type.
func CreateCatWithToys(db *gorm.DB) error {
	cat := Cat{
		Name: "Whiskers",
		Toys: []Toy{
			{Name: "Mouse"},
			{Name: "Feather"},
		},
	}

	result := db.Create(&cat)
	if result.Error != nil {
		return fmt.Errorf("failed to create cat with toys: %w", result.Error)
	}

	// OwnerType will be "cats" (pluralized table name)
	fmt.Printf("Created cat %s with %d toys\n", cat.Name, len(cat.Toys))
	return nil
}

// GetDogWithToys demonstrates preloading polymorphic association.
func GetDogWithToys(db *gorm.DB, dogID uint) (*Dog, error) {
	var dog Dog
	err := db.Preload("Toys").First(&dog, dogID).Error
	if err != nil {
		return nil, fmt.Errorf("failed to get dog with toys: %w", err)
	}

	return &dog, nil
}

// FindToysByOwnerType demonstrates querying polymorphic records by type.
func FindToysByOwnerType(db *gorm.DB, ownerType string) ([]Toy, error) {
	var toys []Toy
	err := db.Where("owner_type = ?", ownerType).Find(&toys).Error
	if err != nil {
		return nil, fmt.Errorf("failed to find toys: %w", err)
	}

	return toys, nil
}

// --- Custom Polymorphic Column Names ---

// Hamster has toys with custom polymorphic configuration
type Hamster struct {
	gorm.Model
	Name string
	Toys []HamsterToy `gorm:"polymorphicType:Kind;polymorphicId:OwnerID;polymorphicValue:master"`
}

// HamsterToy uses custom column name "Kind" instead of "OwnerType"
type HamsterToy struct {
	gorm.Model
	Name    string
	OwnerID uint
	Kind    string // Custom type column
}

// CreateHamsterWithToys demonstrates custom polymorphic configuration.
func CreateHamsterWithToys(db *gorm.DB) error {
	hamster := Hamster{
		Name: "Hammy",
		Toys: []HamsterToy{
			{Name: "Wheel"},
			{Name: "Tunnel"},
		},
	}

	result := db.Create(&hamster)
	if result.Error != nil {
		return fmt.Errorf("failed to create hamster with toys: %w", result.Error)
	}

	// Kind will be "master" (custom polymorphicValue) instead of "hamsters"
	fmt.Printf("Created hamster %s with %d toys\n", hamster.Name, len(hamster.Toys))
	return nil
}

// --- Polymorphic Has One ---

// Company has one polymorphic address
type Company struct {
	gorm.Model
	Name    string
	Address Address `gorm:"polymorphic:Addressable;"`
}

// Person has one polymorphic address
type Person struct {
	gorm.Model
	Name    string
	Address Address `gorm:"polymorphic:Addressable;"`
}

// Address can belong to either Company or Person
type Address struct {
	gorm.Model
	Street          string
	City            string
	AddressableID   uint
	AddressableType string
}

// CreateCompanyWithAddress demonstrates polymorphic has one.
func CreateCompanyWithAddress(db *gorm.DB) error {
	company := Company{
		Name: "Acme Corp",
		Address: Address{
			Street: "123 Business Ave",
			City:   "Tech City",
		},
	}

	return db.Create(&company).Error
}

// CreatePersonWithAddress demonstrates polymorphic has one for a different type.
func CreatePersonWithAddress(db *gorm.DB) error {
	person := Person{
		Name: "John Doe",
		Address: Address{
			Street: "456 Home St",
			City:   "Hometown",
		},
	}

	return db.Create(&person).Error
}

// --- Polymorphic Comments (Common Use Case) ---

// Post has many polymorphic comments
type Post struct {
	gorm.Model
	Title    string
	Content  string
	Comments []Comment `gorm:"polymorphic:Commentable;"`
}

// Video has many polymorphic comments
type Video struct {
	gorm.Model
	Title    string
	URL      string
	Comments []Comment `gorm:"polymorphic:Commentable;"`
}

// Comment can belong to Post, Video, or any other commentable entity
type Comment struct {
	gorm.Model
	Content         string
	CommentableID   uint
	CommentableType string
}

// CreatePostWithComments demonstrates a common polymorphic pattern for comments.
func CreatePostWithComments(db *gorm.DB) error {
	post := Post{
		Title:   "My First Post",
		Content: "Hello World!",
		Comments: []Comment{
			{Content: "Great post!"},
			{Content: "Thanks for sharing."},
		},
	}

	return db.Create(&post).Error
}

// AddCommentToPost demonstrates adding to a polymorphic association.
func AddCommentToPost(db *gorm.DB, postID uint, content string) error {
	var post Post
	if err := db.First(&post, postID).Error; err != nil {
		return fmt.Errorf("post not found: %w", err)
	}

	comment := Comment{Content: content}
	err := db.Model(&post).Association("Comments").Append(&comment)
	if err != nil {
		return fmt.Errorf("failed to add comment: %w", err)
	}

	return nil
}

// GetAllComments demonstrates querying all comments regardless of owner type.
func GetAllComments(db *gorm.DB) ([]Comment, error) {
	var comments []Comment
	err := db.Find(&comments).Error
	if err != nil {
		return nil, fmt.Errorf("failed to get comments: %w", err)
	}

	return comments, nil
}
