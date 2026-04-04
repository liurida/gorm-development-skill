package examples

import (
	"fmt"
	"log"
	"os"
	"strings"
	"time"

	"gorm.io/driver/sqlite"
	"gorm.io/gorm"
	"gorm.io/gorm/logger"
	"gorm.io/gorm/schema"
)

// --- Comprehensive GORM Configuration ---

// GormUser is a simple model for demonstrating configuration effects.
type GormUser struct {
	gorm.Model
	Name string
}

// CreateCustomDB demonstrates initializing GORM with a comprehensive configuration.
func CreateCustomDB() (*gorm.DB, error) {
	// Custom logger configuration
	newLogger := logger.New(
		log.New(os.Stdout, "\r\n", log.LstdFlags),
		logger.Config{
			SlowThreshold:             200 * time.Millisecond,
			LogLevel:                  logger.Info,
			IgnoreRecordNotFoundError: true,
			Colorful:                  true,
		},
	)

	// Open a database connection with a custom gorm.Config
	db, err := gorm.Open(sqlite.Open("file::memory:?cache=shared"), &gorm.Config{
		// Performance optimizations
		SkipDefaultTransaction: true,
		PrepareStmt:            true,

		// Naming strategy for table and column names
		NamingStrategy: schema.NamingStrategy{
			TablePrefix:   "tbl_",   // e.g., `GormUser` becomes `tbl_gorm_users`
			SingularTable: false,  // e.g., `GormUser` becomes `gorm_users` (plural)
			NameReplacer:  strings.NewReplacer("Gorm", ""), // `GormUser` -> `User` -> `users`
		},

		// Logger configuration
		Logger: newLogger,

		// Custom function for `now`
		NowFunc: func() time.Time {
			return time.Now().UTC() // Always use UTC time
		},

		// Behavior flags
		DryRun:                               false, // Set to true to generate SQL without executing
		AllowGlobalUpdate:                    false, // Disallow updates/deletes without a WHERE clause
		DisableAutomaticPing:                 false, // GORM pings the DB on init to check availability
		DisableForeignKeyConstraintWhenMigrating: false, // Create foreign keys during AutoMigrate
	})

	if err != nil {
		return nil, err
	}

	return db, nil
}

// DemonstrateConfigEffects shows the effects of the custom configuration.
func DemonstrateConfigEffects() error {
	db, err := CreateCustomDB()
	if err != nil {
		return err
	}

	// 1. Naming Strategy Effect
	// The table name for GormUser will be `tbl_users`
	// (Gorm removed -> user -> users (plural) -> tbl_users (prefix))
	db.AutoMigrate(&GormUser{})
	if !db.Migrator().HasTable("tbl_users") {
		return fmt.Errorf("expected table 'tbl_users' to exist, but it doesn't")
	}
	fmt.Println("NamingStrategy effect: Table 'tbl_users' created successfully.")

	// 2. NowFunc Effect
	user := GormUser{Name: "test_user"}
	db.Create(&user)
	// Check if the CreatedAt timestamp is in UTC
	if user.CreatedAt.Location() != time.UTC {
		return fmt.Errorf("expected CreatedAt to be in UTC, but it was in %s", user.CreatedAt.Location())
	}
	fmt.Println("NowFunc effect: CreatedAt timestamp is in UTC.")

	// 3. AllowGlobalUpdate Effect
	// This will fail because AllowGlobalUpdate is false
	if err := db.Delete(&GormUser{}).Error; err == nil {
		return fmt.Errorf("expected global delete to fail, but it succeeded")
	}
	fmt.Println("AllowGlobalUpdate effect: Global delete was successfully prevented.")

	// 4. DryRun example (in a separate session)
	dryRunDB := db.Session(&gorm.Session{DryRun: true})
	result := dryRunDB.Delete(&GormUser{}, 1)
	fmt.Printf("DryRun SQL: %s\n", result.Statement.SQL.String())
	if result.RowsAffected != -1 { // DryRun sets RowsAffected to -1
		return fmt.Errorf("expected RowsAffected to be -1 in DryRun mode")
	}

	return nil
}

func main() {
	if err := DemonstrateConfigEffects(); err != nil {
		panic(err)
	}
}
