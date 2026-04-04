# Key Concepts for GORM Models

This document provides key concepts for declaring GORM models.

## Model Declaration

Models are defined using standard Go structs.

```go
type User struct {
  ID           uint
  Name         string
  Email        *string
  Age          uint8
  Birthday     *time.Time
  MemberNumber sql.NullString
  ActivatedAt  sql.NullTime
  CreatedAt    time.Time
  UpdatedAt    time.Time
}
```

## Conventions

- **Primary Key**: GORM uses a field named `ID` as the default primary key.
- **Table Names**: Struct names are converted to `snake_case` and pluralized (e.g., `User` becomes `users`).
- **Column Names**: Field names are converted to `snake_case` (e.g., `MemberNumber` becomes `member_number`).
- **Timestamp Fields**: `CreatedAt` and `UpdatedAt` are used to automatically track creation and update times.

## `gorm.Model`

A predefined struct that includes `ID`, `CreatedAt`, `UpdatedAt`, and `DeletedAt` fields.

```go
type Model struct {
  ID        uint           `gorm:"primaryKey"`
  CreatedAt time.Time
  UpdatedAt time.Time
  DeletedAt gorm.DeletedAt `gorm:"index"`
}
```

## Field Tags

Tags are used to configure field properties.

| Tag Name       | Description                                                  |
|----------------|--------------------------------------------------------------|
| `column`       | Specifies the column's database name.                        |
| `type`         | Specifies the column's data type.                            |
| `size`         | Specifies the column's data size/length.                     |
| `primaryKey`   | Marks a field as the primary key.                            |
| `unique`       | Marks a field as unique.                                     |
| `default`      | Specifies a default value for a field.                       |
| `not null`     | Marks a field as NOT NULL.                                   |
| `autoIncrement`| Enables auto-increment for a field.                          |
