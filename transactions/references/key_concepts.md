# Key Concepts for GORM Transactions

This document provides key concepts for working with database transactions in GORM.

## Automatic Transactions

GORM performs write operations (create, update, delete) inside a transaction by default to ensure data consistency.

## Manual Transactions

You can control transactions manually for a set of operations.

### Using `Transaction` Method

This is the recommended way to handle transactions.

```go
db.Transaction(func(tx *gorm.DB) error {
  // Perform database operations within the transaction
  if err := tx.Create(&User{Name: "Giraffe"}).Error; err != nil {
    return err // Returning an error will rollback the transaction
  }
  return nil // Returning nil will commit the transaction
})
```

### Manual Control

You can also use `Begin`, `Commit`, and `Rollback` for manual control.

```go
tx := db.Begin() // Start a transaction

// Perform operations
tx.Create(&User{Name: "Lion"})

// tx.Rollback() // Rollback the transaction
tx.Commit() // Commit the transaction
```

## Nested Transactions

GORM supports nested transactions. You can rollback a subset of operations.

## `SavePoint` and `RollbackTo`

These methods allow you to create savepoints and rollback to them.

```go
tx := db.Begin()
tx.Create(&user1)

tx.SavePoint("sp1")
tx.Create(&user2)
tx.RollbackTo("sp1") // Rollbacks the creation of user2

tx.Commit() // Commits the creation of user1
```
