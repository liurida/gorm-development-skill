package examples

import (
	"fmt"
	"time"

	"gorm.io/gorm"
)

// --- Basic Many-to-Many ---

// User has and belongs to many Languages
type User struct {
	gorm.Model
	Name      string
	Languages []Language `gorm:"many2many:user_languages;"`
}

// Language has and belongs to many Users
type Language struct {
	gorm.Model
	Name  string
	Users []*User `gorm:"many2many:user_languages;"`
}

// CreateUserWithLanguages demonstrates creating a many-to-many relationship.
func CreateUserWithLanguages(db *gorm.DB) error {
	user := User{
		Name: "John Doe",
		Languages: []Language{
			{Name: "English"},
			{Name: "Spanish"},
			{Name: "French"},
		},
	}

	result := db.Create(&user)
	if result.Error != nil {
		return fmt.Errorf("failed to create user with languages: %w", result.Error)
	}

	fmt.Printf("Created user %d with %d languages\n", user.ID, len(user.Languages))
	return nil
}

// GetUserWithLanguages demonstrates preloading many-to-many associations.
func GetUserWithLanguages(db *gorm.DB, userID uint) (*User, error) {
	var user User
	err := db.Preload("Languages").First(&user, userID).Error
	if err != nil {
		return nil, fmt.Errorf("failed to get user with languages: %w", err)
	}

	return &user, nil
}

// GetLanguageWithUsers demonstrates reverse many-to-many lookup.
func GetLanguageWithUsers(db *gorm.DB, langID uint) (*Language, error) {
	var lang Language
	err := db.Preload("Users").First(&lang, langID).Error
	if err != nil {
		return nil, fmt.Errorf("failed to get language with users: %w", err)
	}

	return &lang, nil
}

// AddLanguageToUser demonstrates appending to a many-to-many association.
func AddLanguageToUser(db *gorm.DB, userID uint, langName string) error {
	var user User
	if err := db.First(&user, userID).Error; err != nil {
		return fmt.Errorf("user not found: %w", err)
	}

	// Check if language exists, create if not
	var lang Language
	result := db.Where("name = ?", langName).FirstOrCreate(&lang, Language{Name: langName})
	if result.Error != nil {
		return fmt.Errorf("failed to find or create language: %w", result.Error)
	}

	err := db.Model(&user).Association("Languages").Append(&lang)
	if err != nil {
		return fmt.Errorf("failed to add language: %w", err)
	}

	return nil
}

// RemoveLanguageFromUser demonstrates removing from a many-to-many association.
func RemoveLanguageFromUser(db *gorm.DB, userID uint, langID uint) error {
	var user User
	if err := db.First(&user, userID).Error; err != nil {
		return fmt.Errorf("user not found: %w", err)
	}

	var lang Language
	if err := db.First(&lang, langID).Error; err != nil {
		return fmt.Errorf("language not found: %w", err)
	}

	// This removes the association from join table, not the language itself
	err := db.Model(&user).Association("Languages").Delete(&lang)
	if err != nil {
		return fmt.Errorf("failed to remove language: %w", err)
	}

	return nil
}

// ReplaceUserLanguages demonstrates replacing all associations.
func ReplaceUserLanguages(db *gorm.DB, userID uint, langIDs []uint) error {
	var user User
	if err := db.First(&user, userID).Error; err != nil {
		return fmt.Errorf("user not found: %w", err)
	}

	var languages []Language
	if err := db.Find(&languages, langIDs).Error; err != nil {
		return fmt.Errorf("failed to find languages: %w", err)
	}

	err := db.Model(&user).Association("Languages").Replace(&languages)
	if err != nil {
		return fmt.Errorf("failed to replace languages: %w", err)
	}

	return nil
}

// --- Self-Referential Many-to-Many (Friends) ---

// Person has many friends (self-referential many-to-many)
type Person struct {
	gorm.Model
	Name    string
	Friends []*Person `gorm:"many2many:person_friends;"`
}

// AddFriend demonstrates self-referential many-to-many.
func AddFriend(db *gorm.DB, personID uint, friendID uint) error {
	var person Person
	if err := db.First(&person, personID).Error; err != nil {
		return fmt.Errorf("person not found: %w", err)
	}

	var friend Person
	if err := db.First(&friend, friendID).Error; err != nil {
		return fmt.Errorf("friend not found: %w", err)
	}

	err := db.Model(&person).Association("Friends").Append(&friend)
	if err != nil {
		return fmt.Errorf("failed to add friend: %w", err)
	}

	return nil
}

// --- Custom Join Table ---

// PersonAddress is a custom join table with additional fields
type PersonAddress struct {
	PersonID  int `gorm:"primaryKey"`
	AddressID int `gorm:"primaryKey"`
	CreatedAt time.Time
	DeletedAt gorm.DeletedAt
	IsPrimary bool // Additional field
}

// PersonM2M has many addresses through custom join table
type PersonM2M struct {
	gorm.Model
	Name      string
	Addresses []AddressM2M `gorm:"many2many:person_addresses;"`
}

// AddressM2M belongs to many persons through custom join table
type AddressM2M struct {
	gorm.Model
	Street string
	City   string
}

// SetupCustomJoinTable demonstrates using a custom join table model.
func SetupCustomJoinTable(db *gorm.DB) error {
	// Must be called before AutoMigrate
	err := db.SetupJoinTable(&PersonM2M{}, "Addresses", &PersonAddress{})
	if err != nil {
		return fmt.Errorf("failed to setup join table: %w", err)
	}

	// Now AutoMigrate will use the custom join table structure
	return db.AutoMigrate(&PersonM2M{}, &AddressM2M{}, &PersonAddress{})
}

// --- Custom Foreign Keys in Many-to-Many ---

// Profile with custom many2many foreign key configuration
type Profile struct {
	gorm.Model
	Refer uint   `gorm:"index:,unique"`
	Name  string
	Tags  []Tag `gorm:"many2many:profile_tags;foreignKey:Refer;joinForeignKey:ProfileReferID;References:TagRefer;joinReferences:TagReferID"`
}

// Tag with custom reference field
type Tag struct {
	gorm.Model
	TagRefer uint   `gorm:"index:,unique"`
	Name     string
}

// CreateProfileWithTags demonstrates custom foreign key configuration.
func CreateProfileWithTags(db *gorm.DB) error {
	profile := Profile{
		Refer: 1001,
		Name:  "Developer Profile",
		Tags: []Tag{
			{TagRefer: 2001, Name: "golang"},
			{TagRefer: 2002, Name: "backend"},
		},
	}

	return db.Create(&profile).Error
}
