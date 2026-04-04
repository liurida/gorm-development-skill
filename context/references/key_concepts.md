
# Key Concepts for GORM Context

This document provides key concepts for working with `context.Context` in GORM.

## Overview

GORM's context support allows you to manage database operations more effectively by integrating with Go's context package. This enables control over cancellations, timeouts, and passing request-scoped values through your application.

## Single Session Mode

For individual database operations, you can pass a context to control its execution.

### Traditional API

Use the `WithContext` method to apply a context to a GORM database object.

```go
ctx, cancel := context.WithTimeout(context.Background(), 2*time.Second)
defer cancel()

var user User
db.WithContext(ctx).First(&user, 1)
```

## Continuous Session Mode

For a series of related operations, you can create a context-aware GORM session.

```go
ctx := context.Background()
tx := db.WithContext(ctx)

var user User
tx.First(&user, 1)
tx.Model(&user).Update("role", "admin")
```

## Context Timeout

You can use context to enforce timeouts on database queries, preventing long-running operations from blocking your application.

```go
ctx, cancel := context.WithTimeout(context.Background(), 100*time.Millisecond)
defer cancel()

var user User
db.WithContext(ctx).First(&user, 1) // This operation will time out if it takes longer than 100ms
```

## Context in Hooks and Callbacks

The context is available within GORM's hooks and callbacks, allowing for contextual logic during the lifecycle of a database operation.

```go
func (u *User) BeforeCreate(tx *gorm.DB) (err error) {
  ctx := tx.Statement.Context
  // You can now use the context, for example, to get request-scoped values
  if value := ctx.Value("request_id"); value != nil {
    fmt.Printf("Processing request: %s\n", value)
  }
  return
}
```

## Integration with Web Servers

GORM's context support integrates seamlessly with web frameworks like Chi, allowing you to pass request-specific contexts to your database operations.

```go
func ListUsers(w http.ResponseWriter, r *http.Request) {
  // Assuming a middleware has added the DB instance to the request context
  db, ok := r.Context().Value("DB").(*gorm.DB)
  if !ok {
    http.Error(w, "Could not get database connection from context", http.StatusInternalServerError)
    return
  }

  var users []User
  // The DB instance from the context is already context-aware
  if err := db.Find(&users).Error; err != nil {
    http.Error(w, err.Error(), http.StatusInternalServerError)
    return
  }

  // ...
}
```

By leveraging `context.Context`, you can build more robust, scalable, and maintainable applications with GORM.
