---
name: gorm-performance
description: Use when optimizing GORM query performance, reducing database latency, or troubleshooting slow queries. Covers prepared statements, batch operations, index hints, field selection, and connection pooling.
---

# Performance

GORM optimizes many things by default, but there are several techniques to improve performance for high-throughput applications.

**Reference:** https://gorm.io/docs/performance.html

## Quick Reference

| Technique | Improvement | Use Case |
|-----------|-------------|----------|
| `SkipDefaultTransaction` | ~30% faster writes | Single-row inserts without rollback needs |
| `PrepareStmt` | Faster repeated queries | High-frequency identical queries |
| `Select` fields | Reduced memory/bandwidth | Large tables with many columns |
| `FindInBatches` | Memory-efficient | Processing large result sets |
| Index Hints | Query optimization | Complex queries with multiple indexes |
| Connection Pool | Connection reuse | High-concurrency applications |

## Disable Default Transaction

GORM wraps write operations in transactions for data consistency. Disable when not needed.

```go
// Global configuration - ~30% performance improvement for writes
db, err := gorm.Open(sqlite.Open("gorm.db"), &gorm.Config{
  SkipDefaultTransaction: true,
})

// Per-session configuration
tx := db.Session(&gorm.Session{SkipDefaultTransaction: true})
tx.Create(&user)
```

**When to use:** Single-row inserts where you handle transactions manually or don't need atomicity.

**When NOT to use:** Multi-table operations, financial transactions, or when data consistency is critical.

## Prepared Statement Caching

Creates prepared statements when executing any SQL and caches them for future calls.

```go
// Global mode - all operations use prepared statements
db, err := gorm.Open(sqlite.Open("gorm.db"), &gorm.Config{
  PrepareStmt: true,
})

// Session mode - specific operations
tx := db.Session(&gorm.Session{PrepareStmt: true})
tx.First(&user, 1)
tx.Find(&users)
tx.Model(&user).Update("Age", 18)

// Works with raw SQL too
db.Raw("select sum(age) from users where role = ?", "admin").Scan(&age)

// Access prepared statement manager
stmtManger, ok := tx.ConnPool.(*PreparedStmtDB)
stmtManger.Close()  // Close statements for current session
```

**MySQL Tip:** Enable `interpolateparams` to reduce roundtrips:
```go
dsn := "user:pass@tcp(127.0.0.1:3306)/db?interpolateParams=true"
```

## Select Only Required Fields

Reduce memory usage and network bandwidth by selecting specific fields.

```go
// Explicit field selection
db.Select("Name", "Age").Find(&users)

// Smart select with smaller struct - GORM auto-selects matching fields
type User struct {
  ID     uint
  Name   string
  Age    int
  Gender string
  // hundreds of fields
}

type APIUser struct {
  ID   uint
  Name string
}

// Automatically selects only `id`, `name`
db.Model(&User{}).Limit(10).Find(&APIUser{})
// SELECT `id`, `name` FROM `users` LIMIT 10
```

## Batch Processing

Process large datasets without loading everything into memory.

```go
// FindInBatches - process records in chunks
result := db.Where("processed = ?", false).FindInBatches(&results, 100, func(tx *gorm.DB, batch int) error {
  for _, result := range results {
    // Process each record
  }
  return nil  // Return error to stop
})

// CreateInBatches - bulk insert with controlled batch size
users := make([]User, 10000)
db.CreateInBatches(users, 100)  // Insert 100 at a time

// Session-level batch size
db.Session(&gorm.Session{CreateBatchSize: 1000}).Create(&users)
```

## Index Hints

Guide the query optimizer to use specific indexes.

```go
import "gorm.io/hints"

// Use specific index
db.Clauses(hints.UseIndex("idx_user_name")).Find(&User{})
// SELECT * FROM `users` USE INDEX (`idx_user_name`)

// Force index for JOIN operations
db.Clauses(hints.ForceIndex("idx_user_name", "idx_user_id").ForJoin()).Find(&User{})
// SELECT * FROM `users` FORCE INDEX FOR JOIN (`idx_user_name`,`idx_user_id`)

// Multiple hints for different operations
db.Clauses(
  hints.ForceIndex("idx_user_name", "idx_user_id").ForOrderBy(),
  hints.IgnoreIndex("idx_user_name").ForGroupBy(),
).Find(&User{})
```

## Connection Pool Configuration

Configure the underlying connection pool for high-concurrency scenarios.

```go
sqlDB, err := db.DB()

// Idle connections - keep warm connections ready
sqlDB.SetMaxIdleConns(10)

// Max open connections - limit concurrent connections
sqlDB.SetMaxOpenConns(100)

// Connection lifetime - prevent stale connections
sqlDB.SetConnMaxLifetime(time.Hour)

// Idle timeout (Go 1.15+)
sqlDB.SetConnMaxIdleTime(10 * time.Minute)
```

**Guidelines:**
- `MaxOpenConns`: Set based on database server capacity (typically 25-100)
- `MaxIdleConns`: Set to ~25% of `MaxOpenConns`
- `ConnMaxLifetime`: Set to less than database server's wait_timeout

## Read/Write Splitting

For high-throughput applications, use DBResolver to split reads and writes.

```go
import "gorm.io/plugin/dbresolver"

db.Use(dbresolver.Register(dbresolver.Config{
  Sources:  []gorm.Dialector{mysql.Open("db_write_dsn")},
  Replicas: []gorm.Dialector{mysql.Open("db_read1_dsn"), mysql.Open("db_read2_dsn")},
  Policy:   dbresolver.RandomPolicy{},
}))
```

See `gorm-dbresolver` skill for details.

## When NOT to Use

- **Premature optimization** - Don't apply these techniques without profiling; measure first to identify actual bottlenecks
- **`SkipDefaultTransaction` with multi-table operations** - Keep transactions for operations requiring atomicity
- **`PrepareStmt` with highly dynamic queries** - Statement caching provides no benefit when every query is unique
- **Connection pool tuning without load testing** - Default settings work well for most applications; tune only under production-like load
- **`FindInBatches` for small datasets** - Adds complexity with no benefit for < 1000 records
- **Index hints without EXPLAIN analysis** - Let the query optimizer do its job unless you've proven it's making poor choices

## Common Mistakes

| Mistake | Fix |
|---------|-----|
| Loading all columns for API responses | Use `Select()` or smaller structs |
| Processing millions of rows at once | Use `FindInBatches` |
| Not setting connection pool limits | Configure `SetMaxOpenConns` |
| Disabling transactions for multi-table ops | Keep transactions enabled |
| Not using indexes | Add `gorm:"index"` tags, use EXPLAIN |
