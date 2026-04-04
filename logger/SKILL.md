---
name: gorm-logger
description: Use when configuring GORM logging, customizing log output, tracking slow queries, or integrating with application logging frameworks.
---

# Logger

GORM provides a configurable logger for SQL debugging, slow query detection, and error tracking.

**Reference:** https://gorm.io/docs/logger.html

## Quick Reference

| Log Level | Description |
|-----------|-------------|
| `Silent` | No logging |
| `Error` | Errors only |
| `Warn` | Errors and warnings |
| `Info` | All SQL statements (verbose) |

## Basic Configuration

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
    SlowThreshold:             time.Second,   // Slow SQL threshold
    LogLevel:                  logger.Warn,   // Log level
    IgnoreRecordNotFoundError: true,          // Ignore ErrRecordNotFound
    ParameterizedQueries:      true,          // Don't log params (security)
    Colorful:                  true,          // Colorized output
  },
)

// Global configuration
db, err := gorm.Open(sqlite.Open("test.db"), &gorm.Config{
  Logger: newLogger,
})

// Session-specific logger
tx := db.Session(&gorm.Session{Logger: newLogger})
```

## Log Levels

```go
// Set globally during init
db, err := gorm.Open(sqlite.Open("test.db"), &gorm.Config{
  Logger: logger.Default.LogMode(logger.Silent),
})

// Change at runtime
db.Logger = logger.Default.LogMode(logger.Info)
```

## Debug Mode

Temporarily enable verbose logging for specific queries.

```go
// Debug single operation
db.Debug().Where("name = ?", "jinzhu").First(&User{})
// Logs the full SQL with parameters

// Debug a chain of operations
db.Debug().
  Preload("Orders").
  Where("status = ?", "active").
  Find(&users)
```

## Slow Query Detection

```go
newLogger := logger.New(
  log.New(os.Stdout, "\r\n", log.LstdFlags),
  logger.Config{
    SlowThreshold: 200 * time.Millisecond, // Warn if query takes > 200ms
    LogLevel:      logger.Warn,
  },
)
```

Output for slow queries:
```
[SLOW SQL >= 200ms] [rows:1000] SELECT * FROM orders WHERE status = 'pending'
```

## Security: Parameterized Queries

Hide sensitive data from logs.

```go
newLogger := logger.New(
  log.New(os.Stdout, "\r\n", log.LstdFlags),
  logger.Config{
    ParameterizedQueries: true, // Don't include params in SQL log
    LogLevel:             logger.Info,
  },
)

// Logs: SELECT * FROM users WHERE email = ?
// Instead of: SELECT * FROM users WHERE email = 'user@example.com'
```

## Custom Logger Implementation

Implement the `logger.Interface` for integration with your logging framework.

```go
type Interface interface {
  LogMode(LogLevel) Interface
  Info(context.Context, string, ...interface{})
  Warn(context.Context, string, ...interface{})
  Error(context.Context, string, ...interface{})
  Trace(ctx context.Context, begin time.Time, fc func() (sql string, rowsAffected int64), err error)
}
```

### Example: Zap Logger Integration

```go
type ZapLogger struct {
  ZapLogger *zap.Logger
  LogLevel  logger.LogLevel
}

func (l *ZapLogger) LogMode(level logger.LogLevel) logger.Interface {
  newLogger := *l
  newLogger.LogLevel = level
  return &newLogger
}

func (l *ZapLogger) Info(ctx context.Context, msg string, data ...interface{}) {
  if l.LogLevel >= logger.Info {
    l.ZapLogger.Sugar().Infof(msg, data...)
  }
}

func (l *ZapLogger) Warn(ctx context.Context, msg string, data ...interface{}) {
  if l.LogLevel >= logger.Warn {
    l.ZapLogger.Sugar().Warnf(msg, data...)
  }
}

func (l *ZapLogger) Error(ctx context.Context, msg string, data ...interface{}) {
  if l.LogLevel >= logger.Error {
    l.ZapLogger.Sugar().Errorf(msg, data...)
  }
}

func (l *ZapLogger) Trace(ctx context.Context, begin time.Time, fc func() (string, int64), err error) {
  elapsed := time.Since(begin)
  sql, rows := fc()
  
  if err != nil {
    l.ZapLogger.Error("SQL Error",
      zap.Error(err),
      zap.String("sql", sql),
      zap.Duration("elapsed", elapsed),
    )
  } else if elapsed > time.Second {
    l.ZapLogger.Warn("Slow SQL",
      zap.String("sql", sql),
      zap.Int64("rows", rows),
      zap.Duration("elapsed", elapsed),
    )
  } else {
    l.ZapLogger.Debug("SQL",
      zap.String("sql", sql),
      zap.Int64("rows", rows),
      zap.Duration("elapsed", elapsed),
    )
  }
}
```

## Context-Aware Logging

Access request context in logs for tracing.

```go
func (l *CustomLogger) Trace(ctx context.Context, begin time.Time, fc func() (string, int64), err error) {
  // Extract trace ID from context
  traceID := ctx.Value("trace_id")
  sql, rows := fc()
  
  l.logger.Info("SQL",
    zap.Any("trace_id", traceID),
    zap.String("sql", sql),
    zap.Int64("rows", rows),
  )
}
```

## Environment-Based Configuration

```go
func NewLogger() logger.Interface {
  logLevel := logger.Warn
  if os.Getenv("DEBUG") == "true" {
    logLevel = logger.Info
  }
  if os.Getenv("SILENT") == "true" {
    logLevel = logger.Silent
  }
  
  return logger.New(
    log.New(os.Stdout, "\r\n", log.LstdFlags),
    logger.Config{
      SlowThreshold:             time.Second,
      LogLevel:                  logLevel,
      IgnoreRecordNotFoundError: true,
      Colorful:                  os.Getenv("NO_COLOR") != "true",
    },
  )
}
```

## When NOT to Use

- **In performance-critical hot paths with `Info` level** - Logging every SQL query adds overhead. Use `Warn` or `Error` in production, or sample logs.
- **`Debug()` in production code** - `Debug()` is for temporary, interactive debugging and should not be committed.
- **Logging sensitive data** - Ensure `ParameterizedQueries: true` is set to avoid logging passwords, PII, or other sensitive information.
- **When a simple logger is sufficient** - GORM's default logger is often enough for development. Only implement a custom logger when you need integration with a specific logging framework.

## Common Mistakes

| Mistake | Fix |
|---------|-----|
| Logging params in production | Set `ParameterizedQueries: true` |
| `Info` level in production | Use `Warn` or `Error` in production |
| Ignoring slow queries | Set appropriate `SlowThreshold` |
| Not using context for tracing | Pass context through `WithContext()` |
