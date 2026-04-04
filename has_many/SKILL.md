---
name: gorm-has-many
description: Use when defining one-to-many relationships where the owner can have zero or many instances of the associated model. Each associated model holds a foreign key referencing the owner.
---

# Has Many

A `has many` association sets up a one-to-many connection with another model. Unlike `has one`, the owner can have zero or many instances of the associated model.

## Basic Definition

The associated model holds the foreign key (e.g., `UserID`). The default foreign key name is the owner's type name plus its primary key field name.

```go
// User has many CreditCards, UserID is the foreign key
type User struct {
  gorm.Model
  CreditCards []CreditCard
}

type CreditCard struct {
  gorm.Model
  Number string
  UserID uint
}
```

## Retrieving with Eager Loading

```go
// Retrieve user list with eager loading credit cards
func GetAll(db *gorm.DB) ([]User, error) {
    var users []User
    err := db.Model(&User{}).Preload("CreditCards").Find(&users).Error
    return users, err
}
```

## Override Foreign Key

Use the `foreignKey` tag to customize:

```go
type User struct {
  gorm.Model
  CreditCards []CreditCard `gorm:"foreignKey:UserRefer"`
}

type CreditCard struct {
  gorm.Model
  Number    string
  UserRefer uint
}
```

## Override References

GORM typically uses the owner's primary key as the foreign key value. Change this with the `references` tag:

```go
type User struct {
  gorm.Model
  MemberNumber string
  CreditCards  []CreditCard `gorm:"foreignKey:UserNumber;references:MemberNumber"`
}

type CreditCard struct {
  gorm.Model
  Number     string
  UserNumber string
}
```

## Self-Referential Has Many

Models can reference themselves for hierarchical structures (e.g., org charts, threaded comments):

```go
type User struct {
  gorm.Model
  Name      string
  ManagerID *uint
  Team      []User `gorm:"foreignkey:ManagerID"`
}
```

## Foreign Key Constraints

Set up `OnUpdate` and `OnDelete` constraints with the `constraint` tag (applied during migration):

```go
type User struct {
  gorm.Model
  CreditCards []CreditCard `gorm:"constraint:OnUpdate:CASCADE,OnDelete:SET NULL;"`
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
| `OnDelete:CASCADE` | Deletes child records when parent is deleted |
| `OnDelete:RESTRICT` | Prevents deletion of parent if children exist |

## CRUD with Has Many

Use Association Mode for managing has many relationships:

```go
var user User
db.First(&user, 1)

// Find all associated credit cards
var cards []CreditCard
db.Model(&user).Association("CreditCards").Find(&cards)

// Append new credit cards
db.Model(&user).Association("CreditCards").Append(&CreditCard{Number: "1234"})
db.Model(&user).Association("CreditCards").Append([]CreditCard{card1, card2})

// Replace all credit cards
db.Model(&user).Association("CreditCards").Replace([]CreditCard{card1, card2})

// Delete specific credit cards (sets FK to NULL, doesn't delete record)
db.Model(&user).Association("CreditCards").Delete(&card1)

// Clear all associations (sets FK to NULL)
db.Model(&user).Association("CreditCards").Clear()

// Count associations
count := db.Model(&user).Association("CreditCards").Count()
```

## Eager Loading

Use `Preload` for eager loading has many associations:

```go
var users []User
db.Preload("CreditCards").Find(&users)
// SELECT * FROM users;
// SELECT * FROM credit_cards WHERE user_id IN (1,2,3,4);
```

**Important**: Do NOT use `Joins` for has many - it creates a Cartesian product:

```go
// BAD: Creates Cartesian product - each user duplicated per credit card
db.Joins("CreditCards").Find(&users)

// GOOD: Use Preload
db.Preload("CreditCards").Find(&users)
```

## Delete with Select

Delete has many associations when deleting the parent:

```go
// Delete user's credit cards when deleting the user
db.Select("CreditCards").Delete(&user)

// Delete multiple association types
db.Select("Orders", "CreditCards").Delete(&user)

// Delete all associations
db.Select(clause.Associations).Delete(&user)
```

## Preload with Conditions

Filter preloaded has many records:

```go
// Only preload active credit cards
db.Preload("CreditCards", "active = ?", true).Find(&users)

// Custom preload with ordering
db.Preload("CreditCards", func(db *gorm.DB) *gorm.DB {
    return db.Order("credit_cards.created_at DESC").Limit(5)
}).Find(&users)
```

## Common Patterns

### Has Many with Polymorphism

```go
type Post struct {
  gorm.Model
  Title    string
  Comments []Comment `gorm:"polymorphic:Commentable;"`
}

type Video struct {
  gorm.Model
  Title    string
  Comments []Comment `gorm:"polymorphic:Commentable;"`
}

type Comment struct {
  gorm.Model
  Content         string
  CommentableID   uint
  CommentableType string
}
```

### Batch Operations

```go
var users = []User{user1, user2, user3}

// Append different cards to different users
db.Model(&users).Association("CreditCards").Append(&card1, &card2, &[]CreditCard{card3, card4})

// The arguments must match the number of users
// user1 gets card1, user2 gets card2, user3 gets card3+card4
```

## Tag Reference

| Tag | Purpose | Example |
|-----|---------|---------|
| `foreignKey` | Specify custom foreign key field | `gorm:"foreignKey:UserRefer"` |
| `references` | Specify which field to reference | `gorm:"references:MemberNumber"` |
| `constraint` | Database FK constraints | `gorm:"constraint:OnDelete:CASCADE"` |
| `polymorphic` | Enable polymorphism | `gorm:"polymorphic:Owner"` |

## When NOT to Use

- **When loading data for table views** - `Joins` is often better for creating flat table structures for display, as `Preload` creates separate queries.
- **For very large numbers of associated records** - Preloading 100,000 child records will consume a lot of memory. Use pagination on the association or query the children separately.
- **If you need to filter the parent based on child data** - `Preload` doesn't filter the parent. Use `Joins` with a `WHERE` clause on the child table instead.
- **Many-to-many relationships** - A `has many` relationship doesn't work for users and roles where a join table is needed. Use `many2many`.

## References

- [GORM Official Documentation: Has Many](https://gorm.io/docs/has_many.html)
