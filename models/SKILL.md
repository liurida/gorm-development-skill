---
name: gorm-models
description: Use when declaring GORM models including struct definition, gorm.Model embedding, field tags, field permissions, timestamp tracking, and embedded structs.
---

# Declaring Models

Models are Go structs that map to database tables. They can contain basic Go types, pointers, aliases, or custom types implementing `Scanner` and `Valuer` interfaces.

## Basic Model

```go
type User struct {
  ID           uint           // Standard field for the primary key
  Name         string         // A regular string field
  Email        *string        // Pointer allows NULL values
  Age          uint8          // Unsigned 8-bit integer
  Birthday     *time.Time     // Pointer to time.Time, can be null
  MemberNumber sql.NullString // sql.NullString for nullable strings
  ActivatedAt  sql.NullTime   // sql.NullTime for nullable time fields
  CreatedAt    time.Time      // Auto-managed by GORM for creation time
  UpdatedAt    time.Time      // Auto-managed by GORM for update time
  ignored      string         // Non-exported fields are ignored
}
```

## gorm.Model

GORM provides a predefined struct with common fields:

```go
type Model struct {
  ID        uint           `gorm:"primaryKey"`
  CreatedAt time.Time
  UpdatedAt time.Time
  DeletedAt gorm.DeletedAt `gorm:"index"` // For soft deletes
}

// Embed in your struct
type User struct {
  gorm.Model
  Name string
}
// Equivalent to:
// type User struct {
//   ID        uint           `gorm:"primaryKey"`
//   CreatedAt time.Time
//   UpdatedAt time.Time
//   DeletedAt gorm.DeletedAt `gorm:"index"`
//   Name      string
// }
```

## Field Tags Reference

| Tag | Description | Example |
|-----|-------------|---------|
| `column` | Custom column name | `gorm:"column:user_name"` |
| `type` | Column data type | `gorm:"type:varchar(100)"` |
| `size` | Column size/length | `gorm:"size:256"` |
| `primaryKey` | Mark as primary key | `gorm:"primaryKey"` |
| `unique` | Unique constraint | `gorm:"unique"` |
| `default` | Default value | `gorm:"default:0"` |
| `not null` | NOT NULL constraint | `gorm:"not null"` |
| `autoIncrement` | Auto increment | `gorm:"autoIncrement"` |
| `index` | Create index | `gorm:"index"` |
| `uniqueIndex` | Create unique index | `gorm:"uniqueIndex"` |
| `embedded` | Embed struct fields | `gorm:"embedded"` |
| `embeddedPrefix` | Prefix for embedded fields | `gorm:"embeddedPrefix:author_"` |
| `serializer` | Serialize/deserialize | `gorm:"serializer:json"` |
| `comment` | Column comment | `gorm:"comment:user email"` |
| `check` | Check constraint | `gorm:"check:age > 13"` |

**Note:** Tags are case-insensitive; `camelCase` preferred. Separate multiple tags with semicolons (`;`).

## Field-Level Permissions

Control read/write permissions with tags:

```go
type User struct {
  Name string `gorm:"<-:create"`          // read and create only
  Name string `gorm:"<-:update"`          // read and update only
  Name string `gorm:"<-"`                 // read and write (create/update)
  Name string `gorm:"<-:false"`           // read only, no write
  Name string `gorm:"->"`                 // read only
  Name string `gorm:"->;<-:create"`       // read and create
  Name string `gorm:"->:false;<-:create"` // create only, no read
  Name string `gorm:"-"`                  // ignore for read/write
  Name string `gorm:"-:all"`              // ignore for read/write/migrate
  Name string `gorm:"-:migration"`        // ignore for migration only
}
```

## Timestamp Tracking

### Auto-managed Timestamps

```go
type User struct {
  CreatedAt time.Time // Set on create if zero
  UpdatedAt time.Time // Set on create/update
}
```

### Unix Timestamps

```go
type User struct {
  CreatedAt time.Time // Default: time.Time
  UpdatedAt int       // Unix seconds
  Updated   int64 `gorm:"autoUpdateTime:nano"`  // Unix nanoseconds
  Updated   int64 `gorm:"autoUpdateTime:milli"` // Unix milliseconds
  Created   int64 `gorm:"autoCreateTime"`       // Unix seconds
}
```

### Custom Field Names

```go
type User struct {
  Created int64 `gorm:"autoCreateTime"` // Custom name with auto-tracking
  Updated int64 `gorm:"autoUpdateTime"` // Custom name with auto-tracking
}
```

### Disable Auto-tracking

```go
type User struct {
  CreatedAt time.Time `gorm:"autoCreateTime:false"`
  UpdatedAt time.Time `gorm:"autoUpdateTime:false"`
}
```

## Embedded Structs

### Anonymous Embedding

```go
type Author struct {
  Name  string
  Email string
}

type Blog struct {
  Author        // Anonymous: fields included directly
  ID      int
  Upvotes int32
}
// Creates columns: name, email, id, upvotes
```

### Tagged Embedding

```go
type Blog struct {
  ID      int
  Author  Author `gorm:"embedded"`
  Upvotes int32
}
// Same result: name, email, id, upvotes
```

### Embedded with Prefix

```go
type Blog struct {
  ID      int
  Author  Author `gorm:"embedded;embeddedPrefix:author_"`
  Upvotes int32
}
// Creates columns: id, author_name, author_email, upvotes
```

## Data Type Mapping

| Go Type | Database Type |
|---------|---------------|
| `bool` | BOOLEAN |
| `int`, `int8`, `int16`, `int32`, `int64` | INTEGER variants |
| `uint`, `uint8`, `uint16`, `uint32`, `uint64` | UNSIGNED INTEGER variants |
| `float32`, `float64` | FLOAT variants |
| `string` | VARCHAR/TEXT |
| `time.Time` | DATETIME/TIMESTAMP |
| `[]byte` | BLOB/BYTEA |
| `*T` (pointer) | Nullable column |
| `sql.NullString`, `sql.NullTime`, etc. | Nullable with explicit NULL handling |

## Comprehensive Example

```go
type User struct {
  gorm.Model                                    // ID, CreatedAt, UpdatedAt, DeletedAt
  
  // Basic fields
  Username  string  `gorm:"size:50;uniqueIndex;not null"`
  Email     *string `gorm:"size:255;unique"`
  Age       uint8   `gorm:"check:age >= 0"`
  
  // Permission-controlled
  Password  string `gorm:"<-:create;size:255"`  // Create only
  Role      string `gorm:"default:user"`
  
  // Embedded
  Profile   Profile `gorm:"embedded;embeddedPrefix:profile_"`
  
  // Custom serialization
  Settings  map[string]interface{} `gorm:"serializer:json"`
  
  // Computed/virtual
  Internal  string `gorm:"-"` // Ignored by GORM
}

type Profile struct {
  Bio       string `gorm:"type:text"`
  AvatarURL string `gorm:"size:512"`
}
```

## Common Mistakes

| Mistake | Fix |
|---------|-----|
| Non-exported fields expected to persist | Use exported fields (uppercase first letter) |
| Using `int` for nullable fields | Use `*int` or `sql.NullInt64` |
| Missing `primaryKey` tag on custom PK | Add `gorm:"primaryKey"` to non-ID primary keys |
| Expecting `TableName()` to be dynamic | Use `db.Table()` or Scopes for dynamic names |

## When NOT to Use

- **Existing database schemas with incompatible naming** - Use `gorm:"column:..."` tags or custom `TableName()` instead of fighting conventions
- **Non-struct data sources** - For dynamic queries or arbitrary JSON, use `db.Table()` with maps instead
- **When you need maximum query flexibility** - Consider raw SQL or the SQL builder for complex reporting queries
- **Legacy tables without primary keys** - GORM requires a primary key; use raw SQL for truly legacy schemas
- **When `gorm.Model` fields aren't needed** - Don't embed `gorm.Model` if you don't want `ID`, `CreatedAt`, `UpdatedAt`, `DeletedAt`

## References

- [Official GORM Documentation: Models](https://gorm.io/docs/models.html)
- [Conventions](https://gorm.io/docs/conventions.html)
- [Serializer](https://gorm.io/docs/serializer.html)
- [Indexes](https://gorm.io/docs/indexes.html)
- [Constraints](https://gorm.io/docs/constraints.html)
