---
name: gorm-write-plugins
description: Use when writing custom GORM plugins, registering callbacks, or extending GORM functionality with reusable components.
---

# Writing Plugins

GORM's plugin system allows for easy extensibility and customization of its core functionalities, enhancing your application's capabilities while maintaining a modular architecture.

## Plugin Interface

A plugin must implement the `gorm.Plugin` interface:

```go
type Plugin interface {
    Name() string
    Initialize(*gorm.DB) error
}
```

- **`Name()`**: Returns a unique string identifier for the plugin.
- **`Initialize(*gorm.DB)`**: Contains the logic to set up the plugin. Called when the plugin is registered with GORM for the first time.

## Registering a Plugin

Once your plugin conforms to the `Plugin` interface, register it with a GORM instance:

```go
db.Use(MyCustomPlugin{})
```

## Accessing Registered Plugins

After registration, plugins are stored in GORM's configuration:

```go
plugin := db.Config.Plugins[pluginName]
```

## Callbacks

GORM leverages `Callbacks` to power its core functionalities. Callbacks provide hooks for various database operations like `Create`, `Query`, `Update`, `Delete`, `Row`, and `Raw`.

**Important**: Callbacks are registered at the global `*gorm.DB` level, not on a session basis. If you need different callback behaviors, initialize a separate `*gorm.DB` instance.

### Registering a Callback

```go
func cropImage(db *gorm.DB) {
    if db.Statement.Schema != nil {
        for _, field := range db.Statement.Schema.Fields {
            switch db.Statement.ReflectValue.Kind() {
            case reflect.Slice, reflect.Array:
                for i := 0; i < db.Statement.ReflectValue.Len(); i++ {
                    if fieldValue, isZero := field.ValueOf(db.Statement.Context, db.Statement.ReflectValue.Index(i)); !isZero {
                        if crop, ok := fieldValue.(CropInterface); ok {
                            crop.Crop()
                        }
                    }
                }
            case reflect.Struct:
                if fieldValue, isZero := field.ValueOf(db.Statement.Context, db.Statement.ReflectValue); !isZero {
                    if crop, ok := fieldValue.(CropInterface); ok {
                        crop.Crop()
                    }
                }
            }
        }
    }
}

db.Callback().Create().Register("crop_image", cropImage)
```

### Callback Operations

| Operation | Description |
|-----------|-------------|
| `Remove("name")` | Removes a callback by name |
| `Replace("name", fn)` | Replaces a callback with a new function |
| `Before("name")` | Registers to execute before specified callback |
| `After("name")` | Registers to execute after specified callback |
| `Before("*")` | Registers to execute before all callbacks |
| `After("*")` | Registers to execute after all callbacks |

### Ordering Callbacks

```go
// Before a specific callback
db.Callback().Create().Before("gorm:create").Register("my_plugin:before_create", beforeCreate)

// After a specific callback
db.Callback().Create().After("gorm:create").Register("my_plugin:after_create", afterCreate)

// Between two callbacks
db.Callback().Create().Before("gorm:create").After("gorm:before_create").Register("my_plugin:middle", middleFn)

// Before/after all callbacks
db.Callback().Create().Before("*").Register("my_plugin:first", firstFn)
db.Callback().Create().After("*").Register("my_plugin:last", lastFn)
```

## Complete Plugin Example

```go
package myplugin

import (
    "fmt"
    "time"
    "gorm.io/gorm"
)

type AuditPlugin struct {
    Enabled bool
}

func (p *AuditPlugin) Name() string {
    return "audit_plugin"
}

func (p *AuditPlugin) Initialize(db *gorm.DB) error {
    // Register before create callback
    db.Callback().Create().Before("gorm:create").Register("audit:before_create", func(db *gorm.DB) {
        if p.Enabled {
            fmt.Printf("[AUDIT] Creating record at %v\n", time.Now())
        }
    })

    // Register after create callback
    db.Callback().Create().After("gorm:create").Register("audit:after_create", func(db *gorm.DB) {
        if p.Enabled && db.Statement.Schema != nil {
            fmt.Printf("[AUDIT] Created %s record\n", db.Statement.Schema.Name)
        }
    })

    // Register query callback
    db.Callback().Query().After("gorm:query").Register("audit:after_query", func(db *gorm.DB) {
        if p.Enabled {
            fmt.Printf("[AUDIT] Query executed: %s\n", db.Statement.SQL.String())
        }
    })

    return nil
}

// Usage
func main() {
    db, _ := gorm.Open(sqlite.Open("test.db"), &gorm.Config{})
    db.Use(&AuditPlugin{Enabled: true})
}
```

## Accessing Schema Information in Callbacks

```go
func myCallback(db *gorm.DB) {
    if db.Statement.Schema != nil {
        // All fields for current model
        fields := db.Statement.Schema.Fields

        // All primary key fields
        primaryFields := db.Statement.Schema.PrimaryFields

        // Prioritized primary key (field with DB name `id` or first defined)
        prioritizedPK := db.Statement.Schema.PrioritizedPrimaryField

        // All relationships
        relationships := db.Statement.Schema.Relationships

        // Find field by name or DB name
        field := db.Statement.Schema.LookUpField("Name")
    }
}
```

## Predefined Callbacks

GORM comes with predefined callbacks that drive its standard features. Review these before creating custom plugins:
- Source: https://github.com/go-gorm/gorm/blob/master/callbacks/callbacks.go

## Official Plugin Example

The Prometheus plugin demonstrates real-world plugin implementation:

```go
import "gorm.io/plugin/prometheus"

db.Use(prometheus.New(prometheus.Config{
    DBName:          "db1",
    RefreshInterval: 15,
    StartServer:     true,
    HTTPServerPort:  8080,
    MetricsCollector: []prometheus.MetricsCollector{
        &prometheus.MySQL{VariableNames: []string{"Threads_running"}},
    },
}))
```

## When NOT to Use

- **For simple, one-off logic** - If you only need to run a piece of code for a single model, a standard GORM hook (`BeforeCreate`, etc.) on the model is simpler than a full plugin.
- **When you need session-specific behavior** - Callbacks registered by plugins are global. If you need different logic for different requests, use `db.Session()` with custom settings or pass data via `context.Context` instead.
- **To replace core GORM logic without understanding it** - Be very careful when replacing or reordering GORM's default callbacks. This can have unintended side effects, such as breaking transaction handling or association management.
- **If an existing plugin already provides the functionality** - Check GORM's official plugins (like Prometheus, DBResolver, Sharding) before writing your own.

## Reference

- Official Docs: https://gorm.io/docs/write_plugins.html
- Predefined Callbacks: https://github.com/go-gorm/gorm/blob/master/callbacks/callbacks.go
- Prometheus Plugin: https://gorm.io/docs/prometheus.html
