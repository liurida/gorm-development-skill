# Key Concepts for Polymorphism

GORM supports polymorphic associations for `has one` and `has many` relationships. Polymorphism allows a model to belong to more than one other model on a single association, storing the owner's table name in a type column and the primary key in an ID column.

## Basic Polymorphic Has Many

Use the `polymorphic` tag to specify the prefix for the type and ID columns.

```go
type Dog struct {
  gorm.Model
  Name string
  Toys []Toy `gorm:"polymorphic:Owner;"`
}

type Cat struct {
  gorm.Model
  Name string
  Toys []Toy `gorm:"polymorphic:Owner;"`
}

type Toy struct {
  gorm.Model
  Name      string
  OwnerID   uint    // Stores the ID of Dog or Cat
  OwnerType string  // Stores "dogs" or "cats" (pluralized table name)
}
```

When creating a Dog with Toys:
- `OwnerID` is set to the Dog's ID
- `OwnerType` is set to "dogs" (pluralized table name)

## Polymorphic Has One

Works the same way for single associations.

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

## Custom Column Names and Values

Three tags customize polymorphic behavior:

| Tag | Purpose | Default |
|-----|---------|---------|
| `polymorphicType` | Column name for type | `<Prefix>Type` |
| `polymorphicId` | Column name for ID | `<Prefix>ID` |
| `polymorphicValue` | Value stored in type column | Pluralized table name |

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
  Kind    string  // Custom type column (stores "master" instead of "dogs")
}
```

## Retrieving with Eager Loading

```go
// Preload polymorphic association
var dog Dog
db.Preload("Toys").First(&dog, dogID)

var cat Cat
db.Preload("Toys").First(&cat, catID)
```

## Querying by Owner Type

Find all toys belonging to a specific type:

```go
var dogToys []Toy
db.Where("owner_type = ?", "dogs").Find(&dogToys)

var catToys []Toy
db.Where("owner_type = ?", "cats").Find(&catToys)
```

## Association Mode Operations

```go
// Append to polymorphic association
db.Model(&dog).Association("Toys").Append(&toy)

// Replace all toys
db.Model(&dog).Association("Toys").Replace(&toys)

// Delete association
db.Model(&dog).Association("Toys").Delete(&toy)

// Clear all
db.Model(&dog).Association("Toys").Clear()

// Count
count := db.Model(&dog).Association("Toys").Count()

// Find all
var toys []Toy
db.Model(&dog).Association("Toys").Find(&toys)
```

## Common Use Cases

### Comments System
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

### Images/Attachments
```go
type User struct {
  gorm.Model
  Avatar Image `gorm:"polymorphic:Imageable;"`
}

type Product struct {
  gorm.Model
  Images []Image `gorm:"polymorphic:Imageable;"`
}

type Image struct {
  gorm.Model
  URL           string
  ImageableID   uint
  ImageableType string
}
```

## Tag Reference

| Tag | Purpose | Example |
|-----|---------|---------|
| `polymorphic` | Enable polymorphism with prefix | `gorm:"polymorphic:Owner"` |
| `polymorphicType` | Custom type column name | `gorm:"polymorphicType:Kind"` |
| `polymorphicId` | Custom ID column name | `gorm:"polymorphicId:OwnerID"` |
| `polymorphicValue` | Custom type value | `gorm:"polymorphicValue:master"` |

## Important Notes

1. **Table Names**: By default, the type column stores the pluralized table name (e.g., "dogs", "cats")
2. **Migration**: GORM automatically handles the polymorphic columns during AutoMigrate
3. **Querying**: When querying, GORM automatically filters by the correct `OwnerType`
4. **Performance**: Consider adding indexes on `OwnerID` and `OwnerType` columns for large tables
