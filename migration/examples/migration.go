package examples

import "gorm.io/gorm"

// User is a basic GORM model for migration examples.
type User struct {
	gorm.Model
	Name string
}

// Product is another GORM model for migration examples.
type Product struct {
	gorm.Model
	Code  string
	Price uint
}

// AutoMigrate demonstrates GORM's auto migration feature.
func AutoMigrate(db *gorm.DB) {
	db.AutoMigrate(&User{}, &Product{})
}

// CreateTable demonstrates creating a table with the migrator.
func CreateTable(db *gorm.DB) {
	db.Migrator().CreateTable(&User{})
}

// DropTable demonstrates dropping a table with the migrator.
func DropTable(db *gorm.DB) {
	db.Migrator().DropTable(&User{})
}

// AddColumn demonstrates adding a column with the migrator.
func AddColumn(db *gorm.DB) {
	db.Migrator().AddColumn(&User{}, "Age")
}
