
# Key Concepts for GORM DBResolver

This document provides key concepts for using the DBResolver plugin in GORM for read/write splitting and multiple database support.

## Overview

The DBResolver plugin allows you to configure multiple database connections and automatically route queries to the appropriate database. This is commonly used for read/write splitting, where write operations go to a primary database and read operations are distributed among replica databases.

## Basic Configuration

To use DBResolver, you register it as a plugin with your GORM DB instance, providing configurations for your data sources and replicas.

```go
import "gorm.io/plugin/dbresolver"

db.Use(dbresolver.Register(dbresolver.Config{
    // Sources are for write operations
    Sources:  []gorm.Dialector{mysql.Open("primary_dsn")},
    // Replicas are for read operations
    Replicas: []gorm.Dialector{mysql.Open("replica1_dsn"), mysql.Open("replica2_dsn")},
    // Policy determines how to select a replica (e.g., Random, RoundRobin)
    Policy:   dbresolver.RandomPolicy{},
}))
```

## Automatic Read/Write Splitting

By default, DBResolver automatically routes queries based on the operation type:
- **Write operations** (`Create`, `Update`, `Delete`, `Save`) are sent to one of the `Sources`.
- **Read operations** (`Find`, `First`, `Take`, raw queries starting with `SELECT`) are sent to one of the `Replicas`.

## Manual Control

You can manually specify which connection to use for a particular query using clauses.

### Forcing Write Mode

To force a read query to use a write connection (source), use `dbresolver.Write`.

```go
var user User
// This will read from a source database instead of a replica
db.Clauses(dbresolver.Write).First(&user, 1)
```

### Specifying a Resolver

If you have multiple named resolver configurations, you can choose a specific one with `dbresolver.Use`.

```go
// Assuming a resolver named "analytics" is configured
db.Clauses(dbresolver.Use("analytics")).Find(&reports)
```

## Transactions

When you start a transaction, all operations within that transaction will be performed on the same connection. DBResolver will not switch between sources and replicas within a transaction.

You can, however, specify which type of connection to use when you begin the transaction.

```go
// Start a transaction on a read replica
tx := db.Clauses(dbresolver.Read).Begin()
// ... perform read-only operations in transaction
tx.Commit()

// Start a transaction on a write source
txWrite := db.Clauses(dbresolver.Write).Begin()
// ... perform write operations
txWrite.Commit()
```

## Multiple Database Sources

DBResolver also allows you to route queries to different databases based on the model being queried. This is useful for sharding or separating data by domain.

```go
db.Use(dbresolver.Register(dbresolver.Config{
    // ... default configuration
}).Register(dbresolver.Config{
    // Configuration for the User and Profile models
    Replicas: []gorm.Dialector{mysql.Open("users_replica_dsn")},
}, &User{}, &Profile{}))
```
In this example, any queries for `User` or `Profile` will use the `users_replica_dsn`, while other models will use the default configuration.
