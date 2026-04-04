---
name: gorm-session
description: Use when creating GORM sessions with custom configurations including DryRun, PrepareStmt, NewDB, SkipHooks, AllowGlobalUpdate, Context, Logger, and batch operations.
---

# Session

GORM's `Session` method creates a new session with specific configurations, allowing fine-grained control over database operations.

## Session Config Struct

```go
type Session struct {
  DryRun                   bool
  PrepareStmt              bool
  NewDB                    bool
  Initialized              bool
  SkipHooks                bool
  SkipDefaultTransaction   bool
  DisableNestedTransaction bool
  AllowGlobalUpdate        bool
  FullSaveAssociations     bool
  QueryFields              bool
  Context                  context.Context
  Logger                   logger.Interface
  NowFunc                  func() time.Time
  CreateBatchSize          int
}
```

## Quick Reference

| Option | Description |
|--------|-------------|
| `DryRun` | Generate SQL without executing |
| `PrepareStmt` | Cache prepared statements |
| `NewDB` | Fresh session without inherited conditions |
| `SkipHooks` | Skip BeforeCreate, AfterCreate, etc. |
| `AllowGlobalUpdate` | Allow UPDATE/DELETE without WHERE |
| `FullSaveAssociations` | Upsert associations when saving |
| `QueryFields` | Select specific fields instead of `*` |
| `CreateBatchSize` | Default batch size for creates |

## DryRun Mode

Generate SQL without executing (useful for debugging and testing):

```go
// Session mode
stmt := db.Session(&gorm.Session{DryRun: true}).First(&user, 1).Statement
stmt.SQL.String() // => SELECT * FROM `users` WHERE `id` = $1 ORDER BY `id`
stmt.Vars         // => []interface{}{1}

// Generate final SQL (NOTE: not safe to execute directly, may have SQL injection)
db.Dialector.Explain(stmt.SQL.String(), stmt.Vars...)
// => SELECT * FROM `users` WHERE `id` = 1
```

## PrepareStmt Mode

Cache prepared statements for performance:

```go
// Session mode
tx := db.Session(&gorm.Session{PrepareStmt: true})
tx.First(&user, 1)
tx.Find(&users)
tx.Model(&user).Update("Age", 18)

// Access prepared statements manager
stmtManger, ok := tx.ConnPool.(*PreparedStmtDB)

// Close prepared statements for current session
stmtManger.Close()

// Get all cached prepared SQL
stmtManger.PreparedSQL // => []string{}

// Access all prepared statements
for sql, stmt := range stmtManger.Stmts {
  sql  // prepared SQL
  stmt // *sql.Stmt
  stmt.Close() // close individual statement
}
```

## NewDB Mode

Create a session without inherited conditions:

```go
tx := db.Where("name = ?", "jinzhu").Session(&gorm.Session{NewDB: true})

tx.First(&user)
// SELECT * FROM users ORDER BY id LIMIT 1
// Note: no WHERE clause inherited

tx.First(&user, "id = ?", 10)
// SELECT * FROM users WHERE id = 10 ORDER BY id

// Without NewDB: conditions are inherited
tx2 := db.Where("name = ?", "jinzhu").Session(&gorm.Session{})
tx2.First(&user)
// SELECT * FROM users WHERE name = "jinzhu" ORDER BY id
```

## SkipHooks Mode

Skip all hooks (BeforeCreate, AfterCreate, etc.):

```go
DB.Session(&gorm.Session{SkipHooks: true}).Create(&user)
DB.Session(&gorm.Session{SkipHooks: true}).Create(&users)
DB.Session(&gorm.Session{SkipHooks: true}).CreateInBatches(users, 100)
DB.Session(&gorm.Session{SkipHooks: true}).Find(&user)
DB.Session(&gorm.Session{SkipHooks: true}).Delete(&user)
DB.Session(&gorm.Session{SkipHooks: true}).Model(User{}).Where("age > ?", 18).Updates(&user)
```

## AllowGlobalUpdate Mode

Allow UPDATE/DELETE without WHERE clause (use with caution):

```go
db.Session(&gorm.Session{AllowGlobalUpdate: true}).Model(&User{}).Update("name", "jinzhu")
// UPDATE users SET `name` = "jinzhu"
```

## DisableNestedTransaction

Disable SAVEPOINT usage for nested transactions:

```go
db.Session(&gorm.Session{DisableNestedTransaction: true}).CreateInBatches(&users, 100)
```

## FullSaveAssociations

Update associations' data using Upsert:

```go
db.Session(&gorm.Session{FullSaveAssociations: true}).Updates(&user)
// INSERT INTO "addresses" ... ON DUPLICATE KEY SET ...
// INSERT INTO "users" ...
// INSERT INTO "emails" ... ON DUPLICATE KEY SET ...
```

## Context

Set context for timeout/cancellation:

```go
timeoutCtx, cancel := context.WithTimeout(context.Background(), time.Second)
defer cancel()

tx := db.Session(&gorm.Session{Context: timeoutCtx})
tx.First(&user) // query with timeout
tx.Model(&user).Update("role", "admin") // update with timeout

// Shortcut method
db.WithContext(ctx).First(&user)
```

## Logger

Customize logging for the session:

```go
newLogger := logger.New(log.New(os.Stdout, "\r\n", log.LstdFlags),
  logger.Config{
    SlowThreshold: time.Second,
    LogLevel:      logger.Silent,
    Colorful:      false,
  })

db.Session(&gorm.Session{Logger: newLogger})

// Set log level
db.Session(&gorm.Session{Logger: logger.Default.LogMode(logger.Silent)})

// Debug shortcut
db.Debug().First(&user) // equivalent to LogMode(logger.Info)
```

## NowFunc

Override timestamp function for the session:

```go
db.Session(&gorm.Session{
  NowFunc: func() time.Time {
    return time.Now().Local()
  },
})
```

## QueryFields

Select specific fields instead of `*`:

```go
db.Session(&gorm.Session{QueryFields: true}).Find(&user)
// SELECT `users`.`name`, `users`.`age`, ... FROM `users`

// Without QueryFields:
// SELECT * FROM `users`
```

## CreateBatchSize

Set default batch size for bulk creates:

```go
users := [5000]User{{Name: "jinzhu", Pets: []Pet{pet1, pet2, pet3}}...}

db.Session(&gorm.Session{CreateBatchSize: 1000}).Create(&users)
// INSERT INTO users xxx (5 batches)
// INSERT INTO pets xxx (15 batches)
```

## Continuous Session Mode

Reuse a session for multiple operations:

```go
tx := db.Session(&gorm.Session{
  SkipDefaultTransaction: true,
  PrepareStmt:            true,
})

tx.First(&user, 1)
tx.Find(&users)
tx.Model(&user).Update("Age", 18)
```

## Initialized Mode

Create a new initialized DB (not Method Chain/Goroutine Safe):

```go
tx := db.Session(&gorm.Session{Initialized: true})
// Use only when you understand the implications
```

## Combining Options

```go
tx := db.Session(&gorm.Session{
  DryRun:                 false,
  PrepareStmt:            true,
  SkipDefaultTransaction: true,
  SkipHooks:              false,
  AllowGlobalUpdate:      false,
  QueryFields:            true,
  CreateBatchSize:        500,
  Context:                ctx,
})
```

## Common Mistakes

| Mistake | Fix |
|---------|-----|
| Using `NewDB` without understanding it clears conditions | Be explicit about which conditions to apply |
| Enabling `AllowGlobalUpdate` globally | Use session-level setting only when needed |
| Forgetting to close prepared statements | Call `stmtManger.Close()` when done |
| Using `Initialized` without understanding thread safety | Avoid unless you have specific requirements |

## When NOT to Use

- **For a single, simple query** - Creating a new session for one `db.First(&user)` call is unnecessary overhead.
- **`AllowGlobalUpdate` in application code** - This should be reserved for specific, controlled scripts or migrations, not general application logic.
- **`SkipHooks` as a default** - Hooks are a key part of GORM's functionality for validation and callbacks. Only skip them when you have a specific performance reason, like a bulk import.
- **`Initialized: true` unless you are an expert** - This mode is not goroutine-safe and breaks the normal chaining behavior. Avoid it unless you have a deep understanding of GORM internals and a specific need for it.

## References

- [Official GORM Documentation: Session](https://gorm.io/docs/session.html)
- [Method Chaining](https://gorm.io/docs/method_chaining.html)
- [Logger](https://gorm.io/docs/logger.html)
- [Context](https://gorm.io/docs/context.html)
