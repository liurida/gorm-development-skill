# Key Concepts for GORM Associations

This document provides key concepts for working with associations in GORM.

## Auto Create/Update

GORM automatically saves associations when creating or updating a record.

```go
user := User{
    Name:            "jinzhu",
    BillingAddress:  Address{Address1: "Billing Address - Address 1"},
    ShippingAddress: Address{Address1: "Shipping Address - Address 1"},
}
db.Create(&user)
```

## Skipping Auto Create/Update

You can use `Omit` to skip auto-saving associations.

```go
db.Omit("BillingAddress").Create(&user)
db.Omit(clause.Associations).Create(&user) // Skips all associations
```

## Association Mode

Association Mode provides helper methods for managing relationships.

```go
// Find associations
db.Model(&user).Association("Languages").Find(&languages)

// Append associations
db.Model(&user).Association("Languages").Append(&languageEN)

// Replace associations
db.Model(&user).Association("Languages").Replace(&languageEN)

// Delete associations
db.Model(&user).Association("Languages").Delete(&languageEN)

// Clear associations
db.Model(&user).Association("Languages").Clear()

// Count associations
db.Model(&user).Association("Languages").Count()
```
