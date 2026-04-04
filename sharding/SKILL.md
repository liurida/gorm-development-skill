---
name: gorm-sharding
description: Use when implementing horizontal database sharding, partitioning large tables across multiple databases, or configuring shard-aware primary key generation in GORM applications.
---

# Sharding

The Sharding plugin partitions large tables into smaller ones using SQL parsing, redirecting queries to the appropriate shard.

**Reference:** https://gorm.io/docs/sharding.html
**Repository:** https://github.com/go-gorm/sharding

## Features

- Non-intrusive design (plugin-based)
- No network middleware overhead
- PostgreSQL and MySQL support
- Built-in primary key generators (Snowflake, PostgreSQL Sequence, custom)

## Quick Reference

| Config | Purpose |
|--------|---------|
| `ShardingKey` | Column used to determine shard |
| `NumberOfShards` | Total number of shards (e.g., 64) |
| `PrimaryKeyGenerator` | ID generation strategy |
| `ShardingAlgorithm` | Custom shard routing function |

## Basic Setup

```go
import (
  "gorm.io/driver/postgres"
  "gorm.io/gorm"
  "gorm.io/sharding"
)

db, err := gorm.Open(postgres.New(postgres.Config{
  DSN: "postgres://localhost:5432/sharding-db?sslmode=disable",
}))

// Register sharding for specific tables
db.Use(sharding.Register(sharding.Config{
  ShardingKey:         "user_id",
  NumberOfShards:      64,
  PrimaryKeyGenerator: sharding.PKSnowflake,
}, "orders", "order_items", "notifications"))
```

## Primary Key Generators

```go
// Snowflake IDs (recommended for distributed systems)
sharding.PKSnowflake

// PostgreSQL Sequence
sharding.PKPGSequence

// Custom generator
sharding.PKCustom
```

## CRUD Operations

**Critical:** All operations on sharded tables MUST include the sharding key.

```go
// CREATE - sharding key determines target shard
db.Create(&Order{UserID: 2})
// INSERT INTO orders_02 ...

// Raw SQL insert
db.Exec("INSERT INTO orders(user_id) VALUES(?)", int64(3))
// INSERT INTO orders_03 ...

// READ - must include sharding key in WHERE clause
var orders []Order
db.Model(&Order{}).Where("user_id", int64(2)).Find(&orders)
// SELECT * FROM orders_02 WHERE user_id = 2

// Raw SQL query
db.Raw("SELECT * FROM orders WHERE user_id = ?", int64(3)).Scan(&orders)
// SELECT * FROM orders_03 WHERE user_id = 3

// UPDATE - must include sharding key
db.Exec("UPDATE orders SET product_id = ? WHERE user_id = ?", 2, int64(3))
// UPDATE orders_03 SET product_id = 2 WHERE user_id = 3

// DELETE - must include sharding key
db.Where("user_id = ?", int64(2)).Delete(&Order{})
// DELETE FROM orders_02 WHERE user_id = 2
```

## Error Handling

Operations without sharding key will fail with `ErrMissingShardingKey`.

```go
// ERROR: Missing sharding key in CREATE
err := db.Create(&Order{Amount: 10, ProductID: 100}).Error
// err = ErrMissingShardingKey

// ERROR: Missing sharding key in WHERE clause
err = db.Model(&Order{}).Where("product_id", "1").Find(&orders).Error
// err = ErrMissingShardingKey

// ERROR: DELETE without sharding key
err = db.Exec("DELETE FROM orders WHERE product_id = 3").Error
// err = ErrMissingShardingKey
```

## Custom Sharding Algorithm

Override the default modulo-based sharding.

```go
db.Use(sharding.Register(sharding.Config{
  ShardingKey:    "user_id",
  NumberOfShards: 64,
  ShardingAlgorithm: func(columnValue interface{}) (suffix string, err error) {
    // Custom logic to determine shard suffix
    userID := columnValue.(int64)
    if userID < 1000 {
      return "_legacy", nil
    }
    return fmt.Sprintf("_%02d", userID%64), nil
  },
  PrimaryKeyGenerator: sharding.PKSnowflake,
}, "orders"))
```

## Multiple Tables with Same Sharding Rule

```go
// All these tables will use the same sharding configuration
db.Use(sharding.Register(sharding.Config{
  ShardingKey:         "user_id",
  NumberOfShards:      64,
  PrimaryKeyGenerator: sharding.PKSnowflake,
}, "orders", "order_items", "payments", "notifications"))
```

## Table Creation

You must create the sharded tables beforehand.

```go
// For 64 shards, create tables: orders_00, orders_01, ..., orders_63
for i := 0; i < 64; i++ {
  tableName := fmt.Sprintf("orders_%02d", i)
  db.Table(tableName).AutoMigrate(&Order{})
}
```

## Common Mistakes

| Mistake | Fix |
|---------|-----|
| Querying without sharding key | Always include sharding key in WHERE clause |
| Cross-shard JOINs | Denormalize data or use application-level joins |
| Forgetting to create shard tables | Create all `table_XX` tables before use |
| Using auto-increment IDs | Use Snowflake or distributed ID generator |
| Changing sharding key after data exists | Plan sharding key carefully upfront |

## When NOT to Use

- **Small to medium datasets** - Sharding adds complexity; single database handles millions of rows well
- **Queries requiring cross-shard JOINs** - Sharding doesn't support cross-shard queries; redesign data model first
- **Rapidly changing sharding requirements** - Re-sharding is extremely difficult; ensure sharding key is stable
- **When vertical scaling is sufficient** - Upgrade hardware before adding horizontal complexity
- **Applications with many ad-hoc queries** - All queries must include sharding key; analytics workloads don't fit
- **Teams without distributed systems experience** - Sharding introduces significant operational complexity

## Design Considerations

1. **Choose sharding key carefully** - Should distribute data evenly and be present in most queries
2. **Avoid cross-shard queries** - Design data model to minimize cross-shard operations
3. **Plan for growth** - Choose `NumberOfShards` with future scale in mind
4. **Test shard distribution** - Verify even distribution before production
