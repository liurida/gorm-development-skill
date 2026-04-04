# Key Concepts for GORM Security

This document provides key concepts for security best practices when using GORM.

## SQL Injection Prevention

GORM uses the `database/sql` package's argument placeholders to construct SQL statements, which automatically escapes arguments to prevent SQL injection.

### Safe Query Methods

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

### Vulnerable Methods

The following methods do NOT escape user input and require whitelist validation:

| Method | Risk Level | Mitigation |
|--------|------------|------------|
| `Select()` | High | Whitelist column names |
| `Distinct()` | High | Whitelist column names |
| `Order()` | High | Whitelist order fields |
| `Group()` | High | Whitelist group fields |
| `Having()` | High | Avoid user input in HAVING |
| `Table()` | High | Whitelist table names |
| `Raw()` | Critical | Avoid or use parameterized |
| `Exec()` | Critical | Avoid or use parameterized |
| `Joins()` | High | Avoid user input in joins |
| `Pluck()` | High | Whitelist column names |

Example of vulnerable usage:

```go
// VULNERABLE: SQL injection possible
db.Select("name; drop table users;").First(&user)
db.Order("name; drop table users;").First(&user)
db.Table("users; drop table users;").Find(&users)
db.Raw("select name from users; drop table users;").First(&user)
```

## Safe Patterns

### Numeric ID Validation

Always validate numeric IDs from user input:

```go
userInputID := "1=1;drop table users;"

// Safe: validate and convert to integer first
id, err := strconv.Atoi(userInputID)
if err != nil {
    return err
}
db.First(&user, id)
```

### Whitelist Validation for Dynamic Fields

For methods that don't escape input, use whitelist validation:

```go
// Whitelist allowed fields
allowedFields := map[string]bool{
    "name":       true,
    "email":      true,
    "created_at": true,
}

if !allowedFields[userField] {
    return fmt.Errorf("invalid field: %s", userField)
}

db.Order(userField).Find(&users)
```

### Raw SQL with Parameters

When using `Raw()` or `Exec()`, always use parameterized queries:

```go
// Safe: parameterized raw query
db.Raw("SELECT * FROM users WHERE name = ?", userInput).Scan(&users)

// Safe: parameterized exec
db.Exec("UPDATE users SET name = ? WHERE id = ?", name, id)
```

## Warning About gorm.Expr

Despite being inside parameterized functions, `gorm.Expr` does NOT escape its content:

```go
// VULNERABLE: gorm.Expr is injectable
db.Exec("UPDATE users SET name = '?' WHERE id = 1", 
    gorm.Expr("alice'; drop table users;-- "))
```

Avoid using `gorm.Expr` with user input.

## Logger Output Warning

The SQL logged by GORM's logger is NOT fully escaped like the actual executed query. Never copy and execute logged SQL directly in a database console without review.

## General Security Rules

1. **Never trust user input** - Always validate and sanitize
2. **Use parameterized queries** - Pass user input as arguments with `?` placeholders
3. **Whitelist validation** - For methods that don't escape, validate against known good values
4. **Type validation** - Convert and validate user input to expected types (int, UUID, etc.)
5. **Principle of least privilege** - Database user should have minimal required permissions
