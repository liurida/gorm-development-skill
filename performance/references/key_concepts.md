
# Key Concepts for GORM Performance

This document provides key concepts and best practices for optimizing GORM performance.

## Disable Default Transaction

GORM performs single create, update, and delete operations within a transaction by default to ensure data consistency. While safe, this adds performance overhead. You can disable it globally:

```go
db, err := gorm.Open(sqlite.Open("gorm.db"), &gorm.Config{
    SkipDefaultTransaction: true,
})
```

**When to disable:**
- High-throughput applications with single, non-critical write operations
- When you are managing transactions manually

**When to keep enabled:**
- When data consistency for every write is paramount

## Prepared Statement Caching

GORM can cache prepared statements to speed up future calls of the same query. This reduces the overhead of parsing and compiling SQL for every execution.

```go
// Enable globally
db, err := gorm.Open(sqlite.Open("gorm.db"), &gorm.Config{
    PrepareStmt: true,
})

// Enable for a specific session
tx := db.Session(&gorm.Session{PrepareStmt: true})
```

**Impact:**
- Reduces CPU load on the database server
- Decreases network round-trip time for subsequent identical queries

## Select Specific Fields

Avoid fetching unnecessary data from the database.

### Using `Select`

Explicitly specify which fields to retrieve:

```go
db.Select("name", "age").Find(&Users{})
```

### Smart Select Fields (using a smaller struct)

Define a smaller struct for API responses or specific use cases. GORM will automatically select only the fields present in that struct.

```go
type APIUser struct {
    ID   uint
    Name string
}

var apiUsers []APIUser
db.Model(&User{}).Limit(10).Find(&apiUsers)
// SELECT `id`, `name` FROM `users` LIMIT 10
```

## Batch Processing

For large datasets, process records in batches to avoid loading everything into memory.

### `FindInBatches`

```go
var results []User
db.Where("processed = ?", false).FindInBatches(&results, 100, func(tx *gorm.DB, batch int) error {
    // Process 100 users at a time
    return nil
})
```

### `Rows` for Iteration

For fine-grained control, iterate over rows one by one:

```go
rows, err := db.Model(&User{}).Rows()
defer rows.Close()

for rows.Next() {
    var user User
    db.ScanRows(rows, &user)
    // Process user
}
```

## Indexing and Hints

### Database Indexes

Ensure your tables have appropriate indexes for columns used in `WHERE`, `ORDER BY`, and `JOIN` clauses. This is the most fundamental database performance optimization.

### Index Hints

Guide the database query optimizer to use a specific index:

```go
import "gorm.io/hints"

// Force the use of idx_user_name
db.Clauses(hints.UseIndex("idx_user_name")).Find(&User{})
```

**When to use:**
- When you know the query optimizer is making a poor choice
- For complex queries where the optimal path is not obvious

## Read/Write Splitting

Distribute database load between a primary write instance and multiple read replicas. See the `dbresolver` skill for a full implementation.

```go
// Conceptual example
db.Use(dbresolver.Register(dbresolver.Config{
    Sources:  []gorm.Dialector{...}, // Write replica
    Replicas: []gorm.Dialector{...}, // Read replicas
}))

// Write operations go to Sources
db.Create(&user)

// Read operations go to Replicas
db.First(&user, 1)
```

## Efficient Updates

### Update with `map` or `struct`

Updating with a map only modifies the specified fields and runs hooks for those fields, which is more efficient than fetching and saving the whole object.

```go
// Updates `name` and `age` only
db.Model(&User{}).Where("id = ?", 1).Updates(map[string]interface{}{"name": "new_name", "age": 30})
```

### Update with SQL Expressions

Perform calculations on the database server to avoid a read-modify-write cycle.

```go
// BAD: Fetches, modifies, saves
db.First(&user, 1)
user.Age += 1
db.Save(&user)

// GOOD: Server-side update
db.Model(&User{}).Where("id = ?", 1).Update("age", gorm.Expr("age + ?", 1))
```

## Performance Anti-Patterns

### N+1 Queries

**Problem:** Accessing associations without preloading.

```go
// BAD: One query for users, then one query PER user for their orders
var users []User
db.Find(&users)
for _, u := range users {
    db.Model(&u).Association("Orders").Find(&u.Orders)
}

// GOOD: Two queries total
db.Preload("Orders").Find(&users)
```

### Updating Full Objects

**Problem:** Fetching and saving a whole object just to change one field.

```go
// INEFFICIENT: Fetches all fields, saves all fields
var user User
db.First(&user, 1)
user.LastLogin = time.Now()
db.Save(&user)

// EFFICIENT: Updates only one field
db.Model(&User{}).Where("id = ?", 1).Update("last_login", time.Now())
```

### Leaking Goroutines in Transactions

**Problem:** If a context is cancelled, a transaction might hang, waiting for a lock that will never be released.

**Solution:** Always check `ctx.Err()` inside transaction loops.

```go
return db.WithContext(ctx).Transaction(func(tx *gorm.DB) error {
    // ... do work ...
    if ctx.Err() != nil {
        return ctx.Err()
    }
    // ... do more work ...
    return nil
})
```

By applying these concepts, you can significantly improve the performance and scalability of your GORM-based applications.
