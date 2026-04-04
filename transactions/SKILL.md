---
name: gorm-transactions
description: Use when performing a set of database operations that must be treated as a single atomic unit, ensuring that all changes are saved or none are.
---

# Transactions

## Basic Transaction

To perform a set of operations within a transaction, you can use the `Transaction` method.

```go
db.Transaction(func(tx *gorm.DB) error {
  // do some database operations in the transaction (use 'tx' from this point, not 'db')
  if err := tx.Create(&Animal{Name: "Giraffe"}).Error; err != nil {
    // return any error will rollback
    return err
  }

  if err := tx.Create(&Animal{Name: "Lion"}).Error; err != nil {
    return err
  }

  // return nil will commit the whole transaction
  return nil
})
```

## Nested Transactions

GORM supports nested transactions.

```go
db.Transaction(func(tx *gorm.DB) error {
  tx.Create(&user1)

  tx.Transaction(func(tx2 *gorm.DB) error {
    tx2.Create(&user2)
    return errors.New("rollback user2") // Rollback user2
  })

  tx.Transaction(func(tx3 *gorm.DB) error {
    tx3.Create(&user3)
    return nil
  })

  return nil
})
```

## Manual Transaction

You can also control the transaction manually.

```go
// begin a transaction
tx := db.Begin()

// do some database operations in the transaction
tx.Create(...)

// rollback the transaction in case of error
tx.Rollback()

// Or commit the transaction
tx.Commit()
```

## When NOT to Use

- **Single independent operations** - Don't wrap single Create/Update/Delete in a transaction; GORM already handles them atomically
- **Read-only queries** - Transactions add overhead with no benefit for SELECT operations
- **Long-running operations** - Transactions hold locks; avoid wrapping external API calls or lengthy processing
- **When `SkipDefaultTransaction` is enabled** - Use manual transactions only when you explicitly need atomicity
- **Cross-database operations** - GORM transactions don't span multiple database connections
