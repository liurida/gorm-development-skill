---
name: gorm-custom-data-types
description: Use when creating custom data types in GORM, implementing Scanner/Valuer interfaces, defining database-specific types, or using SQL expressions for create/update operations.
---

# Custom Data Types

GORM provides interfaces for defining customized data types that integrate seamlessly with database storage and retrieval.

**Reference:** https://gorm.io/docs/data_types.html

## Core Interfaces

Custom data types require implementing interfaces from the `database/sql` package:

| Interface | Package | Purpose |
|-----------|---------|---------|
| `Scanner` | `database/sql` | Receives values from the database |
| `Valuer` | `database/sql/driver` | Saves values to the database |
| `GormDataType` | `gorm.io/gorm/schema` | General data type identifier |
| `GormDBDataType` | `gorm.io/gorm/schema` | Database-specific type for migrations |
| `GormValue` | `gorm.io/gorm` | Create/update from SQL expressions |

## Scanner/Valuer Implementation

The foundation for any custom data type:

```go
type JSON json.RawMessage

// Scan implements sql.Scanner - receives value from database
func (j *JSON) Scan(value interface{}) error {
  bytes, ok := value.([]byte)
  if !ok {
    return errors.New(fmt.Sprint("Failed to unmarshal JSONB value:", value))
  }

  result := json.RawMessage{}
  err := json.Unmarshal(bytes, &result)
  *j = JSON(result)
  return err
}

// Value implements driver.Valuer - saves value to database
func (j JSON) Value() (driver.Value, error) {
  if len(j) == 0 {
    return nil, nil
  }
  return json.RawMessage(j).MarshalJSON()
}
```

## GormDataTypeInterface

Provides a general data type identifier accessible via `schema.Field`:

```go
type GormDataTypeInterface interface {
  GormDataType() string
}

func (JSON) GormDataType() string {
  return "json"
}

// Usage in hooks or plugins
func (user User) BeforeCreate(tx *gorm.DB) {
  field := tx.Statement.Schema.LookUpField("Attrs")
  if field.DataType == "json" {
    // do something
  }
}
```

## GormDBDataTypeInterface

Returns database-specific types during migrations:

```go
type GormDBDataTypeInterface interface {
  GormDBDataType(*gorm.DB, *schema.Field) string
}

func (JSON) GormDBDataType(db *gorm.DB, field *schema.Field) string {
  // Access field tags via field.Tag, field.TagSettings
  switch db.Dialector.Name() {
  case "mysql", "sqlite":
    return "JSON"
  case "postgres":
    return "JSONB"
  }
  return ""
}
```

## GormValuerInterface

Create/update from SQL expressions or context-based values:

```go
type GormValuerInterface interface {
  GormValue(ctx context.Context, db *gorm.DB) clause.Expr
}
```

### Example: Geometry Type with SQL Expression

```go
type Location struct {
  X, Y int
}

func (loc Location) GormDataType() string {
  return "geometry"
}

func (loc Location) GormValue(ctx context.Context, db *gorm.DB) clause.Expr {
  return clause.Expr{
    SQL:  "ST_PointFromText(?)",
    Vars: []interface{}{fmt.Sprintf("POINT(%d %d)", loc.X, loc.Y)},
  }
}

func (loc *Location) Scan(v interface{}) error {
  // Scan a value into struct from database driver
}

// Usage
db.Create(&User{Name: "jinzhu", Location: Location{X: 100, Y: 100}})
// INSERT INTO `users` (`name`,`point`) VALUES ("jinzhu",ST_PointFromText("POINT(100 100)"))
```

### Example: Context-Based Encryption

```go
type EncryptedString struct {
  Value string
}

func (es EncryptedString) GormValue(ctx context.Context, db *gorm.DB) (expr clause.Expr) {
  if encryptionKey, ok := ctx.Value("TenantEncryptionKey").(string); ok {
    return clause.Expr{SQL: "?", Vars: []interface{}{Encrypt(es.Value, encryptionKey)}}
  }
  db.AddError(errors.New("invalid encryption key"))
  return
}
```

## Using Third-Party Types

Types implementing Scanner/Valuer work with GORM:

```go
import (
  "github.com/google/uuid"
  "github.com/lib/pq"
)

type Post struct {
  ID     uuid.UUID      `gorm:"type:uuid;default:uuid_generate_v4()"`
  Title  string
  Tags   pq.StringArray `gorm:"type:text[]"`
}
```

## Clause Expression for Query Helpers

Build custom query helpers by implementing `clause.Expression`:

```go
type Expression interface {
  Build(builder Builder)
}

// Usage with datatypes package
db.Find(&user, datatypes.JSONQuery("attributes").HasKey("role"))
// MySQL: SELECT * FROM `users` WHERE JSON_EXTRACT(`attributes`, '$.role') IS NOT NULL
// PostgreSQL: SELECT * FROM "user" WHERE "attributes"::jsonb ? 'role'

db.Find(&user, datatypes.JSONQuery("attributes").Equals("jinzhu", "name"))
// MySQL: SELECT * FROM `user` WHERE JSON_EXTRACT(`attributes`, '$.name') = "jinzhu"
```

## Type Inference Fallback

When interfaces are not implemented, GORM infers from the first field:

```go
type NullString struct {
  String string // data type inferred from this field
  Valid  bool
}

type User struct {
  Name NullString // data type will be string
}
```

## When NOT to Use

- **For standard Go types** - `string`, `int`, `time.Time`, `bool` are handled automatically and don't need custom types.
- **When a simple alias is enough** - `type UserID int64` works without needing `Scanner`/`Valuer` if it maps to a native database type.
- **If the database has a better native type** - Don't store complex data as JSON if your database has native support for arrays, objects, or other structures.
- **When the logic is too complex** - If your `Scan`/`Value` methods become very complex, it might be a sign that the data should be modeled differently (e.g., as a separate table).
- **If it hurts query performance** - Storing data in a custom format (like encrypted blobs or compressed JSON) can make it impossible to query or index effectively.

## Common Mistakes

| Mistake | Fix |
|---------|-----|
| Forgetting `Scan` receiver is pointer | Use `func (j *JSON) Scan(...)` |
| Not handling nil in `Value()` | Check for empty/nil before processing |
| Wrong database type per dialect | Implement `GormDBDataType` for cross-database support |
| Not returning error from `Scan` | Always propagate unmarshal errors |

## Datatypes Collection

Pre-built custom types: https://github.com/go-gorm/datatypes

Includes: JSON, JSONSlice, JSONType, Date, Time, and more.
