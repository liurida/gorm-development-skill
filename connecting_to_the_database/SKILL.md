---
name: gorm-connecting-to-the-database
description: Use when connecting to MySQL, PostgreSQL, SQLite, SQL Server, Oracle, TiDB, GaussDB, or Clickhouse databases with GORM. Covers DSN configuration, driver customization, existing connections, and connection pooling.
---

# Connecting to a Database

GORM officially supports MySQL, PostgreSQL, GaussDB, SQLite, SQL Server, TiDB, Oracle Database, and Clickhouse.

## Quick Reference

| Database   | Driver Package                      | DSN Format |
|------------|-------------------------------------|------------|
| MySQL      | `gorm.io/driver/mysql`              | `user:pass@tcp(host:port)/dbname?params` |
| PostgreSQL | `gorm.io/driver/postgres`           | `host=x user=x password=x dbname=x port=x` |
| SQLite     | `gorm.io/driver/sqlite`             | `path/to/file.db` or `file::memory:` |
| SQL Server | `gorm.io/driver/sqlserver`          | `sqlserver://user:pass@host:port?database=x` |
| Oracle     | `github.com/oracle-samples/gorm-oracle/oracle` | `user="x" password="x" connectString="host:port/sid"` |

## MySQL

```go
import (
  "gorm.io/driver/mysql"
  "gorm.io/gorm"
)

func main() {
  dsn := "user:pass@tcp(127.0.0.1:3306)/dbname?charset=utf8mb4&parseTime=True&loc=Local"
  db, err := gorm.Open(mysql.Open(dsn), &gorm.Config{})
}
```

**Critical DSN Parameters:**
- `parseTime=True` - Required to handle `time.Time` correctly
- `charset=utf8mb4` - Required for full UTF-8 support (emojis, special chars)

### Advanced MySQL Configuration

```go
db, err := gorm.Open(mysql.New(mysql.Config{
  DSN: "gorm:gorm@tcp(127.0.0.1:3306)/gorm?charset=utf8&parseTime=True&loc=Local",
  DefaultStringSize:         256,  // default size for string fields
  DisableDatetimePrecision:  true, // disable for MySQL < 5.6
  DontSupportRenameIndex:    true, // drop & create for MySQL < 5.7, MariaDB
  DontSupportRenameColumn:   true, // `change` for MySQL < 8, MariaDB
  SkipInitializeWithVersion: false,
}), &gorm.Config{})
```

### Custom MySQL Driver

```go
import _ "example.com/my_mysql_driver"

db, err := gorm.Open(mysql.New(mysql.Config{
  DriverName: "my_mysql_driver",
  DSN: "gorm:gorm@tcp(localhost:9910)/gorm?charset=utf8&parseTime=True&loc=Local",
}), &gorm.Config{})
```

### Using Existing Connection

```go
import "database/sql"

sqlDB, err := sql.Open("mysql", "mydb_dsn")
gormDB, err := gorm.Open(mysql.New(mysql.Config{
  Conn: sqlDB,
}), &gorm.Config{})
```

## PostgreSQL

```go
import (
  "gorm.io/driver/postgres"
  "gorm.io/gorm"
)

dsn := "host=localhost user=gorm password=gorm dbname=gorm port=5432 sslmode=disable TimeZone=Asia/Shanghai"
db, err := gorm.Open(postgres.Open(dsn), &gorm.Config{})
```

GORM uses [pgx](https://github.com/jackc/pgx) as the postgres driver, which enables prepared statement cache by default.

### Disable Prepared Statement Cache

```go
db, err := gorm.Open(postgres.New(postgres.Config{
  DSN: "user=gorm password=gorm dbname=gorm port=5432 sslmode=disable",
  PreferSimpleProtocol: true, // disables implicit prepared statement usage
}), &gorm.Config{})
```

### Cloud SQL Proxy (Custom Driver)

```go
import _ "github.com/GoogleCloudPlatform/cloudsql-proxy/proxy/dialers/postgres"

db, err := gorm.Open(postgres.New(postgres.Config{
  DriverName: "cloudsqlpostgres",
  DSN: "host=project:region:instance user=postgres dbname=postgres password=password sslmode=disable",
}), &gorm.Config{})
```

## SQLite

```go
import (
  "gorm.io/driver/sqlite" // CGO-based driver
  // "github.com/glebarez/sqlite" // Pure-Go alternative
  "gorm.io/gorm"
)

db, err := gorm.Open(sqlite.Open("gorm.db"), &gorm.Config{})
```

**In-Memory Database:**
```go
db, err := gorm.Open(sqlite.Open("file::memory:?cache=shared"), &gorm.Config{})
```

## SQL Server

```go
import (
  "gorm.io/driver/sqlserver"
  "gorm.io/gorm"
)

dsn := "sqlserver://gorm:LoremIpsum86@localhost:1433?database=gorm"
db, err := gorm.Open(sqlserver.Open(dsn), &gorm.Config{})
```

## Oracle Database

**Prerequisite:** Install Oracle Instant Client ([installation guide](https://odpi-c.readthedocs.io/en/latest/user_guide/installation.html))

```go
import (
  "github.com/oracle-samples/gorm-oracle/oracle"
  "gorm.io/gorm"
)

// macOS/Windows: include libDir
dataSourceName := `user="scott" password="tiger" connectString="dbhost:1521/orclpdb1" libDir="/path/to/instantclient"`

// Linux: libDir not supported, use ldconfig instead
dataSourceName := `user="scott" password="tiger" connectString="dbhost:1521/orclpdb1"`

db, err := gorm.Open(oracle.Open(dataSourceName), &gorm.Config{})
```

## TiDB

TiDB is MySQL-compatible. Use the MySQL driver with TiDB-specific features:

```go
type Product struct {
  ID    uint `gorm:"primaryKey;default:auto_random()"` // TiDB AUTO_RANDOM
  Code  string
  Price uint
}

db, err := gorm.Open(mysql.Open("root:@tcp(127.0.0.1:4000)/test"), &gorm.Config{})
```

**TiDB Notes:**
- `SAVEPOINT` supported from v6.2.0
- `FOREIGN KEY` supported from v6.6.0

## Clickhouse

```go
import (
  "gorm.io/driver/clickhouse"
  "gorm.io/gorm"
)

dsn := "tcp://localhost:9000?database=gorm&username=gorm&password=gorm&read_timeout=10&write_timeout=20"
db, err := gorm.Open(clickhouse.Open(dsn), &gorm.Config{})

// Set table engine
db.Set("gorm:table_options", "ENGINE=Distributed(cluster, default, hits)").AutoMigrate(&User{})
```

## Connection Pool

GORM uses `database/sql` to maintain a connection pool.

```go
sqlDB, err := db.DB()

sqlDB.SetMaxIdleConns(10)           // max idle connections
sqlDB.SetMaxOpenConns(100)          // max open connections
sqlDB.SetConnMaxLifetime(time.Hour) // max connection lifetime
```

## Common Mistakes

| Mistake | Fix |
|---------|-----|
| Missing `parseTime=True` for MySQL | Add to DSN for proper `time.Time` handling |
| Using `charset=utf8` instead of `utf8mb4` | Use `utf8mb4` for full Unicode support |
| Not handling connection errors | Always check `err` from `gorm.Open()` |
| Using in-memory SQLite without `cache=shared` | Add `?cache=shared` for concurrent access |

## When NOT to Use

- **When a simpler database library is sufficient** - If you're only running simple queries, `database/sql` might be enough
- **For NoSQL databases** - GORM is for SQL-based relational databases only
- **If you need to manage connections manually** - GORM abstracts away connection handling; use `database/sql` for full control
- **When the official driver has features GORM doesn't expose** - If you need a specific driver feature, using the driver directly might be necessary

## References

- [Official GORM Documentation: Connecting to Database](https://gorm.io/docs/connecting_to_the_database.html)
- [MySQL DSN Parameters](https://github.com/go-sql-driver/mysql#parameters)
- [PostgreSQL pgx Documentation](https://github.com/jackc/pgx)
- [Generic Interface (DB access)](https://gorm.io/docs/generic_interface.html)
