---
name: gorm-security
description: Use when implementing security best practices with GORM, preventing SQL injection, validating user input in database queries, or auditing code for SQL injection vulnerabilities.
---

# Security

GORM uses the `database/sql` package's argument placeholders to construct SQL statements, which automatically escapes arguments to prevent SQL injection.

**Core Principle:** User input should ONLY be used as an argument (with `?` placeholders), never directly embedded in SQL strings.

**Reference:** [GORM Security Documentation](https://gorm.io/docs/security.html)

## Quick Reference

| Pattern | Safe | Example |
|---------|------|---------|
| `Where("name = ?", input)` | Yes | Parameterized |
| `Where(fmt.Sprintf(...))` | NO | SQL injection |
| `First(&user, id)` (int) | Yes | Type-safe |
| `First(&user, stringID)` | NO | Injection risk |
| `Order(userInput)` | NO | Not escaped |
| `Raw(query, args...)` | Partial | Args escaped, query not |

## Safe Query Patterns

Always use parameterized queries with user input:

```go
userInput := "jinzhu;drop table users;"

// Safe: will be escaped
db.Where("name = ?", userInput).First(&user)

// Safe: inline condition with parameter
db.First(&user, "name = ?", userInput)

// Safe: map conditions
db.Where(map[string]interface{}{"name": userInput}).Find(&users)

// Safe: struct conditions
db.Where(&User{Name: userInput}).Find(&users)
```

## Unsafe Patterns to Avoid

**Do not** use string formatting to build queries:

```go
// UNSAFE: SQL injection possible
db.Where(fmt.Sprintf("name = %v", userInput)).First(&user)
// Executes: SELECT * FROM users WHERE name = jinzhu;drop table users;

// SAFE: Use parameterized query
db.Where("name = ?", userInput).First(&user)
// Executes: SELECT * FROM users WHERE name = 'jinzhu;drop table users;'
```

## Numeric ID Validation

Always validate numeric IDs from user input:

```go
userInputID := "1=1;drop table users;"

// Safe: validate and convert to integer first
id, err := strconv.Atoi(userInputID)
if err != nil {
    return err
}
db.First(&user, id)

// UNSAFE: SQL injection
db.First(&user, userInputID)
// Executes: SELECT * FROM users WHERE 1=1;drop table users;
```

## Vulnerable Methods (Complete List)

These methods do NOT escape user input and require whitelist validation:

### Column/Field Methods
```go
// All vulnerable to SQL injection
db.Select("name; drop table users;").First(&user)
db.Distinct("name; drop table users;").First(&user)
db.Model(&user).Pluck("name; drop table users;", &names)
db.Group("name; drop table users;").First(&user)
db.Order("name; drop table users;").First(&user)
```

### Clause Methods
```go
db.Group("name").Having("1 = 1;drop table users;").First(&user)
db.Table("users; drop table users;").Find(&users)
db.Joins("inner join orders; drop table users;").Find(&users)
db.InnerJoins("inner join orders; drop table users;").Find(&users)
```

### Raw SQL Methods
```go
db.Raw("select name from users; drop table users;").First(&user)
db.Exec("select name from users; drop table users;")
```

### Conditional Methods
```go
db.Delete(&User{}, "id=1; drop table users;")
db.Where("id=1").Not("name = 'alice'; drop table users;").Find(&users)
db.Where("id=1").Or("name = 'alice'; drop table users;").Find(&users)
db.Find(&User{}, "name = 'alice'; drop table users;")
```

### gorm.Expr() - Special Case
Even inside parameterized `Exec()`, `gorm.Expr` is injectable:
```go
// UNSAFE: gorm.Expr bypasses parameterization
db.Exec("UPDATE users SET name = '?' WHERE id = 1", 
    gorm.Expr("alice'; drop table users;-- "))
```

### Blind SQL Injection Methods
These can only be exploited via blind SQL injection:
```go
db.First(&users, "2 or 1=1-- ")
db.FirstOrCreate(&users, "2 or 1=1-- ")
db.FirstOrInit(&users, "2 or 1=1-- ")
db.Last(&users, "2 or 1=1-- ")
db.Take(&users, "2 or 1=1-- ")
```

## Whitelist Validation Pattern

Use whitelist validation for all vulnerable methods:

```go
// Define allowed values
var allowedOrderFields = map[string]bool{
    "name":       true,
    "email":      true,
    "created_at": true,
}

var allowedDirections = map[string]bool{
    "asc":  true,
    "desc": true,
    "ASC":  true,
    "DESC": true,
}

// Validate before use
func SafeOrder(db *gorm.DB, field, direction string) *gorm.DB {
    if !allowedOrderFields[field] {
        return db // or return error
    }
    if !allowedDirections[direction] {
        direction = "asc"
    }
    return db.Order(field + " " + direction)
}

// Usage
SafeOrder(db, userField, userDirection).Find(&users)
```

## Safe Dynamic Query Builder

For complex dynamic queries with user input:

```go
func BuildSafeQuery(db *gorm.DB, filters map[string]interface{}) *gorm.DB {
    allowedFilters := map[string]bool{
        "name": true, "email": true, "status": true,
    }
    
    for field, value := range filters {
        if !allowedFilters[field] {
            continue // Skip unknown fields
        }
        // Use parameterized query - value is escaped
        db = db.Where(field+" = ?", value)
    }
    return db
}
```

## Logger Warning

SQL logged by GORM's logger is NOT fully escaped like the executed SQL. Never copy and execute logged SQL directly in a database console - it may contain unescaped values that could cause data loss.

## When NOT to Use

- **Trusted internal systems only** - Even internal services should use parameterized queries; skip only if performance profiling proves necessity
- **Static queries with no user input** - Hardcoded queries without external data don't need extra validation (but parameterized queries are still cleaner)
- **Database migrations** - Schema changes typically don't involve user input; use raw DDL statements

Note: The techniques in this skill should be applied to ALL GORM code handling user input. There is almost never a reason to skip SQL injection prevention.

## Security Checklist

Before deploying GORM code:
- [ ] All user input uses `?` placeholders
- [ ] No `fmt.Sprintf` or string concatenation in queries
- [ ] Numeric IDs validated with `strconv.Atoi`
- [ ] Vulnerable methods use whitelist validation
- [ ] No direct use of `gorm.Expr` with user input
- [ ] Raw SQL reviewed for injection points
