---
name: gorm-context
description: Use when managing request timeouts, implementing cancellation, passing tracing information through database operations, or integrating GORM with web middleware.
---

# Context

GORM supports `context.Context` for timeout control, cancellation, and request tracing across database operations.

**Reference:** https://gorm.io/docs/context.html

## Quick Reference

| Method | Purpose |
|--------|---------|
| `db.WithContext(ctx)` | Attach context to operations |
| `gorm.G[T](db).Find(ctx)` | Generics API with context |
| `tx.Statement.Context` | Access context in hooks |

## Single Session Mode

Execute individual operations with context control.

### Traditional API

```go
ctx, cancel := context.WithTimeout(context.Background(), 2*time.Second)
defer cancel()

// Query with timeout
db.WithContext(ctx).Find(&users)

// If context expires, returns context.DeadlineExceeded error
```

### Generics API

```go
ctx, cancel := context.WithTimeout(context.Background(), 2*time.Second)
defer cancel()

// Context as first parameter
users, err := gorm.G[User](db).Find(ctx)
```

## Continuous Session Mode

Maintain context across multiple related operations.

```go
ctx, cancel := context.WithTimeout(context.Background(), 5*time.Second)
defer cancel()

// Create session with context
tx := db.WithContext(ctx)

// All operations share the same context
tx.First(&user, 1)           // Uses ctx
tx.Model(&user).Update("role", "admin")  // Uses ctx
tx.Find(&orders, "user_id = ?", 1)       // Uses ctx
```

## Timeout Management

Prevent long-running queries from blocking.

```go
func GetUser(id uint) (*User, error) {
  ctx, cancel := context.WithTimeout(context.Background(), 2*time.Second)
  defer cancel()
  
  var user User
  err := db.WithContext(ctx).First(&user, id).Error
  if err != nil {
    if errors.Is(err, context.DeadlineExceeded) {
      return nil, fmt.Errorf("query timed out after 2s")
    }
    return nil, err
  }
  return &user, nil
}
```

## Cancellation

Support graceful shutdown and request cancellation.

```go
func ProcessOrders(ctx context.Context) error {
  // Long-running batch operation
  return db.WithContext(ctx).FindInBatches(&orders, 100, func(tx *gorm.DB, batch int) error {
    select {
    case <-ctx.Done():
      return ctx.Err() // Stop processing if cancelled
    default:
      // Process batch
      for _, order := range orders {
        processOrder(order)
      }
      return nil
    }
  }).Error
}
```

## Context in Hooks

Access context within GORM lifecycle hooks.

```go
func (u *User) BeforeCreate(tx *gorm.DB) error {
  ctx := tx.Statement.Context
  
  // Check for cancellation
  select {
  case <-ctx.Done():
    return ctx.Err()
  default:
  }
  
  // Access context values (e.g., request ID for audit)
  if requestID, ok := ctx.Value("request_id").(string); ok {
    u.CreatedByRequest = requestID
  }
  
  return nil
}

func (u *User) AfterCreate(tx *gorm.DB) error {
  ctx := tx.Statement.Context
  
  // Use context for external service calls
  return notifyService.UserCreated(ctx, u.ID)
}
```

## Web Framework Integration

### Chi Middleware

```go
func SetDBMiddleware(next http.Handler) http.Handler {
  return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
    // Create timeout context from request
    timeoutCtx, cancel := context.WithTimeout(r.Context(), 5*time.Second)
    defer cancel()
    
    // Attach DB with context to request
    ctx := context.WithValue(r.Context(), "DB", db.WithContext(timeoutCtx))
    next.ServeHTTP(w, r.WithContext(ctx))
  })
}

// Router setup
r := chi.NewRouter()
r.Use(SetDBMiddleware)

// Handler
r.Get("/users/{id}", func(w http.ResponseWriter, r *http.Request) {
  db := r.Context().Value("DB").(*gorm.DB)
  // db already has request context with timeout
  db.First(&user, chi.URLParam(r, "id"))
})
```

### Gin Middleware

```go
func DBMiddleware() gin.HandlerFunc {
  return func(c *gin.Context) {
    // Use Gin's request context (supports cancellation)
    ctx, cancel := context.WithTimeout(c.Request.Context(), 5*time.Second)
    defer cancel()
    
    c.Set("db", db.WithContext(ctx))
    c.Next()
  }
}

// Usage in handler
func GetUser(c *gin.Context) {
  db := c.MustGet("db").(*gorm.DB)
  var user User
  if err := db.First(&user, c.Param("id")).Error; err != nil {
    c.JSON(500, gin.H{"error": err.Error()})
    return
  }
  c.JSON(200, user)
}
```

## Request Tracing

Pass trace information through database operations.

```go
type contextKey string
const TraceIDKey contextKey = "trace_id"

// Middleware adds trace ID
func TracingMiddleware(next http.Handler) http.Handler {
  return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
    traceID := r.Header.Get("X-Trace-ID")
    if traceID == "" {
      traceID = uuid.New().String()
    }
    
    ctx := context.WithValue(r.Context(), TraceIDKey, traceID)
    next.ServeHTTP(w, r.WithContext(ctx))
  })
}

// Custom logger uses trace ID
func (l *TracingLogger) Trace(ctx context.Context, begin time.Time, fc func() (string, int64), err error) {
  traceID, _ := ctx.Value(TraceIDKey).(string)
  sql, rows := fc()
  
  l.logger.Info("SQL",
    "trace_id", traceID,
    "sql", sql,
    "rows", rows,
    "duration", time.Since(begin),
  )
}
```

## Goroutine Safety

`WithContext` is goroutine-safe.

```go
// Safe: each goroutine gets its own session
var wg sync.WaitGroup
for i := 0; i < 10; i++ {
  wg.Add(1)
  go func(id int) {
    defer wg.Done()
    ctx, cancel := context.WithTimeout(context.Background(), time.Second)
    defer cancel()
    
    db.WithContext(ctx).First(&user, id)
  }(i)
}
wg.Wait()
```

## When NOT to Use

- **Background jobs without deadlines** - Long-running batch jobs may not need timeouts; use `context.Background()`
- **Database migrations** - Schema changes should complete regardless of request timeouts
- **When parent context is already configured** - Don't add redundant `WithContext` calls; the context propagates
- **Extremely short timeouts** - Timeouts under 100ms can cause false failures under normal load
- **When you don't handle cancellation** - Adding context without checking `ctx.Err()` doesn't help

## Common Mistakes

| Mistake | Fix |
|---------|-----|
| Not calling `cancel()` | Always `defer cancel()` after `WithTimeout` |
| Ignoring context errors | Check for `context.DeadlineExceeded` |
| Timeout too short for batches | Account for total operation time |
| No timeout in production | Always set reasonable timeouts |
| Sharing session across goroutines | Create new session per goroutine |
