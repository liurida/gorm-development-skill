package examples

import (
	"fmt"
	"time"

	"gorm.io/driver/sqlite"
	"gorm.io/gorm"
)

// --- Getting the Generic *sql.DB Object ---

// GetSqlDB demonstrates how to get the underlying *sql.DB object from GORM.
func GetSqlDB(db *gorm.DB) error {
	// db.DB() returns the underlying sql.DB object.
	sqlDB, err := db.DB()
	if err != nil {
		return fmt.Errorf("failed to get sql.DB object: %w", err)
	}

	// You can now use methods from the standard library's sql.DB.
	// Ping verifies the connection to the database is still alive.
	if err := sqlDB.Ping(); err != nil {
		return fmt.Errorf("failed to ping database: %w", err)
	}

	fmt.Println("Successfully pinged database using the generic *sql.DB object.")

	// sqlDB.Close() would close the database connection for the GORM instance.

	return nil
}

// --- Configuring the Connection Pool ---

// ConfigureConnectionPool shows how to set connection pool parameters.
func ConfigureConnectionPool(db *gorm.DB) error {
	sqlDB, err := db.DB()
	if err != nil {
		return err
	}

	// SetMaxIdleConns sets the maximum number of connections in the idle connection pool.
	// If n <= 0, no idle connections are retained.
	sqlDB.SetMaxIdleConns(10)

	// SetMaxOpenConns sets the maximum number of open connections to the database.
	// If n <= 0, then there is no limit on the number of open connections.
	sqlDB.SetMaxOpenConns(100)

	// SetConnMaxLifetime sets the maximum amount of time a connection may be reused.
	// Expired connections may be closed lazily before reuse.
	// If d <= 0, connections are not closed due to a connection's age.
	sqlDB.SetConnMaxLifetime(time.Hour)

	fmt.Println("Connection pool configured successfully.")
	return nil
}

// --- Checking Database Statistics ---

// ShowDBStats demonstrates how to retrieve database connection statistics.
func ShowDBStats(db *gorm.DB) error {
	sqlDB, err := db.DB()
	if err != nil {
		return err
	}

	stats := sqlDB.Stats()

	fmt.Println("--- Database Statistics ---")
	fmt.Printf("Open Connections: %d\n", stats.OpenConnections) // Number of open connections to the database.
	fmt.Printf("In Use: %d\n", stats.InUse)             // The number of connections currently in use.
	fmt.Printf("Idle: %d\n", stats.Idle)               // The number of idle connections.
	fmt.Printf("Wait Count: %d\n", stats.WaitCount)         // The total number of connections waited for.
	fmt.Printf("Wait Duration: %v\n", stats.WaitDuration)   // The total time blocked waiting for a new connection.
	fmt.Printf("Max Idle Closed: %d\n", stats.MaxIdleClosed) // The total number of connections closed due to SetMaxIdleConns.
	fmt.Printf("Max Lifetime Closed: %d\n", stats.MaxLifetimeClosed) // The total number of connections closed due to SetConnMaxLifetime.
	fmt.Println("---------------------------")

	return nil
}

// --- Transaction Caveat ---

// TransactionDB illustrates that db.DB() returns an error within a transaction.
func TransactionDB(db *gorm.DB) error {
	return db.Transaction(func(tx *gorm.DB) error {
		// Inside a transaction, tx does not represent the main connection pool.
		// Therefore, tx.DB() will return an error.
		sqlDB, err := tx.DB()
		if err != nil {
			fmt.Printf("Expected error inside transaction: %v\n", err)
			return nil // Returning nil to commit the transaction for this example.
		}

		// This part will not be reached.
		_ = sqlDB
		return errors.New("should have failed to get *sql.DB in transaction")
	})
}

func main() {
	db, err := gorm.Open(sqlite.Open("file::memory:?cache=shared"), &gorm.Config{})
	if err != nil {
		panic("failed to connect database")
	}

	if err := GetSqlDB(db); err != nil {
		panic(err)
	}

	if err := ConfigureConnectionPool(db); err != nil {
		panic(err)
	}

	if err := ShowDBStats(db); err != nil {
		panic(err)
	}

	if err := TransactionDB(db); err != nil {
		panic(err)
	}
}
