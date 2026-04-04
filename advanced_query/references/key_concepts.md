# Key Concepts for GORM Advanced Query

This document provides detailed explanations of advanced querying techniques in GORM.

## Smart Select Fields

GORM automatically selects only the fields present in the destination struct, reducing data transfer.

```go
type APIUser struct {
  ID   uint
  Name string
}

// Only selects id and name columns
db.Model(&User{}).Limit(10).Find(&APIUser{})
```

**QueryFields Mode**: When enabled, GORM uses explicit field names in SELECT statements.

```go
db, err := gorm.Open(sqlite.Open("gorm.db"), &gorm.Config{
  QueryFields: true,
})
// SQL: SELECT `users`.`name`, `users`.`age`, ... FROM `users`
```

## Locking Strategies

### FOR UPDATE
Locks rows for update, preventing other transactions from modifying them.

```go
db.Clauses(clause.Locking{Strength: "UPDATE"}).Find(&users)
```

### FOR SHARE
Allows other transactions to read but not modify locked rows.

```go
db.Clauses(clause.Locking{Strength: "SHARE"}).Find(&users)
```

### Locking Options
- **NOWAIT**: Fails immediately if lock unavailable
- **SKIP LOCKED**: Skips locked rows (useful for job queues)

## SubQueries

### In WHERE Clause
```go
db.Where("amount > (?)", db.Table("orders").Select("AVG(amount)")).Find(&orders)
```

### In FROM Clause
```go
db.Table("(?) as u", db.Model(&User{}).Select("name", "age")).Where("age = ?", 18).Find(&User{})
```

### Combining Multiple Subqueries
```go
subQuery1 := db.Model(&User{}).Select("name")
subQuery2 := db.Model(&Pet{}).Select("name")
db.Table("(?) as u, (?) as p", subQuery1, subQuery2).Find(&result)
```

## Group Conditions

Build complex boolean logic with nested `Where` and `Or` calls.

```go
// (A AND (B OR C)) OR (D AND E)
db.Where(
  db.Where("A").Where(db.Where("B").Or("C")),
).Or(
  db.Where("D").Where("E"),
).Find(&results)
```

## IN with Multiple Columns

```go
db.Where("(name, age, role) IN ?", [][]interface{}{
  {"jinzhu", 18, "admin"},
  {"jinzhu2", 19, "user"},
}).Find(&users)
```

## Named Arguments

Three ways to use named arguments:

1. **sql.Named**: `sql.Named("name", "value")`
2. **Map**: `map[string]interface{}{"name": "value"}`
3. **Struct**: Fields become `@FieldName` placeholders

## FirstOrInit vs FirstOrCreate

| Method | Database Write | Use Case |
|--------|----------------|----------|
| FirstOrInit | No | Initialize struct without saving |
| FirstOrCreate | Yes | Get or create record |

### Attrs vs Assign

- **Attrs**: Only used when creating (ignored if record found)
- **Assign**: Always applied to struct (updates if record found)

## FindInBatches

Processes large datasets in chunks to manage memory.

```go
db.FindInBatches(&results, 100, func(tx *gorm.DB, batch int) error {
  // batch is the current batch number (1-indexed)
  // tx.RowsAffected is the count in current batch
  return nil // return error to stop processing
})
```

## Scopes

Reusable query conditions defined as functions.

### Basic Scope
```go
func ActiveUsers(db *gorm.DB) *gorm.DB {
  return db.Where("active = ?", true)
}
```

### Parameterized Scope
```go
func OlderThan(age int) func(db *gorm.DB) *gorm.DB {
  return func(db *gorm.DB) *gorm.DB {
    return db.Where("age > ?", age)
  }
}
```

## Iteration with Rows

For processing records one at a time (memory efficient).

```go
rows, _ := db.Model(&User{}).Rows()
defer rows.Close()

for rows.Next() {
  var user User
  db.ScanRows(rows, &user)
}
```

## Query Hooks

The `AfterFind` hook runs after each record is retrieved.

```go
func (u *User) AfterFind(tx *gorm.DB) error {
  if u.Role == "" {
    u.Role = "user" // Set default
  }
  return nil
}
```

## Optimizer and Index Hints

```go
import "gorm.io/hints"

// Optimizer hint
db.Clauses(hints.New("MAX_EXECUTION_TIME(10000)")).Find(&users)

// Index hint
db.Clauses(hints.UseIndex("idx_user_name")).Find(&users)

// Force index for JOIN
db.Clauses(hints.ForceIndex("idx_user_name").ForJoin()).Find(&users)
```

## Pluck

Extracts a single column into a slice.

```go
var names []string
db.Model(&User{}).Pluck("name", &names)

// With Distinct
db.Model(&User{}).Distinct().Pluck("name", &names)
```

## Count

```go
var count int64

// Basic count
db.Model(&User{}).Count(&count)

// Count distinct
db.Model(&User{}).Distinct("name").Count(&count)

// Count with group (returns number of groups)
db.Model(&User{}).Group("role").Count(&count)
```
