
# Key Concepts for GORM Composite Primary Keys

This document provides detailed explanations of using composite primary keys in GORM.

## Overview

A composite primary key is a primary key that consists of two or more columns. The combination of these columns must be unique for each row in the table. GORM provides straightforward support for composite primary keys using struct tags.

## Defining a Composite Primary Key

To define a composite primary key, add the `gorm:"primaryKey"` tag to each field that is part of the key.

```go
// The primary key for this model is the combination of (ID, LanguageCode).
type Product struct {
    ID           string `gorm:"primaryKey"`
    LanguageCode string `gorm:"primaryKey"`
    Name         string
    Description  string
}
```

When GORM migrates this struct, it will create a table with a primary key constraint on both the `id` and `language_code` columns.

## Working with Integer Composite Keys

When using integer types for a composite primary key, you typically want to disable the `autoIncrement` behavior, as the combination of keys is what matters, not a single auto-incrementing value.

```go
type OrderItem struct {
    OrderID   uint `gorm:"primaryKey;autoIncrement:false"`
    ProductID uint `gorm:"primaryKey;autoIncrement:false"`
    Quantity  int
}
```

## CRUD Operations with Composite Keys

When performing operations like `First`, `Update`, or `Delete`, you can pass a struct literal to specify the composite key.

### Creating

Creating a record is straightforward. Just populate all the primary key fields.

```go
product := Product{ID: "p-123", LanguageCode: "en", Name: "GORM Book"}
db.Create(&product)
```

### Querying

GORM automatically uses the primary key fields when you pass a struct to a query method.

```go
var foundProduct Product
// GORM will build the WHERE clause based on the primary key fields provided:
// WHERE id = 'p-123' AND language_code = 'en'
db.First(&foundProduct, Product{ID: "p-123", LanguageCode: "en"})
```

You can also provide the values in order as arguments, but using a struct is clearer and safer.

### Updating and Deleting

Similar to querying, you can use a struct to specify the record to update or delete.

```go
// Update
db.Model(&Product{ID: "p-123", LanguageCode: "en"}).Update("Name", "New Name")

// Delete
db.Delete(&Product{ID: "p-123", LanguageCode: "en"})
```

## Composite Keys and Associations

Composite primary keys are fully supported in associations, such as many-to-many relationships.

When you define a many-to-many relationship with a model that has a composite key, GORM will automatically create the correct foreign keys in the join table.

```go
type Tag struct {
    ID     uint   `gorm:"primaryKey"`
    Locale string `gorm:"primaryKey"`
    Value  string
}

type Blog struct {
    gorm.Model
    Title string
    Tags  []Tag `gorm:"many2many:blog_tags;"`
}
```

In this example, the `blog_tags` join table will have four columns:
- `blog_id`
- `tag_id`
- `tag_locale`
- A primary key on (`blog_id`, `tag_id`, `tag_locale`)

This allows a blog to be associated with a tag that is unique for a specific locale (e.g., the tag with `ID: 1, Locale: "en"` is different from the tag with `ID: 1, Locale: "fr"`).

By using the `gorm:"primaryKey"` tag, you can effectively model complex data relationships that require multi-column uniqueness constraints.
