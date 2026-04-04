---
name: gorm-serializer
description: Use when serializing complex Go types to database storage, implementing custom serialization logic, or using built-in JSON/GOB/UnixTime serializers in GORM.
---

# Serializer

Serializers provide an extensible interface to customize how data is serialized and deserialized between Go types and database storage.

**Reference:** https://gorm.io/docs/serializer.html

## Built-in Serializers

| Serializer | Purpose | Database Storage |
|------------|---------|------------------|
| `json` | JSON encoding | String/text |
| `gob` | Go binary encoding | Bytes |
| `unixtime` | Unix timestamp | Datetime |

## Basic Usage

Use the `serializer` tag to specify serialization:

```go
type User struct {
  Name        []byte                 `gorm:"serializer:json"`
  Roles       Roles                  `gorm:"serializer:json"`
  Contracts   map[string]interface{} `gorm:"serializer:json"`
  JobInfo     Job                    `gorm:"type:bytes;serializer:gob"`
  CreatedTime int64                  `gorm:"serializer:unixtime;type:time"`
}

type Roles []string

type Job struct {
  Title    string
  Location string
  IsIntern bool
}
```

### Example: Create and Query

```go
createdAt := time.Date(2020, 1, 1, 0, 8, 0, 0, time.UTC)
data := User{
  Name:        []byte("jinzhu"),
  Roles:       []string{"admin", "owner"},
  Contracts:   map[string]interface{}{"name": "jinzhu", "age": 10},
  CreatedTime: createdAt.Unix(),
  JobInfo: Job{
    Title:    "Developer",
    Location: "NY",
    IsIntern: false,
  },
}

DB.Create(&data)
// INSERT INTO `users` (`name`,`roles`,`contracts`,`job_info`,`created_time`) VALUES
//   ("\"amluemh1\"","[\"admin\",\"owner\"]","{\"age\":10,\"name\":\"jinzhu\"}",<gob binary>,"2020-01-01 00:08:00")

var result User
DB.First(&result, "id = ?", data.ID)
// Data is automatically deserialized back to Go types
```

## SerializerInterface

Implement this interface to create custom serializers:

```go
import "gorm.io/gorm/schema"

type SerializerInterface interface {
  Scan(ctx context.Context, field *schema.Field, dst reflect.Value, dbValue interface{}) error
  SerializerValuerInterface
}

type SerializerValuerInterface interface {
  Value(ctx context.Context, field *schema.Field, dst reflect.Value, fieldValue interface{}) (interface{}, error)
}
```

### Method Parameters

| Parameter | Description |
|-----------|-------------|
| `ctx` | Request-scoped context values |
| `field` | Field metadata including GORM settings and struct tags |
| `dst` | Current model value |
| `dbValue` (Scan) | Value from database |
| `fieldValue` (Value) | Go value to serialize |

## Default JSONSerializer Implementation

```go
type JSONSerializer struct{}

func (JSONSerializer) Scan(ctx context.Context, field *schema.Field, dst reflect.Value, dbValue interface{}) (err error) {
  fieldValue := reflect.New(field.FieldType)

  if dbValue != nil {
    var bytes []byte
    switch v := dbValue.(type) {
    case []byte:
      bytes = v
    case string:
      bytes = []byte(v)
    default:
      return fmt.Errorf("failed to unmarshal JSONB value: %#v", dbValue)
    }
    err = json.Unmarshal(bytes, fieldValue.Interface())
  }

  field.ReflectValueOf(ctx, dst).Set(fieldValue.Elem())
  return
}

func (JSONSerializer) Value(ctx context.Context, field *schema.Field, dst reflect.Value, fieldValue interface{}) (interface{}, error) {
  return json.Marshal(fieldValue)
}
```

## Registering Custom Serializers

```go
schema.RegisterSerializer("myserializer", MySerializer{})
```

Then use with tag:

```go
type User struct {
  Data []byte `gorm:"serializer:myserializer"`
}
```

## Custom Serializer Type (Field-Level)

Create a type that implements `SerializerInterface` directly:

```go
type EncryptedString string

func (es *EncryptedString) Scan(ctx context.Context, field *schema.Field, dst reflect.Value, dbValue interface{}) (err error) {
  switch value := dbValue.(type) {
  case []byte:
    *es = EncryptedString(bytes.TrimPrefix(value, []byte("hello")))
  case string:
    *es = EncryptedString(strings.TrimPrefix(value, "hello"))
  default:
    return fmt.Errorf("unsupported data %#v", dbValue)
  }
  return nil
}

func (es EncryptedString) Value(ctx context.Context, field *schema.Field, dst reflect.Value, fieldValue interface{}) (interface{}, error) {
  return "hello" + string(es), nil
}

type User struct {
  gorm.Model
  Password EncryptedString
}
```

### Usage

```go
data := User{Password: EncryptedString("pass")}
DB.Create(&data)
// INSERT INTO `users` (`password`) VALUES ("hellopass")

var result User
DB.First(&result, "id = ?", data.ID)
// result.Password equals "pass" (decrypted)

DB.Where(User{Password: EncryptedString("pass")}).Take(&result)
// SELECT * FROM `users` WHERE `users`.`password` = "hellopass"
```

## Serializer vs Custom Data Type

| Feature | Serializer | Custom Data Type |
|---------|------------|------------------|
| Use via tag | Yes (`gorm:"serializer:json"`) | No (field type) |
| Reusable across types | Yes (register once) | Per-type implementation |
| Access to field metadata | Yes (`schema.Field`) | Limited |
| SQL expression support | No | Yes (`GormValue`) |
| Database type control | No | Yes (`GormDBDataType`) |

**Use serializers when:** You need consistent serialization across multiple fields/types.

**Use custom data types when:** You need database-specific behavior or SQL expressions.

## When NOT to Use

- **For simple data types** - Don't use a serializer for `string`, `int`, etc. that map directly to database types.
- **When you need to query the serialized data** - Storing data as JSON or GOB makes it un-queryable at the database level. If you need to filter by fields within the serialized data, it should be a separate column or table.
- **For large objects** - Storing large binary objects (like images or files) in the database via a serializer is often inefficient. Use a dedicated file storage solution (like S3) and store a reference.
- **When portability is key** - A custom GOB serializer is Go-specific. If other applications need to read the data, use a language-agnostic format like JSON.

## Common Mistakes

| Mistake | Fix |
|---------|-----|
| Forgetting pointer receiver on `Scan` | Use `func (es *EncryptedString) Scan(...)` |
| Not handling both `[]byte` and `string` | Database drivers return either type |
| Modifying dst directly | Use `field.ReflectValueOf(ctx, dst).Set(...)` |
| Not returning errors | Always propagate deserialization errors |
