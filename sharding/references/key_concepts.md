
# Key Concepts for GORM Sharding

This document provides key concepts for using the GORM sharding plugin to partition large tables.

## Overview

The `gorm.io/sharding` plugin allows you to implement database sharding, which is the process of splitting a large table into smaller, more manageable pieces (shards). This is a common technique for horizontally scaling databases.

## How It Works

The sharding plugin works by intercepting GORM queries. It parses the SQL to identify the target table and looks for a **sharding key** in the `WHERE` clause. Based on the value of the sharding key, it dynamically rewrites the query to target a specific shard table (e.g., `orders_01`, `orders_02`, etc.).

## Configuration

You configure the sharding plugin by registering it with GORM and providing a `sharding.Config`.

```go
import "gorm.io/sharding"

db.Use(sharding.Register(sharding.Config{
    // The field name to be used as the sharding key
    ShardingKey:         "user_id",
    // The total number of shards
    NumberOfShards:      64,
    // Optional: A primary key generator for new records
    PrimaryKeyGenerator: sharding.PKSnowflake,
}, "orders")) // The list of tables to apply sharding to
```

- **`ShardingKey`**: This is the most important setting. It's the name of the column that will be used to determine which shard a row belongs to.
- **`NumberOfShards`**: The total number of shard tables you will have.
- **`PrimaryKeyGenerator`**: Since data is distributed, you can't rely on auto-incrementing primary keys. The plugin provides generators like Snowflake to create unique IDs across all shards.

## Usage

After configuring the plugin, you use GORM as you normally would. However, there's a critical requirement:

**All queries (reads and writes) on a sharded table MUST include the sharding key in the `WHERE` clause.**

If you don't provide the sharding key, the plugin cannot determine which shard to route the query to and will return an `ErrMissingShardingKey` error.

### Correct Usage (with Sharding Key)

```go
// Create will be routed to the correct shard based on UserID
db.Create(&Order{UserID: 123, Amount: 100})

// Find will query the specific shard for UserID 123
var orders []Order
db.Where("user_id = ?", 123).Find(&orders)
```

### Incorrect Usage (without Sharding Key)

```go
// This will return ErrMissingShardingKey
db.Create(&Order{Amount: 100})

// This will also return ErrMissingShardingKey
db.Where("amount > ?", 50).Find(&orders)
```

## Transactions

Transactions on sharded tables are complex. A single GORM transaction will operate on a single shard, determined by the first operation in the transaction that includes a sharding key. Cross-shard transactions are not supported out of the box and require a distributed transaction manager (which is outside the scope of the GORM sharding plugin).

Sharding is a powerful but advanced technique. It introduces complexity, so it should only be used when you have very large tables that are causing performance bottlenecks.
