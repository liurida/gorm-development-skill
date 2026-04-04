---
name: gorm-associations
description: Use when managing GORM associations including auto create/update, association mode operations (find, append, replace, delete, clear, count), batch handling, and association tags.
---

# Associations

GORM automates the saving of associations and their references when creating or updating records, using an upsert technique that primarily updates foreign key references for existing associations.

## Auto Create/Update

When you create a new record, GORM automatically saves its associated data, inserting data into related tables and managing foreign key references.

```go
user := User{
  Name:            "jinzhu",
  BillingAddress:  Address{Address1: "Billing Address - Address 1"},
  ShippingAddress: Address{Address1: "Shipping Address - Address 1"},
  Emails:          []Email{
    {Email: "jinzhu@example.com"},
    {Email: "jinzhu-2@example.com"},
  },
  Languages:       []Language{
    {Name: "ZH"},
    {Name: "EN"},
  },
}

db.Create(&user)
// BEGIN TRANSACTION;
// INSERT INTO "addresses" (address1) VALUES ("Billing Address - Address 1"), ("Shipping Address - Address 1") ON DUPLICATE KEY DO NOTHING;
// INSERT INTO "users" (name,billing_address_id,shipping_address_id) VALUES ("jinzhu", 1, 2);
// INSERT INTO "emails" (user_id,email) VALUES (111, "jinzhu@example.com"), (111, "jinzhu-2@example.com") ON DUPLICATE KEY DO NOTHING;
// INSERT INTO "languages" ("name") VALUES ('ZH'), ('EN') ON DUPLICATE KEY DO NOTHING;
// INSERT INTO "user_languages" ("user_id","language_id") VALUES (111, 1), (111, 2) ON DUPLICATE KEY DO NOTHING;
// COMMIT;
```

### Updating with FullSaveAssociations

For full updates of associated data (not just foreign key references):

```go
db.Session(&gorm.Session{FullSaveAssociations: true}).Updates(&user)
```

## Skip Auto Create/Update

### Using Select

Include only specific fields:

```go
db.Select("Name").Create(&user)
// SQL: INSERT INTO "users" (name) VALUES ("jinzhu");
```

### Using Omit

Exclude fields or associations:

```go
// Skip creating BillingAddress
db.Omit("BillingAddress").Create(&user)

// Skip all associations
db.Omit(clause.Associations).Create(&user)
```

**For many-to-many associations**, skip upserting with `.*`:

```go
// Skip upserting Languages associations
db.Omit("Languages.*").Create(&user)

// Skip creating both association and references
db.Omit("Languages").Create(&user)
```

### Select/Omit Association Fields

Target specific fields within associations:

```go
// Include only specific fields of BillingAddress
db.Select("BillingAddress.Address1", "BillingAddress.Address2").Create(&user)

// Exclude specific fields of BillingAddress
db.Omit("BillingAddress.Address2", "BillingAddress.CreatedAt").Create(&user)
```

## Delete Associations

Delete associated relationships when deleting a primary record:

```go
// Delete user's account
db.Select("Account").Delete(&user)

// Delete Orders and CreditCards
db.Select("Orders", "CreditCards").Delete(&user)

// Delete all has one, has many, and many2many associations
db.Select(clause.Associations).Delete(&user)
```

**Important**: Associations are deleted only if the primary key is not zero:

```go
// WRONG - won't delete accounts
db.Select("Account").Where("name = ?", "jinzhu").Delete(&User{})

// CORRECT - specify ID
db.Select("Account").Where("name = ?", "jinzhu").Delete(&User{ID: 1})
```

## Association Mode

Start association mode by specifying source model and relationship field:

```go
var user User
db.Model(&user).Association("Languages")

// Check for errors
error := db.Model(&user).Association("Languages").Error
```

### Finding Associations

```go
// Simple find
db.Model(&user).Association("Languages").Find(&languages)

// Find with conditions
codes := []string{"zh-CN", "en-US", "ja-JP"}
db.Model(&user).Where("code IN ?", codes).Association("Languages").Find(&languages)
```

### Appending Associations

Adds new associations for many-to-many/has-many, or replaces for has-one/belongs-to:

```go
db.Model(&user).Association("Languages").Append([]Language{languageZH, languageEN})
db.Model(&user).Association("Languages").Append(&Language{Name: "DE"})
db.Model(&user).Association("CreditCard").Append(&CreditCard{Number: "411111111111"})
```

### Replacing Associations

```go
db.Model(&user).Association("Languages").Replace([]Language{languageZH, languageEN})
db.Model(&user).Association("Languages").Replace(Language{Name: "DE"}, languageEN)
```

### Deleting Associations

Removes relationship between source and arguments (only deletes reference):

```go
db.Model(&user).Association("Languages").Delete([]Language{languageZH, languageEN})
db.Model(&user).Association("Languages").Delete(languageZH, languageEN)
```

### Clearing Associations

```go
db.Model(&user).Association("Languages").Clear()
```

### Counting Associations

```go
db.Model(&user).Association("Languages").Count()

// With conditions
codes := []string{"zh-CN", "en-US", "ja-JP"}
db.Model(&user).Where("code IN ?", codes).Association("Languages").Count()
```

## Batch Data Handling

Handle relationships for multiple records:

```go
// Find associations for multiple users
db.Model(&users).Association("Role").Find(&roles)

// Delete across multiple records
db.Model(&users).Association("Team").Delete(&userA)

// Count for batch
db.Model(&users).Association("Team").Count()

// Append/Replace - argument lengths must match data
var users = []User{user1, user2, user3}
db.Model(&users).Association("Team").Append(&userA, &userB, &[]User{userA, userB, userC})
db.Model(&users).Association("Team").Replace(&userA, &userB, &[]User{userA, userB, userC})
```

## Delete Association Record

`Replace`, `Delete`, and `Clear` methods update foreign key to null but don't delete actual records.

### Using Unscoped

**Soft Delete** (sets deleted_at):
```go
db.Model(&user).Association("Languages").Unscoped().Clear()
```

**Permanent Delete**:
```go
db.Unscoped().Model(&user).Association("Languages").Unscoped().Clear()
```

## Association Tags

| Tag | Description |
|-----|-------------|
| `foreignKey` | Column name of current model used as foreign key in join table |
| `references` | Column name in reference table that foreign key maps to |
| `polymorphic` | Defines polymorphic type (typically model name) |
| `polymorphicValue` | Sets polymorphic value (usually table name) |
| `many2many` | Names the join table for many-to-many relationships |
| `joinForeignKey` | Foreign key column in join table mapping to current model |
| `joinReferences` | Foreign key column in join table linking to reference model |
| `constraint` | Specifies relational constraints like `OnUpdate`, `OnDelete` |

## When NOT to Use

- **Bulk data imports** - Auto-save associations adds overhead; use `Omit(clause.Associations)` or raw inserts
- **Complex association graphs** - Deeply nested creates can be hard to debug; build associations explicitly
- **When you need fine-grained control** - Auto-upsert may not match your exact requirements; use explicit operations
- **Performance-critical paths** - Association handling adds queries; consider denormalization for hot paths
- **Many-to-many with extra join table fields** - Use explicit join table model instead of GORM's automatic handling
- **When association order matters** - GORM doesn't guarantee association insert order; handle explicitly if needed

## References

- [GORM Official Documentation: Associations](https://gorm.io/docs/associations.html)
