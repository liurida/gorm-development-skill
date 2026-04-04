# Key Concepts for GORM Configuration (`gorm.Config`)

This document provides detailed explanations of the `gorm.Config` struct used for initializing GORM.

## Overview

The `gorm.Config` struct is passed during GORM initialization (`gorm.Open`) to control a wide range of behaviors, from performance settings to naming conventions.

```go
import "gorm.io/gorm"

db, err := gorm.Open(sqlite.Open("test.db"), &gorm.Config{
    // ... configuration options here ...
})
```

## Key Configuration Fields

### Performance Tuning

- **`SkipDefaultTransaction` (bool)**
  - If `true`, GORM will not wrap single create, update, or delete operations in a transaction. This can improve performance for write-heavy workloads but sacrifices the default data consistency guarantee for single operations.

- **`PrepareStmt` (bool)**
  - If `true`, GORM will cache prepared statements for every query, which can significantly speed up repeated queries by avoiding the overhead of SQL compilation.

### Naming Conventions

- **`NamingStrategy` (schema.Namer)**
  - Allows you to define custom rules for how GORM maps struct and field names to database table and column names.

  ```go
  import "gorm.io/gorm/schema"

  &gorm.Config{
      NamingStrategy: schema.NamingStrategy{
          TablePrefix:   "tbl_",   // Add a prefix to all table names
          SingularTable: true,     // Use singular table names (e.g., "user" instead of "users")
          NameReplacer:  strings.NewReplacer("CID", "Cid"), // Replace names before conversion
      },
  }
  ```

### Logging

- **`Logger` (logger.Interface)**
  - Allows you to provide a custom logger. You can control log levels, output formats, and slow query thresholds. See the `logger` skill for more details.

### Time and Functions

- **`NowFunc` (func() time.Time)**
  - Overrides the default function for getting the current time. This is useful for ensuring consistent timezones (e.g., always using UTC) for `CreatedAt` and `UpdatedAt` fields.

  ```go
  &gorm.Config{
      NowFunc: func() time.Time {
          return time.Now().UTC()
      },
  }
  ```

### Behavior Control

- **`DryRun` (bool)**
  - If `true`, GORM will generate the SQL for an operation but will not execute it. This is useful for testing, debugging, or preparing SQL statements.

- **`AllowGlobalUpdate` (bool)**
  - If `false` (the default), GORM will return an error if you attempt to perform an update or delete operation without a `WHERE` clause, preventing accidental modification of all rows in a table.

- **`DisableNestedTransaction` (bool)**
  - If `true`, GORM will not use database `SAVEPOINT`s to simulate nested transactions.

- **`DisableAutomaticPing` (bool)**
  - If `true`, GORM will not automatically ping the database upon initialization to check for connectivity.

### Migration Settings

- **`DisableForeignKeyConstraintWhenMigrating` (bool)**
  - If `true`, GORM will not create foreign key constraints in the database when running `AutoMigrate`. This can be useful for certain database architectures or sharding strategies.

## Example of a Comprehensive Configuration

```go
db, err := gorm.Open(sqlite.Open("test.db"), &gorm.Config{
    SkipDefaultTransaction: true,
    PrepareStmt:            true,
    NamingStrategy: schema.NamingStrategy{
        TablePrefix:   "t_",
        SingularTable: true,
    },
    Logger: logger.Default.LogMode(logger.Info),
    NowFunc: func() time.Time {
        return time.Now().UTC()
    },
    AllowGlobalUpdate: false,
})
```

By carefully configuring these options, you can tailor GORM's behavior to meet the specific needs of your application, balancing performance, safety, and convention.
