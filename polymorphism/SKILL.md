---
name: gorm-polymorphism
description: Use when defining a relationship where a single model can belong to multiple different owner models on a single association, such as comments belonging to both posts and videos, or addresses for both users and companies.
---

# Polymorphism

GORM supports polymorphic associations for `has one` and `has many` relationships. This allows a model to belong to more than one other model on a single association, using a type column to distinguish the owner.

## Core Concept

Polymorphic associations store:
- **Owner ID**: The primary key of the owning record
- **Owner Type**: The table name of the owning model (pluralized by default)

## Basic Has One Polymorphism

```go
type Cat struct {
  gorm.Model
  Name string
  Toy  Toy `gorm:"polymorphic:Owner;"`
}

type Dog struct {
  gorm.Model
  Name string
  Toy  Toy `gorm:"polymorphic:Owner;"`
}

type Toy struct {
  gorm.Model
  Name      string
  OwnerID   uint
  OwnerType string
}

db.Create(&Dog{Name: "dog1", Toy: Toy{Name: "toy1"}})
// INSERT INTO `dogs` (`name`) VALUES ("dog1")
// INSERT INTO `toys` (`name`,`owner_id`,`owner_type`) VALUES ("toy1",1,"dogs")

db.Create(&Cat{Name: "cat1", Toy: Toy{Name: "toy2"}})
// INSERT INTO `cats` (`name`) VALUES ("cat1")
// INSERT INTO `toys` (`name`,`owner_id`,`owner_type`) VALUES ("toy2",1,"cats")
```

## Basic Has Many Polymorphism

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

db.Create(&Post{
  Title: "My Post",
  Comments: []Comment{
    {Content: "Great post!"},
    {Content: "Thanks for sharing"},
  },
})
// INSERT INTO `posts` (`title`) VALUES ("My Post")
// INSERT INTO `comments` (`content`,`commentable_id`,`commentable_type`) 
//   VALUES ("Great post!",1,"posts"), ("Thanks for sharing",1,"posts")
```

## Customizing Polymorphic Columns

Use separate tags to customize column names and values:

| Tag | Purpose | Default |
|-----|---------|---------|
| `polymorphicType` | Column for storing the type | `<Prefix>Type` |
| `polymorphicId` | Column for storing the ID | `<Prefix>ID` |
| `polymorphicValue` | Value stored in type column | Table name (pluralized) |

```go
type Dog struct {
  gorm.Model
  Name string
  Toys []Toy `gorm:"polymorphicType:Kind;polymorphicId:OwnerID;polymorphicValue:master"`
}

type Toy struct {
  gorm.Model
  Name    string
  OwnerID uint
  Kind    string  // Custom type column name
}

db.Create(&Dog{Name: "dog1", Toys: []Toy{{Name: "toy1"}, {Name: "toy2"}}})
// INSERT INTO `dogs` (`name`) VALUES ("dog1")
// INSERT INTO `toys` (`name`,`owner_id`,`kind`) VALUES ("toy1",1,"master"), ("toy2",1,"master")
```

## Eager Loading Polymorphic Associations

Use `Preload` to load polymorphic associations:

```go
// Load all dogs with their toys
var dogs []Dog
db.Preload("Toy").Find(&dogs)
// SELECT * FROM dogs;
// SELECT * FROM toys WHERE owner_type = 'dogs' AND owner_id IN (1,2,3);

// Load all posts with their comments
var posts []Post
db.Preload("Comments").Find(&posts)
// SELECT * FROM posts;
// SELECT * FROM comments WHERE commentable_type = 'posts' AND commentable_id IN (1,2,3);
```

### Preload with Conditions

```go
// Only load approved comments
db.Preload("Comments", "approved = ?", true).Find(&posts)

// Order comments by creation time
db.Preload("Comments", func(db *gorm.DB) *gorm.DB {
    return db.Order("comments.created_at DESC")
}).Find(&posts)
```

## CRUD with Polymorphic Associations

Use Association Mode for managing polymorphic relationships:

```go
var dog Dog
db.First(&dog, 1)

// Find the associated toy (has one)
var toy Toy
db.Model(&dog).Association("Toy").Find(&toy)

// Replace the toy
db.Model(&dog).Association("Toy").Replace(&Toy{Name: "new toy"})

// Clear the association (sets OwnerID/OwnerType to null)
db.Model(&dog).Association("Toy").Clear()

// For has many (e.g., Post with Comments)
var post Post
db.First(&post, 1)

// Find all comments
var comments []Comment
db.Model(&post).Association("Comments").Find(&comments)

// Append new comments
db.Model(&post).Association("Comments").Append(&Comment{Content: "New comment"})

// Delete specific comments
db.Model(&post).Association("Comments").Delete(&comment1)

// Count comments
count := db.Model(&post).Association("Comments").Count()
```

## Delete with Select

Delete polymorphic associations when deleting the parent:

```go
// Delete dog's toy when deleting the dog
db.Select("Toy").Delete(&dog)

// Delete post's comments when deleting the post
db.Select("Comments").Delete(&post)

// Delete all associations
db.Select(clause.Associations).Delete(&post)
```

## Querying Across Polymorphic Types

Query the associated table directly when you need cross-type queries:

```go
// Find all toys regardless of owner type
var toys []Toy
db.Find(&toys)

// Find all toys belonging to dogs
db.Where("owner_type = ?", "dogs").Find(&toys)

// Find all comments on posts
var comments []Comment
db.Where("commentable_type = ?", "posts").Find(&comments)
```

## Foreign Key Constraints

Polymorphic associations typically **do not use database foreign key constraints** because the foreign key references multiple tables. The `OwnerType` column determines which table to look up.

If you need referential integrity, consider:
1. Application-level validation
2. Database triggers
3. Separate foreign key columns per type (non-polymorphic approach)

## Common Patterns

### Address for Multiple Entity Types

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
  City            string
  AddressableID   uint
  AddressableType string
}
```

### Taggable Pattern

```go
type Article struct {
  gorm.Model
  Title string
  Tags  []Tag `gorm:"polymorphic:Taggable;"`
}

type Photo struct {
  gorm.Model
  URL  string
  Tags []Tag `gorm:"polymorphic:Taggable;"`
}

type Tag struct {
  gorm.Model
  Name         string
  TaggableID   uint
  TaggableType string
}
```

### Activity Feed / Audit Log

```go
type User struct {
  gorm.Model
  Activities []Activity `gorm:"polymorphic:Actor;"`
}

type System struct {
  gorm.Model
  Activities []Activity `gorm:"polymorphic:Actor;"`
}

type Activity struct {
  gorm.Model
  Action    string
  ActorID   uint
  ActorType string
  CreatedAt time.Time
}
```

## Tag Reference

| Tag | Purpose | Example |
|-----|---------|---------|
| `polymorphic` | Enable polymorphism with prefix | `gorm:"polymorphic:Owner"` |
| `polymorphicType` | Custom type column name | `gorm:"polymorphicType:Kind"` |
| `polymorphicId` | Custom ID column name | `gorm:"polymorphicId:OwnerID"` |
| `polymorphicValue` | Custom value for type column | `gorm:"polymorphicValue:master"` |

## Polymorphism vs Separate Tables

| Aspect | Polymorphic | Separate Tables |
|--------|-------------|-----------------|
| Schema | Single table with type column | Multiple join/reference tables |
| Foreign Keys | No DB-level FK constraints | Full FK support |
| Queries | Single table scan | Joins or separate queries |
| Use When | Simple associations, flexible schema | Strict integrity, complex relationships |

## Common Mistakes

| Mistake | Fix |
|---------|-----|
| Missing type/ID columns on associated model | Add `OwnerID uint` and `OwnerType string` fields |
| Expecting DB foreign keys | Polymorphism doesn't create FK constraints; validate in application |
| Wrong column naming | Column names default to `<Prefix>ID` and `<Prefix>Type` |
| Querying without type filter | Include `WHERE owner_type = ?` when querying associated table directly |

## When NOT to Use

- **When database-level referential integrity is required** - Polymorphic associations cannot use foreign key constraints. If you need the database to enforce these rules, use separate tables and standard `has one`/`has many` relationships.
- **If the associated models have very different life cycles or access patterns** - It might be cleaner to have separate tables (e.g., `post_comments`, `video_comments`) if they are queried and managed differently.
- **When performance for a specific owner type is critical** - A polymorphic table must be indexed on both `(OwnerID, OwnerType)`. A dedicated table with a simple index on `OwnerID` can be more performant.
- **If the number of owning types is very large and queries often target a single type** - The `OwnerType` column might have low cardinality, making indexes less effective. Separate tables might be better in this case.

## References

- [GORM Official Documentation: Polymorphism](https://gorm.io/docs/polymorphism.html)
