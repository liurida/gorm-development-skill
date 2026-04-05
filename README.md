# GORM Development Skills

A structured repository of GORM best practices and patterns optimized for AI agents and LLMs.

## Structure

- `*/SKILL.md` - Individual skill files (one per topic)
- `SKILL.md` - Main router skill with overview and navigation
- `*/examples/` - Code examples for each skill

## Skills

### Getting Started (CRITICAL)

- `connecting_to_the_database/` - Connect to PostgreSQL, MySQL, SQLite, SQL Server
- `gorm_config/` - Configure GORM with `gorm.Config` during initialization
- `conventions/` - Default conventions: primary keys, table names, timestamps

### Core Concepts (HIGH)

- `models/` - Define models with `gorm.Model` and field tags
- `custom_data_types/` - Create types with `Scanner` and `Valuer` interfaces
- `serializer/` - Use custom serializers for field data

### CRUD Operations (HIGH)

- `create/` - Create records (single and batch)
- `query/` - Query with conditions, ordering, pagination
- `advanced_query/` - Complex queries with joins, subqueries
- `update/` - Update records, columns, selected fields
- `delete/` - Delete records, soft deletes
- `raw_sql/` - Raw SQL and SQL Builder for complex queries
- `sql_builder/` - Build SQL queries programmatically

### Associations (HIGH)

- `associations/` - Overview of managing relationships
- `belongs_to/` - One-to-one where model belongs to another
- `has_one/` - One-to-one in opposite direction
- `has_many/` - One-to-many relationships
- `many_to_many/` - Many-to-many with join tables
- `polymorphism/` - Polymorphic associations
- `preload/` - Eager loading to avoid N+1 queries

### Schema & Database Design (MEDIUM)

- `indexes/` - Single and composite indexes
- `constraints/` - Check and foreign key constraints
- `composite_primary_key/` - Multi-column primary keys
- `migration/` - Auto-migrate database schemas

### Advanced Topics (MEDIUM)

- `transactions/` - Ensure data integrity with transactions
- `hooks/` - Intercept CRUD operations with callbacks
- `scopes/` - Reusable query scopes
- `session/` - Sessions with specific configurations
- `method_chaining/` - Build queries with chained methods
- `context/` - Use `context.Context` for cancellation/timeouts
- `settings/` - Pass values between code and hooks
- `error_handling/` - Handle errors including `ErrRecordNotFound`
- `security/` - Prevent SQL injection, validate input

### Performance & Monitoring (MEDIUM)

- `performance/` - Tips for improving query performance
- `logger/` - Configure GORM's logger
- `prometheus/` - Export metrics to Prometheus
- `hints/` - Add optimizer hints to queries
- `generic_interface/` - Access `*sql.DB`, configure connection pools

### Plugins & Scaling (LOW)

- `dbresolver/` - Read/write splitting across replicas
- `sharding/` - Partition data across databases
- `write_plugins/` - Write custom GORM plugins
- `write_driver/` - Write custom database drivers

### Experimental

- `generics_api/` - Generics-based API (experimental)
- `gorm-generics/` - Generic helper utilities

## Using These Skills

Each skill follows this structure:

```markdown
---
name: skill-name
description: When to use this skill
---

# Skill Title

Reference: [Official Docs](https://gorm.io/docs/...)

## Quick Reference
| Method | Purpose |
|--------|---------|

## Basic Usage
[Code examples]

## When NOT to Use
[Anti-patterns and alternatives]

## Common Mistakes
| Mistake | Fix |
|---------|-----|
```

## Impact Levels

- `CRITICAL` - Required for any GORM project
- `HIGH` - Core functionality used in most projects
- `MEDIUM` - Important for production applications
- `LOW` - Advanced features for specific use cases

## Quick Start

```go
import (
  "gorm.io/gorm"
  "gorm.io/driver/postgres"
)

// Connect
dsn := "host=localhost user=gorm password=gorm dbname=gorm port=9920"
db, err := gorm.Open(postgres.Open(dsn), &gorm.Config{})

// Define model
type User struct {
  gorm.Model
  Name string
  Age  int
}

// Auto-migrate
db.AutoMigrate(&User{})

// CRUD
db.Create(&User{Name: "John", Age: 30})
db.First(&user, 1)
db.Model(&user).Update("Age", 31)
db.Delete(&user, 1)
```

## References

- [GORM Official Documentation](https://gorm.io/docs/)
- [GORM GitHub Repository](https://github.com/go-gorm/gorm)
