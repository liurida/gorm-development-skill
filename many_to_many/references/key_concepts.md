# Key Concepts for Many To Many

A `many to many` association is used when a record can be associated with multiple records of another model, and vice-versa. GORM uses a **join table** to manage these relationships.

## Basic Declaration

Use the `many2many` tag to specify the join table name.

```go
// User has and belongs to many Languages
type User struct {
  gorm.Model
  Languages []Language `gorm:"many2many:user_languages;"`
}

type Language struct {
  gorm.Model
  Name string
}
```

GORM's `AutoMigrate` automatically creates the `user_languages` join table.

## Bidirectional (Back-Reference)

Both models can reference each other using the same join table.

```go
type User struct {
  gorm.Model
  Languages []*Language `gorm:"many2many:user_languages;"`
}

type Language struct {
  gorm.Model
  Name  string
  Users []*User `gorm:"many2many:user_languages;"`
}
```

## Self-Referential Many-to-Many

Models can reference themselves for relationships like friends or followers.

```go
type User struct {
  gorm.Model
  Friends []*User `gorm:"many2many:user_friends"`
}
```

## Custom Foreign Keys

Four tags control join table foreign key configuration:

| Tag | Purpose |
|-----|---------|
| `foreignKey` | Source model's reference field |
| `joinForeignKey` | Join table's column referencing source |
| `References` | Target model's reference field |
| `joinReferences` | Join table's column referencing target |

```go
type User struct {
  gorm.Model
  Profiles []Profile `gorm:"many2many:user_profiles;foreignKey:Refer;joinForeignKey:UserReferID;References:UserRefer;joinReferences:ProfileRefer"`
  Refer    uint      `gorm:"index:,unique"`
}

type Profile struct {
  gorm.Model
  UserRefer uint `gorm:"index:,unique"`
}
```

**Important**: Fields referenced by foreign keys should have unique index tags.

## Custom Join Table

Join tables can be full models with soft delete, hooks, and additional fields.

```go
type PersonAddress struct {
  PersonID  int `gorm:"primaryKey"`
  AddressID int `gorm:"primaryKey"`
  CreatedAt time.Time
  DeletedAt gorm.DeletedAt
  IsPrimary bool  // Additional custom field
}

// Setup must be called before AutoMigrate
db.SetupJoinTable(&Person{}, "Addresses", &PersonAddress{})
```

**Note**: Custom join tables require composite primary keys or composite unique index for foreign keys.

## Retrieving with Eager Loading

```go
// From User side
db.Model(&User{}).Preload("Languages").Find(&users)

// From Language side
db.Model(&Language{}).Preload("Users").Find(&languages)
```

## Association Mode Operations

```go
// Append new associations (adds to join table)
db.Model(&user).Association("Languages").Append(&language)

// Replace all associations
db.Model(&user).Association("Languages").Replace(&languages)

// Delete associations (removes from join table, not the record itself)
db.Model(&user).Association("Languages").Delete(&language)

// Clear all associations
db.Model(&user).Association("Languages").Clear()

// Count associations
count := db.Model(&user).Association("Languages").Count()

// Find all associations
var languages []Language
db.Model(&user).Association("Languages").Find(&languages)
```

## Foreign Key Constraints

Configure `OnUpdate` and `OnDelete` behavior using the `constraint` tag.

```go
type User struct {
  gorm.Model
  Languages []Language `gorm:"many2many:user_languages;constraint:OnUpdate:CASCADE,OnDelete:CASCADE;"`
}
```

## Composite Foreign Keys

With composite primary keys, GORM enables composite foreign keys by default. Specify multiple keys by separating with commas.

```go
Tags []Tag `gorm:"many2many:blog_tags;ForeignKey:id,locale;References:id"`
```

## Tag Reference

| Tag | Purpose | Example |
|-----|---------|---------|
| `many2many` | Specify join table name | `gorm:"many2many:user_languages"` |
| `foreignKey` | Source reference field | `gorm:"foreignKey:Refer"` |
| `joinForeignKey` | Join table source column | `gorm:"joinForeignKey:UserReferID"` |
| `References` | Target reference field | `gorm:"References:UserRefer"` |
| `joinReferences` | Join table target column | `gorm:"joinReferences:ProfileRefer"` |
| `constraint` | FK constraints | `gorm:"constraint:OnDelete:CASCADE"` |
