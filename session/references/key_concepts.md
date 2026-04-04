# Key Concepts for GORM Session

GORM's `Session` method creates a new session with specific configurations, allowing fine-grained control over database operations.

## Session Configuration Options

### DryRun Mode

Generate SQL without executing it. Useful for debugging and testing queries.

```go
db.Session(&gorm.Session{DryRun: true}).First(&user, 1)
// Generates SQL but does not execute it
```

### PrepareStmt Mode

Caches prepared statements for performance improvement.

```go
db.Session(&gorm.Session{PrepareStmt: true})
```

### NewDB Mode

Creates a fresh session without any inherited conditions, clauses, or context.

```go
db.Session(&gorm.Session{NewDB: true})
```

### SkipHooks Mode

Skips all hooks (BeforeCreate, AfterCreate, etc.) during operations.

```go
db.Session(&gorm.Session{SkipHooks: true})
```

### AllowGlobalUpdate Mode

By default, GORM does not allow global UPDATE/DELETE without conditions. This mode allows it.

```go
db.Session(&gorm.Session{AllowGlobalUpdate: true}).Model(&User{}).Update("Name", "jinzhu")
```

### SkipDefaultTransaction Mode

Disables the default transaction wrapper for write operations, improving performance.

```go
db.Session(&gorm.Session{SkipDefaultTransaction: true})
```

## Continuous Session Mode

You can reuse a session for multiple operations:

```go
tx := db.Session(&gorm.Session{SkipDefaultTransaction: true})
tx.First(&user, 1)
tx.Find(&users)
tx.Model(&user).Update("Age", 18)
```

## Important Notes

- Sessions inherit settings from the parent `*gorm.DB` unless explicitly overridden
- Use `NewDB: true` when you need a clean slate without inherited conditions
- Combine multiple session options as needed for your use case
