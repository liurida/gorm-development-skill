
# Key Concepts for GORM Indexes

This document provides key concepts for creating and using database indexes with GORM.

## Overview

Indexes are crucial for database performance, as they speed up data retrieval operations. GORM allows you to define indexes directly in your model structs using tags.

## Basic Index

You can create a basic index on a field using the `index` tag.

```go
type User struct {
    Name string `gorm:"index"`
}
```

## Unique Index

To create a unique index, use the `uniqueIndex` tag. This ensures that all values in the column are distinct.

```go
type User struct {
    Email string `gorm:"uniqueIndex"`
}
```

## Composite Indexes

A composite index is an index on multiple columns. You can create one by using the same index name for multiple fields.

```go
type User struct {
    FirstName string `gorm:"index:idx_name"`
    LastName  string `gorm:"index:idx_name"`
}
```

To create a unique composite index, add the `unique` option to the tag.

```go
type User struct {
    FirstName string `gorm:"index:idx_name,unique"`
    LastName  string `gorm:"index:idx_name,unique"`
}
```

### Index Priority

The order of columns in a composite index matters. You can control the order using the `priority` option.

```go
type User struct {
    FirstName string `gorm:"index:idx_name,priority:2"`
    LastName  string `gorm:"index:idx_name,priority:1"`
}
// The index will be created on (last_name, first_name)
```

## Index Options

GORM supports several options for indexes, which can be specified in the tag.

- **`class`**: The index class (e.g., `FULLTEXT`).
- **`type`**: The index type (e.g., `btree`, `hash`).
- **`where`**: A partial index condition.
- **`sort`**: The sort order (`asc` or `desc`).
- **`length`**: The length of the indexed prefix.

```go
type User struct {
    Name string `gorm:"index:idx_name,class:FULLTEXT,type:btree,length:10"`
    Age  int    `gorm:"index:idx_age,where:age > 18"`
}
```

## Applying Indexes

The indexes defined in your models are created when you run `db.AutoMigrate()`.

```go
db.AutoMigrate(&User{})
```

It's important to design your indexes based on your application's query patterns to ensure optimal database performance.
