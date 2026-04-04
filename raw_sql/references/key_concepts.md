# Key Concepts for GORM Raw SQL & SQL Builder

This document provides key concepts for using raw SQL and the SQL builder in GORM.

## Executing Raw SQL

- `Exec`: Executes a raw SQL statement.
- `Raw`: Executes a raw SQL query and scans the result into a struct.

```go
// Exec
db.Exec("UPDATE users SET name = ? WHERE id = ?", "jinzhu", 1)

// Raw and Scan
var result Result
db.Raw("SELECT name, age FROM users WHERE id = ?", 1).Scan(&result)
```

## Named Arguments

GORM supports named arguments, which can make your queries more readable.

```go
import "database/sql"

db.Where("name = @name", sql.Named("name", "jinzhu")).First(&user)
```

## DryRun Mode

In `DryRun` mode, GORM generates the SQL statement without executing it. This is useful for testing and debugging.

```go
stmt := db.Session(&gorm.Session{DryRun: true}).First(&user, 1).Statement
sql := stmt.SQL.String() // Get the generated SQL
vars := stmt.Vars        // Get the query variables
```

## `ToSQL`

`ToSQL` is another way to get the generated SQL without executing the query.

```go
sql := db.ToSQL(func(tx *gorm.DB) *gorm.DB {
  return tx.Model(&User{}).Where("id = ?", 100).Limit(10).Find(&[]User{})
})
```

## `Row` & `Rows`

You can get the result as a `*sql.Row` or `*sql.Rows` to have more control over scanning the data.

```go
// Get a single row
row := db.Table("users").Where("name = ?", "jinzhu").Select("name, age").Row()
var name string
var age int
row.Scan(&name, &age)

// Get multiple rows
rows, _ := db.Model(&User{}).Where("name = ?", "jinzhu").Rows()
defer rows.Close()
for rows.Next() {
  var user User
  db.ScanRows(rows, &user)
  // ...
}
```
