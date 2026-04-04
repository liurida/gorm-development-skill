---
name: gorm-delete
description: Use when deleting records from the database with GORM. Covers single record delete, delete by primary key, batch delete, soft delete, permanent delete, and hooks.
---

# Delete

Reference: [GORM Delete Documentation](https://gorm.io/docs/delete.html)

## Quick Reference

| Method | Purpose |
|--------|---------|
| `Delete(&record)` | Delete a single record (requires primary key). |
| `Delete(&User{}, 10)` | Delete by primary key. |
| `Where(...).Delete(&Model{})` | Batch delete matching records. |
| `Unscoped()` | Required to find or permanently delete soft-deleted records. |

## Delete a Record

When deleting a record, the object must have a primary key, otherwise GORM will perform a batch delete.

```go
// Delete a record with a primary key
// email's ID is 10
db.Delete(&email)
// DELETE from emails where id = 10;

// Delete with additional conditions
db.Where("name = ?", "jinzhu").Delete(&email)
// DELETE from emails where id = 10 AND name = "jinzhu";
```

## Delete with Primary Key

GORM allows deleting records by primary key directly.

```go
db.Delete(&User{}, 10)
// DELETE FROM users WHERE id = 10;

db.Delete(&User{}, "10")
// DELETE FROM users WHERE id = 10;

db.Delete(&users, []int{1,2,3})
// DELETE FROM users WHERE id IN (1,2,3);
```

## Batch Delete

If the value passed to `Delete` does not have a primary key, GORM performs a batch delete.

```go
db.Where("email LIKE ?", "%jinzhu%").Delete(&Email{})
// DELETE from emails where email LIKE "%jinzhu%";

db.Delete(&Email{}, "email LIKE ?", "%jinzhu%")
// DELETE from emails where email LIKE "%jinzhu%";
```

To efficiently delete a large number of records, pass a slice of primary keys:
```go
var users = []User{{ID: 1}, {ID: 2}, {ID: 3}}
db.Delete(&users)
// DELETE FROM users WHERE id IN (1,2,3);
```

### Block Global Delete

By default, GORM prevents global deletes without a `WHERE` clause and returns `ErrMissingWhereClause`.

```go
db.Delete(&User{}) // Returns gorm.ErrMissingWhereClause

// To perform a global delete, you must provide a condition or use a special session.
// 1. Add a condition
db.Where("1 = 1").Delete(&User{})
// DELETE FROM `users` WHERE 1=1

// 2. Use Raw SQL
db.Exec("DELETE FROM users")

// 3. Use AllowGlobalUpdate session mode
db.Session(&gorm.Session{AllowGlobalUpdate: true}).Delete(&User{})
// DELETE FROM users
```

## Soft Delete

If your model includes `gorm.Model` or a `gorm.DeletedAt` field, it automatically gets soft delete capability.

```go
type User struct {
  gorm.Model
  Name string
}

// Soft delete a user
db.Delete(&user) // user's ID is 111
// UPDATE users SET deleted_at="2013-10-29 10:23" WHERE id = 111;

// Batch soft delete
db.Where("age = ?", 20).Delete(&User{})
// UPDATE users SET deleted_at="2013-10-29 10:23" WHERE age = 20;
```

Soft-deleted records are automatically excluded from queries.

### Find Soft-Deleted Records

Use `Unscoped` to find soft-deleted records.

```go
db.Unscoped().Where("age = 20").Find(&users)
// SELECT * FROM users WHERE age = 20;
```

### Delete Permanently

Use `Unscoped` to perform a hard delete.

```go
db.Unscoped().Delete(&order)
// DELETE FROM orders WHERE id=10;
```

### Custom Soft Delete Flag

GORM supports custom soft delete flags using the `gorm.io/plugin/soft_delete` plugin, allowing you to use unix timestamps or a boolean flag.

**Unix Second:**
```go
import "gorm.io/plugin/soft_delete"

type User struct {
  ID        uint
  Name      string
  DeletedAt soft_delete.DeletedAt
}
// Query: SELECT * FROM users WHERE deleted_at = 0;
// Delete: UPDATE users SET deleted_at = /* current unix second */ WHERE ID = 1;
```

**Boolean Flag (1/0):**
```go
import "gorm.io/plugin/soft_delete"

type User struct {
  ID    uint
  Name  string
  IsDel soft_delete.DeletedAt `gorm:"softDelete:flag"`
}
// Query: SELECT * FROM users WHERE is_del = 0;
// Delete: UPDATE users SET is_del = 1 WHERE ID = 1;
```

## Delete Hooks

GORM allows hooks for delete operations: `BeforeDelete` and `AfterDelete`.

```go
func (u *User) BeforeDelete(tx *gorm.DB) (err error) {
  if u.Role == "admin" {
    return errors.New("admin user not allowed to delete")
  }
  return
}
```

## Returning Deleted Data

For databases that support it (like PostgreSQL), you can return the deleted data.

```go
var users []User
db.Clauses(clause.Returning{}).Where("role = ?", "admin").Delete(&users)
// DELETE FROM `users` WHERE role = "admin" RETURNING *
// users => []User{{...}, ...}
```

## Common Mistakes

| Mistake | Fix |
|---------|-----|
| Deleting without primary key | This triggers a batch delete. Ensure your record has a PK for single deletes. |
| Trying to find soft-deleted records | Use `db.Unscoped()` to include them in queries. |
| Accidental global delete | GORM blocks this by default. Be explicit with `Where` or `AllowGlobalUpdate`. |
| Using soft delete with unique index | Default `*time.Time` can cause issues. Use a unix timestamp or flag from the `soft_delete` plugin. |

## When NOT to Use

- **When you need to archive data** - Soft delete is for marking records as unusable, not for historical archiving. Move records to a separate archive table for that.
- **For GDPR or data privacy requirements** - Soft delete does not permanently remove user data. Use hard deletes (`Unscoped`) for compliance.
- **In high-performance systems where storage is a concern** - Soft-deleted records remain in your tables, which can impact index size and query performance.
- **When you need to re-use a unique key** - A soft-deleted record still holds its unique keys, preventing new records from using them.
- **Cascading deletes in complex graphs** - GORM's soft delete doesn't automatically cascade. You need to handle this manually in hooks or your application logic.

## Related Topics

- [Hooks](https://gorm.io/docs/hooks.html) - Intercepting operations.
- [Transactions](https://gorm.io/docs/transactions.html) - Ensuring data integrity.
- [Query](https://gorm.io/docs/query.html) - Finding records.
