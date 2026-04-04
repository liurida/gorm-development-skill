# Key Concepts for GORM Create

This document provides key concepts for creating records with GORM.

## Create a Record

Use the `Create` method to insert a new record into the database.

```go
user := User{Name: "Jinzhu", Age: 18, Birthday: time.Now()}
result := db.Create(&user)
```

- The `user.ID` will be populated with the primary key of the inserted record.
- `result.Error` will contain any error that occurred.
- `result.RowsAffected` will return the number of inserted records.

## Create with Selected Fields

Use `Select` to specify which fields to include in the `INSERT` statement.

```go
db.Select("Name", "Age", "CreatedAt").Create(&user)
```

## Batch Insert

Pass a slice to the `Create` method to insert multiple records at once. GORM will generate a single SQL statement to insert all the data.

```go
var users = []User{{Name: "jinzhu1"}, {Name: "jinzhu2"}, {Name: "jinzhu3"}}
db.Create(&users)
```

## Upsert / On Conflict

GORM provides `clause.OnConflict` for handling "upsert" operations.

```go
import "gorm.io/gorm/clause"

// Do nothing on conflict
db.Clauses(clause.OnConflict{DoNothing: true}).Create(&user)

// Update all columns to new value on conflict
db.Clauses(clause.OnConflict{
  UpdateAll: true,
}).Create(&users)
```
