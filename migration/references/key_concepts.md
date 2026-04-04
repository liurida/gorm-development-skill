# Key Concepts for GORM Migration

This document provides key concepts for database migrations in GORM.

## Auto Migration

GORM's `AutoMigrate` function is used to automatically migrate your schema to keep it up-to-date.

```go
db.AutoMigrate(&User{}, &Product{})
```

**Important Notes:**
- `AutoMigrate` will only create tables, add missing columns, and add missing indexes.
- It **will not** delete unused columns or indexes to protect your data.
- It will change a column's type if the size or precision has changed.

## Migrator Interface

For more control over migrations, GORM provides a `Migrator` interface.

### Table Operations

```go
// Create table
db.Migrator().CreateTable(&User{})

// Drop table
db.Migrator().DropTable(&User{})

// Check if table exists
hasTable := db.Migrator().HasTable(&User{})

// Rename table
db.Migrator().RenameTable(&User{}, &UserInfo{})
```

### Column Operations

```go
// Add column
db.Migrator().AddColumn(&User{}, "Name")

// Drop column
db.Migrator().DropColumn(&User{}, "Name")

// Alter column
db.Migrator().AlterColumn(&User{}, "Name")

// Check if column exists
hasColumn := db.Migrator().HasColumn(&User{}, "Name")

// Rename column
db.Migrator().RenameColumn(&User{}, "Name", "NewName")
```
