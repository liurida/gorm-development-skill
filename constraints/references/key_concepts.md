# Key Concepts for GORM Constraints

This document provides detailed explanations of using database constraints in GORM.

## Overview

GORM allows you to define database-level constraints directly in your model structs using tags. These constraints are automatically created when you use GORM's auto-migration feature, helping to ensure data integrity at the database layer.

## Check Constraints

A `CHECK` constraint is used to specify a condition that must be true for each row in a table. It's a powerful way to enforce business rules.

### Defining a Check Constraint

Use the `gorm:"check:<constraint_name>,<expression>"` tag.

```go
type User struct {
    gorm.Model
    // Named CHECK constraint: ensures age is 18 or greater
    Age  int `gorm:"check:age_checker,age >= 18"`

    // Unnamed CHECK constraint: ensures role is one of the allowed values
    Role string `gorm:"check:role IN ('admin', 'user', 'guest')"`
}
```

- The constraint name (e.g., `age_checker`) is optional but recommended for easier management.
- If you try to create or update a record that violates a check constraint, the database will return an error, which GORM will propagate.

## Foreign Key Constraints

GORM automatically creates foreign key constraints for `belongs to` and `has one/many` relationships. You can customize the behavior of these constraints using the `constraint` tag.

### Defining Foreign Key Constraints

Use the `gorm:"constraint:OnUpdate:<action>,OnDelete:<action>;"` tag on the association field.

```go
type User struct {
    gorm.Model
    CompanyID int
    Company   Company `gorm:"constraint:OnUpdate:CASCADE,OnDelete:SET NULL;"`
}

type Company struct {
    ID   int `gorm:"primaryKey"`
    Name string
}
```

### Foreign Key Actions

| Action | `OnUpdate` Behavior | `OnDelete` Behavior |
|---|---|---|
| `CASCADE` | When the referenced key is updated, the foreign key in the child table is also updated. | When the referenced row is deleted, all corresponding rows in the child table are also deleted. |
| `SET NULL` | When the referenced key is updated, the foreign key in the child table is set to `NULL`. (The FK column must be nullable). | When the referenced row is deleted, the foreign key in the child table is set to `NULL`. (The FK column must be nullable). |
| `RESTRICT` | Prevents the update of the referenced key if there are dependent rows. | Prevents the deletion of the referenced row if there are dependent rows. |
| `NO ACTION` | Similar to `RESTRICT`, the check is deferred until the end of the transaction in some databases. | Similar to `RESTRICT`, the check is deferred. |

**Important:** For `SET NULL` to work correctly, the foreign key field in your model should be a pointer type (e.g., `*int`) to be nullable.

## Index Constraints

While indexes are primarily for performance, `UNIQUE` indexes also serve as a constraint. See the `indexes` skill for more details.

```go
// Ensures that no two users can have the same email.
type User struct {
    Email string `gorm:"unique"`
}
```

## Disabling Foreign Key Constraints

In some scenarios (e.g., during complex data migrations or with certain database sharding strategies), you may want to prevent GORM from creating foreign key constraints automatically.

You can disable this feature globally during initialization:

```go
db, err := gorm.Open(sqlite.Open("gorm.db"), &gorm.Config{
    DisableForeignKeyConstraintWhenMigrating: true,
})
```

This setting only affects GORM's auto-migration. It does not disable existing foreign key constraints in the database.

By leveraging these constraint tags, you can enforce data integrity rules directly at the database level, creating a more robust and reliable application.
