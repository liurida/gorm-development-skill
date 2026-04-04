# Key Concepts for GORM Logger

This document provides detailed explanations of the logger in GORM.

## Overview

GORM includes a flexible and customizable logger that provides insights into database operations. It can be configured to control log levels, output formats, and slow query detection.

## Default Logger

GORM's default logger prints slow queries and errors. You can customize its behavior during initialization.

```go
import (
    "log"
    "os"
    "time"
    "gorm.io/gorm/logger"
)

newLogger := logger.New(
    log.New(os.Stdout, "\r\n", log.LstdFlags), // io writer
    logger.Config{
        SlowThreshold:             time.Second,
        LogLevel:                  logger.Info,
        IgnoreRecordNotFoundError: true,
        ParameterizedQueries:      false, // Log SQL with params
        Colorful:                  true,
    },
)

db, err := gorm.Open(sqlite.Open("test.db"), &gorm.Config{
    Logger: newLogger,
})
```

### Logger Configuration Options

| Option | Description |
|--------|-------------|
| `SlowThreshold` | Duration to consider a query "slow" |
| `LogLevel` | Log level: `Silent`, `Error`, `Warn`, `Info` |
| `IgnoreRecordNotFoundError` | If `true`, `ErrRecordNotFound` errors are not logged |
| `ParameterizedQueries` | If `true`, logs SQL without inline values |
| `Colorful` | Enable/disable colorized log output |

## Log Levels

- `Silent`: No logs
- `Error`: Only logs errors
- `Warn`: Logs warnings and errors
- `Info`: Logs all SQL queries, warnings, and errors

Change the log level:

```go
// Globally
db, err := gorm.Open(..., &gorm.Config{
    Logger: logger.Default.LogMode(logger.Warn),
})

// Session-level
tx := db.Session(&gorm.Session{Logger: logger.Default.LogMode(logger.Info)})
```

## Debugging

Use `Debug()` to log a single operation at the `Info` level:

```go
db.Debug().Where("name = ?", "jinzhu").First(&User{})
```

## Custom Logger

You can implement a custom logger by satisfying the `logger.Interface`.

### Interface Definition

```go
type Interface interface {
    LogMode(LogLevel) Interface
    Info(context.Context, string, ...interface{})
    Warn(context.Context, string, ...interface{})
    Error(context.Context, string, ...interface{})
    Trace(ctx context.Context, begin time.Time, fc func() (sql string, rowsAffected int64), err error)
}
```

### Trace Method

The `Trace` method is the core of the logger. It is called after every database operation and receives:
- `ctx`: The context of the operation
- `begin`: The start time of the operation
- `fc`: A function that returns the SQL query and rows affected
- `err`: Any error that occurred

**Example Trace Implementation**
```go
func (l *CustomLogger) Trace(ctx context.Context, begin time.Time, fc func() (string, int64), err error) {
    duration := time.Since(begin)
    sql, rows := fc()

    // Log the query
    log.Printf("SQL: %s | Rows: %d | Duration: %v", sql, rows, duration)

    // Log slow queries
    if l.SlowThreshold > 0 && duration > l.SlowThreshold {
        log.Printf("SLOW: %s | Duration: %v", sql, duration)
    }

    // Log errors
    if err != nil && !errors.Is(err, gorm.ErrRecordNotFound) {
        log.Printf("ERROR: %v | SQL: %s", err, sql)
    }
}
```

## Contextual Logging

Since the logger interface methods accept a `context.Context`, you can create loggers that include request-scoped information, such as a request ID:

```go
func (l *ContextualLogger) Trace(ctx context.Context, ...) {
    if requestID := ctx.Value("request_id"); requestID != nil {
        log.Printf("[RequestID: %v] ...", requestID)
    }
}

// Usage
ctx := context.WithValue(context.Background(), "request_id", "some-id")
db.WithContext(ctx).First(&user)
```

## Best Practices

1. **Use `logger.Info` for development**, `logger.Warn` or `logger.Error` for production
2. **Set a reasonable `SlowThreshold`** to identify performance bottlenecks
3. **Set `IgnoreRecordNotFoundError` to `true`** in production to reduce noise
4. **Use contextual logging** to trace queries back to specific requests
5. **Consider structured logging** (e.g., JSON) in your custom logger for better observability

## Common Use Cases

### Logging to a File

```go
file, err := os.OpenFile("gorm.log", os.O_CREATE|os.O_WRONLY|os.O_APPEND, 0666)
if err != nil {
    log.Fatal(err)
}

newLogger := logger.New(
    log.New(file, "\r\n", log.LstdFlags),
    ...
)
```

### Integrating with a Third-Party Logger (e.g., Zap, Logrus)

Create a custom logger that wraps your preferred logging library:

```go
type ZapLogger struct {
    zap *zap.Logger
}

func (l *ZapLogger) Trace(ctx context.Context, begin time.Time, fc func() (string, int64), err error) {
    // ... implement with zap fields
}
```

By customizing the logger, you can integrate GORM's database logging seamlessly into your application's overall logging strategy.
