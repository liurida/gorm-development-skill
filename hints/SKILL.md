---
name: gorm-hints
description: Use when adding optimizer hints, index hints, or SQL comments to GORM queries for query optimization, performance tuning, or database routing.
---

# Hints

GORM supports optimizer hints, index hints, and comment hints for query optimization and SQL routing.

**Reference:** https://gorm.io/docs/hints.html
**Repository:** https://github.com/go-gorm/hints

## Quick Reference

| Hint Type | Purpose |
|-----------|---------|
| Optimizer Hints | Guide query optimizer (timeouts, parallelism) |
| Index Hints | Force, suggest, or ignore specific indexes |
| Comment Hints | Add comments for routing/debugging |

## Optimizer Hints

Guide the database query optimizer with database-specific directives.

```go
import "gorm.io/hints"

// MySQL MAX_EXECUTION_TIME
db.Clauses(hints.New("MAX_EXECUTION_TIME(10000)")).Find(&User{})
// SELECT /*+ MAX_EXECUTION_TIME(10000) */ * FROM `users`

// PostgreSQL parallel query
db.Clauses(hints.New("parallel(users, 4)")).Find(&User{})
// SELECT /*+ parallel(users, 4) */ FROM "users"

// Multiple hints
db.Clauses(
  hints.New("MAX_EXECUTION_TIME(5000)"),
  hints.New("NO_INDEX_MERGE(users)"),
).Find(&User{})
```

## Index Hints

Control which indexes the query optimizer uses.

### USE INDEX

Suggest indexes for the optimizer to consider.

```go
db.Clauses(hints.UseIndex("idx_user_name")).Find(&User{})
// SELECT * FROM `users` USE INDEX (`idx_user_name`)

// Multiple indexes
db.Clauses(hints.UseIndex("idx_user_name", "idx_user_email")).Find(&User{})
// SELECT * FROM `users` USE INDEX (`idx_user_name`,`idx_user_email`)
```

### FORCE INDEX

Require the optimizer to use specific indexes.

```go
// Force index for entire query
db.Clauses(hints.ForceIndex("idx_user_name")).Find(&User{})
// SELECT * FROM `users` FORCE INDEX (`idx_user_name`)

// Force index for JOIN operations only
db.Clauses(hints.ForceIndex("idx_user_name", "idx_user_id").ForJoin()).Find(&User{})
// SELECT * FROM `users` FORCE INDEX FOR JOIN (`idx_user_name`,`idx_user_id`)

// Force index for ORDER BY
db.Clauses(hints.ForceIndex("idx_created_at").ForOrderBy()).Find(&User{})
// SELECT * FROM `users` FORCE INDEX FOR ORDER BY (`idx_created_at`)
```

### IGNORE INDEX

Prevent the optimizer from using specific indexes.

```go
// Ignore index for GROUP BY
db.Clauses(hints.IgnoreIndex("idx_user_name").ForGroupBy()).Find(&User{})
// SELECT * FROM `users` IGNORE INDEX FOR GROUP BY (`idx_user_name`)
```

### Combined Index Hints

```go
db.Clauses(
  hints.ForceIndex("idx_user_name", "idx_user_id").ForOrderBy(),
  hints.IgnoreIndex("idx_user_name").ForGroupBy(),
).Find(&User{})
// SELECT * FROM `users` FORCE INDEX FOR ORDER BY (`idx_user_name`,`idx_user_id`) IGNORE INDEX FOR GROUP BY (`idx_user_name`)
```

## Comment Hints

Add SQL comments for routing, debugging, or query identification.

```go
// Comment after SELECT
db.Clauses(hints.Comment("select", "master")).Find(&User{})
// SELECT /*master*/ * FROM `users`

// Comment before statement
db.Clauses(hints.CommentBefore("insert", "node2")).Create(&user)
// /*node2*/ INSERT INTO `users` ...

// Comment after WHERE clause
db.Clauses(hints.CommentAfter("where", "hint")).Find(&User{}, "id = ?", 1)
// SELECT * FROM `users` WHERE id = ? /* hint */
```

### Use Cases for Comments

**Database Routing (ProxySQL, Vitess)**
```go
// Route to master
db.Clauses(hints.Comment("select", "master")).Find(&user)

// Route to specific shard
db.Clauses(hints.CommentBefore("select", "/*shard:users_1*/")).Find(&user)
```

**Query Identification**
```go
// Tag queries for monitoring
db.Clauses(hints.Comment("select", "source:user-service,endpoint:get-user")).
  First(&user, id)
```

## Performance Optimization Examples

### Force Covering Index

```go
// Use covering index to avoid table lookups
db.Clauses(hints.ForceIndex("idx_user_name_email")).
  Select("name", "email").
  Find(&users)
```

### Index for Sorting

```go
// Avoid filesort by forcing index
db.Clauses(hints.ForceIndex("idx_created_at").ForOrderBy()).
  Order("created_at DESC").
  Limit(100).
  Find(&users)
```

### Index for Joins

```go
// Optimize join with index hint
db.Clauses(hints.ForceIndex("idx_user_id").ForJoin()).
  Joins("JOIN orders ON orders.user_id = users.id").
  Find(&users)
```

## Combining with Other Features

```go
// With DBResolver for read/write splitting
db.Clauses(
  dbresolver.Write,
  hints.Comment("select", "after-write-read"),
).First(&user)

// With preloading
db.Clauses(hints.UseIndex("idx_user_status")).
  Preload("Orders", func(db *gorm.DB) *gorm.DB {
    return db.Clauses(hints.UseIndex("idx_order_user_id"))
  }).
  Where("status = ?", "active").
  Find(&users)
```

## When NOT to Use

- **Without profiling and `EXPLAIN` analysis** - Don't use hints as a first resort. Trust the query optimizer until you can prove it's making a mistake.
- **On small tables** - The overhead of parsing the hint may outweigh the benefit for tables with few records.
- **When the underlying data distribution changes frequently** - A hint that works today might be counterproductive after a large data import or deletion. Re-evaluate hints periodically.
- **If your application needs to be database-agnostic** - Optimizer and index hints are often database-specific and will break portability.
- **As a substitute for proper indexing** - If a query is slow, the first step is to ensure the correct indexes exist, not to add a hint.

## Common Mistakes

| Mistake | Fix |
|---------|-----|
| Using hints with non-existent indexes | Verify index exists first |
| Force index on small tables | Let optimizer decide for small tables |
| Ignoring EXPLAIN output | Always verify with EXPLAIN |
| Hint syntax varies by database | Test hints on target database |
| Over-using hints | Only use when optimizer makes poor choices |
