# Key Concepts for Connecting to a Database

This document provides key concepts and examples for connecting to various databases using GORM.

## Supported Databases

GORM officially supports the following databases:
- MySQL
- PostgreSQL
- GaussDB
- SQLite
- SQL Server
- TiDB
- Oracle Database

## DSN (Data Source Name)

The DSN is a string that contains the connection information for the database. The format of the DSN varies depending on the database driver.

### MySQL DSN

`user:pass@tcp(127.0.0.1:3306)/dbname?charset=utf8mb4&parseTime=True&loc=Local`

**Important Parameters:**
- `parseTime=True`: Required to handle `time.Time` correctly.
- `charset=utf8mb4`: For full UTF-8 support.

### PostgreSQL DSN

`host=localhost user=gorm password=gorm dbname=gorm port=9920 sslmode=disable TimeZone=Asia/Shanghai`

## Connection Pool

GORM uses the `database/sql` package to maintain a connection pool. You can configure the connection pool to optimize performance.

```go
sqlDB, err := db.DB()

// SetMaxIdleConns sets the maximum number of connections in the idle connection pool.
sqlDB.SetMaxIdleConns(10)

// SetMaxOpenConns sets the maximum number of open connections to the database.
sqlDB.SetMaxOpenConns(100)

// SetConnMaxLifetime sets the maximum amount of time a connection may be reused.
sqlDB.SetConnMaxLifetime(time.Hour)
```
