---
name: gorm-development
description: Use when developing with GORM in Go. Router to 40+ atomic sub-skills covering CRUD, associations, transactions, security, performance, and plugins.
---

# GORM Development Skills

This skill is the main entry point for all GORM development sub-skills. Each sub-skill is self-contained with its own examples and references.

## Getting Started

| Skill | Description |
|-------|-------------|
| [Connecting to a Database](./connecting_to_the_database/SKILL.md) | Connect to PostgreSQL, MySQL, SQLite, SQL Server |
| [GORM Config](./gorm_config/SKILL.md) | Configure with `gorm.Config` during initialization |
| [Conventions](./conventions/SKILL.md) | Default conventions: primary keys, table names, timestamps |

## Core Concepts

| Skill | Description |
|-------|-------------|
| [Models](./models/SKILL.md) | Define models with `gorm.Model` and field tags |
| [Custom Data Types](./custom_data_types/SKILL.md) | Create types with `Scanner` and `Valuer` interfaces |
| [Serializer](./serializer/SKILL.md) | Use custom serializers for field data |

## CRUD Operations

| Skill | Description |
|-------|-------------|
| [Create](./create/SKILL.md) | Create records (single and batch) |
| [Query](./query/SKILL.md) | Query with conditions, ordering, pagination |
| [Update](./update/SKILL.md) | Update records, columns, selected fields |
| [Delete](./delete/SKILL.md) | Delete records, soft deletes |
| [Raw SQL & SQL Builder](./raw_sql/SKILL.md) | Raw SQL for complex queries |

## Associations (Relationships)

| Skill | Description |
|-------|-------------|
| [Associations](./associations/SKILL.md) | Overview of managing relationships |
| [Belongs To](./belongs_to/SKILL.md) | One-to-one where model belongs to another |
| [Has One](./has_one/SKILL.md) | One-to-one in opposite direction |
| [Has Many](./has_many/SKILL.md) | One-to-many relationships |
| [Many To Many](./many_to_many/SKILL.md) | Many-to-many with join tables |
| [Polymorphism](./polymorphism/SKILL.md) | Polymorphic associations |
| [Preload](./preload/SKILL.md) | Eager loading to avoid N+1 queries |

## Schema & Database Design

| Skill | Description |
|-------|-------------|
| [Indexes](./indexes/SKILL.md) | Single and composite indexes |
| [Constraints](./constraints/SKILL.md) | Check and foreign key constraints |
| [Composite Primary Key](./composite_primary_key/SKILL.md) | Multi-column primary keys |
| [Migration](./migration/SKILL.md) | Auto-migrate database schemas |

## Advanced Topics

| Skill | Description |
|-------|-------------|
| [Transactions](./transactions/SKILL.md) | Ensure data integrity |
| [Hooks](./hooks/SKILL.md) | Intercept CRUD with callbacks |
| [Scopes](./scopes/SKILL.md) | Reusable query scopes |
| [Session](./session/SKILL.md) | Sessions with specific configurations |
| [Method Chaining](./method_chaining/SKILL.md) | Build queries with chained methods |
| [Context](./context/SKILL.md) | Use `context.Context` for cancellation/timeouts |
| [Settings](./settings/SKILL.md) | Pass values between code and hooks |
| [Error Handling](./error_handling/SKILL.md) | Handle errors including `ErrRecordNotFound` |
| [Security](./security/SKILL.md) | Prevent SQL injection, validate input |

## Performance & Monitoring

| Skill | Description |
|-------|-------------|
| [Performance](./performance/SKILL.md) | Tips for improving performance |
| [Logger](./logger/SKILL.md) | Configure GORM's logger |
| [Prometheus](./prometheus/SKILL.md) | Export metrics to Prometheus |
| [Hints](./hints/SKILL.md) | Add optimizer hints to queries |
| [Generic Interface](./generic_interface/SKILL.md) | Access `*sql.DB`, configure connection pools |

## Plugins & Scaling

| Skill | Description |
|-------|-------------|
| [DBResolver](./dbresolver/SKILL.md) | Read/write splitting |
| [Sharding](./sharding/SKILL.md) | Partition data across databases |
| [Write Plugins](./write_plugins/SKILL.md) | Write custom GORM plugins |
| [Write Driver](./write_driver/SKILL.md) | Write custom database drivers |

## Quick Reference

**Common patterns:**
- `db.Create(&user)` - Insert record
- `db.First(&user, 1)` - Find by primary key
- `db.Where("name = ?", name).Find(&users)` - Query with conditions
- `db.Model(&user).Update("name", "newname")` - Update field
- `db.Delete(&user, 1)` - Delete record
- `db.Preload("Orders").Find(&users)` - Eager load associations
- `db.Transaction(func(tx *gorm.DB) error { ... })` - Transaction block
