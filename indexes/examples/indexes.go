
package examples

import (
	"gorm.io/gorm"
)

// IndexUser model demonstrates different ways to define indexes.
type IndexUser struct {
	gorm.Model

	// Basic index
	Name string `gorm:"index"`

	// Unique index
	Email string `gorm:"uniqueIndex"`

	// Composite index on (FirstName, LastName)
	FirstName string `gorm:"index:idx_name"`
	LastName  string `gorm:"index:idx_name"`

	// Unique composite index with priority
	// The index will be on (Country, City)
	City    string `gorm:"index:idx_location,unique,priority:2"`
	Country string `gorm:"index:idx_location,unique,priority:1"`

	// Index with options
	Description string `gorm:"index:idx_desc,class:FULLTEXT,type:btree"`

	// Partial index (for supported databases like PostgreSQL)
	Active bool `gorm:"index:idx_active_users,where:active = true"`
}

// MultiIndexUser demonstrates multiple indexes on a single field.
type MultiIndexUser struct {
	gorm.Model
	OID          int64  `gorm:"index:idx_id;index:idx_oid,unique"`
	MemberNumber string `gorm:"index:idx_id"`
}

// ExpressionIndexUser demonstrates expression-based indexes.
type ExpressionIndexUser struct {
	gorm.Model
	Age int64 `gorm:"index:,expression:ABS(age)"`
}

// SharedIndexBase is embedded to demonstrate shared composite indexes.
// Using `composite` allows multiple embedding structs to have their own index names.
type SharedIndexBase struct {
	IndexA int `gorm:"index:,unique,composite:shared_idx"`
	IndexB int `gorm:"index:,unique,composite:shared_idx"`
}

// ProductWithSharedIndex embeds SharedIndexBase.
// The composite index will be named: idx_product_with_shared_index_shared_idx
type ProductWithSharedIndex struct {
	gorm.Model
	SharedIndexBase
	Name string
}

// OrderWithSharedIndex also embeds SharedIndexBase.
// The composite index will be named: idx_order_with_shared_index_shared_idx
type OrderWithSharedIndex struct {
	gorm.Model
	SharedIndexBase
	Total float64
}

// CreateIndexes demonstrates that indexes are created during auto-migration.
func CreateIndexes(db *gorm.DB) error {
	// AutoMigrate will create the table and all the defined indexes.
	return db.AutoMigrate(&IndexUser{})
}

// UsingCompositeIndex shows an example of a query that would benefit from a composite index.
func UsingCompositeIndex(db *gorm.DB) (*IndexUser, error) {
	var user IndexUser
	// This query can efficiently use the composite index 'idx_location'.
	err := db.Where("country = ? AND city = ?", "USA", "New York").First(&user).Error
	return &user, err
}

// CreateSharedIndexTables demonstrates tables with shared composite indexes.
func CreateSharedIndexTables(db *gorm.DB) error {
	// Each table gets its own uniquely named index from the embedded struct
	return db.AutoMigrate(&ProductWithSharedIndex{}, &OrderWithSharedIndex{})
}
