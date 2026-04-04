
# Key Concepts for Writing GORM Plugins

This document provides key concepts for extending GORM by writing custom plugins and using callbacks.

## Overview

GORM's plugin system is built upon its callback mechanism. Callbacks are functions that are executed at specific points during CRUD operations. By creating custom plugins, you can register new callbacks to add or modify functionality.

## The Plugin Interface

A GORM plugin is a struct that implements the `gorm.Plugin` interface.

```go
type Plugin interface {
    Name() string
    Initialize(*gorm.DB) error
}
```

- `Name()`: This method must return a unique string that identifies the plugin.
- `Initialize(*gorm.DB)`: This method is called when the plugin is registered with a GORM DB instance. This is where you register your callbacks.

## Registering a Plugin

You register a plugin using the `Use` method on a `*gorm.DB` instance.

```go
db.Use(&MyCustomPlugin{})
```

## Callbacks

Callbacks are the heart of the plugin system. You can hook into various stages of an operation (`Create`, `Query`, `Update`, `Delete`).

### Registering a Callback

You can register a function to be called at a specific point in an operation's lifecycle.

```go
func (p *MyCustomPlugin) Initialize(db *gorm.DB) error {
    // Register a function to run before the main create operation
    db.Callback().Create().Before("gorm:create").Register("my_plugin:before_create", myBeforeCreateFunction)
    return nil
}

func myBeforeCreateFunction(db *gorm.DB) {
    // Your custom logic here
    // For example, you can inspect the object being created via db.Statement.ReflectValue
    fmt.Println("About to create a new record!")
}
```

### Callback Ordering

You can control the execution order of your callbacks relative to GORM's default callbacks or other plugins.

- `Before(name)`: Register the callback to run before the callback named `name`.
- `After(name)`: Register the callback to run after the callback named `name`.
- `Replace(name, newFunc)`: Replace an existing callback with a new function.
- `Remove(name)`: Remove a callback entirely.

### Accessing Model Data in Callbacks

Within a callback, you have access to the `*gorm.DB` instance for the current operation. You can get information about the model, the SQL statement, and the data being manipulated through the `db.Statement` field.

```go
func myBeforeCreateFunction(db *gorm.DB) {
    if db.Statement.Schema != nil {
        // Get the model's table name
        tableName := db.Statement.Schema.Table

        // Get the value of the model being created
        modelValue := db.Statement.ReflectValue

        fmt.Printf("Creating a record for table: %s\n", tableName)
        
        // You can even modify the data before it's saved
        if modelValue.Kind() == reflect.Struct {
            if nameField := modelValue.FieldByName("Name"); nameField.IsValid() && nameField.CanSet() {
                 if name, ok := nameField.Interface().(string); ok {
                     nameField.SetString(strings.ToUpper(name))
                 }
            }
        }
    }
}
```

By understanding and using plugins and callbacks, you can greatly extend GORM's functionality to fit the specific needs of your application.
