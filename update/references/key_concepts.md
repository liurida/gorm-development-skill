# Key Concepts for GORM Update

This document provides key concepts for updating records with GORM.

## Save All Fields

`Save` will save all fields when performing an update. If the record does not exist, it will be created.

```go
var user User
db.First(&user)
user.Name = "jinzhu 2"
user.Age = 100
db.Save(&user)
```

## Update Single Column

Use `Model` and `Update` to update a single column. You must provide a condition.

```go
db.Model(&User{}).Where("active = ?", true).Update("name", "hello")
```

## Update Multiple Columns

Use `Updates` to update multiple columns. You can pass a `struct` or a `map`.

- When using a `struct`, only non-zero fields are updated.
- When using a `map`, all key-value pairs are used for the update.

```go
// With struct
db.Model(&user).Updates(User{Name: "hello", Age: 18})

// With map
db.Model(&user).Updates(map[string]interface{}{"name": "hello", "age": 18, "active": false})
```

## Update Selected Fields

Use `Select` to specify which fields to update.

```go
db.Model(&user).Select("name").Updates(map[string]interface{}{"name": "hello", "age": 18, "active": false})
```

Use `Omit` to specify which fields to ignore.

```go
db.Model(&user).Omit("name").Updates(map[string]interface{}{"name": "hello", "age": 18, "active": false})
```
