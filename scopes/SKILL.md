---
name: gorm-scopes
description: Use when creating reusable query scopes in GORM, building composable query conditions, implementing pagination, multi-tenancy, or dynamic table routing.
---

# Scopes

Scopes allow you to re-use commonly used logic. The shared logic needs to be defined as type `func(*gorm.DB) *gorm.DB`.

## Basic Query Scopes

```go
func AmountGreaterThan1000(db *gorm.DB) *gorm.DB {
    return db.Where("amount > ?", 1000)
}

func PaidWithCreditCard(db *gorm.DB) *gorm.DB {
    return db.Where("pay_mode = ?", "card")
}

func PaidWithCod(db *gorm.DB) *gorm.DB {
    return db.Where("pay_mode = ?", "cod")
}

// Usage: combine multiple scopes
db.Scopes(AmountGreaterThan1000, PaidWithCreditCard).Find(&orders)
// SELECT * FROM orders WHERE amount > 1000 AND pay_mode = 'card'

db.Scopes(AmountGreaterThan1000, PaidWithCod).Find(&orders)
// SELECT * FROM orders WHERE amount > 1000 AND pay_mode = 'cod'
```

## Parameterized Scopes

Scopes can accept parameters by returning a closure:

```go
func OrderStatus(status []string) func(db *gorm.DB) *gorm.DB {
    return func(db *gorm.DB) *gorm.DB {
        return db.Where("status IN ?", status)
    }
}

func AmountGreaterThan(amount int) func(db *gorm.DB) *gorm.DB {
    return func(db *gorm.DB) *gorm.DB {
        return db.Where("amount > ?", amount)
    }
}

// Combine parameterized scopes
db.Scopes(AmountGreaterThan(1000), OrderStatus([]string{"paid", "shipped"})).Find(&orders)
// SELECT * FROM orders WHERE amount > 1000 AND status IN ('paid', 'shipped')
```

## Pagination Scope

A common use case for scopes is implementing pagination:

```go
func Paginate(r *http.Request) func(db *gorm.DB) *gorm.DB {
    return func(db *gorm.DB) *gorm.DB {
        q := r.URL.Query()
        page, _ := strconv.Atoi(q.Get("page"))
        if page <= 0 {
            page = 1
        }

        pageSize, _ := strconv.Atoi(q.Get("page_size"))
        switch {
        case pageSize > 100:
            pageSize = 100
        case pageSize <= 0:
            pageSize = 10
        }

        offset := (page - 1) * pageSize
        return db.Offset(offset).Limit(pageSize)
    }
}

// Usage across different models
db.Scopes(Paginate(r)).Find(&users)
db.Scopes(Paginate(r)).Find(&articles)
db.Scopes(Paginate(r)).Find(&orders)
```

## Dynamic Table Scope

Use scopes to dynamically set the query table for table sharding or multi-database scenarios:

```go
// Table sharding by year
func TableOfYear(user *User, year int) func(db *gorm.DB) *gorm.DB {
    return func(db *gorm.DB) *gorm.DB {
        tableName := user.TableName() + strconv.Itoa(year)
        return db.Table(tableName)
    }
}

db.Scopes(TableOfYear(user, 2023)).Find(&users)
// SELECT * FROM users_2023

db.Scopes(TableOfYear(user, 2024)).Find(&users)
// SELECT * FROM users_2024

// Table from different database/schema
func TableOfOrg(user *User, dbName string) func(db *gorm.DB) *gorm.DB {
    return func(db *gorm.DB) *gorm.DB {
        tableName := dbName + "." + user.TableName()
        return db.Table(tableName)
    }
}

db.Scopes(TableOfOrg(user, "org1")).Find(&users)
// SELECT * FROM org1.users

db.Scopes(TableOfOrg(user, "org2")).Find(&users)
// SELECT * FROM org2.users
```

## Update/Delete Scopes

Scopes work with update and delete operations too:

```go
func CurOrganization(r *http.Request) func(db *gorm.DB) *gorm.DB {
    return func(db *gorm.DB) *gorm.DB {
        org := r.URL.Query().Get("org")
        if org != "" {
            var organization Organization
            if db.Session(&gorm.Session{}).First(&organization, "name = ?", org).Error == nil {
                return db.Where("org_id = ?", organization.ID)
            }
        }
        db.AddError(errors.New("invalid organization"))
        return db
    }
}

// Update with scope
db.Model(&Article{}).Scopes(CurOrganization(r)).Update("Name", "name 1")
// UPDATE articles SET name = 'name 1' WHERE org_id = 111

// Delete with scope
db.Scopes(CurOrganization(r)).Delete(&Article{})
// DELETE FROM articles WHERE org_id = 111
```

## Multi-Tenancy Scope

Implement row-level multi-tenancy:

```go
func TenantScope(tenantID uint) func(db *gorm.DB) *gorm.DB {
    return func(db *gorm.DB) *gorm.DB {
        return db.Where("tenant_id = ?", tenantID)
    }
}

// Apply to all queries for a tenant
tenantDB := db.Scopes(TenantScope(currentTenantID))

tenantDB.Find(&users)    // Only users for this tenant
tenantDB.Find(&orders)   // Only orders for this tenant
tenantDB.Create(&newOrder) // tenant_id must be set in model
```

## Soft Delete Scope

Custom soft delete logic:

```go
func NotDeleted(db *gorm.DB) *gorm.DB {
    return db.Where("deleted_at IS NULL")
}

func IncludeDeleted(db *gorm.DB) *gorm.DB {
    return db.Unscoped()
}

func OnlyDeleted(db *gorm.DB) *gorm.DB {
    return db.Unscoped().Where("deleted_at IS NOT NULL")
}

db.Scopes(NotDeleted).Find(&users)     // Active users only
db.Scopes(IncludeDeleted).Find(&users) // All users
db.Scopes(OnlyDeleted).Find(&users)    // Deleted users only
```

## Composing Scopes

Scopes can call other scopes:

```go
func ActivePremiumUsers(db *gorm.DB) *gorm.DB {
    return db.Scopes(NotDeleted).Where("subscription = ?", "premium")
}

func RecentActiveUsers(days int) func(db *gorm.DB) *gorm.DB {
    return func(db *gorm.DB) *gorm.DB {
        return db.Scopes(NotDeleted).Where("last_login > ?", time.Now().AddDate(0, 0, -days))
    }
}

db.Scopes(ActivePremiumUsers, RecentActiveUsers(30)).Find(&users)
```

## When NOT to Use

- **One-off queries** - Don't create a scope for a condition used only once; inline it
- **Simple conditions** - `Where("active = ?", true)` is clearer than `Scopes(ActiveOnly)` for trivial filters
- **When scope logic is complex** - If a scope has many side effects or branches, consider a repository method instead
- **Cross-cutting concerns better suited for middleware** - Use plugins or callbacks for logging, metrics, etc.
- **When readability suffers** - Multiple nested scopes can obscure what SQL is actually generated

## Quick Reference

| Pattern | Use Case |
|---------|----------|
| Simple scope | Reusable static conditions |
| Parameterized scope | Conditions with dynamic values |
| Pagination scope | Consistent pagination logic |
| Dynamic table scope | Table sharding, multi-database |
| Organization scope | Multi-tenancy, row-level security |
| Composed scope | Complex reusable query patterns |

## Reference

- Official Docs: https://gorm.io/docs/scopes.html
- Advanced Query: https://gorm.io/docs/advanced_query.html
