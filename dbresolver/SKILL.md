---
name: gorm-dbresolver
description: Use when implementing read/write splitting, multiple database sources, automatic connection switching, or load balancing across database replicas in GORM applications.
---

# Database Resolver

DBResolver adds multiple database support to GORM for read/write splitting and connection management.

**Reference:** https://gorm.io/docs/dbresolver.html
**Repository:** https://github.com/go-gorm/dbresolver

## Features

- Multiple sources and replicas
- Read/Write splitting (automatic)
- Connection switching by table/struct
- Manual connection switching
- Load balancing across replicas
- Works with raw SQL and transactions

## Quick Reference

| Operation | Connection Used |
|-----------|----------------|
| `Create`, `Update`, `Delete` | Sources |
| `Find`, `First`, `Scan` | Replicas |
| `SELECT ... FOR UPDATE` | Sources |
| Raw `SELECT` statements | Replicas |
| Inside transaction | Same connection throughout |

## Basic Setup

```go
import (
  "gorm.io/gorm"
  "gorm.io/plugin/dbresolver"
  "gorm.io/driver/mysql"
)

db, err := gorm.Open(mysql.Open("db1_dsn"), &gorm.Config{})

db.Use(dbresolver.Register(dbresolver.Config{
  // Write operations go to sources
  Sources:  []gorm.Dialector{mysql.Open("db2_dsn")},
  // Read operations go to replicas
  Replicas: []gorm.Dialector{mysql.Open("db3_dsn"), mysql.Open("db4_dsn")},
  // Load balancing policy
  Policy: dbresolver.RandomPolicy{},
  // Log which connection is used
  TraceResolverMode: true,
}))
```

## Multiple Resolver Configuration

Configure different databases for different tables/models.

```go
db.Use(dbresolver.Register(dbresolver.Config{
  // Global resolver: db2 for writes, db3/db4 for reads
  Sources:  []gorm.Dialector{mysql.Open("db2_dsn")},
  Replicas: []gorm.Dialector{mysql.Open("db3_dsn"), mysql.Open("db4_dsn")},
  Policy:   dbresolver.RandomPolicy{},
}).Register(dbresolver.Config{
  // User/Address use default connection for writes, db5 for reads
  Replicas: []gorm.Dialector{mysql.Open("db5_dsn")},
}, &User{}, &Address{}).Register(dbresolver.Config{
  // Orders/Product use db6/db7 for writes, db8 for reads
  Sources:  []gorm.Dialector{mysql.Open("db6_dsn"), mysql.Open("db7_dsn")},
  Replicas: []gorm.Dialector{mysql.Open("db8_dsn")},
}, "orders", &Product{}, "secondary"))
```

## Automatic Connection Switching

DBResolver automatically routes queries based on operation type.

```go
// User Resolver Examples
db.Table("users").Rows()                            // replicas (db5)
db.Model(&User{}).Find(&AdvancedUser{})             // replicas (db5)
db.Exec("update users set name = ?", "jinzhu")      // sources (db1)
db.Raw("select name from users").Row().Scan(&name)  // replicas (db5)
db.Create(&user)                                    // sources (db1)
db.Delete(&User{}, "name = ?", "jinzhu")            // sources (db1)

// Global Resolver Examples
db.Find(&Pet{})  // replicas (db3 or db4)
db.Save(&Pet{})  // sources (db2)

// Orders Resolver Examples
db.Find(&Order{})                    // replicas (db8)
db.Table("orders").Find(&Report{})   // replicas (db8)
```

## Manual Connection Switching

Override automatic routing when needed.

```go
// Force read from source (e.g., for consistency after write)
db.Clauses(dbresolver.Write).First(&user)

// Use specific resolver by name
db.Clauses(dbresolver.Use("secondary")).First(&user)

// Combine resolver selection with write mode
db.Clauses(dbresolver.Use("secondary"), dbresolver.Write).First(&user)
```

## Transactions

Transactions maintain a single connection throughout.

```go
// Start transaction on default replicas
tx := db.Clauses(dbresolver.Read).Begin()

// Start transaction on default sources
tx := db.Clauses(dbresolver.Write).Begin()

// Start transaction on specific resolver's sources
tx := db.Clauses(dbresolver.Use("secondary"), dbresolver.Write).Begin()

// All operations in transaction use the same connection
tx.Create(&order)
tx.Create(&orderItems)
tx.Commit()
```

## Connection Pool Configuration

Configure pool settings per resolver.

```go
db.Use(
  dbresolver.Register(dbresolver.Config{ /* config */ }).
  SetConnMaxIdleTime(time.Hour).
  SetConnMaxLifetime(24 * time.Hour).
  SetMaxIdleConns(100).
  SetMaxOpenConns(200),
)
```

## Custom Load Balancing Policy

Implement the `Policy` interface for custom routing.

```go
type Policy interface {
  Resolve([]gorm.ConnPool) gorm.ConnPool
}

// Example: Round-robin policy
type RoundRobinPolicy struct {
  counter uint64
}

func (p *RoundRobinPolicy) Resolve(pools []gorm.ConnPool) gorm.ConnPool {
  n := atomic.AddUint64(&p.counter, 1)
  return pools[n%uint64(len(pools))]
}

// Usage
db.Use(dbresolver.Register(dbresolver.Config{
  Replicas: []gorm.Dialector{mysql.Open("db1"), mysql.Open("db2")},
  Policy:   &RoundRobinPolicy{},
}))
```

## When NOT to Use

- **Single database deployments** - DBResolver adds complexity with no benefit if you don't have replicas
- **Strong consistency requirements** - Replicas have replication lag; reads immediately after writes may return stale data
- **Simple applications** - The overhead isn't worth it unless you have genuine scalability needs
- **When all queries need latest data** - If most queries require `dbresolver.Write`, you're not benefiting from read replicas
- **Cross-region deployments with latency concerns** - DBResolver doesn't handle geo-routing; use a dedicated solution
- **Without proper monitoring** - You need to track replica lag and connection health

## Common Mistakes

| Mistake | Fix |
|---------|-----|
| Reading immediately after write expects new data | Use `dbresolver.Write` for consistent reads |
| Different resolvers for related tables | Keep related tables on same resolver |
| Not configuring connection pools | Set pool limits per resolver |
| Forgetting replication lag exists | Design for eventual consistency |
