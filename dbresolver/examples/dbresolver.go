
package examples

import (
	"fmt"

	"gorm.io/driver/sqlite"
	"gorm.io/gorm"
	"gorm.io/plugin/dbresolver"
)

// Product model for dbresolver examples
type Product struct {
	gorm.Model
	Name  string
	Price float64
}

// SetupDBResolver sets up a GORM DB instance with the dbresolver plugin.
func SetupDBResolver() (*gorm.DB, error) {
	db, err := gorm.Open(sqlite.Open("file::memory:?cache=shared"), &gorm.Config{})
	if err != nil {
		return nil, err
	}

	// Let's assume we have a replica database (we'll use another in-memory db for this example)
	replicaDB := sqlite.Open("file::memory:?cache=shared")

	db.Use(dbresolver.Register(dbresolver.Config{
		// Source database for writes
		Sources: []gorm.Dialector{sqlite.Open("file::memory:?cache=shared")},
		// Replica database for reads
		Replicas: []gorm.Dialector{replicaDB},
		Policy:   dbresolver.RandomPolicy{},
	}))

	// Auto-migrate the schema for both databases
	if err := db.AutoMigrate(&Product{}); err != nil {
		return nil, err
	}

	return db, nil
}

// CreateProduct demonstrates a write operation.
func CreateProduct(db *gorm.DB) error {
	// This will use a source connection
	return db.Create(&Product{Name: "Laptop", Price: 1200.00}).Error
}

// GetProducts demonstrates a read operation.
func GetProducts(db *gorm.DB) ([]Product, error) {
	var products []Product
	// This will use a replica connection
	err := db.Find(&products).Error
	return products, err
}

// ForceWriteRead demonstrates forcing a read from the source database.
func ForceWriteRead(db *gorm.DB) (*Product, error) {
	var product Product
	// This will use a source connection because of Clauses(dbresolver.Write)
	err := db.Clauses(dbresolver.Write).First(&product).Error
	return &product, err
}

// ReadWriteInTransaction demonstrates that all operations in a transaction use the same connection.
func ReadWriteInTransaction(db *gorm.DB) error {
	return db.Transaction(func(tx *gorm.DB) error {
		// All operations inside this transaction will use the source connection
		if err := tx.Create(&Product{Name: "Mouse", Price: 25.00}).Error; err != nil {
			return err
		}

		var product Product
		if err := tx.First(&product).Error; err != nil {
			return err
		}

		fmt.Printf("Read product in transaction: %v\n", product.Name)
		return nil
	})
}
