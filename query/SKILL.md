---
name: gorm-query
description: Use when querying records with GORM. Covers First/Take/Last/Find, Where conditions, Select, Order, Limit/Offset, Group/Having, Joins, and Scan.
---

# Query

Reference: [GORM Query](https://gorm.io/docs/query.html)

## Quick Reference

| Method | Purpose |
|--------|---------|
| `First` | First by PK, `ErrRecordNotFound` if none |
| `Take` | One record (no order), `ErrRecordNotFound` if none |
| `Last` | Last by PK desc, `ErrRecordNotFound` if none |
| `Find` | All matching records (slice) |
| `Scan` | Scan into struct (works with Raw SQL) |

## Single Record

```go
db.First(&user)           // ORDER BY id LIMIT 1
db.Take(&user)            // LIMIT 1 (no order)
db.Last(&user)            // ORDER BY id DESC LIMIT 1
db.First(&user, 10)       // WHERE id = 10
db.First(&user, "id = ?", "uuid-string") // String PK

// Check error
errors.Is(result.Error, gorm.ErrRecordNotFound)

// Avoid ErrRecordNotFound
db.Limit(1).Find(&user)   // Empty struct if not found
```

## Multiple Records

```go
db.Find(&users)                    // All records
db.Find(&users, []int{1,2,3})      // WHERE id IN (1,2,3)
```

## Where Conditions

### String

```go
db.Where("name = ?", "jinzhu").First(&user)
db.Where("name IN ?", []string{"a", "b"}).Find(&users)
db.Where("name LIKE ?", "%jin%").Find(&users)
db.Where("name = ? AND age >= ?", "jinzhu", 22).Find(&users)
db.Where("created_at BETWEEN ? AND ?", lastWeek, today).Find(&users)
```

### Struct & Map

```go
// Struct - ignores zero values
db.Where(&User{Name: "jinzhu", Age: 20}).First(&user)

// Map - includes zero values
db.Where(map[string]interface{}{"name": "jinzhu", "age": 0}).Find(&users)

// Specify struct fields to include zero values
db.Where(&User{Name: "jinzhu"}, "name", "Age").Find(&users)
```

### Inline Conditions

```go
db.Find(&user, "name = ?", "jinzhu")
db.Find(&users, User{Age: 20})
```

### Not & Or

```go
db.Not("name = ?", "jinzhu").First(&user)
db.Not([]int64{1,2,3}).First(&user)

db.Where("role = ?", "admin").Or("role = ?", "super").Find(&users)
db.Where("name = ?", "a").Or(User{Name: "b", Age: 18}).Find(&users)
```

## Select Fields

```go
db.Select("name", "age").Find(&users)
db.Table("users").Select("COALESCE(age,?)", 42).Rows()
```

## Order, Limit, Offset

```go
db.Order("age desc, name").Find(&users)
db.Limit(10).Offset(5).Find(&users)
db.Limit(-1).Find(&users) // Cancel limit
```

## Group & Having

```go
db.Model(&User{}).Select("name, sum(age) as total").Group("name").Having("total > ?", 100).Scan(&results)
```

## Distinct

```go
db.Distinct("name", "age").Find(&results)
```

## Joins

```go
// Manual join
db.Model(&User{}).Select("users.name, emails.email").
  Joins("left join emails on emails.user_id = users.id").Scan(&result)

// Joins preloading (association name)
db.Joins("Company").Find(&users)
db.InnerJoins("Company").Find(&users)
db.Joins("Company", db.Where(&Company{Alive: true})).Find(&users)
```

## Scan

```go
var result struct{ Name string; Age int }
db.Table("users").Select("name", "age").Where("id = ?", 1).Scan(&result)
db.Raw("SELECT name, age FROM users WHERE id = ?", 1).Scan(&result)
```

## Common Mistakes

| Mistake | Fix |
|---------|-----|
| `Find` without limit for single record | Use `First`, `Take`, or `Limit(1).Find` |
| Struct conditions ignore zero values | Use map or specify fields |
| Not checking `ErrRecordNotFound` | `errors.Is(result.Error, gorm.ErrRecordNotFound)` |
| Not closing `Rows()` | `defer rows.Close()` |

## When NOT to Use

- Complex analytical queries with CTEs/window functions - use raw SQL
- Bulk exports of millions of rows - use `db.Raw()` with streaming
- Database-specific optimizations - use raw SQL
- Performance-critical simple queries - raw SQL can be faster

## Advanced Patterns

See [Advanced Query](https://gorm.io/docs/advanced_query.html): SubQueries, Locking, FirstOrInit/FirstOrCreate, FindInBatches, Scopes, Pluck, Count.
