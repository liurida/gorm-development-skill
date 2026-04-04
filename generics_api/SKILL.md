---
name: gorm-generics-api
description: Use when working with GORM's type-safe generics API (v1.30.0+) for CRUD operations, joins, preloads, and code generation.
---

# GORM Generics API

The generics API (GORM >= v1.30.0) provides type safety and reduces SQL pollution from reused `gorm.DB` instances.

## Basic CRUD Operations

```go
ctx := context.Background()

// Create
gorm.G[User](db).Create(ctx, &User{Name: "Alice"})
gorm.G[User](db).CreateInBatches(ctx, users, 10)

// Query
user, err := gorm.G[User](db).Where("name = ?", "Jinzhu").First(ctx)
users, err := gorm.G[User](db).Where("age <= ?", 18).Find(ctx)

// Update
gorm.G[User](db).Where("id = ?", u.ID).Update(ctx, "age", 18)
gorm.G[User](db).Where("id = ?", u.ID).Updates(ctx, User{Name: "Jinzhu", Age: 18})

// Delete
gorm.G[User](db).Where("id = ?", u.ID).Delete(ctx)
```

## Advanced Options

Pass clauses and hints as optional parameters.

```go
// OnConflict handling
err := gorm.G[Language](db, clause.OnConflict{DoNothing: true}).Create(ctx, &lang)

// Execution hints
err := gorm.G[User](db,
  hints.New("MAX_EXECUTION_TIME(100)"),
  hints.New("USE_INDEX(t1, idx1)"),
).Find(ctx)
// SQL: SELECT /*+ MAX_EXECUTION_TIME(100) USE_INDEX(t1, idx1) */ * FROM `users`

// DB Resolver - read from master
err := gorm.G[User](db, dbresolver.Write).Find(ctx)

// Get result metadata
result := gorm.WithResult()
err := gorm.G[User](db, result).CreateInBatches(ctx, &users, 2)
// result.RowsAffected, result.Result.LastInsertId()
```

## Enhanced Joins

```go
// Load users who have a company
users, err := gorm.G[User](db).Joins(clause.Has("Company"), nil).Find(ctx)

// Left Join with custom filter
user, err = gorm.G[User](db).Joins(clause.LeftJoin.Association("Company"), func(db gorm.JoinBuilder, joinTable, curTable clause.Table) error {
    db.Where(map[string]any{"name": company.Name})
    return nil
}).Where(map[string]any{"name": user.Name}).First(ctx)
```

## Enhanced Preload

```go
// Basic preload with conditions
users, err := gorm.G[User](db).Preload("Friends", func(db gorm.PreloadBuilder) error {
    db.Where("age > ?", 14)
    return nil
}).Where("age > ?", 18).Find(ctx)

// Nested preload with per-record limit
users, err = gorm.G[User](db).Preload("Friends", func(db gorm.PreloadBuilder) error {
    db.Select("id", "name").Order("age desc")
    return nil
}).Preload("Friends.Pets", func(db gorm.PreloadBuilder) error {
    db.LimitPerRecord(2)
    return nil
}).Find(ctx)
```

## Raw SQL with Generics

```go
users, err := gorm.G[User](db).Raw("SELECT name FROM users WHERE id = ?", user.ID).Find(ctx)
```

## Code Generator (Recommended)

For type-safe raw queries, use the GORM CLI tool.

```bash
# Install
go install gorm.io/cli/gorm@latest

# Generate from interface
gorm gen -i ./examples/example.go -o query
```

Define query interface:

```go
type Query[T any] interface {
    // SELECT * FROM @@table WHERE id=@id
    GetByID(id int) (T, error)
}
```

Use generated code:

```go
import "your_project/query"

company, err := query.Query[Company](db).GetByID(ctx, 10)
// SELECT * FROM `companies` WHERE id=10
```

## When NOT to Use

- **Projects not using Go 1.18+** - The generics API relies on Go generics, which were introduced in Go 1.18.
- **When you need `FirstOrCreate` or `Save`** - These methods were intentionally removed. Use `First` followed by `Create` for `FirstOrCreate`, and `Create`/`Updates` for `Save`.
- **If you prefer the chained method style** - The generics API is more functional. The traditional API might feel more fluent if you prefer long chains.
- **When working with a codebase that heavily uses the traditional API** - While they are compatible, mixing styles can lead to inconsistency. Stick to one style within a project or module.

## Key Differences from Traditional API

- Context is required for all operations
- Results are returned directly (not via pointer parameter)
- `FirstOrCreate` and `Save` removed (ambiguity/concurrency issues)
- Full backward compatibility with traditional API
