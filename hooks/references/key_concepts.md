# Key Concepts for GORM Hooks

This document provides key concepts for using hooks in GORM.

## Overview

Hooks are functions that are automatically called before or after CRUD operations (Create, Read, Update, Delete). They allow you to add custom logic at different points in an object's lifecycle.

## Available Hooks

### Create
- `BeforeCreate`
- `AfterCreate`

### Update
- `BeforeUpdate`
- `AfterUpdate`

### Delete
- `BeforeDelete`
- `AfterDelete`

### Query
- `AfterFind`

### General
- `BeforeSave` (called for both Create and Update)
- `AfterSave` (called for both Create and Update)

## Hook Function Signature

Hook functions must have the following signature:

```go
func (u *User) HookMethodName(tx *gorm.DB) (err error)
```

- The method is defined on the model struct.
- It receives a `*gorm.DB` instance, which is the current database transaction.
- Returning an error will cause GORM to rollback the transaction.

## Example: `BeforeCreate`

```go
import "github.com/google/uuid"

func (u *User) BeforeCreate(tx *gorm.DB) (err error) {
  u.UUID = uuid.New().String()
  return
}
```

This hook sets a new UUID for the `User` record before it is created.
