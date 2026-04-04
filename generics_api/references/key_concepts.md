# Key Concepts for GORM Generics API

This document provides detailed explanations of GORM's generics API introduced in v1.30.0.

## Overview

The generics API provides:
- **Type safety**: Compile-time type checking
- **Reduced SQL pollution**: Prevents issues from reusing `*gorm.DB` instances
- **Context-first design**: All operations require `context.Context`
- **Backward compatibility**: Can be mixed with traditional API

## The `gorm.G[T]` Function

The entry point for generics API is `gorm.G[T](db, ...options)`.

```go
// Basic usage
gorm.G[User](db)

// With options (clauses, hints, etc.)
gorm.G[User](db, clause.OnConflict{DoNothing: true})
```

## CRUD Operations

### Create
```go
ctx := context.Background()

// Single record
gorm.G[User](db).Create(ctx, &User{Name: "Alice"})

// Multiple records in batches
gorm.G[User](db).CreateInBatches(ctx, users, 100)
```

### Read
```go
// First record matching condition
user, err := gorm.G[User](db).Where("name = ?", "Alice").First(ctx)

// All records matching condition
users, err := gorm.G[User](db).Where("age > ?", 18).Find(ctx)
```

**Key difference**: Results are returned directly, not via pointer parameter.

### Update
```go
// Single field
gorm.G[User](db).Where("id = ?", 1).Update(ctx, "name", "Bob")

// Multiple fields
gorm.G[User](db).Where("id = ?", 1).Updates(ctx, User{Name: "Bob", Age: 30})
```

### Delete
```go
gorm.G[User](db).Where("id = ?", 1).Delete(ctx)
```

## Removed APIs

The following APIs are intentionally removed in the generics version due to ambiguity or concurrency issues:

- **FirstOrCreate**: Ambiguous behavior in concurrent scenarios
- **Save**: Unclear whether it creates or updates

Use explicit Create/Update operations instead.

## Options Parameter

Pass clauses and plugins as variadic options.

### OnConflict
```go
// Do nothing on conflict
gorm.G[User](db, clause.OnConflict{DoNothing: true}).Create(ctx, &user)

// Update specific columns on conflict
gorm.G[User](db, clause.OnConflict{
  Columns:   []clause.Column{{Name: "email"}},
  DoUpdates: clause.AssignmentColumns([]string{"name", "age"}),
}).Create(ctx, &user)
```

### Hints
```go
import "gorm.io/hints"

gorm.G[User](db, hints.New("MAX_EXECUTION_TIME(100)")).Find(ctx)
```

### DB Resolver
```go
import "gorm.io/plugin/dbresolver"

// Force read from master
gorm.G[User](db, dbresolver.Write).Find(ctx)
```

### Result Metadata
```go
result := gorm.WithResult()
gorm.G[User](db, result).Create(ctx, &user)

rowsAffected := result.RowsAffected
lastID, _ := result.Result.LastInsertId()
```

## Enhanced Joins

### Basic Join
```go
// Inner join - only users with companies
users, _ := gorm.G[User](db).Joins(clause.Has("Company"), nil).Find(ctx)
```

### Join Types
```go
// Left join
gorm.G[User](db).Joins(clause.LeftJoin.Association("Company"), nil).Find(ctx)

// With custom filter on joined table
gorm.G[User](db).Joins(clause.LeftJoin.Association("Company"), func(db gorm.JoinBuilder, joinTable, curTable clause.Table) error {
    db.Where(map[string]any{"name": "ACME"})
    return nil
}).Find(ctx)
```

### Join with Subquery
```go
gorm.G[User](db).Joins(
  clause.LeftJoin.AssociationFrom("Company", gorm.G[Company](db).Select("Name")).As("t"),
  func(db gorm.JoinBuilder, joinTable, curTable clause.Table) error {
    db.Where("?.name = ?", joinTable, "ACME")
    return nil
  },
).Find(ctx)
```

## Enhanced Preload

### Basic Preload
```go
gorm.G[User](db).Preload("Company", nil).Find(ctx)
```

### Preload with Conditions
```go
gorm.G[User](db).Preload("Friends", func(db gorm.PreloadBuilder) error {
    db.Where("age > ?", 18)
    return nil
}).Find(ctx)
```

### Nested Preload
```go
gorm.G[User](db).Preload("Friends.Pets", nil).Find(ctx)
```

### LimitPerRecord

Limits related records per parent record (useful for "top N per group").

```go
gorm.G[User](db).Preload("Posts", func(db gorm.PreloadBuilder) error {
    db.Order("created_at DESC").LimitPerRecord(5) // Latest 5 posts per user
    return nil
}).Find(ctx)
```

## Raw SQL

```go
// Query into typed result
users, _ := gorm.G[User](db).Raw("SELECT * FROM users WHERE id = ?", 1).Find(ctx)

// Query into primitive type
count, _ := gorm.G[int](db).Raw("SELECT COUNT(*) FROM users").Find(ctx)
```

## Code Generator

For type-safe raw queries, use the GORM CLI tool.

### Installation
```bash
go install gorm.io/cli/gorm@latest
```

### Define Query Interface
```go
type Query[T any] interface {
    // SELECT * FROM @@table WHERE id=@id
    GetByID(id int) (T, error)

    // SELECT * FROM @@table
    //   {{if user.Name != ""}}
    //       WHERE name=@user.Name
    //   {{end}}
    QueryByUser(user User) ([]T, error)
}
```

### Generate Code
```bash
gorm gen -i ./query.go -o ./generated
```

### Use Generated Code
```go
import "project/generated"

user, _ := generated.Query[User](db).GetByID(ctx, 1)
```

## Template Syntax

The code generator supports template syntax:

| Syntax | Purpose |
|--------|---------|
| `@@table` | Current model's table name |
| `@@column` | Dynamic column name |
| `@value` | Parameter value |
| `{{if}}` | Conditional SQL |
| `{{for}}` | Loop for IN clauses |
| `{{where}}` | Smart WHERE clause |
| `{{set}}` | Smart SET clause for updates |

## Migration from Traditional API

| Traditional | Generics |
|-------------|----------|
| `db.First(&user)` | `user, _ := gorm.G[User](db).First(ctx)` |
| `db.Find(&users)` | `users, _ := gorm.G[User](db).Find(ctx)` |
| `db.Create(&user)` | `gorm.G[User](db).Create(ctx, &user)` |
| `db.Save(&user)` | Use explicit Create or Updates |
| `db.FirstOrCreate(&user, ...)` | Use explicit First + Create |

## Best Practices

1. **Always pass context**: Enables timeout and cancellation handling
2. **Use explicit operations**: Prefer Create/Updates over Save
3. **Leverage type safety**: Let the compiler catch type mismatches
4. **Use code generator**: For complex raw queries needing type safety
5. **Mix APIs as needed**: Traditional API still works alongside generics
