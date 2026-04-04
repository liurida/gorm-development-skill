---
name: gorm-belongs-to
description: Use when defining one-to-one relationships where the declaring model holds the foreign key. The child model "belongs to" one instance of the parent model.
---

# Belongs To

A `belongs to` association sets up a one-to-one connection with another model, such that each instance of the declaring model "belongs to" one instance of the other model.

## Basic Definition

The declaring model holds the foreign key (e.g., `CompanyID`). By default, the foreign key name uses the owner's type name plus its primary key field name.

```go
// `User` belongs to `Company`, `CompanyID` is the foreign key
type User struct {
  gorm.Model
  Name      string
  CompanyID int
  Company   Company
}

type Company struct {
  ID   int
  Name string
}
```

## Override Foreign Key

Customize the foreign key using the `foreignKey` tag:

```go
type User struct {
  gorm.Model
  Name         string
  CompanyRefer int
  Company      Company `gorm:"foreignKey:CompanyRefer"`
  // use CompanyRefer as foreign key
}

type Company struct {
  ID   int
  Name string
}
```

## Override References

By default, GORM uses the owner's primary field as the foreign key's value. Change it with the `references` tag:

```go
type User struct {
  gorm.Model
  Name      string
  CompanyID string
  Company   Company `gorm:"references:Code"` // use Code as references
}

type Company struct {
  ID   int
  Code string
  Name string
}
```

**Important**: If the override foreign key name already exists in owner's type, GORM may guess the relationship as `has one`. Use `references` to clarify:

```go
type User struct {
  gorm.Model
  Name      string
  CompanyID int
  Company   Company `gorm:"references:CompanyID"` // use Company.CompanyID as references
}

type Company struct {
  CompanyID   int
  Code        string
  Name        string
}
```

## Foreign Key Constraints

Set up `OnUpdate` and `OnDelete` constraints with the `constraint` tag (applied during migration):

```go
type User struct {
  gorm.Model
  Name      string
  CompanyID int
  Company   Company `gorm:"constraint:OnUpdate:CASCADE,OnDelete:SET NULL;"`
}

type Company struct {
  ID   int
  Name string
}
```

| Constraint Option | Behavior |
|-------------------|----------|
| `OnUpdate:CASCADE` | Updates foreign key when parent key changes |
| `OnDelete:SET NULL` | Sets foreign key to NULL when parent is deleted |
| `OnDelete:CASCADE` | Deletes child records when parent is deleted |
| `OnDelete:RESTRICT` | Prevents deletion of parent if children exist |

## Eager Loading

Use `Preload` or `Joins` to load the associated company:

```go
// Using Preload (separate query)
var users []User
db.Preload("Company").Find(&users)
// SELECT * FROM users;
// SELECT * FROM companies WHERE id IN (1,2,3,4);

// Using Joins (single query with LEFT JOIN)
db.Joins("Company").Find(&users)
// SELECT users.*, Company.* FROM users LEFT JOIN companies AS Company ON ...
```

## CRUD with Belongs To

Use Association Mode for managing belongs to relationships:

```go
var user User
db.First(&user, 1)

// Find the associated company
var company Company
db.Model(&user).Association("Company").Find(&company)

// Replace the company (belongs to replaces, doesn't append)
db.Model(&user).Association("Company").Append(&newCompany)

// Clear the association (sets CompanyID to null)
db.Model(&user).Association("Company").Clear()
```

## Belongs To vs Has One

| Aspect | Belongs To | Has One |
|--------|-----------|---------|
| Foreign Key Location | On the declaring model | On the associated model |
| Example | User has CompanyID | User's CreditCard has UserID |
| Ownership | User belongs to Company | User owns CreditCard |

## Common Patterns

### Optional Belongs To (Nullable Foreign Key)

```go
type User struct {
  gorm.Model
  Name      string
  CompanyID *int     // Pointer makes it nullable
  Company   *Company // Pointer allows nil
}
```

### Belongs To with Self-Reference

```go
type Employee struct {
  gorm.Model
  Name      string
  ManagerID *uint
  Manager   *Employee
}
```

## Tag Reference

| Tag | Purpose | Example |
|-----|---------|---------|
| `foreignKey` | Specify custom foreign key field | `gorm:"foreignKey:CompanyRefer"` |
| `references` | Specify which field to reference | `gorm:"references:Code"` |
| `constraint` | Database FK constraints | `gorm:"constraint:OnDelete:CASCADE"` |

## When NOT to Use

- **When the other model should hold the foreign key** - Use `has one` if the `User` owns the `CreditCard` (credit card has a `user_id`)
- **One-to-many relationships** - A user can't belong to multiple companies; use `has many` for a user that has many posts
- **Many-to-many relationships** - Use `many2many` for users and languages where a join table is needed
- **When the relationship is optional and you want a separate table** - For optional one-to-one, a nullable foreign key is standard; a separate table is over-engineering

## References

- [GORM Official Documentation: Belongs To](https://gorm.io/docs/belongs_to.html)
