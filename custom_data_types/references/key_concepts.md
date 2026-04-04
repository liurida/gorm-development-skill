# Key Concepts for GORM Custom Data Types

This document provides detailed explanations for creating and using custom data types in GORM.

## Overview

GORM is highly extensible, allowing you to define your own data types. This is useful for mapping complex Go types to database columns, such as JSON objects, encrypted strings, or custom value objects.

To create a custom data type, you need to implement a set of interfaces that tell GORM how to handle your type during database operations.

## The Core Interfaces

### 1. `sql.Scanner` and `driver.Valuer`

These are the fundamental interfaces from Go's standard `database/sql` package for custom type mapping.

- **`driver.Valuer`**: Defines how to convert your custom Go type into a primitive type that the database driver can understand (e.g., `string`, `[]byte`, `int64`).
- **`sql.Scanner`**: Defines how to scan a primitive value from the database driver and convert it back into your custom Go type.

```go
import (
    "database/sql/driver"
    "encoding/json"
)

// JSONB is a custom type for JSON data.
type JSONB json.RawMessage

// Value implements the driver.Valuer interface.
func (j JSONB) Value() (driver.Value, error) {
    if len(j) == 0 {
        return nil, nil
    }
    return json.RawMessage(j).MarshalJSON()
}

// Scan implements the sql.Scanner interface.
func (j *JSONB) Scan(value interface{}) error {
    bytes, ok := value.([]byte)
    if !ok {
        return errors.New("Scan source was not []byte")
    }
    *j = JSONB(bytes)
    return nil
}
```

### 2. `GormDataTypeInterface`

This interface tells GORM the general, high-level data type of your custom type. This is useful for plugins and hooks that might need to know about the type.

```go
type GormDataTypeInterface interface {
    GormDataType() string
}

func (j JSONB) GormDataType() string {
    return "json"
}
```

### 3. `GormDBDataTypeInterface`

This interface allows you to specify a different database column type depending on the dialect (e.g., MySQL, PostgreSQL, SQLite).

```go
func (j MyJSON) GormDBDataType(db *gorm.DB, field *schema.Field) string {
    switch db.Dialector.Name() {
    case "mysql":
        return "JSON"
    case "postgres":
        return "JSONB"
    default:
        return "TEXT"
    }
}
```

This is more powerful than `GormDataType` for defining schema during migrations.

### 4. `GormValuerInterface`

This interface provides a way to use SQL expressions when creating or updating records with your custom type. It is particularly useful for types that should be handled by database functions (e.g., spatial data types).

```go
// GormValue returns a clause.Expr which GORM uses to build the SQL.
func (loc Location) GormValue(ctx context.Context, db *gorm.DB) gorm.Clause {
    return gorm.Expr("ST_PointFromText(?)", fmt.Sprintf("POINT(%d %d)", loc.X, loc.Y))
}

// Usage:
db.Create(&User{Location: Location{X: 100, Y: 200}})
// INSERT INTO `users` (`location`) VALUES (ST_PointFromText('POINT(100 200)'))
```

## Example: Encrypted String

A common use case is encrypting data before writing it to the database and decrypting it after reading.

```go
type EncryptedString string

// Value encrypts the string before saving.
func (s EncryptedString) Value() (driver.Value, error) {
    return encrypt(string(s))
}

// Scan decrypts the string after reading.
func (s *EncryptedString) Scan(value interface{}) error {
    encrypted, ok := value.([]byte)
    if !ok {
        return errors.New("invalid type for encrypted string")
    }
    decrypted, err := decrypt(encrypted)
    if err != nil {
        return err
    }
    *s = EncryptedString(decrypted)
    return nil
}
```

By implementing these interfaces, you can seamlessly integrate your own complex data types into GORM models, making your code more expressive and type-safe.
