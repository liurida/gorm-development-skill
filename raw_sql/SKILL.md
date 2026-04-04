---
name: gorm-raw-sql
description: Use when working with raw SQL, the SQL builder, named arguments, and other advanced SQL features in GORM.
---

# Raw SQL & SQL Builder

Reference: [GORM SQL Builder Documentation](https://gorm.io/docs/sql_builder.html)

## Raw SQL

GORM provides `Raw`, `Exec`, and `Scan` for working directly with SQL.

### `Raw` with `Scan`

Query with raw SQL and scan the results into a struct or primitive.

```go
type Result struct {
  ID   int
  Name string
  Age  int
}

var result Result
db.Raw("SELECT id, name, age FROM users WHERE id = ?", 3).Scan(&result)

var age int
db.Raw("SELECT SUM(age) FROM users WHERE role = ?", "admin").Scan(&age)

var users []User
db.Raw("UPDATE users SET name = ? WHERE age = ? RETURNING id, name", "jinzhu", 20).Scan(&users) // PostgreSQL
```

### `Exec`

Execute raw SQL commands for operations that don't return rows (e.g., UPDATE, DELETE, DROP).

```go
db.Exec("DROP TABLE users")

db.Exec("UPDATE orders SET shipped_at = ? WHERE id IN ?", time.Now(), []int64{1, 2, 3})

// Exec with SQL Expression
db.Exec("UPDATE users SET money = ? WHERE name = ?", gorm.Expr("money * ? + ?", 10000, 1), "jinzhu")
```

## Named Arguments

GORM supports named arguments using `sql.Named`, `map[string]interface{}`, or a struct. This improves readability for queries with multiple parameters.

```go
// Using sql.NamedArg
db.Where("name1 = @name OR name2 = @name", sql.Named("name", "jinzhu")).Find(&user)
// SELECT * FROM `users` WHERE name1 = "jinzhu" OR name2 = "jinzhu"

// Using a map
db.Where("name1 = @name OR name2 = @name", map[string]interface{}{"name": "jinzhu2"}).First(&user)
// SELECT * FROM `users` WHERE name1 = "jinzhu2" OR name2 = "jinzhu2" ORDER BY `users`.`id` LIMIT 1

// Raw SQL with named arguments
db.Raw("SELECT * FROM users WHERE name1 = @name OR name2 = @name2",
   sql.Named("name", "jinzhu1"), sql.Named("name2", "jinzhu2")).Find(&user)

// Using a struct
type NamedArgument struct {
	Name string
	Name2 string
}
db.Raw("SELECT * FROM users WHERE name1 = @Name AND name2 = @Name2",
	 NamedArgument{Name: "jinzhu", Name2: "jinzhu2"}).Find(&user)
```

## DryRun Mode

Generate the SQL and its arguments without executing the query. This is useful for testing and debugging.

```go
stmt := db.Session(&gorm.Session{DryRun: true}).First(&user, 1).Statement

sql := stmt.SQL.String() // => SELECT * FROM `users` WHERE `id` = $1 ORDER BY `id`
vars := stmt.Vars        // => []interface{}{1}
```

## ToSQL

`ToSQL` provides a convenient way to get the generated SQL string for a GORM operation without executing it. It returns the interpolated SQL string, which is useful for debugging.

**Warning:** The generated SQL from `ToSQL` is for debugging only and does not provide the same safety guarantees against SQL injection as GORM's normal execution path.

```go
sql := db.ToSQL(func(tx *gorm.DB) *gorm.DB {
  return tx.Model(&User{}).Where("id = ?", 100).Limit(10).Order("age desc").Find(&[]User{})
})
// sql => SELECT * FROM "users" WHERE id = 100 AND "users"."deleted_at" IS NULL ORDER BY age desc LIMIT 10
```

## `Row` & `Rows`

Get results as `*sql.Row` or `*sql.Rows` for low-level access.

### `Row`

```go
var name string
var age int

row := db.Table("users").Where("name = ?", "jinzhu").Select("name", "age").Row()
row.Scan(&name, &age)
```

### `Rows`

```go
rows, err := db.Model(&User{}).Where("name = ?", "jinzhu").Select("name, age, email").Rows()
if err != nil {
    // handle error
}
defer rows.Close()

for rows.Next() {
  var user User
  // ScanRows scans a row into the user struct
  db.ScanRows(rows, &user)
  // do something with user
}
```

## Connection

Run multiple SQL statements within the same database TCP connection (but not in a transaction).

```go
db.Connection(func(tx *gorm.DB) error {
  // Set a session variable
  tx.Exec("SET my.role = ?", "admin")

  // Query using the same connection
  tx.First(&User{})

  return nil
})
```

## Clauses

GORM uses a clause-based SQL builder. You can interact with these clauses directly for advanced query construction.

```go
import "gorm.io/gorm/clause"

// Example: Using an INSERT clause with a modifier
db.Clauses(clause.Insert{Modifier: "IGNORE"}).Create(&user)
// INSERT IGNORE INTO users ...
```

### Optimizer & Index Hints

You can use clauses to pass hints to the database optimizer.

```go
import "gorm.io/hints"

// Optimizer Hint
db.Clauses(hints.New("MAX_EXECUTION_TIME(10000)")).Find(&User{})
// SELECT * /*+ MAX_EXECUTION_TIME(10000) */ FROM `users`

// Index Hint
db.Clauses(hints.UseIndex("idx_user_name")).Find(&User{})
// SELECT * FROM `users` USE INDEX (`idx_user_name`)

// Force Index
db.Clauses(hints.ForceIndex("idx_user_name")).Find(&User{})
// SELECT * FROM `users` FORCE INDEX (`idx_user_name`)
```

## When NOT to Use

- **Standard CRUD operations** - Use GORM's built-in methods; they're safer and more maintainable
- **When struct mapping works** - Prefer `Find`, `First`, `Create` over raw SQL for typical operations
- **Queries that GORM handles well** - `Where`, `Joins`, `Preload` cover most use cases without raw SQL
- **When you need hooks** - Raw SQL bypasses model hooks; use GORM methods if you need lifecycle callbacks
- **Cross-database portability** - Raw SQL may use database-specific syntax; GORM abstracts dialect differences

Use raw SQL only for: complex reporting queries, database-specific features, bulk operations, or when GORM's query builder is insufficient.

## Common Mistakes

| Mistake | Fix |
|---------|-----|
| SQL Injection Risk | Always use `?` for placeholders in raw SQL. Avoid `fmt.Sprintf` to build queries. |
| Forgetting to `defer rows.Close()` | This can lead to connection leaks. Always defer the close after getting `*sql.Rows`. |
| Using `ToSQL` in production | `ToSQL` is for debugging. It doesn't provide the same safety guarantees. |
| `DryRun` modifying state | `DryRun` does not execute, so it cannot return data or change the database. |

## Related Topics

- [Security](https://gorm.io/docs/security.html) - Best practices for avoiding SQL injection.
- [Performance](https://gorm.io/docs/performance.html) - Caching prepared statements.
- [Session](https://gorm.io/docs/session.html) - DryRun mode and other session configurations.
