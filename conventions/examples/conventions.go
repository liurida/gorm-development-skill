package examples

import (
	"fmt"
	"time"

	"gorm.io/driver/sqlite"
	"gorm.io/gorm"
	"gorm.io/gorm/schema"
)

// --- Default Conventions ---

// UserModel demonstrates GORM's default conventions.
type UserModel struct {
	ID        uint // `ID` is the default primary key.
	UserName  string
	CreatedAt time.Time // Automatically tracked
	UpdatedAt time.Time // Automatically tracked
}

// After migrating this struct, the table will be named `user_models`.
// The `UserName` field will be a column named `user_name`.

// --- Overriding Table Name ---

// AdminUser demonstrates overriding the table name by implementing the Tabler interface.
type AdminUser struct {
	gorm.Model
	Email string
}

// TableName overrides the default table name for the AdminUser struct.
func (AdminUser) TableName() string {
	return "admin_portal_users"
}

// DemonstrateTableName shows how table names are generated and overridden.
func DemonstrateTableName(db *gorm.DB) {
	db.AutoMigrate(&UserModel{}, &AdminUser{})

	// Check if tables exist
	if db.Migrator().HasTable(&UserModel{}) {
		fmt.Println("Table for UserModel exists (default: user_models).")
	}
	if db.Migrator().HasTable("user_models") {
		fmt.Println("Table 'user_models' confirmed to exist.")
	}

	if db.Migrator().HasTable(&AdminUser{}) {
		fmt.Println("Table for AdminUser exists (overridden: admin_portal_users).")
	}
	if db.Migrator().HasTable("admin_portal_users") {
		fmt.Println("Table 'admin_portal_users' confirmed to exist.")
	}
}

// --- Overriding Column Name ---

// Profile demonstrates overriding column names with the `column` tag.
type Profile struct {
	gorm.Model
	UserUUID   string `gorm:"column:user_id"`
	Bio        string `gorm:"column:biography"`
	Subscribed bool   // will be `subscribed` by default
}

// DemonstrateColumnName shows how column names are generated and overridden.
func DemonstrateColumnName(db *gorm.DB) {
	db.AutoMigrate(&Profile{})

	if db.Migrator().HasColumn(&Profile{}, "user_id") {
		fmt.Println("Column 'user_id' exists in profiles table.")
	}
	if db.Migrator().HasColumn(&Profile{}, "biography") {
		fmt.Println("Column 'biography' exists in profiles table.")
	}
}

// --- Disabling Timestamp Tracking ---

// Post demonstrates disabling automatic timestamp tracking.
type Post struct {
	gorm.Model
	Title     string
	Content   string
	CreatedAt time.Time `gorm:"autoCreateTime:false"`
	UpdatedAt int64     `gorm:"autoUpdateTime:false"` // Can be different type
}

// DemonstrateTimestampTracking shows default and disabled timestamp tracking.
func DemonstrateTimestampTracking(db *gorm.DB) {
	db.AutoMigrate(&UserModel{}, &Post{})

	// Default behavior: CreatedAt and UpdatedAt are set automatically.
	user := UserModel{UserName: "default_user"}
	db.Create(&user)
	fmt.Printf("Default user CreatedAt: %v\n", user.CreatedAt)

	// Disabled behavior: Fields are not set by GORM.
	post := Post{Title: "My First Post"}
	db.Create(&post)
	if post.CreatedAt.IsZero() {
		fmt.Println("Disabled autoCreateTime: CreatedAt is zero.")
	}
}

// --- Using a Custom Naming Strategy ---

// DemonstrateCustomNamingStrategy shows how to use a different naming convention.
func DemonstrateCustomNamingStrategy() {
	db, err := gorm.Open(sqlite.Open("file::memory:"), &gorm.Config{
		NamingStrategy: schema.NamingStrategy{
			TablePrefix:   "tbl_", // add table prefix
			SingularTable: true,   // use singular table name, e.g., "user" instead of "users"
		},
	})
	if err != nil {
		panic("failed to connect database")
	}

	db.AutoMigrate(&UserModel{})

	// The table name will now be "tbl_user_model"
	if db.Migrator().HasTable("tbl_user_model") {
		fmt.Println("Custom naming strategy applied: table 'tbl_user_model' exists.")
	}
}

func main() {
	db, err := gorm.Open(sqlite.Open("file::memory:?cache=shared"), &gorm.Config{})
	if err != nil {
		panic("failed to connect database")
	}

	fmt.Println("--- Testing Table Name Conventions ---")
	DemonstrateTableName(db)

	fmt.Println("\n--- Testing Column Name Conventions ---")
	DemonstrateColumnName(db)

	fmt.Println("\n--- Testing Timestamp Tracking ---")
	DemonstrateTimestampTracking(db)

	fmt.Println("\n--- Testing Custom Naming Strategy ---")
	DemonstrateCustomNamingStrategy()
}
