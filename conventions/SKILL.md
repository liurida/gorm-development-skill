---
name: gorm-conventions
description: Use when understanding GORM's default conventions for primary keys, table names, column names, timestamp tracking, and how to override them with custom naming strategies.
---

# Conventions

GORM follows "convention over configuration" to minimize explicit setup. Understanding these conventions helps write concise models.

## Quick Reference

| Convention | Default Behavior | Override Method |
|------------|------------------|-----------------|
| Primary Key | Field named `ID` | `gorm:"primaryKey"` tag |
| Table Name | Pluralized snake_case | `TableName()` method or `db.Table()` |
| Column Name | snake_case | `gorm:"column:name"` tag |
| CreatedAt | Auto-set on create | `gorm:"autoCreateTime:false"` |
| UpdatedAt | Auto-set on create/update | `gorm:"autoUpdateTime:false"` |

## Primary Key

### Default Convention

GORM uses a field named `ID` as the default primary key:

```go
type User struct {
  ID   string // Used as primary key by default
  Name string
}
```

### Custom Primary Key

Use the `primaryKey` tag on any field:

```go
type Animal struct {
  ID     int64
  UUID   string `gorm:"primaryKey"` // UUID is now the primary key
  Name   string
  Age    int64
}
```

See the `composite_primary_key` skill for composite keys.

## Table Names

### Default Convention

Struct names are converted to `snake_case` and pluralized:

| Struct | Table Name |
|--------|------------|
| `User` | `users` |
| `UserProfile` | `user_profiles` |
| `APIToken` | `api_tokens` |

### Override with Tabler Interface

Implement `TableName()` to specify a custom table name:

```go
type Tabler interface {
  TableName() string
}

type User struct {
  gorm.Model
  Name string
}

func (User) TableName() string {
  return "profiles" // Maps to `profiles` table
}
```

**Note:** `TableName()` result is cached. For dynamic names, use Scopes:

```go
func UserTable(user User) func(tx *gorm.DB) *gorm.DB {
  return func(tx *gorm.DB) *gorm.DB {
    if user.Admin {
      return tx.Table("admin_users")
    }
    return tx.Table("users")
  }
}

db.Scopes(UserTable(user)).Create(&user)
```

### Temporary Table Override

Use `Table()` method for one-off operations:

```go
// Create table with different name
db.Table("deleted_users").AutoMigrate(&User{})

// Query from different table
var deletedUsers []User
db.Table("deleted_users").Find(&deletedUsers)
// SELECT * FROM deleted_users;

// Delete from different table
db.Table("deleted_users").Where("name = ?", "jinzhu").Delete(&User{})
// DELETE FROM deleted_users WHERE name = 'jinzhu';
```

## Column Names

### Default Convention

Field names are converted to `snake_case`:

```go
type User struct {
  ID        uint      // column: id
  Name      string    // column: name
  Birthday  time.Time // column: birthday
  CreatedAt time.Time // column: created_at
  UserName  string    // column: user_name
}
```

### Custom Column Name

Use the `column` tag:

```go
type Animal struct {
  AnimalID int64     `gorm:"column:beast_id"`         // beast_id
  Birthday time.Time `gorm:"column:day_of_the_beast"` // day_of_the_beast
  Age      int64     `gorm:"column:age_of_the_beast"` // age_of_the_beast
}
```

## Timestamp Tracking

### CreatedAt

Automatically set to current time when creating a record (if value is zero):

```go
db.Create(&user) // Sets CreatedAt to current time

// Provide explicit value to skip auto-setting
user2 := User{Name: "jinzhu", CreatedAt: time.Now()}
db.Create(&user2) // CreatedAt is NOT changed

// Update CreatedAt explicitly
db.Model(&user).Update("CreatedAt", time.Now())
```

### UpdatedAt

Automatically set on create and update:

```go
db.Save(&user) // Sets UpdatedAt to current time

db.Model(&user).Update("name", "jinzhu") // Updates UpdatedAt

// UpdateColumn does NOT update UpdatedAt
db.Model(&user).UpdateColumn("name", "jinzhu") // UpdatedAt unchanged
```

### Disable Auto-tracking

```go
type User struct {
  CreatedAt time.Time `gorm:"autoCreateTime:false"`
  UpdatedAt time.Time `gorm:"autoUpdateTime:false"`
}
```

### Unix Timestamps

Track as Unix seconds, milliseconds, or nanoseconds:

```go
type User struct {
  Created int64 `gorm:"autoCreateTime"`      // Unix seconds
  Updated int64 `gorm:"autoUpdateTime:milli"` // Unix milliseconds
  Modified int64 `gorm:"autoUpdateTime:nano"` // Unix nanoseconds
}
```

## NamingStrategy

For complete control, configure a custom `NamingStrategy`:

```go
import "gorm.io/gorm/schema"

db, err := gorm.Open(sqlite.Open("test.db"), &gorm.Config{
  NamingStrategy: schema.NamingStrategy{
    TablePrefix:   "prod_",  // Add prefix: User -> prod_users
    SingularTable: true,     // Singular names: User -> prod_user
    NoLowerCase:   true,     // Skip snake_case conversion
    NameReplacer:  strings.NewReplacer("CID", "Cid"), // Custom replacements
  },
})
```

**NamingStrategy Options:**

| Option | Description | Example |
|--------|-------------|---------|
| `TablePrefix` | Add prefix to all tables | `"t_"` -> `t_users` |
| `SingularTable` | Use singular table names | `true` -> `user` (not `users`) |
| `NoLowerCase` | Skip snake_case conversion | `UserName` stays `UserName` |
| `NameReplacer` | Custom string replacements | Replace `API` with `Api` |

## Comprehensive Example

```go
type User struct {
  gorm.Model               // ID, CreatedAt, UpdatedAt, DeletedAt
  Username  string         // column: username
  Email     string         // column: email
  CreatedAt time.Time      // auto-tracked (from gorm.Model)
  UpdatedAt time.Time      // auto-tracked (from gorm.Model)
}

type AdminUser struct {
  ID        uint      `gorm:"primaryKey"`
  AdminName string    `gorm:"column:name"` // custom column name
  Created   int64     `gorm:"autoCreateTime:milli"` // Unix milliseconds
}

func (AdminUser) TableName() string {
  return "admin_portal_users" // custom table name
}
```

## Common Mistakes

| Mistake | Fix |
|---------|-----|
| Expecting dynamic `TableName()` | Use `db.Table()` or Scopes for dynamic names |
| Using `UpdateColumn` and expecting `UpdatedAt` to change | Use `Update` instead |
| Forgetting `TableName()` is cached | Use Scopes for dynamic table routing |

## When NOT to Use

- **When working with a legacy database with inconsistent naming** - Instead of fighting conventions, use explicit tags (`gorm:"column:..."`, `TableName()`) for everything.
- **If your team's coding standards conflict with GORM's** - Use a custom `NamingStrategy` to enforce your team's conventions project-wide.
- **When table or column names must be dynamic** - Conventions are static; use `db.Table()` or Scopes for dynamic names based on runtime data.
- **If you prefer explicit configuration over implicit magic** - For maximum clarity, you can choose to define all names and relationships explicitly, ignoring conventions.

## References

- [Official GORM Documentation: Conventions](https://gorm.io/docs/conventions.html)
- [GORM Config: NamingStrategy](https://gorm.io/docs/gorm_config.html#naming_strategy)
- [Models: Time Tracking](https://gorm.io/docs/models.html#time_tracking)
- [Composite Primary Key](https://gorm.io/docs/composite_primary_key.html)
