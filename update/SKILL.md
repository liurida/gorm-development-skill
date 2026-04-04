---
name: gorm-update
description: Use when updating records in the database with GORM. Covers Save (upsert), Update, Updates, updating with struct/map, selecting/omitting fields, batch updates, SQL expressions, and hooks.
---

# Update

Reference: [GORM Update Documentation](https://gorm.io/docs/update.html)

## Quick Reference

| Method | Purpose |
|--------|---------|
| `Save(&record)` | Upsert: Updates all fields if PK exists, otherwise creates. |
| `Update("column", value)` | Update a single column. Requires a `Where` clause. |
| `Updates(struct)` | Update multiple columns based on non-zero fields of a struct. |
| `Updates(map)` | Update multiple columns based on a map, including zero values. |
| `Select(...)` | Specify which fields to update. |
| `Omit(...)` | Specify which fields to exclude from update. |
| `UpdateColumn(s)`| Update without running hooks or tracking update time. |

## Save (Upsert)

`Save` updates all fields if a primary key is provided and exists, otherwise it creates a new record. **It has been removed from the Generics API.**

```go
// Record exists: Update all fields
db.First(&user)
user.Name = "jinzhu 2"
user.Age = 100
db.Save(&user)
// UPDATE users SET name='jinzhu 2', age=100, birthday='2016-01-01', updated_at = ... WHERE id=111;

// No primary key: Create
db.Save(&User{Name: "jinzhu", Age: 100})
// INSERT INTO `users` ...

// Primary key not found: Create
// (if db.Save(&User{ID: 1, ...}) finds no user with ID 1, it will INSERT)
db.Save(&User{ID: 1, Name: "jinzhu", Age: 100})
```

**Warning:** `Save` with `Model()` is undefined behavior. To prevent unintended creations, use `Select("*").Updates()`.

## Update Single Column

Requires a condition to avoid accidental global updates (returns `ErrMissingWhereClause` otherwise).

```go
// Update with conditions
db.Model(&User{}).Where("active = ?", true).Update("name", "hello")
// UPDATE users SET name='hello', updated_at='...' WHERE active=true;

// If model has primary key, it's used as condition
// User's ID is `111`:
db.Model(&user).Update("name", "hello")
// UPDATE users SET name='hello', updated_at='...' WHERE id=111;
```

## Update Multiple Columns

### Updates with Struct

Updates only non-zero fields by default.

```go
// user's ID is 111
db.Model(&user).Updates(User{Name: "hello", Age: 18, Active: false})
// UPDATE users SET name='hello', age=18, updated_at = '...' WHERE id = 111;
// Active: false is a zero-value and is IGNORED!
```

### Updates with Map

Includes all key-values, even zero-values.

```go
// user's ID is 111
db.Model(&user).Updates(map[string]interface{}{"name": "hello", "age": 18, "active": false})
// UPDATE users SET name='hello', age=18, active=false, updated_at='...' WHERE id=111;
// Active: false IS included
```

## Update Selected Fields

Use `Select` and `Omit` for fine-grained control.

```go
// Update only the 'name' field
db.Model(&user).Select("name").Updates(map[string]interface{}{"name": "hello", "age": 18, "active": false})
// UPDATE users SET name='hello' WHERE id=111;

// Exclude the 'name' field from update
db.Model(&user).Omit("name").Updates(map[string]interface{}{"name": "hello", "age": 18, "active": false})
// UPDATE users SET age=18, active=false, updated_at='...' WHERE id=111;

// Select specific fields from a struct, including zero-values
db.Model(&user).Select("Name", "Age").Updates(User{Name: "new_name", Age: 0})
// UPDATE users SET name='new_name', age=0 WHERE id=111;

// Select ALL fields from a struct, including zero-values
db.Model(&user).Select("*").Updates(User{Name: "jinzhu", Role: "admin", Age: 0})
```

## Batch Updates

If `Model` is used without a primary key, GORM performs a batch update.

```go
// Update with struct (non-zero fields)
db.Model(User{}).Where("role = ?", "admin").Updates(User{Name: "hello", Age: 18})
// UPDATE users SET name='hello', age=18 WHERE role = 'admin';

// Update with map
db.Table("users").Where("id IN ?", []int{10, 11}).Updates(map[string]interface{}{"name": "hello", "age": 18})
// UPDATE users SET name='hello', age=18 WHERE id IN (10, 11);
```

### Block Global Updates

By default, GORM prevents global updates without a `WHERE` clause. To override:

```go
db.Model(&User{}).Update("name", "jinzhu") // returns gorm.ErrMissingWhereClause

// 1. Add a condition
db.Model(&User{}).Where("1 = 1").Update("name", "jinzhu")

// 2. Use raw SQL
db.Exec("UPDATE users SET name = ?", "jinzhu")

// 3. Use a session with AllowGlobalUpdate
db.Session(&gorm.Session{AllowGlobalUpdate: true}).Model(&User{}).Update("name", "jinzhu")
```

### Get Rows Affected

```go
result := db.Model(User{}).Where("role = ?", "admin").Updates(User{Name: "hello", Age: 18})

result.RowsAffected // returns updated records count
result.Error        // returns updating error
```

## Update with SQL Expression

Use `gorm.Expr` to update with SQL expressions.

```go
// product's ID is `3`
db.Model(&product).Update("price", gorm.Expr("price * ? + ?", 2, 100))
// UPDATE "products" SET "price" = price * 2 + 100, ... WHERE "id" = 3;

db.Model(&product).Updates(map[string]interface{}{"price": gorm.Expr("price * ? + ?", 2, 100)})
// UPDATE "products" SET "price" = price * 2 + 100, ... WHERE "id" = 3;

db.Model(&product).UpdateColumn("quantity", gorm.Expr("quantity - ?", 1))
// UPDATE "products" SET "quantity" = quantity - 1 WHERE "id" = 3;
```

## Update from SubQuery

```go
db.Model(&user).Update("company_name", db.Model(&Company{}).Select("name").Where("companies.id = users.company_id"))
// UPDATE "users" SET "company_name" = (SELECT name FROM companies WHERE companies.id = users.company_id);

db.Table("users as u").Where("name = ?", "jinzhu").Update("company_name", db.Table("companies as c").Select("name").Where("c.id = u.company_id"))
```

## Skip Hooks and Time Tracking

Use `UpdateColumn` and `UpdateColumns` to skip hooks and `updated_at` tracking.

```go
// Update single column without hooks
db.Model(&user).UpdateColumn("name", "hello")
// UPDATE users SET name='hello' WHERE id = 111;

// Update multiple columns without hooks
db.Model(&user).UpdateColumns(User{Name: "hello", Age: 18})
// UPDATE users SET name='hello', age=18 WHERE id = 111;
```

## Returning Modified Data

For databases that support it (like PostgreSQL), you can return the modified data.

```go
var users []User
db.Model(&users).Clauses(clause.Returning{}).Where("role = ?", "admin").Update("salary", gorm.Expr("salary * ?", 2))
// UPDATE `users` SET `salary`=salary * 2, ... WHERE role = "admin" RETURNING *
// users => []User{{...}, ...}
```

## Common Mistakes

| Mistake | Fix |
|---------|-----|
| `Updates(struct)` not updating zero values | Use `Updates(map[...])` or `Select` to include zero values. |
| Accidental global update | Always use `Where`. Enable `AllowGlobalUpdate` only if necessary. |
| Using `Save` expecting only an update | `Save` is an upsert. Use `Updates` for pure update logic. |
| `BeforeUpdate` hook not firing | `UpdateColumn(s)` bypasses hooks. Use `Update(s)` instead. |

## When NOT to Use

- **When you want to create a record if it doesn't exist** - `Update` will fail if the record is not found. Use `Save` (in the traditional API) for upsert logic, or `FirstOrCreate` followed by an `Update`.
- **`Save` when you only want to update** - `Save` will create a new record if the primary key is zero or not found. This can lead to unexpected new rows. Use `Updates` for pure update logic.
- **When you need to update a specific, hardcoded set of fields** - `Update` and `Updates` are flexible. If you have a dedicated operation (e.g., `ActivateUser`), a raw SQL `EXEC` can sometimes be clearer and more explicit.
- **For large, complex data transformations** - If an update requires complex logic based on other tables, a raw SQL `UPDATE ... FROM ...` query might be more efficient and readable.

## Related Topics

- [Hooks](https://gorm.io/docs/hooks.html) - Intercepting operations.
- [Transactions](https://gorm.io/docs/transactions.html) - Grouping operations.
