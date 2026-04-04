
# Key Concepts for Writing GORM Drivers

This document provides key concepts for creating custom database drivers for GORM.

## Overview

While GORM has built-in support for many popular databases, you might need to integrate with a database that isn't supported out of the box. You can do this by creating a custom driver that implements GORM's `gorm.Dialector` interface.

## The Dialector Interface

The `Dialector` interface is the heart of a GORM driver. It defines how GORM communicates with a specific database dialect.

```go
type Dialector interface {
    Name() string
    Initialize(*DB) error
    Migrator(db *DB) Migrator
    DataTypeOf(*schema.Field) string
    DefaultValueOf(*schema.Field) clause.Expression
    BindVarTo(writer clause.Writer, stmt *Statement, v interface{})
    QuoteTo(clause.Writer, string)
    Explain(sql string, vars ...interface{}) string
}
```

### Key Methods

- **`Name()`**: Returns the name of the dialect (e.g., `"mysql"`, `"postgres"`).
- **`Initialize(*DB)`**: Called when the GORM DB object is initialized. This is where you can set up dialect-specific configurations, like custom clause builders.
- **`DataTypeOf(*schema.Field)`**: Maps a Go type from a model field to the corresponding database column type (e.g., `string` -> `VARCHAR(255)`).
- **`BindVarTo(clause.Writer, *Statement, interface{})`**: Defines how bind variables (like `?` or `$1`) are handled in SQL queries.
- **`QuoteTo(clause.Writer, string)`**: Defines how to quote identifiers (like table and column names).
- **`Explain(string, ...interface{})`**: Returns a human-readable version of a SQL query with its variables interpolated.

## Nested Transaction Support

If your database supports savepoints, your dialect can implement the `SavePointerDialectorInterface` to enable nested transaction support in GORM.

```go
type SavePointerDialectorInterface interface {
    SavePoint(tx *DB, name string) error
    RollbackTo(tx *DB, name string) error
}
```

## Custom Clause Builders

Different databases have different SQL syntax. You can define custom clause builders to handle these differences. For example, if your database uses a non-standard `LIMIT` clause, you can provide a custom builder for it.

1.  **Define the builder function**:

    ```go
    func MyCustomLimitBuilder(c clause.Clause, builder clause.Builder) {
        // Custom logic to build the LIMIT clause
    }
    ```

2.  **Register the builder** in your dialect's `Initialize` method:

    ```go
    func (d *MyDialector) Initialize(db *gorm.DB) error {
        db.ClauseBuilders["LIMIT"] = MyCustomLimitBuilder
        return nil
    }
    ```

## Creating a Custom Driver

To create a full-fledged driver, you typically need:

1.  A struct that represents your dialect (e.g., `MyDialector`).
2.  An `Open(dsn string) gorm.Dialector` function that users will call to connect.
3.  Implementation of all the methods in the `gorm.Dialector` interface.

Creating a custom driver is an advanced task. It's often helpful to look at GORM's existing drivers (like the [MySQL driver](https://github.com/go-gorm/mysql)) as a reference.
