package examples

import (
	"context"
	"database/sql/driver"
	"encoding/json"
	"errors"
	"fmt"

	"gorm.io/driver/sqlite"
	"gorm.io/gorm"
	"gorm.io/gorm/schema"
)

// --- 1. Basic Custom Type using sql.Scanner and driver.Valuer ---

// JSONB is a custom type for storing JSON data.
// It implements sql.Scanner and driver.Valuer to tell GORM how to read/write it.
type JSONB json.RawMessage

// Scan implements the sql.Scanner interface.
func (j *JSONB) Scan(value interface{}) error {
	bytes, ok := value.([]byte)
	if !ok {
		return errors.New("Scan source was not []byte")
	}
	*j = JSONB(bytes)
	return nil
}

// Value implements the driver.Valuer interface.
func (j JSONB) Value() (driver.Value, error) {
	if len(j) == 0 {
		return nil, nil
	}
	return json.RawMessage(j).MarshalJSON()
}

// GormDataType tells GORM what the general data type is.
func (j JSONB) GormDataType() string {
	return "json"
}

// UserWithJSONB has a field of our custom JSONB type.
type UserWithJSONB struct {
	gorm.Model
	Name       string
	Attributes JSONB
}

// DemonstrateScannerValuer shows how to create and read a record with a custom type.
func DemonstrateScannerValuer(db *gorm.DB) error {
	db.AutoMigrate(&UserWithJSONB{})

	// Create a user with JSON attributes
	attributes := JSONB(`{"role": "admin", "permissions": ["read", "write"]}`)
	user := UserWithJSONB{Name: "json_user", Attributes: attributes}
	if err := db.Create(&user).Error; err != nil {
		return err
	}

	// Find the user and verify the attributes
	var foundUser UserWithJSONB
	db.First(&foundUser, user.ID)

	var decoded map[string]interface{}
	json.Unmarshal(foundUser.Attributes, &decoded)

	fmt.Printf("Found user with role: %s\n", decoded["role"])
	return nil
}

// --- 2. Dialect-Specific Data Type ---

// MyJSON is another custom JSON type that returns different DB types based on the dialect.
type MyJSON json.RawMessage

// Scan and Value are similar to JSONB...
func (j *MyJSON) Scan(value interface{}) error { /* ... */ return nil }
func (j MyJSON) Value() (driver.Value, error)  { /* ... */ return nil, nil }

// GormDBDataType returns a database-specific type.
func (j MyJSON) GormDBDataType(db *gorm.DB, field *schema.Field) string {
	switch db.Dialector.Name() {
	case "mysql":
		return "JSON"
	case "postgres":
		return "JSONB"
	case "sqlite":
		return "TEXT" // SQLite doesn't have a native JSON type
	}
	return ""
}

// --- 3. Custom Type with SQL Expression ---

// Location is a custom type that uses a SQL function to be stored.
type Location struct {
	X, Y int
}

// GormDataType for Location.
func (loc Location) GormDataType() string {
	return "geometry"
}

// GormValue implements GormValuerInterface to generate a SQL expression.
func (loc Location) GormValue(ctx context.Context, db *gorm.DB) gorm.Clause {
	return gorm.Expr("ST_PointFromText(?)", fmt.Sprintf("POINT(%d %d)", loc.X, loc.Y))
}

// Scan for Location (simplified for example).
func (loc *Location) Scan(v interface{}) error {
	// In a real application, you would parse the WKT format from the database.
	return nil
}

// UserWithLocation has a field of the custom Location type.
type UserWithLocation struct {
	gorm.Model
	Name     string
	Location Location
}

// DemonstrateSQLExpression shows how a custom type can be saved using a SQL function.
func DemonstrateSQLExpression(db *gorm.DB) error {
	// Note: This requires a database with GIS extensions like PostGIS.
	// The following is a conceptual demonstration.
	dryRunDB := db.Session(&gorm.Session{DryRun: true})
	user := UserWithLocation{Name: "gis_user", Location: Location{X: 100, Y: 200}}
	result := dryRunDB.Create(&user)

	fmt.Println("SQL generated for custom type with GormValue:")
	fmt.Println(result.Statement.SQL.String())

	return nil
}

func main() {
	db, err := gorm.Open(sqlite.Open("file::memory:?cache=shared"), &gorm.Config{})
	if err != nil {
		panic("failed to connect database")
	}

	fmt.Println("--- Testing sql.Scanner and driver.Valuer ---")
	if err := DemonstrateScannerValuer(db); err != nil {
		panic(err)
	}

	fmt.Println("\n--- Testing SQL Expression with GormValuer ---")
	if err := DemonstrateSQLExpression(db); err != nil {
		panic(err)
	}
}
