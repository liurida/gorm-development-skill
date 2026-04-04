---
name: gorm-create
description: Use when inserting records with GORM. Covers single/batch insert, Select/Omit, CreateInBatches, Map creation, Upsert/OnConflict, associations, and defaults.
---

# Create

Reference: [GORM Create](https://gorm.io/docs/create.html)

## Quick Reference

| Method | Purpose |
|--------|---------|
| `Create(&record)` | Insert single record |
| `Create(&records)` | Batch insert (slice) |
| `CreateInBatches(&records, size)` | Insert in batches |
| `Select(...).Create(...)` | Insert selected fields only |
| `Omit(...).Create(...)` | Exclude fields from insert |
| `Clauses(clause.OnConflict{...}).Create(...)` | Upsert |

## Basic Usage

```go
user := User{Name: "Jinzhu", Age: 18}
result := db.Create(&user) // Must pass pointer

user.ID             // inserted primary key
result.Error        // error if any
result.RowsAffected // count

// Batch insert
users := []*User{{Name: "A"}, {Name: "B"}}
db.Create(users)
```

## Select/Omit Fields

```go
db.Select("Name", "Age").Create(&user)
db.Omit("Name", "Age").Create(&user)
```

## Batch Insert

```go
// Large datasets - use CreateInBatches
db.CreateInBatches(users, 100)

// Global config
db, _ := gorm.Open(sqlite.Open("gorm.db"), &gorm.Config{CreateBatchSize: 1000})

// Per session
db.Session(&gorm.Session{CreateBatchSize: 1000}).Create(&users)
```

## Hooks

```go
func (u *User) BeforeCreate(tx *gorm.DB) error {
  u.UUID = uuid.New()
  return nil
}
// Available: BeforeSave, BeforeCreate, AfterSave, AfterCreate

// Skip hooks
db.Session(&gorm.Session{SkipHooks: true}).Create(&user)
```

## Create From Map

Hooks not invoked, associations not saved, PK not backfilled:
```go
db.Model(&User{}).Create(map[string]interface{}{"Name": "jinzhu", "Age": 18})
```

## SQL Expressions

```go
db.Model(User{}).Create(map[string]interface{}{
  "Name":     "jinzhu",
  "Location": clause.Expr{SQL: "ST_PointFromText(?)", Vars: []interface{}{"POINT(100 100)"}},
})
```

## Associations

Non-zero associations are upserted automatically:
```go
db.Create(&User{Name: "jinzhu", CreditCard: CreditCard{Number: "411111111111"}})

// Skip associations
db.Omit("CreditCard").Create(&user)
db.Omit(clause.Associations).Create(&user) // skip all
```

## Default Values

```go
type User struct {
  Name string `gorm:"default:galeone"`
  Age  *int   `gorm:"default:18"` // Use pointer for zero-value support
}
```

Zero values (`0`, `""`, `false`) won't save for fields with defaults. Use pointer/sql.Null types.

## Upsert (OnConflict)

```go
import "gorm.io/gorm/clause"

// Do nothing on conflict
db.Clauses(clause.OnConflict{DoNothing: true}).Create(&user)

// Update specific columns on conflict
db.Clauses(clause.OnConflict{
  Columns:   []clause.Column{{Name: "id"}},
  DoUpdates: clause.AssignmentColumns([]string{"name", "age"}),
}).Create(&users)

// Update all columns on conflict
db.Clauses(clause.OnConflict{UpdateAll: true}).Create(&users)
```

## When NOT to Use

- **Bulk imports from external sources** - Use `CreateInBatches` or raw SQL for large datasets
- **When you need conditional insert** - Use `FirstOrCreate` instead of `Create` to avoid duplicates
- **Insert with complex subqueries** - Use raw SQL when GORM's Create doesn't support your needs
- **When hooks add too much overhead** - Use `Session{SkipHooks: true}` for bulk operations
- **Creating from maps when you need hooks** - Maps bypass hooks; use structs if hooks are required
- **When you don't need the inserted ID** - Consider `Exec` with raw SQL for fire-and-forget inserts

## Common Mistakes

| Mistake | Fix |
|---------|-----|
| `db.Create(user)` | Use `db.Create(&user)` - pointer required |
| Zero values not saving | Use pointer types or sql.Null |
| Large batch without batching | Use `CreateInBatches` for 1000+ records |
| Expecting hooks with map | Use struct for hooks |
