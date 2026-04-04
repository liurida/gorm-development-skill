---
name: gorm-settings
description: Use when you need to pass data between GORM hooks or to subsequent operations within a single session, such as passing an audit user ID to a BeforeCreate hook.
---

# Settings

GORM provides `Set`, `Get`, `InstanceSet`, `InstanceGet` methods to pass values to hooks or other methods.

## Set / Get

Use `Set()` to pass values available in all hooks, including nested associations:

```go
// Set a value before create
db.Set("my_value", 123).Create(&User{})

// Read in hook
func (u *User) BeforeCreate(tx *gorm.DB) error {
    myValue, ok := tx.Get("my_value")
    if ok {
        val := myValue.(int) // val = 123
    }
    return nil
}
```

Values set with `Set()` propagate to association hooks:

```go
func (c *CreditCard) BeforeCreate(tx *gorm.DB) error {
    myValue, ok := tx.Get("my_value")
    // ok => true, values propagate to associations
    return nil
}
```

## InstanceSet / InstanceGet

Use `InstanceSet()` for values scoped to the current statement only. Unlike `Set()`, these values do NOT propagate to association hooks:

```go
db.InstanceSet("my_value", 123).Create(&User{})

// In User hook: ok => true
// In CreditCard hook: ok => false (new Statement created)
```

## Table Options

Set table options during migrations:

```go
// Set MySQL/MariaDB table engine
db.Set("gorm:table_options", "ENGINE=InnoDB").AutoMigrate(&User{})

// Set table comment
db.Set("gorm:table_options", "COMMENT='User accounts'").AutoMigrate(&User{})
```

## When NOT to Use

- **For passing general application data** - GORM settings are for controlling database operations. Use `context.Context` for passing request-scoped data like trace IDs or user information through your application layers.
- **As a replacement for proper session management** - While you can store a user ID, don't use it as your primary means of authentication or authorization within hooks. That logic belongs in middleware or a service layer.
- **When a value needs to be passed between separate GORM operations** - `Set` and `InstanceSet` are scoped to a single chain of GORM calls. The values will not persist to the next independent `db.Create()` or `db.Find()` call.

## Common Use Cases

- Conditional hook logic (skip audit, soft delete bypass)
- Passing audit context (user ID, action)
- Multi-tenant isolation (tenant ID)
- Migration table options
