# Key Concepts for GORM Query

This document provides key concepts for querying records with GORM.

## Retrieving a Single Object

- `First`: Retrieves the first record, ordered by primary key.
- `Take`: Retrieves one record, without a specified order.
- `Last`: Retrieves the last record, ordered by primary key desc.

```go
var user User
db.First(&user) // Get the first record
db.Take(&user)  // Get one record
db.Last(&user)   // Get the last record
```

## Retrieving All Objects

- `Find`: Retrieves all records that match the given conditions.

```go
var users []User
db.Find(&users) // Get all records
```

## Conditions

- **String Conditions**: Use a string with placeholders for arguments.
  ```go
  db.Where("name = ?", "jinzhu").First(&user)
  ```
- **Struct & Map Conditions**: Use a struct or map to build the query. GORM will only query with non-zero fields for structs.
  ```go
  db.Where(&User{Name: "jinzhu", Age: 20}).First(&user)
  db.Where(map[string]interface{}{"name": "jinzhu", "age": 20}).Find(&users)
  ```

## Other Common Methods

- `Select`: Specifies the fields to be retrieved.
- `Order`: Specifies the order of the results.
- `Limit`: Specifies the maximum number of records to retrieve.
- `Offset`: Specifies the number of records to skip.
- `Group`: Groups the results by a specified column.
- `Having`: Adds a `HAVING` clause to the query.
- `Joins`: Specifies a `JOIN` clause.
