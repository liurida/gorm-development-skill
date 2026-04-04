---
name: gorm-write-driver
description: Use when writing a custom GORM database driver, implementing the Dialector interface, or adding support for unsupported databases.
---

# Writing a Database Driver

GORM offers built-in support for `SQLite`, `MySQL`, `Postgres`, `SQLServer`, and `ClickHouse`. To integrate GORM with unsupported databases or databases with unique features, create a custom driver by implementing the `Dialector` interface.

## Dialector Interface

```go
type Dialector interface {
    Name() string                      // Returns the name of the database dialect
    Initialize(*DB) error              // Initializes the database connection
    Migrator(db *DB) Migrator          // Provides the database migration tool
    DataTypeOf(*schema.Field) string   // Determines the data type for a schema field
    DefaultValueOf(*schema.Field) clause.Expression // Provides default value for a field
    BindVarTo(writer clause.Writer, stmt *Statement, v interface{}) // Handles variable binding
    QuoteTo(clause.Writer, string)     // Manages quoting of identifiers
    Explain(sql string, vars ...interface{}) string // Formats SQL statements with variables
}
```

## Compatibility with MySQL or Postgres

For databases that closely resemble `MySQL` or `Postgres`, you can often use those dialects directly. However, if your database significantly deviates or offers additional features, develop a custom driver.

## Nested Transaction Support (SavePoints)

If your database supports savepoints, implement `SavePointerDialectorInterface`:

```go
type SavePointerDialectorInterface interface {
    SavePoint(tx *DB, name string) error     // Saves a savepoint within a transaction
    RollbackTo(tx *DB, name string) error    // Rolls back to the specified savepoint
}
```

## Custom Clause Builders

Custom clause builders extend query capabilities for database-specific operations.

### Step 1: Define Custom Clause Builder Function

```go
func MyCustomLimitBuilder(c clause.Clause, builder clause.Builder) {
    if limit, ok := c.Expression.(clause.Limit); ok {
        // Access limit values via limit.Limit and limit.Offset
        if limit.Limit != nil && *limit.Limit >= 0 {
            builder.WriteString("FETCH FIRST ")
            builder.AddVar(nil, *limit.Limit)
            builder.WriteString(" ROWS ONLY")
        }
        if limit.Offset > 0 {
            builder.WriteString(" OFFSET ")
            builder.AddVar(nil, limit.Offset)
        }
    }
}
```

### Step 2: Register the Clause Builder

```go
func (d *MyDialector) Initialize(db *gorm.DB) error {
    db.ClauseBuilders["LIMIT"] = MyCustomLimitBuilder
    // ... other initialization
    return nil
}
```

### Step 3: Use in Queries

```go
query := db.Model(&User{}).Limit(10).Offset(5)
result := query.Find(&results)
// SQL: SELECT * FROM users FETCH FIRST 10 ROWS ONLY OFFSET 5
```

## Complete Driver Example

```go
package mydriver

import (
    "database/sql"
    "gorm.io/gorm"
    "gorm.io/gorm/clause"
    "gorm.io/gorm/migrator"
    "gorm.io/gorm/schema"
)

type Dialector struct {
    DSN string
}

func Open(dsn string) gorm.Dialector {
    return &Dialector{DSN: dsn}
}

func (d *Dialector) Name() string {
    return "mydb"
}

func (d *Dialector) Initialize(db *gorm.DB) error {
    // Open the underlying SQL connection
    sqlDB, err := sql.Open("mydb", d.DSN)
    if err != nil {
        return err
    }
    db.ConnPool = sqlDB

    // Register custom clause builders if needed
    db.ClauseBuilders["LIMIT"] = MyCustomLimitBuilder

    // Register callbacks if needed
    // ...

    return nil
}

func (d *Dialector) Migrator(db *gorm.DB) gorm.Migrator {
    return migrator.Migrator{Config: migrator.Config{
        DB:        db,
        Dialector: d,
    }}
}

func (d *Dialector) DataTypeOf(field *schema.Field) string {
    switch field.DataType {
    case schema.Bool:
        return "BOOLEAN"
    case schema.Int, schema.Uint:
        return "INTEGER"
    case schema.Float:
        return "DOUBLE"
    case schema.String:
        return "VARCHAR(255)"
    case schema.Time:
        return "TIMESTAMP"
    case schema.Bytes:
        return "BLOB"
    default:
        return "VARCHAR(255)"
    }
}

func (d *Dialector) DefaultValueOf(field *schema.Field) clause.Expression {
    return clause.Expr{SQL: "NULL"}
}

func (d *Dialector) BindVarTo(writer clause.Writer, stmt *gorm.Statement, v interface{}) {
    writer.WriteByte('?')
}

func (d *Dialector) QuoteTo(writer clause.Writer, str string) {
    writer.WriteByte('"')
    writer.WriteString(str)
    writer.WriteByte('"')
}

func (d *Dialector) Explain(sql string, vars ...interface{}) string {
    return gorm.ExplainSQL(sql, nil, `"`, vars...)
}

// SavePoint support (optional)
func (d *Dialector) SavePoint(tx *gorm.DB, name string) error {
    return tx.Exec("SAVEPOINT " + name).Error
}

func (d *Dialector) RollbackTo(tx *gorm.DB, name string) error {
    return tx.Exec("ROLLBACK TO SAVEPOINT " + name).Error
}
```

## Usage

```go
import (
    "gorm.io/gorm"
    "mydriver"
)

func main() {
    db, err := gorm.Open(mydriver.Open("connection-string"), &gorm.Config{})
    if err != nil {
        panic(err)
    }

    // Use GORM as normal
    db.AutoMigrate(&User{})
    db.Create(&User{Name: "jinzhu"})
}
```

## When NOT to Use

- **If an official or community driver already exists** - Always check for an existing driver before writing your own. The official drivers are well-tested and maintained.
- **For minor SQL syntax differences** - If your database is mostly compatible with MySQL or PostgreSQL, you might be able to use their drivers directly, perhaps with some custom callbacks, instead of writing a full driver.
- **If you only need a few custom features** - It might be simpler to use raw SQL for specific, non-standard queries rather than implementing a full driver.
- **Without a deep understanding of the target database** - Writing a driver requires in-depth knowledge of the database's SQL dialect, data types, and transaction behavior.

## Reference

- Official Docs: https://gorm.io/docs/write_driver.html
- MySQL Driver Source: https://github.com/go-gorm/mysql
- PostgreSQL Driver Source: https://github.com/go-gorm/postgres
- SQLite Driver Source: https://github.com/go-gorm/sqlite
- SQL Server Driver Source: https://github.com/go-gorm/sqlserver
- ClickHouse Driver Source: https://github.com/go-gorm/clickhouse
