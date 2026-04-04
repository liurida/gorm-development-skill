---
name: gorm-gorm-config
description: Use when initializing GORM with custom configuration including performance tuning, naming strategies, logging, transactions, and migration settings.
---

# GORM Config

GORM provides a `gorm.Config` struct for initialization configuration.

## Config Struct

```go
type Config struct {
  SkipDefaultTransaction                   bool
  NamingStrategy                           schema.Namer
  Logger                                   logger.Interface
  NowFunc                                  func() time.Time
  DryRun                                   bool
  PrepareStmt                              bool
  DisableNestedTransaction                 bool
  AllowGlobalUpdate                        bool
  DisableAutomaticPing                     bool
  DisableForeignKeyConstraintWhenMigrating bool
}
```

## Quick Reference

| Field | Default | Description |
|-------|---------|-------------|
| `SkipDefaultTransaction` | `false` | Skip transaction wrapper for write operations |
| `PrepareStmt` | `false` | Cache prepared statements |
| `DryRun` | `false` | Generate SQL without executing |
| `AllowGlobalUpdate` | `false` | Allow UPDATE/DELETE without WHERE |
| `DisableAutomaticPing` | `false` | Skip DB ping on initialization |
| `DisableForeignKeyConstraintWhenMigrating` | `false` | Skip FK creation in AutoMigrate |

## Performance Tuning

### SkipDefaultTransaction

GORM wraps write operations (create/update/delete) in transactions by default. Disable for better performance when not needed:

```go
db, err := gorm.Open(sqlite.Open("gorm.db"), &gorm.Config{
  SkipDefaultTransaction: true,
})
```

### PrepareStmt

Cache prepared statements to speed up repeated queries:

```go
db, err := gorm.Open(sqlite.Open("gorm.db"), &gorm.Config{
  PrepareStmt: true,
})
```

## NamingStrategy

Customize how GORM maps struct/field names to table/column names:

```go
db, err := gorm.Open(sqlite.Open("gorm.db"), &gorm.Config{
  NamingStrategy: schema.NamingStrategy{
    TablePrefix:   "t_",                              // table prefix (User -> t_users)
    SingularTable: true,                              // singular names (User -> user)
    NoLowerCase:   true,                              // skip snake_casing
    NameReplacer:  strings.NewReplacer("CID", "Cid"), // custom replacements
  },
})
```

**Namer Interface:**
```go
type Namer interface {
  TableName(table string) string
  SchemaName(table string) string
  ColumnName(table, column string) string
  JoinTableName(table string) string
  RelationshipFKName(Relationship) string
  CheckerName(table, column string) string
  IndexName(table, column string) string
}
```

## Logger

Configure GORM's logging behavior:

```go
db, err := gorm.Open(sqlite.Open("gorm.db"), &gorm.Config{
  Logger: logger.Default.LogMode(logger.Silent), // Silent, Error, Warn, Info
})
```

See the `logger` skill for custom logger implementation.

## NowFunc

Override the timestamp function for `CreatedAt`/`UpdatedAt` fields:

```go
db, err := gorm.Open(sqlite.Open("gorm.db"), &gorm.Config{
  NowFunc: func() time.Time {
    return time.Now().UTC() // Always use UTC
  },
})
```

## DryRun

Generate SQL without executing (useful for testing/debugging):

```go
db, err := gorm.Open(sqlite.Open("gorm.db"), &gorm.Config{
  DryRun: true,
})

// All operations generate SQL but don't execute
stmt := db.Find(&user, 1).Statement
stmt.SQL.String() // => SELECT * FROM `users` WHERE `id` = ?
```

## Transaction Control

### DisableNestedTransaction

Disable SAVEPOINT support for nested transactions:

```go
db, err := gorm.Open(sqlite.Open("gorm.db"), &gorm.Config{
  DisableNestedTransaction: true,
})
```

### AllowGlobalUpdate

Allow UPDATE/DELETE without WHERE clause (dangerous, use with caution):

```go
db, err := gorm.Open(sqlite.Open("gorm.db"), &gorm.Config{
  AllowGlobalUpdate: true, // Not recommended for production
})
```

## Database Initialization

### DisableAutomaticPing

Skip the automatic database ping on initialization:

```go
db, err := gorm.Open(sqlite.Open("gorm.db"), &gorm.Config{
  DisableAutomaticPing: true,
})
```

### DisableForeignKeyConstraintWhenMigrating

Skip foreign key creation during AutoMigrate:

```go
db, err := gorm.Open(sqlite.Open("gorm.db"), &gorm.Config{
  DisableForeignKeyConstraintWhenMigrating: true,
})
```

## Comprehensive Example

```go
import (
  "log"
  "os"
  "strings"
  "time"

  "gorm.io/driver/sqlite"
  "gorm.io/gorm"
  "gorm.io/gorm/logger"
  "gorm.io/gorm/schema"
)

func initDB() (*gorm.DB, error) {
  newLogger := logger.New(
    log.New(os.Stdout, "\r\n", log.LstdFlags),
    logger.Config{
      SlowThreshold:             200 * time.Millisecond,
      LogLevel:                  logger.Info,
      IgnoreRecordNotFoundError: true,
      Colorful:                  true,
    },
  )

  return gorm.Open(sqlite.Open("app.db"), &gorm.Config{
    SkipDefaultTransaction: true,
    PrepareStmt:            true,
    NamingStrategy: schema.NamingStrategy{
      TablePrefix:   "app_",
      SingularTable: false,
    },
    Logger: newLogger,
    NowFunc: func() time.Time {
      return time.Now().UTC()
    },
    DisableForeignKeyConstraintWhenMigrating: false,
  })
}
```

## Common Mistakes

| Mistake | Fix |
|---------|-----|
| Enabling `AllowGlobalUpdate` in production | Use session-level setting instead |
| Using `DryRun` and expecting data changes | DryRun only generates SQL |
| Not considering `PrepareStmt` memory usage | Monitor connection pool with many unique queries |

## When NOT to Use

- **For session-specific settings** - Use `db.Session(&gorm.Session{...})` to apply settings like `DryRun` or `SkipHooks` to a specific set of operations, not globally.
- **`AllowGlobalUpdate` in production code** - This is a dangerous setting that should only be used in controlled scripts or migrations, not in application code.
- **`SkipDefaultTransaction` for operations that must be atomic** - Only disable default transactions if you are certain the operation is safe to run without transactional integrity or if you are managing transactions manually.
- **`DryRun` in production logic** - `DryRun` is for testing, debugging, and SQL generation, not for live application logic.

## References

- [Official GORM Documentation: Config](https://gorm.io/docs/gorm_config.html)
- [Session Configuration](https://gorm.io/docs/session.html)
- [Logger Configuration](https://gorm.io/docs/logger.html)
