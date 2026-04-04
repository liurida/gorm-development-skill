---
name: gorm-error-handling
description: Use when handling errors in GORM, checking for specific error types, translating database errors, or implementing robust error handling patterns with GORM operations.
---

# Error Handling

Effective error handling is critical when working with GORM. GORM offers two API styles with different error handling approaches.

**Reference:** [GORM Error Handling Documentation](https://gorm.io/docs/error_handling.html)

## Quick Reference

| Error | When Returned | Check Pattern |
|-------|---------------|---------------|
| `ErrRecordNotFound` | First/Last/Take finds nothing | `errors.Is(err, gorm.ErrRecordNotFound)` |
| `ErrDuplicatedKey` | Unique constraint violation | `errors.Is(err, gorm.ErrDuplicatedKey)` |
| `ErrForeignKeyViolated` | FK constraint violation | `errors.Is(err, gorm.ErrForeignKeyViolated)` |
| `ErrMissingWhereClause` | Delete/Update without WHERE | `errors.Is(err, gorm.ErrMissingWhereClause)` |
| `ErrInvalidTransaction` | Bad Commit/Rollback | `errors.Is(err, gorm.ErrInvalidTransaction)` |

## Basic Error Handling

### Generics API (Go 1.18+)

Errors are returned directly, following Go's standard pattern:

```go
ctx := context.Background()

// Error handling with direct return values
user, err := gorm.G[User](db).Where("name = ?", "jinzhu").First(ctx)
if err != nil {
    return fmt.Errorf("failed to find user: %w", err)
}

// For operations that don't return a result
err := gorm.G[User](db).Where("id = ?", 1).Delete(ctx)
if err != nil {
    return fmt.Errorf("failed to delete user: %w", err)
}
```

### Traditional API

Errors are stored in the `*gorm.DB` instance's `Error` field:

```go
// Pattern 1: Check Error field directly
if err := db.Where("name = ?", "jinzhu").First(&user).Error; err != nil {
    return fmt.Errorf("failed to find user: %w", err)
}

// Pattern 2: Store result for additional checks
result := db.Where("name = ?", "jinzhu").First(&user)
if result.Error != nil {
    return fmt.Errorf("failed to find user: %w", result.Error)
}
// Can also check result.RowsAffected
```

## ErrRecordNotFound

GORM returns `ErrRecordNotFound` when `First`, `Last`, or `Take` methods find no records.

```go
// Generics API
ctx := context.Background()
user, err := gorm.G[User](db).First(ctx)
if errors.Is(err, gorm.ErrRecordNotFound) {
    // Handle record not found - this is often expected, not an error
    return nil, nil
}
if err != nil {
    return nil, fmt.Errorf("database error: %w", err)
}

// Traditional API
err := db.First(&user, 100).Error
if errors.Is(err, gorm.ErrRecordNotFound) {
    // Handle record not found
}
```

**Note:** `Find` does NOT return `ErrRecordNotFound` for empty results - it returns an empty slice.

## Complete Error Types

All GORM error types from `gorm/errors.go`:

```go
// Record errors
gorm.ErrRecordNotFound         // First/Last/Take found nothing
gorm.ErrEmptySlice             // Empty slice provided

// Transaction errors
gorm.ErrInvalidTransaction     // Invalid Commit/Rollback

// Query/Model errors
gorm.ErrMissingWhereClause     // Delete/Update without WHERE
gorm.ErrPrimaryKeyRequired     // Primary key required
gorm.ErrModelValueRequired     // Model value required
gorm.ErrModelAccessibleFieldsRequired // Accessible fields required
gorm.ErrSubQueryRequired       // Subquery required

// Data errors
gorm.ErrInvalidData            // Unsupported data type
gorm.ErrInvalidField           // Invalid field
gorm.ErrInvalidValue           // Should be pointer to struct/slice
gorm.ErrInvalidValueOfLength   // Association values length mismatch

// Constraint errors (require TranslateError)
gorm.ErrDuplicatedKey          // Unique constraint violation
gorm.ErrForeignKeyViolated     // Foreign key violation
gorm.ErrCheckConstraintViolated // Check constraint violation

// Driver/System errors
gorm.ErrUnsupportedDriver      // Unsupported driver
gorm.ErrUnsupportedRelation    // Unsupported relations
gorm.ErrNotImplemented         // Not implemented
gorm.ErrRegistered             // Already registered
gorm.ErrInvalidDB              // Invalid database
gorm.ErrDryRunModeUnsupported  // Dry run not supported
gorm.ErrPreloadNotAllowed      // Preload with count
```

## Dialect Translated Errors

Enable `TranslateError` for unified error handling across databases:

```go
db, err := gorm.Open(postgres.Open(dsn), &gorm.Config{
    TranslateError: true,
})
```

This translates database-specific errors into common GORM error types:

```go
result := db.Create(&user)

// Check for duplicate key (e.g., unique email)
if errors.Is(result.Error, gorm.ErrDuplicatedKey) {
    return fmt.Errorf("user with this email already exists")
}

// Check for foreign key violation
if errors.Is(result.Error, gorm.ErrForeignKeyViolated) {
    return fmt.Errorf("referenced record does not exist")
}

// Check for check constraint violation
if errors.Is(result.Error, gorm.ErrCheckConstraintViolated) {
    return fmt.Errorf("data validation failed")
}
```

## Database-Specific Error Codes

For fine-grained control, parse database-specific errors:

### MySQL
```go
import "github.com/go-sql-driver/mysql"

result := db.Create(&record)
if result.Error != nil {
    if mysqlErr, ok := result.Error.(*mysql.MySQLError); ok {
        switch mysqlErr.Number {
        case 1062: // Duplicate entry
            return fmt.Errorf("duplicate entry: %w", result.Error)
        case 1452: // Foreign key constraint fails
            return fmt.Errorf("foreign key violation: %w", result.Error)
        case 1406: // Data too long
            return fmt.Errorf("data too long: %w", result.Error)
        default:
            return fmt.Errorf("mysql error %d: %w", mysqlErr.Number, result.Error)
        }
    }
    return result.Error
}
```

### PostgreSQL
```go
import "github.com/lib/pq"

result := db.Create(&record)
if result.Error != nil {
    if pqErr, ok := result.Error.(*pq.Error); ok {
        switch pqErr.Code {
        case "23505": // unique_violation
            return fmt.Errorf("duplicate entry: %w", result.Error)
        case "23503": // foreign_key_violation
            return fmt.Errorf("foreign key violation: %w", result.Error)
        case "23514": // check_violation
            return fmt.Errorf("check constraint violation: %w", result.Error)
        }
    }
    return result.Error
}
```

## Error Handling Patterns

### Pattern 1: Service Layer Error Handling

```go
func (s *UserService) GetByID(ctx context.Context, id uint) (*User, error) {
    var user User
    err := s.db.WithContext(ctx).First(&user, id).Error
    
    if errors.Is(err, gorm.ErrRecordNotFound) {
        return nil, ErrUserNotFound // Custom domain error
    }
    if err != nil {
        return nil, fmt.Errorf("failed to get user %d: %w", id, err)
    }
    return &user, nil
}

func (s *UserService) Create(ctx context.Context, user *User) error {
    err := s.db.WithContext(ctx).Create(user).Error
    
    if errors.Is(err, gorm.ErrDuplicatedKey) {
        return ErrUserAlreadyExists
    }
    if err != nil {
        return fmt.Errorf("failed to create user: %w", err)
    }
    return nil
}
```

### Pattern 2: Checking Affected Rows

```go
result := db.Model(&user).Where("id = ?", id).Update("name", "new_name")
if result.Error != nil {
    return fmt.Errorf("update failed: %w", result.Error)
}
if result.RowsAffected == 0 {
    return ErrUserNotFound // No rows updated
}
```

### Pattern 3: Transaction Error Handling

```go
err := db.Transaction(func(tx *gorm.DB) error {
    if err := tx.Create(&user).Error; err != nil {
        return err // Rollback on any error
    }
    if err := tx.Create(&profile).Error; err != nil {
        return err
    }
    return nil // Commit
})
if err != nil {
    return fmt.Errorf("transaction failed: %w", err)
}
```

## Common Mistakes

| Mistake | Fix |
|---------|-----|
| Ignoring errors | Always check `.Error` or returned `err` |
| Using `Find` and expecting `ErrRecordNotFound` | Use `First`/`Last`/`Take` or check slice length |
| Not enabling `TranslateError` | Enable for portable constraint error handling |
| String comparison for errors | Use `errors.Is()` for proper error chain checking |
| Not wrapping errors | Use `fmt.Errorf("context: %w", err)` |

## When NOT to Use

- **When `ErrRecordNotFound` is a normal condition** - Don't log `ErrRecordNotFound` as an error; it's an expected outcome when checking for existence.
- **Database-specific error codes without a fallback** - Relying only on specific codes makes your code less portable. Use `TranslateError` and `errors.Is` first.
- **`RowsAffected == 0` as a definitive error** - An update that doesn't change data might affect 0 rows and still be successful. Check this only when you expect a change.
- **Ignoring errors in background jobs** - All database operations, even in goroutines or background tasks, need robust error handling and logging.

## Error Handling Checklist

- [ ] All database operations check for errors
- [ ] `ErrRecordNotFound` handled appropriately (often not a real error)
- [ ] Constraint violations return user-friendly messages
- [ ] Errors wrapped with context using `%w`
- [ ] `TranslateError` enabled for cross-database compatibility
- [ ] `RowsAffected` checked for update/delete operations
