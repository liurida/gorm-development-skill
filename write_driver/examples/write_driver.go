
package examples

import (
	"gorm.io/gorm"
	"gorm.io/gorm/clause"
	"gorm.io/gorm/schema"
	"io"
)

// This example demonstrates the key components of a custom GORM driver.
// It is a simplified, non-functional example for illustrative purposes.

// 1. Define the Dialector struct
type MyDialector struct {
	// Driver-specific configurations would go here
}

// 2. Implement the Open function
func Open(dsn string) gorm.Dialector {
	return &MyDialector{}
}

// 3. Implement the gorm.Dialector interface

func (d MyDialector) Name() string {
	return "mydialect"
}

func (d MyDialector) Initialize(db *gorm.DB) error {
	// Register custom clause builders if needed
	db.ClauseBuilders["LIMIT"] = myCustomLimitBuilder

	// Register callbacks or perform other initializations
	return nil
}

func (d MyDialector) Migrator(db *gorm.DB) gorm.Migrator {
	// Return a custom migrator if necessary, or a default one
	return nil // Simplified for example
}

func (d MyDialector) DataTypeOf(field *schema.Field) string {
	// Map Go types to custom database types
	switch field.DataType {
	case schema.String:
		return "MY_TEXT_TYPE"
	case schema.Int:
		return "MY_INTEGER_TYPE"
	default:
		return string(field.DataType)
	}
}

func (d MyDialector) DefaultValueOf(field *schema.Field) clause.Expression {
	// Return default values for columns if the database has special keywords
	return clause.Expr{SQL: "MY_DEFAULT_KEYWORD"}
}

func (d MyDialector) BindVarTo(writer clause.Writer, stmt *gorm.Statement, v interface{}) {
	// Define how bind variables are written (e.g., ?, $1, @p1)
	writer.WriteString("?")
}

func (d MyDialector) QuoteTo(writer clause.Writer, str string) {
	// Define how to quote identifiers
	writer.WriteByte('`')
	writer.WriteString(str)
	writer.WriteByte('`')
}

func (d MyDialector) Explain(sql string, vars ...interface{}) string {
	// Return a human-readable version of the query
	return gorm.Explain(sql, vars...)
}

// 4. Implement custom clause builders if necessary

func myCustomLimitBuilder(c clause.Clause, builder clause.Builder) {
	if limit, ok := c.Expression.(clause.Limit); ok {
		if stmt, ok := builder.(*gorm.Statement); ok {
			// Example of a non-standard LIMIT clause: LIMIT <limit> AT <offset>
			if limit.Limit > 0 {
				stmt.WriteString("LIMIT ")
				stmt.AddVar(stmt, limit.Limit)
			}
			if limit.Offset > 0 {
				stmt.WriteString(" AT ")
				stmt.AddVar(stmt, limit.Offset)
			}
		}
	}
}

// This function shows how you might use the custom driver (conceptually).
func UseCustomDriver() {
	// The following code is conceptual and won't run without a real database driver.
	// db, err := gorm.Open(Open("my_dsn"), &gorm.Config{})
	// if err != nil {
	// 	panic("failed to connect database")
	// }
	//
	// // GORM will now use the custom dialect for all operations.
	// db.AutoMigrate(&MyModel{})
	// db.Create(&MyModel{...})
}

type MyModel struct {
	gorm.Model
	Name string
}
