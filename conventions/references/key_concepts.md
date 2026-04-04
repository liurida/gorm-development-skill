# Key Concepts for GORM Conventions

This document provides detailed explanations of GORM's default conventions and how to override them.

## Overview

GORM follows the principle of "convention over configuration." It uses a set of default rules for mapping Go structs to database tables, which minimizes the amount of explicit configuration required.

## Primary Key

- **Convention**: A field named `ID` of an integer or string type is used as the primary key.
- **Override**: Use the `gorm:"primaryKey"` tag on any other field to designate it as the primary key.

```go
// `ID` is the primary key by default
type User struct {
    ID   uint
    Name string
}

// `UUID` is explicitly set as the primary key
type Product struct {
    ID   uint
    UUID string `gorm:"primaryKey"`
}
```

## Table Names

- **Convention**: The struct name is converted to `snake_case` and pluralized. For example, `UserModel` becomes `user_models`.
- **Override**: Implement the `Tabler` interface on your struct.

```go
// This struct will map to the `admin_users` table.
type AdminUser struct {
    gorm.Model
}

func (AdminUser) TableName() string {
    return "admin_users"
}
```

- **Dynamic Table Names**: Use `db.Table("table_name")` or `Scopes` for dynamic table names, as the `TableName()` method result is cached.

## Column Names

- **Convention**: The struct field name is converted to `snake_case`. For example, `UserName` becomes `user_name`.
- **Override**: Use the `gorm:"column:<name>"` tag.

```go
type User struct {
    // This field will map to the `user_email_address` column.
    EmailAddress string `gorm:"column:user_email_address"`
}
```

## Timestamp Tracking

GORM uses `CreatedAt` and `UpdatedAt` fields to automatically track creation and update times.

### `CreatedAt`

- **Convention**: A field named `CreatedAt` of type `time.Time` or `int` (for Unix timestamp) is automatically set to the current time when a record is first created.
- **Override**: You can disable this behavior with the `gorm:"autoCreateTime:false"` tag.

### `UpdatedAt`

- **Convention**: A field named `UpdatedAt` of type `time.Time` or `int` is automatically set to the current time whenever a record is created or updated via `db.Save()` or `db.Update()`.
- **Override**: You can disable this with the `gorm:"autoUpdateTime:false"` tag.

**Note**: `UpdateColumn` and `UpdateColumns` will **not** update the `UpdatedAt` field.

## `gorm.Model`

`gorm.Model` is a predefined struct that includes the most common convention fields:

```go
type Model struct {
    ID        uint `gorm:"primaryKey"`
    CreatedAt time.Time
    UpdatedAt time.Time
    DeletedAt gorm.DeletedAt `gorm:"index"` // For soft delete
}
```

Embedding `gorm.Model` in your structs is a quick way to include these conventional fields.

## Naming Strategy

For complete control over naming conventions, you can provide a custom `NamingStrategy` in the GORM config. This allows you to define rules for table names, column names, join table names, and more.

```go
import "gorm.io/gorm/schema"

db, err := gorm.Open(sqlite.Open("test.db"), &gorm.Config{
    NamingStrategy: schema.NamingStrategy{
        TablePrefix:   "prod_",  // Add a prefix to all table names
        SingularTable: true,     // Use singular table names (e.g., "user" instead of "users")
    },
})
```

This would map a `User` struct to a table named `prod_user`.

By understanding these conventions, you can write concise and readable GORM models, only adding explicit configuration when you need to deviate from the defaults.
