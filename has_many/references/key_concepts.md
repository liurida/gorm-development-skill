# Key Concepts for Has Many

A `has many` association creates a one-to-many connection between models. The owner can have zero or many instances of the associated model.

## Basic Declaration

The foreign key is automatically inferred from the owner's type name plus its primary key field name.

```go
// User has many CreditCards, UserID is the foreign key
type User struct {
  gorm.Model
  CreditCards []CreditCard
}

type CreditCard struct {
  gorm.Model
  Number string
  UserID uint  // Foreign key (convention: OwnerTypeName + ID)
}
```

## Retrieving with Eager Loading

Use `Preload` to fetch associated records and avoid N+1 queries.

```go
var users []User
db.Model(&User{}).Preload("CreditCards").Find(&users)
```

## Custom Foreign Key

Override the default foreign key name with the `foreignKey` tag.

```go
type User struct {
  gorm.Model
  CreditCards []CreditCard `gorm:"foreignKey:UserRefer"`
}

type CreditCard struct {
  gorm.Model
  Number    string
  UserRefer uint  // Custom foreign key field
}
```

## Custom References

By default, GORM uses the owner's primary key. Use `references` to specify a different field.

```go
type User struct {
  gorm.Model
  MemberNumber string
  CreditCards  []CreditCard `gorm:"foreignKey:UserNumber;references:MemberNumber"`
}

type CreditCard struct {
  gorm.Model
  Number     string
  UserNumber string  // References User.MemberNumber instead of User.ID
}
```

## Self-Referential Has Many

Models can reference themselves for hierarchical structures like org charts or threaded comments.

```go
type User struct {
  gorm.Model
  Name      string
  ManagerID *uint
  Team      []User `gorm:"foreignkey:ManagerID"`
}
```

## Foreign Key Constraints

Configure database-level constraints with the `constraint` tag.

```go
type User struct {
  gorm.Model
  CreditCards []CreditCard `gorm:"constraint:OnUpdate:CASCADE,OnDelete:SET NULL;"`
}
```

| Constraint Option | Behavior |
|-------------------|----------|
| `OnUpdate:CASCADE` | Updates foreign key when parent key changes |
| `OnDelete:SET NULL` | Sets foreign key to NULL when parent is deleted |
| `OnDelete:CASCADE` | Deletes child records when parent is deleted |

## Association Mode Operations

GORM provides helper methods to manage has many relationships.

```go
// Append new associations
db.Model(&user).Association("CreditCards").Append(&card)

// Replace all associations
db.Model(&user).Association("CreditCards").Replace(&cards)

// Delete specific associations (sets FK to NULL)
db.Model(&user).Association("CreditCards").Delete(&card)

// Clear all associations (sets FK to NULL)
db.Model(&user).Association("CreditCards").Clear()

// Count associations
count := db.Model(&user).Association("CreditCards").Count()

// Find all associations
var cards []CreditCard
db.Model(&user).Association("CreditCards").Find(&cards)
```

## Tag Reference

| Tag | Purpose | Example |
|-----|---------|---------|
| `foreignKey` | Specify custom foreign key field | `gorm:"foreignKey:UserRefer"` |
| `references` | Specify which field to reference | `gorm:"references:MemberNumber"` |
| `constraint` | Database FK constraints | `gorm:"constraint:OnDelete:CASCADE"` |
