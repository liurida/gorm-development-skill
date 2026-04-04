# Key Concepts for GORM Settings

This document provides key concepts for using GORM's settings system to pass values between code and hooks.

## Overview

GORM provides four methods for passing values between different parts of your application:

| Method | Scope | Use Case |
|--------|-------|----------|
| `Set(key, value)` | All hooks, including associations | Pass values to any hook |
| `Get(key)` | Read values from Set() | Access values in hooks |
| `InstanceSet(key, value)` | Current statement only | Values for immediate model only |
| `InstanceGet(key)` | Read values from InstanceSet() | Access statement-scoped values |

## Set / Get

Use `Set()` to pass values that need to be available in all hooks, including those of associated models.

### Setting Values

```go
myValue := 123
db.Set("my_value", myValue).Create(&User{})
```

### Reading Values in Hooks

```go
func (u *User) BeforeCreate(tx *gorm.DB) error {
    myValue, ok := tx.Get("my_value")
    if ok {
        // myValue is interface{}, type assert as needed
        val := myValue.(int) // val = 123
    }
    return nil
}
```

### Propagation to Associations

Values set with `Set()` are available in ALL hooks, including nested associations:

```go
type User struct {
    gorm.Model
    CreditCard CreditCard
}

func (c *CreditCard) BeforeCreate(tx *gorm.DB) error {
    // Values from Set() ARE available here
    myValue, ok := tx.Get("my_value")
    // ok => true
    // myValue => 123
    return nil
}
```

## InstanceSet / InstanceGet

Use `InstanceSet()` when you only need the value in the immediate model's hooks, not in associated model hooks.

### Setting Values

```go
myValue := 123
db.InstanceSet("my_value", myValue).Create(&User{})
```

### Important Limitation

When creating associations, GORM creates a new `*Statement`, so `InstanceSet` values do NOT propagate:

```go
func (u *User) BeforeCreate(tx *gorm.DB) error {
    // Works: InstanceSet values available
    myValue, ok := tx.InstanceGet("my_value")
    // ok => true
    return nil
}

func (c *CreditCard) BeforeCreate(tx *gorm.DB) error {
    // Does NOT work: new Statement for associations
    myValue, ok := tx.InstanceGet("my_value")
    // ok => false
    // myValue => nil
    return nil
}
```

## Common Use Cases

### Table Options During Migration

```go
// Set MySQL/MariaDB table engine
db.Set("gorm:table_options", "ENGINE=InnoDB").AutoMigrate(&User{})

// Set table comment
db.Set("gorm:table_options", "COMMENT='User accounts'").AutoMigrate(&User{})
```

### Conditional Hook Logic

```go
func (o *Order) BeforeCreate(tx *gorm.DB) error {
    skipAudit, ok := tx.Get("skip_audit")
    if ok && skipAudit.(bool) {
        return nil // Skip audit logging
    }
    // Perform audit logging...
    return nil
}

// Usage
db.Set("skip_audit", true).Create(&order)
```

### Passing Context to Hooks

```go
// Pass audit information
db.Set("audit_user", currentUser.ID).
   Set("audit_action", "create").
   Create(&record)

// In hook
func (r *Record) BeforeCreate(tx *gorm.DB) error {
    userID, _ := tx.Get("audit_user")
    action, _ := tx.Get("audit_action")
    // Log audit trail...
    return nil
}
```

### Tenant Isolation in Multi-tenant Apps

```go
// Set tenant context
db.Set("tenant_id", tenantID).Create(&resource)

// In hook, ensure tenant isolation
func (r *Resource) BeforeCreate(tx *gorm.DB) error {
    tenantID, ok := tx.Get("tenant_id")
    if ok {
        r.TenantID = tenantID.(uint)
    }
    return nil
}
```

## When to Use Set vs InstanceSet

| Use `Set()` when | Use `InstanceSet()` when |
|------------------|--------------------------|
| Value needed in association hooks | Value only needed in immediate model |
| Value should propagate through nested creates | Performance optimization (smaller scope) |
| Building audit trails | Temporary computation values |
| Multi-tenant context | Single-model transformations |

## Best Practices

1. **Use descriptive keys** - Prefix with your domain (e.g., `"audit:user_id"`)
2. **Check ok value** - Always check if the value exists before using
3. **Type assert carefully** - Values are `interface{}`, handle type assertions safely
4. **Prefer Set() for safety** - Use `InstanceSet()` only when you're certain about scope
5. **Document settings** - Comment which settings your hooks expect
