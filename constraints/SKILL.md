---
name: gorm-constraints
description: Use when defining database rules to enforce data integrity, such as CHECK constraints for value validation or foreign key constraints to manage relationships between tables with ON UPDATE/ON DELETE actions.
---

# Constraints

GORM allows creating database constraints with tags. Constraints are created during [AutoMigrate or CreateTable](migration.html).

## CHECK Constraint

Create CHECK constraints with the `check` tag. You can specify a constraint name or let GORM generate one:

```go
type UserIndex struct {
    Name  string `gorm:"check:name_checker,name <> 'jinzhu'"` // Named constraint
    Name2 string `gorm:"check:name <> 'jinzhu'"`              // Unnamed constraint
    Name3 string `gorm:"check:,name <> 'jinzhu'"`             // Unnamed (explicit empty name)
}
```

**Syntax:** `check:constraint_name,expression` or `check:expression`

## Index Constraint

For index constraints, see the [Indexes](indexes.html) skill.

## Foreign Key Constraint

GORM creates foreign key constraints for associations automatically. You can configure `OnDelete` and `OnUpdate` actions with the `constraint` tag:

```go
type User struct {
    gorm.Model
    CompanyID  int
    Company    Company    `gorm:"constraint:OnUpdate:CASCADE,OnDelete:SET NULL;"`
    CreditCard CreditCard `gorm:"constraint:OnUpdate:CASCADE,OnDelete:SET NULL;"`
}

type CreditCard struct {
    gorm.Model
    Number string
    UserID uint
}

type Company struct {
    ID   int
    Name string
}
```

### Available Actions

| Action | Description |
|--------|-------------|
| `CASCADE` | Update/delete child rows when parent changes |
| `SET NULL` | Set foreign key to NULL when parent is deleted |
| `SET DEFAULT` | Set foreign key to default value |
| `RESTRICT` | Prevent parent modification if children exist |
| `NO ACTION` | Similar to RESTRICT (database-dependent) |

## Disabling Foreign Key Constraints

To disable foreign key constraint creation during migration:

```go
db, err := gorm.Open(sqlite.Open("gorm.db"), &gorm.Config{
    DisableForeignKeyConstraintWhenMigrating: true,
})
```

**Use cases:**
- Performance during bulk migrations
- Circular dependencies between tables
- Legacy databases without constraint support

## When NOT to Use

- **When application-level logic is preferred** - If you need complex, multi-step validation, handle it in your service layer or hooks instead of database constraints.
- **For temporary or staging data** - Constraints can hinder the process of loading incomplete or intermediate data. Disable them during such operations.
- **In a shared database with multiple applications** - If other applications don't respect the constraints, it can lead to issues. Ensure all connected systems are aware.
- **When database performance is paramount** - Constraints add overhead to write operations. For high-throughput ingest systems, you might defer validation to a separate process.
- **For databases that don't support them well** - Not all databases (especially older versions or some NoSQL-like systems) have robust support for all constraint types.

## Quick Reference

| Tag | Example | Description |
|-----|---------|-------------|
| `check` | `gorm:"check:age > 18"` | Unnamed CHECK constraint |
| `check:name,expr` | `gorm:"check:age_check,age > 18"` | Named CHECK constraint |
| `constraint:OnUpdate` | `gorm:"constraint:OnUpdate:CASCADE"` | Foreign key update action |
| `constraint:OnDelete` | `gorm:"constraint:OnDelete:SET NULL"` | Foreign key delete action |
| Combined | `gorm:"constraint:OnUpdate:CASCADE,OnDelete:SET NULL;"` | Both actions |

## Reference

- Official Docs: https://gorm.io/docs/constraints.html
