---
name: gorm-sql-builder
description: Use when executing raw SQL, building custom clauses, or working with low-level SQL operations in GORM.
---

# SQL Builder

## Raw SQL with Scan

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
```

## Exec for Write Operations

```go
db.Exec("DROP TABLE users")
db.Exec("UPDATE orders SET shipped_at = ? WHERE id IN ?", time.Now(), []int64{1, 2, 3})

// With SQL Expression
db.Exec("UPDATE users SET money = ? WHERE name = ?", gorm.Expr("money * ? + ?", 10000, 1), "jinzhu")
```

## Named Arguments

```go
// Using sql.Named
db.Where("name1 = @name OR name2 = @name", sql.Named("name", "jinzhu")).Find(&user)

// Using map
db.Where("name1 = @name OR name2 = @name", map[string]interface{}{"name": "jinzhu"}).Find(&user)

// Using struct
type NamedArg struct {
    Name  string
    Name2 string
}
db.Raw("SELECT * FROM users WHERE name1 = @Name AND name2 = @Name2",
    NamedArg{Name: "jinzhu", Name2: "jinzhu2"}).Find(&user)
```

## DryRun Mode

Generate SQL without executing.

```go
stmt := db.Session(&gorm.Session{DryRun: true}).First(&user, 1).Statement
stmt.SQL.String() //=> SELECT * FROM `users` WHERE `id` = $1 ORDER BY `id`
stmt.Vars         //=> []interface{}{1}
```

## ToSQL

Returns generated SQL for debugging (not escaped, use only for debugging).

```go
sql := db.ToSQL(func(tx *gorm.DB) *gorm.DB {
  return tx.Model(&User{}).Where("id = ?", 100).Limit(10).Order("age desc").Find(&[]User{})
})
// sql => SELECT * FROM "users" WHERE id = 100 AND "users"."deleted_at" IS NULL ORDER BY age desc LIMIT 10
```

## Row & Rows

```go
// Single row
row := db.Table("users").Where("name = ?", "jinzhu").Select("name", "age").Row()
row.Scan(&name, &age)

// Multiple rows
rows, err := db.Model(&User{}).Where("name = ?", "jinzhu").Select("name, age, email").Rows()
defer rows.Close()
for rows.Next() {
  rows.Scan(&name, &age, &email)
}

// ScanRows into struct
for rows.Next() {
  db.ScanRows(rows, &user)
}
```

## Connection (Same TCP Connection)

Run multiple SQL statements in the same connection (not a transaction).

```go
db.Connection(func(tx *gorm.DB) error {
  tx.Exec("SET my.role = ?", "admin")
  tx.First(&User{})
  return nil
})
```

## Clauses

GORM builds SQL using clauses. You can add custom clauses.

```go
// Insert with modifier
db.Clauses(clause.Insert{Modifier: "IGNORE"}).Create(&user)
// INSERT IGNORE INTO users (name,age...) VALUES ("jinzhu",18...);
```

## StatementModifier / Hints

```go
import "gorm.io/hints"

db.Clauses(hints.New("hint")).Find(&User{})
// SELECT * /*+ hint */ FROM `users`
```

## When NOT to Use

- **For standard CRUD operations** - GORM's high-level methods (`Create`, `Find`, `Update`) are safer, more readable, and more maintainable for common tasks.
- **When you need portability** - Raw SQL can be database-specific. GORM's standard methods abstract away these differences.
- **If you are not confident in writing secure SQL** - GORM's methods automatically handle SQL injection prevention. With raw SQL, you are responsible for properly parameterizing all user input.
- **When you need GORM hooks** - Raw SQL operations bypass GORM's model hooks (`BeforeCreate`, `AfterUpdate`, etc.).

## Clause Builder

Different databases may generate different SQL for the same operation.

```go
db.Offset(10).Limit(5).Find(&users)
// SQL Server: SELECT * FROM "users" OFFSET 10 ROW FETCH NEXT 5 ROWS ONLY
// MySQL:      SELECT * FROM `users` LIMIT 5 OFFSET 10
```
