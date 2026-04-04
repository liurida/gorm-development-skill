# Key Concepts for GORM Delete

This document provides key concepts for deleting records with GORM.

## Delete a Record

Use the `Delete` method to delete a record. The record must have a primary key.

```go
var email Email
db.First(&email, 10)
db.Delete(&email)
```

## Delete with Primary Key

You can also delete a record by its primary key directly.

```go
db.Delete(&User{}, 10)
```

## Batch Delete

If the value passed to `Delete` does not have a primary key, GORM performs a batch delete.

```go
db.Where("email LIKE ?", "%jinzhu%").Delete(&Email{})
```

## Soft Delete

If a model has a `gorm.DeletedAt` field, it will be soft-deleted. The record will not be removed from the database, but a timestamp will be set in the `deleted_at` column.

```go
db.Delete(&user)
```

- **Finding Soft-Deleted Records**: Use `Unscoped` to find soft-deleted records.
  ```go
  db.Unscoped().Where("age = 20").Find(&users)
  ```
- **Permanent Deletion**: Use `Unscoped` to permanently delete a record.
  ```go
  db.Unscoped().Delete(&order)
  ```
