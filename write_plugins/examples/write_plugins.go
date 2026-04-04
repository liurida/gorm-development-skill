
package examples

import (
	"fmt"
	"reflect"
	"strings"

	"gorm.io/gorm"
)

// UppercaseNamePlugin is a GORM plugin that converts the Name field of a model to uppercase before creation.	ype UppercaseNamePlugin struct{}

// Name returns the plugin's name.
func (p *UppercaseNamePlugin) Name() string {
	return "uppercaseNamePlugin"
}

// Initialize registers the callbacks for the plugin.
func (p *UppercaseNamePlugin) Initialize(db *gorm.DB) error {
	// Register the beforeCreate callback to run before the main GORM create operation.
	db.Callback().Create().Before("gorm:create").Register("my_plugin:uppercase_name", uppercaseName)
	return nil
}

// uppercaseName is the callback function that performs the logic.
func uppercaseName(db *gorm.DB) {
	// Check if the statement and schema are valid
	if db.Statement == nil || db.Statement.Schema == nil {
		return
	}

	// Iterate through fields of the model to find the 'Name' field
	for _, field := range db.Statement.Schema.Fields {
		if field.Name == "Name" {
			// Use reflection to get and set the field value
			switch db.Statement.ReflectValue.Kind() {
			case reflect.Slice, reflect.Array:
				for i := 0; i < db.Statement.ReflectValue.Len(); i++ {
					// Get the model instance from the slice
					model := db.Statement.ReflectValue.Index(i)
					// Get the field from the model instance
					nameField := model.FieldByName("Name")
					if nameField.IsValid() && nameField.CanSet() {
						if name, ok := nameField.Interface().(string); ok {
							nameField.SetString(strings.ToUpper(name))
						}
					}
				}
			case reflect.Struct:
				nameField := db.Statement.ReflectValue.FieldByName("Name")
				if nameField.IsValid() && nameField.CanSet() {
					if name, ok := nameField.Interface().(string); ok {
						nameField.SetString(strings.ToUpper(name))
					}
				}
			}
			break // Found the 'Name' field, no need to check others
		}
	}
}

// PluginUser is a user model for the plugin example.
type PluginUser struct {
	gorm.Model
	Name string
}

// UsePlugin demonstrates how to use a custom GORM plugin.
func UsePlugin(db *gorm.DB) error {
	// Register the custom plugin
	if err := db.Use(&UppercaseNamePlugin{}); err != nil {
		return err
	}

	// Create a user. The plugin will automatically convert the name to uppercase.
	user := PluginUser{Name: "john doe"}
	if err := db.Create(&user).Error; err != nil {
		return err
	}

	fmt.Printf("User created with name: %s\n", user.Name) // Should be "JOHN DOE"

	// Verify the name was stored in uppercase
	var foundUser PluginUser
	db.First(&foundUser, user.ID)
	fmt.Printf("Found user with name: %s\n", foundUser.Name) // Should be "JOHN DOE"

	return nil
}
