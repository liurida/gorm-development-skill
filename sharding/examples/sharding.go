
package examples

import (
	"fmt"

	"gorm.io/driver/sqlite"
	"gorm.io/gorm"
	"gorm.io/sharding"
)

// ShardOrder model for sharding examples. The actual table will be orders_0, orders_1, etc.	ype ShardOrder struct {
	gorm.Model
	UserID int64
	Amount float64
}

// SetupSharding configures the GORM sharding middleware.
func SetupSharding() (*gorm.DB, error) {
	// We use an in-memory sqlite database for this example.
	// In a real application, this would be your actual database connection.
	db, err := gorm.Open(sqlite.Open("file::memory:?cache=shared"), &gorm.Config{})
	if err != nil {
		return nil, err
	}

	// Configure the sharding middleware
	// We'll shard the 'orders' table into 4 shards based on the 'user_id'.
	shardingMiddleware := sharding.Register(sharding.Config{
		ShardingKey:    "user_id",
		NumberOfShards: 4, // Keep it small for the example
	}, "orders")      // Specify the table name to shard

	if err := db.Use(shardingMiddleware); err != nil {
		return nil, err
	}

	// You need to manually create the sharded tables in your database.
	// GORM's AutoMigrate won't do this automatically for sharded tables.
	for i := 0; i < 4; i++ {
		if err := db.Exec(fmt.Sprintf("CREATE TABLE orders_%d (id INTEGER PRIMARY KEY, user_id BIGINT, amount REAL, created_at DATETIME, updated_at DATETIME, deleted_at DATETIME);", i)).Error; err != nil {
			return nil, err
		}
	}

	return db, nil
}

// CreateShardedOrder demonstrates creating a record in a sharded table.
func CreateShardedOrder(db *gorm.DB, userID int64, amount float64) error {
	// The sharding middleware will hash the userID and determine the correct shard.
	// For example, if userID is 1, it might go to orders_1 (1 % 4).
	order := ShardOrder{UserID: userID, Amount: amount}
	return db.Table("orders").Create(&order).Error
}

// GetShardedOrdersForUser demonstrates querying a sharded table.
func GetShardedOrdersForUser(db *gorm.DB, userID int64) ([]ShardOrder, error) {
	var orders []ShardOrder
	// The WHERE clause MUST contain the sharding key ('user_id').
	// The middleware will route this query to the correct shard table.
	err := db.Table("orders").Where("user_id = ?", userID).Find(&orders).Error
	return orders, err
}

// FailWithoutShardingKey demonstrates that queries will fail without the sharding key.
func FailWithoutShardingKey(db *gorm.DB) error {
	var orders []ShardOrder
	// This query will fail because the sharding key ('user_id') is not in the WHERE clause.
	// The middleware doesn't know which shard to query.
	return db.Table("orders").Where("amount > ?", 100).Find(&orders).Error
}
