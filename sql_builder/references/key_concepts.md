# Key Concepts for GORM SQL Builder

This document provides detailed explanations of GORM's SQL builder and raw SQL capabilities.

## Raw SQL Queries

### Basic Raw Query with Scan

```go
type Result struct {
  ID   int
  Name string
  Age  int
}

var result Result
db.Raw("SELECT id, name, age FROM users WHERE id = ?", 3).Scan(&result)
```

### Scanning into Different Types

```go
// Single value
var count int
db.Raw("SELECT COUNT(*) FROM users").Scan(&count)

// Slice
var results []Result
db.Raw("SELECT * FROM users WHERE age > ?", 18).Scan(&results)

// Map
var result map[string]interface{}
db.Model(&User{}).First(&result)
```

## Exec for Write Operations

Use `Exec` for INSERT, UPDATE, DELETE, and DDL statements.

```go
// Basic exec
db.Exec("UPDATE users SET name = ? WHERE id = ?", "new_name", 1)

// With slice parameter (IN clause)
db.Exec("DELETE FROM users WHERE id IN ?", []int64{1, 2, 3})
```

### SQL Expressions

Use `gorm.Expr()` for database expressions.

```go
// Arithmetic expression
db.Exec("UPDATE users SET balance = ? WHERE id = ?",
    gorm.Expr("balance + ?", 100), 1)

// Complex expression
db.Exec("UPDATE products SET price = ? WHERE category = ?",
    gorm.Expr("price * ? * (1 - ?)", 1.1, 0.05), "electronics")
```

## Named Arguments

### sql.Named
```go
db.Where("name = @name OR email = @name", sql.Named("name", "jinzhu")).Find(&user)
```

### Map
```go
db.Where("name = @name AND age >= @age", map[string]interface{}{
    "name": "jinzhu",
    "age":  18,
}).Find(&user)
```

### Struct
```go
type NamedArgs struct {
    Name string
    Age  int
}

db.Raw("SELECT * FROM users WHERE name = @Name AND age = @Age",
    NamedArgs{Name: "jinzhu", Age: 18}).Scan(&users)
```

**Note**: Struct field names become `@FieldName` placeholders.

## DryRun Mode

Generate SQL without executing it.

```go
stmt := db.Session(&gorm.Session{DryRun: true}).First(&user, 1).Statement

sql := stmt.SQL.String()  // Generated SQL
vars := stmt.Vars         // Bound parameters
```

### Use Cases
- SQL debugging
- Query validation
- Generating SQL for logging
- Testing query generation

## ToSQL

Returns the complete SQL statement for debugging purposes.

```go
sql := db.ToSQL(func(tx *gorm.DB) *gorm.DB {
    return tx.Model(&User{}).Where("id = ?", 100).Find(&[]User{})
})
```

**Warning**: ToSQL escapes arguments for display but does NOT provide SQL injection protection. Use only for debugging.

## Row and Rows

### Single Row
```go
row := db.Table("users").Where("name = ?", "jinzhu").Select("name", "age").Row()
row.Scan(&name, &age)
```

### Multiple Rows
```go
rows, err := db.Model(&User{}).Where("age > ?", 18).Rows()
if err != nil {
    return err
}
defer rows.Close()

for rows.Next() {
    var name string
    var age int
    rows.Scan(&name, &age)
}
```

### ScanRows

Scan `*sql.Rows` into GORM models.

```go
rows, _ := db.Model(&User{}).Rows()
defer rows.Close()

for rows.Next() {
    var user User
    db.ScanRows(rows, &user)
    // user is now populated
}
```

## Connection

Run multiple SQL statements in the same TCP connection (not a transaction).

```go
db.Connection(func(tx *gorm.DB) error {
    // All operations use the same connection
    tx.Exec("SET @my_var = ?", "value")
    tx.Raw("SELECT @my_var").Scan(&result)
    return nil
})
```

### Use Cases
- Setting session variables
- Using temporary tables
- Connection-specific settings

## Clauses

GORM builds SQL using a clause system. Each API call adds clauses to a `*gorm.Statement`.

### How Clauses Work

```go
// db.First(&user, 1) adds these clauses:
clause.Select{Columns: []clause.Column{{Name: "*"}}}
clause.From{Tables: []clause.Table{{Name: "users"}}}
clause.Limit{Limit: 1}
clause.OrderBy{Columns: []clause.OrderByColumn{{Column: clause.PrimaryKey}}}
```

### Custom Clauses

```go
// INSERT IGNORE
db.Clauses(clause.Insert{Modifier: "IGNORE"}).Create(&user)

// INSERT ... ON DUPLICATE KEY UPDATE
db.Clauses(clause.OnConflict{
    Columns:   []clause.Column{{Name: "id"}},
    DoUpdates: clause.AssignmentColumns([]string{"name", "age"}),
}).Create(&user)
```

## Clause Builder

Different databases generate different SQL for the same operation.

```go
db.Offset(10).Limit(5).Find(&users)
```

| Database | Generated SQL |
|----------|--------------|
| MySQL | `SELECT * FROM users LIMIT 5 OFFSET 10` |
| SQL Server | `SELECT * FROM users OFFSET 10 ROW FETCH NEXT 5 ROWS ONLY` |
| PostgreSQL | `SELECT * FROM users LIMIT 5 OFFSET 10` |

### Custom Clause Builder

Drivers can register custom clause builders:

```go
// Example from sqlserver driver
func (dialector Dialector) ClauseBuilders() map[string]clause.ClauseBuilder {
    return map[string]clause.ClauseBuilder{
        "LIMIT": limitClauseBuilder,
    }
}
```

## StatementModifier

Interface for modifying statements before execution.

```go
type StatementModifier interface {
    ModifyStatement(*Statement)
}
```

### Hints Plugin

```go
import "gorm.io/hints"

// Optimizer hint
db.Clauses(hints.New("MAX_EXECUTION_TIME(10000)")).Find(&users)
// SQL: SELECT * /*+ MAX_EXECUTION_TIME(10000) */ FROM users

// Index hint
db.Clauses(hints.UseIndex("idx_name")).Find(&users)
// SQL: SELECT * FROM users USE INDEX (idx_name)

// Comment
db.Clauses(hints.Comment("request_id", "abc123")).Find(&users)
// SQL: SELECT * FROM users /* request_id:abc123 */
```

## Prepared Statements

GORM can cache prepared statements for performance.

```go
db, err := gorm.Open(mysql.Open(dsn), &gorm.Config{
    PrepareStmt: true,
})
```

### Benefits
- Reduced parsing overhead
- Improved query execution time
- Better database resource utilization

### Considerations
- Increases connection memory usage
- May cause issues with connection pooling
- Not all queries benefit equally

## Best Practices

1. **Use parameterized queries**: Always use `?` placeholders, never string concatenation
2. **Close Rows**: Always `defer rows.Close()` when using `Rows()`
3. **Check errors**: Raw SQL operations can fail silently
4. **Use DryRun for debugging**: Verify generated SQL before execution
5. **Prefer GORM API**: Use raw SQL only when necessary
6. **Use Connection for session state**: When you need connection-specific settings
