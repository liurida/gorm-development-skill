
# Key Concepts for GORM Serializer

This document provides key concepts for using custom serializers in GORM.

## Overview

GORM's serializer feature allows you to customize how data is stored in and read from the database. This is particularly useful for fields that are complex types (like structs, slices, or maps) and need to be stored in a single database column (e.g., as JSON, Gob, or another format).

## Built-in Serializers

GORM comes with a few built-in serializers that you can use out of the box:

- **`json`**: Serializes the field to JSON before storing it.
- **`gob`**: Serializes the field using Go's built-in `gob` encoder.
- **`unixtime`**: Converts a `time.Time` object to a Unix timestamp (integer) for storage.

```go
type User struct {
    Roles     []string               `gorm:"serializer:json"`
    JobInfo   Job                    `gorm:"type:bytes;serializer:gob"`
    CreatedAt int64                  `gorm:"serializer:unixtime;type:time"`
}
```

## Custom Serializers

You can create your own custom serializers by implementing the `serializer.SerializerInterface`.

```go
import "gorm.io/gorm/schema"

type SerializerInterface interface {
    Scan(ctx context.Context, field *schema.Field, dst reflect.Value, dbValue interface{}) error
    Value(ctx context.Context, field *schema.Field, dst reflect.Value, fieldValue interface{}) (interface{}, error)
}
```

- **`Scan`**: This method is called when reading data from the database. It takes the raw database value (`dbValue`) and must unmarshal it into the appropriate Go type for the field (`dst`).
- **`Value`**: This method is called when writing data to the database. It takes the Go field value (`fieldValue`) and must marshal it into a format that can be stored in the database (e.g., `[]byte` or `string`).

### Registering a Custom Serializer

Once you have implemented your custom serializer, you need to register it with GORM.

```go
schema.RegisterSerializer("my_serializer", MyCustomSerializer{})
```

Then you can use it in your model structs:

```go
type User struct {
    Preferences MyPreferences `gorm:"serializer:my_serializer"`
}
```

### Field-Level Custom Serializer

Alternatively, you can have a field's type itself implement the `serializer.SerializerInterface`. In this case, you don't need to register it globally. GORM will automatically use the `Scan` and `Value` methods defined on the type.

```go
type EncryptedString string

func (es *EncryptedString) Scan(ctx context.Context, field *schema.Field, dst reflect.Value, dbValue interface{}) error {
    // Decryption logic here
}

func (es EncryptedString) Value(ctx context.Context, field *schema.Field, dst reflect.Value, fieldValue interface{}) (interface{}, error) {
    // Encryption logic here
}

// In your model:
type User struct {
    Password EncryptedString
}
```

This approach is very powerful for creating self-contained, reusable data types with custom serialization logic.
