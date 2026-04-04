---
name: gorm-generic-interface
description: Use when working with GORM's generic database interface to get the underlying `*sql.DB` object and configure the connection pool.
---

# GORM Generic Database Interface

GORM provides a generic database interface that allows you to get the underlying `*sql.DB` object from a GORM `*gorm.DB` instance. This is useful for working with other Go-based tools and libraries that expect a `*sql.DB` object, and for configuring the connection pool.

## Getting the Generic DB Object

You can get the underlying `*sql.DB` object using the `DB()` method.

```go
sqlDB, err := db.DB()
if err != nil {
    // handle error
}

// You can now use sqlDB with other libraries that expect a *sql.DB
```

## Connection Pool Configuration

Once you have the `*sql.DB` object, you can configure the connection pool for better performance and resource management.

```go
// SetMaxIdleConns sets the maximum number of connections in the idle connection pool.
sqlDB.SetMaxIdleConns(10)

// SetMaxOpenConns sets the maximum number of open connections to the database.
sqlDB.SetMaxOpenConns(100)

// SetConnMaxLifetime sets the maximum amount of time a connection may be reused.
sqlDB.SetConnMaxLifetime(time.Hour)
```

- **`SetMaxIdleConns`**: Use this to keep a certain number of idle connections in the pool, which can improve performance by avoiding the cost of establishing a new connection for every operation.
- **`SetMaxOpenConns`**: This limits the total number of open connections to the database, which is crucial for preventing your application from overwhelming the database with too many connections.
- **`SetConnMaxLifetime`**: This is useful for environments where connections can become stale, such as when there are network devices that close idle connections.

## When NOT to Use

- **For standard GORM operations** - Use the `*gorm.DB` instance directly for `Create`, `Find`, `Update`, etc. You don't need the underlying `*sql.DB` for normal usage.
- **When inside a GORM transaction** - GORM manages the transaction on its `*gorm.DB` object. Using the raw `*sql.DB` will execute queries outside the transaction.
- **When using plugins like DBResolver or Sharding** - Getting the `*sql.DB` will bypass the plugin's logic (e.g., read/write splitting). Configure the pool through the plugin's interface if available.
- **If GORM provides a higher-level configuration** - For features like prepared statement caching, use GORM's config (`PrepareStmt: true`) instead of managing it at a lower level.
