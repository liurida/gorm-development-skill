---
name: gorm-many-to-many
description: Use when defining many-to-many relationships where records can be associated with multiple records of another model via a join table, such as users and languages, posts and tags, or students and courses.
---

# Many To Many

A `many to many` association adds a join table between two models. Each instance of either model can be associated with multiple instances of the other.

## Basic Definition

The `many2many` tag specifies the join table name. GORM creates this table automatically during `AutoMigrate`.

```go
// User has and belongs to many Languages, `user_languages` is the join table
type User struct {
  gorm.Model
  Languages []Language `gorm:"many2many:user_languages;"`
}

type Language struct {
  gorm.Model
  Name string
}
```

## Back-Reference (Bidirectional)

For bidirectional navigation, declare the association on both models with the same join table:

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

### Retrieving Bidirectional Associations

```go
// Retrieve users with their languages
func GetAllUsers(db *gorm.DB) ([]User, error) {
    var users []User
    err := db.Model(&User{}).Preload("Languages").Find(&users).Error
    return users, err
}

// Retrieve languages with their users
func GetAllLanguages(db *gorm.DB) ([]Language, error) {
    var languages []Language
    err := db.Model(&Language{}).Preload("Users").Find(&languages).Error
    return languages, err
}
```

## Override Foreign Key

The join table owns two foreign keys referencing both models. Override them with `foreignKey`, `references`, `joinForeignKey`, `joinReferences`:

```go
type User struct {
  gorm.Model
  Profiles []Profile `gorm:"many2many:user_profiles;foreignKey:Refer;joinForeignKey:UserReferID;References:UserRefer;joinReferences:ProfileRefer"`
  Refer    uint      `gorm:"index:,unique"`
}

type Profile struct {
  gorm.Model
  Name      string
  UserRefer uint `gorm:"index:,unique"`
}

// Creates join table: user_profiles
//   foreign key: user_refer_id, reference: users.refer
//   foreign key: profile_refer, reference: profiles.user_refer
```

**Important**: When creating database foreign keys during migration, the referenced field must have a unique index.

| Tag | Purpose | Example |
|-----|---------|---------|
| `foreignKey` | Field on current model to use | `gorm:"foreignKey:Refer"` |
| `references` | Field on related model to reference | `gorm:"References:UserRefer"` |
| `joinForeignKey` | Join table column for current model | `gorm:"joinForeignKey:UserReferID"` |
| `joinReferences` | Join table column for related model | `gorm:"joinReferences:ProfileRefer"` |

## Self-Referential Many2Many

Models can reference themselves for relationships like friends or followers:

```go
type User struct {
  gorm.Model
  Friends []*User `gorm:"many2many:user_friends"`
}

// Creates join table: user_friends
//   foreign key: user_id, reference: users.id
//   foreign key: friend_id, reference: users.id
```

## Customize Join Table

The join table can be a full-featured model with `Soft Delete`, `Hooks`, and additional fields:

```go
type Person struct {
  ID        int
  Name      string
  Addresses []Address `gorm:"many2many:person_addresses;"`
}

type Address struct {
  ID   uint
  Name string
}

type PersonAddress struct {
  PersonID  int `gorm:"primaryKey"`
  AddressID int `gorm:"primaryKey"`
  CreatedAt time.Time
  DeletedAt gorm.DeletedAt
}

func (PersonAddress) BeforeCreate(db *gorm.DB) error {
  // Custom hook logic
  return nil
}

// Setup the custom join table
err := db.SetupJoinTable(&Person{}, "Addresses", &PersonAddress{})
```

**Important**: Custom join table foreign keys must be composite primary keys or have a composite unique index.

## Foreign Key Constraints

Set up `OnUpdate` and `OnDelete` constraints:

```go
type User struct {
  gorm.Model
  Languages []Language `gorm:"many2many:user_speaks;"`
}

type Language struct {
  Code string `gorm:"primarykey"`
  Name string
}
// Creates constraints on user_speaks table:
// FOREIGN KEY (user_id) REFERENCES users(id) ON DELETE SET NULL ON UPDATE CASCADE
// FOREIGN KEY (language_code) REFERENCES languages(code) ON DELETE SET NULL ON UPDATE CASCADE
```

| Constraint Option | Behavior |
|-------------------|----------|
| `OnUpdate:CASCADE` | Updates foreign key when parent key changes |
| `OnDelete:SET NULL` | Sets foreign key to NULL when parent is deleted |
| `OnDelete:CASCADE` | Deletes join table records when parent is deleted |
| `OnDelete:RESTRICT` | Prevents deletion of parent if join records exist |

## CRUD with Many2Many

Use Association Mode for managing many-to-many relationships:

```go
var user User
db.First(&user, 1)

// Find all associated languages
var languages []Language
db.Model(&user).Association("Languages").Find(&languages)

// Append new languages
db.Model(&user).Association("Languages").Append(&Language{Name: "DE"})
db.Model(&user).Association("Languages").Append([]Language{lang1, lang2})

// Replace all languages
db.Model(&user).Association("Languages").Replace([]Language{lang1, lang2})

// Delete specific languages (removes join table entry, not the language record)
db.Model(&user).Association("Languages").Delete(&lang1)

// Clear all associations (removes all join table entries)
db.Model(&user).Association("Languages").Clear()

// Count associations
count := db.Model(&user).Association("Languages").Count()
```

## Eager Loading

Use `Preload` for eager loading many-to-many associations:

```go
var users []User
db.Preload("Languages").Find(&users)
// SELECT * FROM users;
// SELECT * FROM user_languages WHERE user_id IN (1,2,3,4);
// SELECT * FROM languages WHERE id IN (1,2,3);
```

### Preload with Conditions

```go
// Only preload specific languages
db.Preload("Languages", "name IN ?", []string{"EN", "ZH"}).Find(&users)

// Custom preload with ordering
db.Preload("Languages", func(db *gorm.DB) *gorm.DB {
    return db.Order("languages.name ASC")
}).Find(&users)
```

## Delete with Select

Delete many-to-many associations when deleting the parent:

```go
// Delete user's language associations (join table entries)
db.Select("Languages").Delete(&user)

// Delete all associations
db.Select(clause.Associations).Delete(&user)
```

## Composite Foreign Keys

With composite primary keys, GORM enables composite foreign keys by default:

```go
type Tag struct {
  ID     uint   `gorm:"primaryKey"`
  Locale string `gorm:"primaryKey"`
  Value  string
}

type Blog struct {
  ID         uint   `gorm:"primaryKey"`
  Locale     string `gorm:"primaryKey"`
  Subject    string
  Body       string
  Tags       []Tag `gorm:"many2many:blog_tags;"`
  LocaleTags []Tag `gorm:"many2many:locale_blog_tags;ForeignKey:id,locale;References:id"`
  SharedTags []Tag `gorm:"many2many:shared_blog_tags;ForeignKey:id;References:id"`
}

// Join Table: blog_tags (all keys)
//   foreign key: blog_id, blog_locale -> blogs.id, blogs.locale
//   foreign key: tag_id, tag_locale -> tags.id, tags.locale

// Join Table: locale_blog_tags (partial keys)
//   foreign key: blog_id, blog_locale -> blogs.id, blogs.locale
//   foreign key: tag_id -> tags.id

// Join Table: shared_blog_tags (single keys)
//   foreign key: blog_id -> blogs.id
//   foreign key: tag_id -> tags.id
```

## Skip Auto-Create During Save

```go
// Skip upserting Languages (don't create/update language records)
db.Omit("Languages.*").Create(&user)

// Skip both association and join table entries
db.Omit("Languages").Create(&user)
```

## Tag Reference

| Tag | Purpose | Example |
|-----|---------|---------|
| `many2many` | Specify join table name | `gorm:"many2many:user_languages"` |
| `foreignKey` | Current model field for FK | `gorm:"foreignKey:Refer"` |
| `references` | Related model field to reference | `gorm:"References:UserRefer"` |
| `joinForeignKey` | Join table column for current model | `gorm:"joinForeignKey:UserReferID"` |
| `joinReferences` | Join table column for related model | `gorm:"joinReferences:ProfileRefer"` |
| `constraint` | Database FK constraints | `gorm:"constraint:OnDelete:CASCADE"` |

## Common Mistakes

| Mistake | Fix |
|---------|-----|
| Different join table names on each model | Use identical `many2many:table_name` on both sides |
| Missing unique index on custom FK fields | Add `gorm:"index:,unique"` to referenced fields |
| Using `Joins` instead of `Preload` | Always use `Preload` for many-to-many to avoid Cartesian products |
| Forgetting `SetupJoinTable` for custom join tables | Call `db.SetupJoinTable()` before using the association |

## When NOT to Use

- **When the relationship is one-to-many** - If a `User` can have many `Posts`, but a `Post` belongs to only one `User`, use `has many` and `belongs to`.
- **If the join table needs many custom fields** - While GORM supports custom join tables, if the join table becomes a core entity in your domain with its own complex logic, it's often better to model it as a first-class struct with two separate `belongs to` relationships.
- **For simple one-to-one relationships** - Don't use `many2many` for a `User` and `Profile`; use `has one`.

## References

- [GORM Official Documentation: Many To Many](https://gorm.io/docs/many_to_many.html)
