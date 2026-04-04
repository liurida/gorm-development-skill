
# Key Concepts for GORM Error Handling

This document provides key concepts for handling errors in GORM.

## Overview

GORM follows standard Go error handling patterns. Most GORM methods return a `*gorm.DB` object, and you can access any error that occurred via the `Error` field.

## Basic Error Checking

Always check the `Error` field after a chain of GORM methods, especially after a finisher method like `First`, `Find`, `Create`, `Update`, or `Delete`.

```go
var user User
err := db.First(&user, 1).Error
if err != nil {
    // Handle error
}
```

## Record Not Found

A very common error is `gorm.ErrRecordNotFound`, which is returned by `First`, `Last`, and `Take` when no record is found. It's important to handle this case specifically, as it often isn't a "true" error in application logic.

```go
import "errors"

err := db.First(&user, "non_existent_id").Error
if errors.Is(err, gorm.ErrRecordNotFound) {
    fmt.Println("User not found, which is an expected case.")
} else if err != nil {
    // Handle other database errors
}
```

**Note**: The `Find` method does **not** return `ErrRecordNotFound`. It will return an empty slice and a `nil` error if no records are found.

## Dialect-Specific Errors

Databases can return specific error codes for things like unique constraint violations or foreign key violations. You can inspect these errors by type-asserting the error to the specific driver's error type.

### Translated Errors

To simplify this, GORM can translate common database errors into generic GORM errors if you enable the `TranslateError` config.

```go
db, err := gorm.Open(postgres.Open(dsn), &gorm.Config{
    TranslateError: true,
})
```

This provides two useful errors:

- **`gorm.ErrDuplicatedKey`**: Returned when a unique constraint is violated.
- **`gorm.ErrForeignKeyViolated`**: Returned when a foreign key constraint is violated.

```go
err := db.Create(&User{Name: "existing_user"}).Error
if errors.Is(err, gorm.ErrDuplicatedKey) {
    fmt.Println("This user already exists!")
}
```

## Transaction Errors

GORM's `Transaction` method makes error handling simple. If the function passed to `Transaction` returns an error, the transaction is automatically rolled back. If it returns `nil`, the transaction is committed.

```go
db.Transaction(func(tx *gorm.DB) error {
    if err := tx.Create(&User{Name: "test"}).Error; err != nil {
        // The transaction will be rolled back
        return err
    }

    if someConditionFails {
        // The transaction will be rolled back
        return errors.New("something went wrong")
    }

    // The transaction will be committed
    return nil
})
```

By consistently checking for errors and handling specific cases like `ErrRecordNotFound`, you can build reliable and robust database interactions with GORM.
