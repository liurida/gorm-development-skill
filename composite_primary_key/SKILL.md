---
name: gorm-composite-primary-key
description: Use when a single column is not sufficient to uniquely identify a record and a combination of multiple columns is required as the primary key, including string-based keys, integer-based keys with autoIncrement disabled, or composite keys in associations.
---

# Composite Primary Key

GORM supports composite primary keys by setting multiple fields with the `primaryKey` tag.

## Basic Usage

```go
type Product struct {
    ID           string `gorm:"primaryKey"`
    LanguageCode string `gorm:"primaryKey"`
    Code         string
    Name         string
}
```

The combination of (ID, LanguageCode) must be unique for each record.

## Integer Primary Keys

**Important:** Integer fields marked as `PrioritizedPrimaryField` enable `AutoIncrement` by default. For composite keys with integers, you must disable auto-increment:

```go
type Product struct {
    CategoryID uint64 `gorm:"primaryKey;autoIncrement:false"`
    TypeID     uint64 `gorm:"primaryKey;autoIncrement:false"`
}
```

Without `autoIncrement:false`, the database will try to auto-generate values for integer primary key fields.

## Querying with Composite Keys

When querying records with composite keys, provide all key fields:

```go
var product Product

// Using struct with key fields
db.First(&product, Product{ID: "product-001", LanguageCode: "en-US"})

// Using Where clause
db.Where("id = ? AND language_code = ?", "product-001", "en-US").First(&product)
```

## Associations with Composite Keys

Composite keys work with associations. GORM automatically uses all primary key fields in join tables:

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

The join table `blog_tags` will include both `tag_id` and `tag_locale` columns.

## Quick Reference

| Scenario | Tag |
|----------|-----|
| String primary key | `gorm:"primaryKey"` |
| Integer primary key (no auto-increment) | `gorm:"primaryKey;autoIncrement:false"` |
| Multiple fields | Apply `primaryKey` to each field |

## When NOT to Use

- **When a single unique ID is sufficient** - Prefer a single `ID` or `UUID` field for simplicity if it can uniquely identify records
- **When you need simple URL routing** - `/products/1` is easier to handle than `/products/product-001/en-US`
- **If joins become too complex** - Joining on multiple keys can be less performant and harder to reason about
- **With ORMs that have poor support** - GORM handles them well, but other tools might not, limiting portability
- **When one part of the key changes frequently** - Primary keys should be stable; frequent updates to key columns are an anti-pattern

## Reference

- Official Docs: https://gorm.io/docs/composite_primary_key.html
