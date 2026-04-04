---
name: gorm-generics
description: Use when working with GORM to leverage type-safe, generic APIs for common database operations like Create, Read, Update, and Delete.
---

# GORM Generics API

## Overview

GORM has introduced a new Generics API that provides enhanced type safety and a more fluent interface for database operations. This skill provides a reference for using the new API with examples.

The new API is accessed through `gorm.G[T](db)`, where `T` is your model type.

## Basic CRUD Operations

### Create

Create a single record:

```go
import "gorm.io/gorm"

user := User{Name: "TestGenericsCreate", Age: 18}
err := gorm.G[User](DB).Create(context.Background(), &user)
```

Create records in batches:

```go
import "gorm.io/gorm"

users := []User{
    {Name: "GenericsCreateInBatches1"},
    {Name: "GenericsCreateInBatches2"},
    {Name: "GenericsCreateInBatches3"},
}

err := gorm.G[User](DB).CreateInBatches(context.Background(), &users, 2)
```

### Read

Find a single record:

```go
import "gorm.io/gorm"

// Get the first record ordered by primary key
user, err := gorm.G[User](DB).Where("name = ?", "Jinzhu").First(context.Background())

// Get one record, no specified order
user, err := gorm.G[User](DB).Where("name = ?", "Jinzhu").Take(context.Background())

// Get the last record ordered by primary key
user, err := gorm.G[User](DB).Where("name = ?", "Jinzhu").Last(context.Background())
```

Find multiple records:

```go
import "gorm.io/gorm"

users, err := gorm.G[User](DB).Where("age <= ?", 18).Find(context.Background())
```

### Update

Update a single column:

```go
import "gorm.io/gorm"

rows, err := gorm.G[User](DB).Where("id = ?", u.ID).Update(context.Background(), "age", 18)
```

Update multiple columns:

```go
import "gorm.io/gorm"

rows, err := gorm.G[User](DB).Where("id = ?", u.ID).Updates(context.Background(), User{Name: "Jinzhu", Age: 18})
```

### Delete

Delete records:

```go
import "gorm.io/gorm"

rows, err := gorm.G[User](DB).Where("id = ?", u.ID).Delete(context.Background())
```

## Advanced Queries

### Joins

```go
import (
    "gorm.io/gorm"
    "gorm.io/gorm/clause"
)

// Inner Join
result, err := gorm.G[User](DB).Joins(clause.Has("Company"), func(db gorm.JoinBuilder, joinTable clause.Table, curTable clause.Table) error {
    db.Where("?.name = ?", joinTable, "MyCompany")
    return nil
}).First(context.Background())

// Left Join
result, err = gorm.G[User](DB).Joins(clause.LeftJoin.Association("Company"), nil).Where("name = ?", "jinzhu").First(context.Background())
```

### Preload

Eager load associations:

```go
import "gorm.io/gorm"

// Preload Company and Pets
users, err := gorm.G[User](DB).Preload("Company", nil).Preload("Pets", nil).Where("name = ?", "jinzhu").First(context.Background())

// Preload with conditions
users, err = gorm.G[User](DB).Preload("Pets", func(db gorm.PreloadBuilder) error {
    db.Where("age > ?", 10)
    return nil
}).Find(context.Background())
```

### Upsert (OnConflict)

```go
import (
    "gorm.io/gorm"
    "gorm.io/gorm/clause"
)

// Do nothing on conflict
err := gorm.G[Language](DB, clause.OnConflict{DoNothing: true}).Create(ctx, &lang)

// Update on conflict
err := gorm.G[Language](DB, clause.OnConflict{
  Columns:   []clause.Column{{Name: "code"}},
  DoUpdates: clause.Assignments(map[string]interface{}{"name": "upsert-new"}),
}).Create(ctx, &lang3)
```

## Raw SQL

You can still use Raw SQL with the generics API.

```go
import "gorm.io/gorm"

users, err := gorm.G[User](DB).Raw("SELECT name FROM users WHERE id = ?", 1).Find(context.Background())
```
