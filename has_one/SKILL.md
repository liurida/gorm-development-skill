---
name: gorm-has-one
description: Use when defining one-to-one relationships where the associated model holds the foreign key. The owner "has one" instance of the associated model.
---

# Has One

A `has one` association sets up a one-to-one connection with another model. This association indicates that each instance of a model contains or possesses one instance of another model.

## Basic Definition

The associated model holds the foreign key (e.g., `UserID`). The foreign key name is typically generated from the owner's type name plus its primary key.

```go
// User has one CreditCard, UserID is the foreign key
type User struct {
  gorm.Model
  CreditCard CreditCard
}

type CreditCard struct {
  gorm.Model
  Number string
  UserID uint
}
```

## Retrieving with Eager Loading

```go
// Retrieve user list with eager loading credit card
func GetAll(db *gorm.DB) ([]User, error) {
    var users []User
    err := db.Model(&User{}).Preload("CreditCard").Find(&users).Error
    return users, err
}
```

## Override Foreign Key

Use the `foreignKey` tag to specify a different field:

```go
type User struct {
  gorm.Model
  CreditCard CreditCard `gorm:"foreignKey:UserName"`
  // use UserName as foreign key
}

type CreditCard struct {
  gorm.Model
  Number   string
  UserName string
}
```

## Override References

By default, the owned entity saves the owner's primary key. Change this with the `references` tag:

```go
type User struct {
  gorm.Model
  Name       string     `gorm:"index"`
  CreditCard CreditCard `gorm:"foreignKey:UserName;references:Name"`
}

type CreditCard struct {
  gorm.Model
  Number   string
  UserName string
}
```

## Self-Referential Has One

Models can reference themselves:

```go
type User struct {
  gorm.Model
  Name      string
  ManagerID *uint
  Manager   *User
}
```

## Foreign Key Constraints

Set up `OnUpdate` and `OnDelete` constraints with the `constraint` tag:

```go
type User struct {
  gorm.Model
  CreditCard CreditCard `gorm:"constraint:OnUpdate:CASCADE,OnDelete:SET NULL;"`
}

type CreditCard struct {
  gorm.Model
  Number string
  UserID uint
}
```

| Constraint Option | Behavior |
|-------------------|----------|
| `OnUpdate:CASCADE` | Updates foreign key when parent key changes |
| `OnDelete:SET NULL` | Sets foreign key to NULL when parent is deleted |
| `OnDelete:CASCADE` | Deletes child record when parent is deleted |
| `OnDelete:RESTRICT` | Prevents deletion of parent if child exists |

## CRUD with Has One

Use Association Mode for managing has one relationships:

```go
var user User
db.First(&user, 1)

// Find the associated credit card
var card CreditCard
db.Model(&user).Association("CreditCard").Find(&card)

// Replace the credit card (has one replaces existing)
db.Model(&user).Association("CreditCard").Append(&newCard)

// Clear the association (sets UserID to null on CreditCard)
db.Model(&user).Association("CreditCard").Clear()
```

## Eager Loading

Use `Preload` or `Joins` for eager loading:

```go
// Using Preload (separate query)
var users []User
db.Preload("CreditCard").Find(&users)
// SELECT * FROM users;
// SELECT * FROM credit_cards WHERE user_id IN (1,2,3,4);

// Using Joins (single query with LEFT JOIN) - more efficient for has one
db.Joins("CreditCard").Find(&users)
// SELECT users.*, CreditCard.* FROM users LEFT JOIN credit_cards AS CreditCard ON ...
```

## Delete with Select

Delete the has one association when deleting the parent:

```go
// Delete user's credit card when deleting the user
db.Select("CreditCard").Delete(&user)
```

## Has One vs Belongs To

| Aspect | Has One | Belongs To |
|--------|---------|-----------|
| Foreign Key Location | On the associated model | On the declaring model |
| Example | CreditCard has UserID | User has CompanyID |
| Ownership | User owns CreditCard | User belongs to Company |

## Common Patterns

### Optional Has One (Nullable)

```go
type User struct {
  gorm.Model
  Profile *Profile // Pointer allows nil
}

type Profile struct {
  gorm.Model
  Bio    string
  UserID uint
}
```

### Has One with Polymorphism

```go
type Company struct {
  gorm.Model
  Name    string
  Address Address `gorm:"polymorphic:Addressable;"`
}

type Person struct {
  gorm.Model
  Name    string
  Address Address `gorm:"polymorphic:Addressable;"`
}

type Address struct {
  gorm.Model
  Street          string
  AddressableID   uint
  AddressableType string
}
```

## Tag Reference

| Tag | Purpose | Example |
|-----|---------|---------|
| `foreignKey` | Specify custom foreign key field | `gorm:"foreignKey:UserName"` |
| `references` | Specify which field to reference | `gorm:"references:Name"` |
| `constraint` | Database FK constraints | `gorm:"constraint:OnDelete:CASCADE"` |
| `polymorphic` | Enable polymorphism | `gorm:"polymorphic:Owner"` |

## When NOT to Use

- **When the declaring model should hold the foreign key** - Use `belongs to` if the `User` should have a `CompanyID`.
- **One-to-many relationships** - A `has one` relationship restricts the owner to a single associated record. Use `has many` if a user can have multiple credit cards.
- **When the data is part of the same logical entity** - If a `Profile` is always accessed with a `User` and has no independent existence, consider embedding it directly in the `User` struct instead of using a `has one` relationship.

## References

- [GORM Official Documentation: Has One](https://gorm.io/docs/has_one.html)
