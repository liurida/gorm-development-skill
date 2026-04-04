---
name: gorm-hooks
description: Use when working with GORM hooks to add custom logic to CRUD operations, intercept database operations, or implement model-level validation and transformation.
---

# Hooks

Hooks are functions called before or after creation/querying/updating/deletion. If you define specified methods for a model, GORM calls them automatically. If any callback returns an error, GORM stops future operations and rolls back the current transaction.

**Hook method signature**: `func(*gorm.DB) error`

## Object Life Cycle

### Creating an Object

```go
// begin transaction
BeforeSave
BeforeCreate
// save before associations
// insert into database
// save after associations
AfterCreate
AfterSave
// commit or rollback transaction
```

**Example**:

```go
func (u *User) BeforeCreate(tx *gorm.DB) (err error) {
    u.UUID = uuid.New()

    if !u.IsValid() {
        err = errors.New("can't save invalid data")
    }
    return
}

func (u *User) AfterCreate(tx *gorm.DB) (err error) {
    if u.ID == 1 {
        tx.Model(u).Update("role", "admin")
    }
    return
}
```

### Updating an Object

```go
// begin transaction
BeforeSave
BeforeUpdate
// save before associations
// update database
// save after associations
AfterUpdate
AfterSave
// commit or rollback transaction
```

**Example**:

```go
func (u *User) BeforeUpdate(tx *gorm.DB) (err error) {
    if u.readonly() {
        err = errors.New("read only user")
    }
    return
}

func (u *User) AfterUpdate(tx *gorm.DB) (err error) {
    if u.Confirmed {
        tx.Model(&Address{}).Where("user_id = ?", u.ID).Update("verified", true)
    }
    return
}
```

### Deleting an Object

```go
// begin transaction
BeforeDelete
// delete from database
AfterDelete
// commit or rollback transaction
```

**Example**:

```go
func (u *User) BeforeDelete(tx *gorm.DB) (err error) {
    if u.Role == "admin" {
        return errors.New("admin users cannot be deleted")
    }
    return
}

func (u *User) AfterDelete(tx *gorm.DB) (err error) {
    if u.Confirmed {
        tx.Model(&Address{}).Where("user_id = ?", u.ID).Update("invalid", true)
    }
    return
}
```

### Querying an Object

```go
// load data from database
// Preloading (eager loading)
AfterFind
```

**Example**:

```go
func (u *User) AfterFind(tx *gorm.DB) (err error) {
    if u.MemberShip == "" {
        u.MemberShip = "user"
    }
    return
}
```

## Modify Current Operation

Within hooks, you can modify the current operation through `tx.Statement`:

```go
func (u *User) BeforeCreate(tx *gorm.DB) error {
    // Select specific fields for insertion
    tx.Statement.Select("Name", "Age")

    // Add ON CONFLICT clause
    tx.Statement.AddClause(clause.OnConflict{DoNothing: true})

    // tx is a new session mode with NewDB option
    // Operations based on it run in the same transaction
    // but without any current conditions
    var role Role
    err := tx.First(&role, "name = ?", user.Role).Error
    // SELECT * FROM roles WHERE name = "admin"

    return err
}
```

## Transaction Behavior

**Important**: Save/Delete operations in GORM run in transactions by default. Changes made in that transaction are not visible until committed. If you return any error in your hooks, the change will be rolled back.

```go
func (u *User) AfterCreate(tx *gorm.DB) (err error) {
    if !u.IsValid() {
        return errors.New("rollback invalid user")
    }
    return nil
}
```

## Complete Example

```go
package models

import (
    "errors"
    "time"
    "github.com/google/uuid"
    "gorm.io/gorm"
    "gorm.io/gorm/clause"
)

type User struct {
    ID        uint      `gorm:"primaryKey"`
    UUID      uuid.UUID `gorm:"type:uuid"`
    Name      string
    Email     string
    Age       int
    Role      string
    CreatedAt time.Time
    UpdatedAt time.Time
}

// BeforeSave runs before both Create and Update
func (u *User) BeforeSave(tx *gorm.DB) (err error) {
    if u.Email == "" {
        return errors.New("email is required")
    }
    return
}

// BeforeCreate runs before Create only
func (u *User) BeforeCreate(tx *gorm.DB) (err error) {
    u.UUID = uuid.New()
    u.Role = "member"
    return
}

// AfterCreate runs after Create only
func (u *User) AfterCreate(tx *gorm.DB) (err error) {
    // Create audit log
    return tx.Create(&AuditLog{
        Action:   "user_created",
        UserID:   u.ID,
        Metadata: fmt.Sprintf("User %s created", u.Name),
    }).Error
}

// BeforeUpdate runs before Update only
func (u *User) BeforeUpdate(tx *gorm.DB) (err error) {
    // Prevent role changes on normal updates
    if tx.Statement.Changed("Role") {
        return errors.New("role cannot be changed through normal update")
    }
    return
}

// AfterFind runs after queries
func (u *User) AfterFind(tx *gorm.DB) (err error) {
    // Set default role if empty
    if u.Role == "" {
        u.Role = "guest"
    }
    return
}

// BeforeDelete runs before Delete
func (u *User) BeforeDelete(tx *gorm.DB) (err error) {
    if u.Role == "admin" {
        return errors.New("cannot delete admin users")
    }
    return
}
```

## When NOT to Use

- **Performance-critical batch operations** - Hooks run for each record; use `SkipHooks: true` for bulk imports
- **Simple derived fields** - Compute in application code or use database defaults instead of hooks
- **External API calls** - Hooks run inside transactions; external calls can cause long-held locks and timeouts
- **Complex business logic** - Keep hooks simple; use service layer for multi-step workflows
- **When creating from maps** - Hooks aren't invoked for `db.Model(&User{}).Create(map[string]interface{}{...})`
- **Audit logging with high volume** - Consider database triggers or async event queues instead

## Quick Reference

| Hook | Triggered On | Use Cases |
|------|--------------|-----------|
| `BeforeSave` | Create, Update | Validation, normalization |
| `BeforeCreate` | Create | UUID generation, defaults |
| `AfterCreate` | Create | Audit logs, notifications |
| `AfterSave` | Create, Update | Post-save processing |
| `BeforeUpdate` | Update | Validation, change prevention |
| `AfterUpdate` | Update | Cascading updates, audit |
| `BeforeDelete` | Delete | Soft delete, authorization |
| `AfterDelete` | Delete | Cleanup, audit logs |
| `AfterFind` | Query | Default values, computed fields |

## Reference

- Official Docs: https://gorm.io/docs/hooks.html
- Callbacks (Plugin): https://gorm.io/docs/write_plugins.html
