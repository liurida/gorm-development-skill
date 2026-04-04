---
name: gorm-indexes
description: Use when creating single or composite database indexes to improve query performance by allowing the database to find rows faster, including unique indexes, partial indexes, expression indexes, or database-specific index options.
---

# Indexes

GORM allows creating database indexes with tags `index` and `uniqueIndex`. Indexes are created during [AutoMigrate or CreateTable](migration.html).

## Index Tag

GORM accepts many index settings: `class`, `type`, `where`, `comment`, `expression`, `sort`, `collate`, `option`.

```go
type User struct {
    Name  string `gorm:"index"`
    Name2 string `gorm:"index:idx_name,unique"`
    Name3 string `gorm:"index:,sort:desc,collate:utf8,type:btree,length:10,where:name3 != 'jinzhu'"`
    Name4 string `gorm:"uniqueIndex"`
    Age   int64  `gorm:"index:,class:FULLTEXT,comment:hello \\, world,where:age > 10"`
    Age2  int64  `gorm:"index:,expression:ABS(age)"`
}

// MySQL option
type User struct {
    Name string `gorm:"index:,class:FULLTEXT,option:WITH PARSER ngram INVISIBLE"`
}

// PostgreSQL option
type User struct {
    Name string `gorm:"index:,option:CONCURRENTLY"`
}
```

### uniqueIndex

The `uniqueIndex` tag works like `index` and equals `index:,unique`:

```go
type User struct {
    Name1 string `gorm:"uniqueIndex"`
    Name2 string `gorm:"uniqueIndex:idx_name,sort:desc"`
}
```

**Note:** This does not work for unique composite indexes (see below).

## Composite Indexes

Use the same index name for multiple fields to create composite indexes:

```go
// Create composite index `idx_member` with columns `name`, `number`
type User struct {
    Name   string `gorm:"index:idx_member"`
    Number string `gorm:"index:idx_member"`
}
```

For a unique composite index:

```go
// Create unique composite index `idx_member` with columns `name`, `number`
type User struct {
    Name   string `gorm:"index:idx_member,unique"`
    Number string `gorm:"index:idx_member,unique"`
}
```

### Fields Priority

Column order impacts performance. Use `priority` to control order (default is `10`). Lower priority values come first:

```go
type User struct {
    Name   string `gorm:"index:idx_member"`
    Number string `gorm:"index:idx_member"`
}
// column order: name, number

type User struct {
    Name   string `gorm:"index:idx_member,priority:2"`
    Number string `gorm:"index:idx_member,priority:1"`
}
// column order: number, name

type User struct {
    Name   string `gorm:"index:idx_member,priority:12"`
    Number string `gorm:"index:idx_member"`
}
// column order: number, name
```

### Shared Composite Indexes

When embedding a struct multiple times, explicit index names cause duplicates. Use `composite` to let GORM generate names via NamingStrategy:

```go
type Foo struct {
    IndexA int `gorm:"index:,unique,composite:myname"`
    IndexB int `gorm:"index:,unique,composite:myname"`
}
// Table Foo gets index: idx_foo_myname

type Bar0 struct {
    Foo
}
// Table Bar0 gets index: idx_bar0_myname

type Bar1 struct {
    Foo
}
// Table Bar1 gets index: idx_bar1_myname
```

**Note:** `composite` only works when no explicit index name is specified.

## Multiple Indexes

A field can have multiple `index` or `uniqueIndex` tags:

```go
type UserIndex struct {
    OID          int64  `gorm:"index:idx_id;index:idx_oid,unique"`
    MemberNumber string `gorm:"index:idx_id"`
}
```

## Quick Reference

| Setting | Example | Description |
|---------|---------|-------------|
| `index` | `gorm:"index"` | Basic index |
| `index:name` | `gorm:"index:idx_name"` | Named index |
| `uniqueIndex` | `gorm:"uniqueIndex"` | Unique index |
| `unique` | `gorm:"index:idx_name,unique"` | Unique constraint |
| `type` | `gorm:"index:,type:btree"` | Index type (btree, hash, etc.) |
| `class` | `gorm:"index:,class:FULLTEXT"` | Index class |
| `sort` | `gorm:"index:,sort:desc"` | Sort order |
| `collate` | `gorm:"index:,collate:utf8"` | Collation |
| `length` | `gorm:"index:,length:10"` | Prefix length |
| `where` | `gorm:"index:,where:active = true"` | Partial index condition |
| `comment` | `gorm:"index:,comment:my comment"` | Index comment |
| `expression` | `gorm:"index:,expression:ABS(age)"` | Expression index |
| `option` | `gorm:"index:,option:CONCURRENTLY"` | Database-specific option |
| `priority` | `gorm:"index:idx,priority:1"` | Column order in composite |
| `composite` | `gorm:"index:,composite:myname"` | Shared composite index |

## When NOT to Use

- **On every column** - Indexes slow down write operations (INSERT, UPDATE, DELETE). Only index columns that are frequently used in `WHERE`, `JOIN`, or `ORDER BY` clauses.
- **On columns with low cardinality** - Indexing a boolean `is_active` column is often not useful, as the index will not be selective enough for the optimizer to use it.
- **On very small tables** - For tables with only a few hundred rows, a full table scan is often faster than an index lookup.
- **For columns that are rarely queried** - Don't add an index for a query that is run once a day if it slows down a high-frequency write operation.
- **When the index is too large** - Indexing very large `TEXT` or `BLOB` columns can create huge indexes and is generally not effective. Use full-text search engines instead.

## Reference

- Official Docs: https://gorm.io/docs/indexes.html
