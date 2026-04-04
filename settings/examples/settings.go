package examples

import (
	"gorm.io/gorm"
)

// User is a basic GORM model for settings examples.
type User struct {
	gorm.Model
	Name       string
	CreditCard CreditCard
}

// CreditCard demonstrates nested associations with settings.
type CreditCard struct {
	gorm.Model
	Number string
	UserID uint
}

// SetGetExample demonstrates using Set and Get to pass values to hooks.
// Values set with Set() are available in all hooks, including nested associations.
func SetGetExample(db *gorm.DB) error {
	// Set a value that will be available in all hooks
	myValue := 123
	result := db.Set("my_value", myValue).Create(&User{Name: "John"})
	return result.Error
}

// BeforeCreateWithGet demonstrates reading values in a BeforeCreate hook.
// This hook receives values set with db.Set().
func (u *User) BeforeCreate(tx *gorm.DB) error {
	myValue, ok := tx.Get("my_value")
	if ok {
		// myValue is available here as interface{}
		// Type assert as needed: myValue.(int)
		_ = myValue
	}
	return nil
}

// BeforeCreateCreditCard shows that Set() values are available in association hooks.
func (c *CreditCard) BeforeCreate(tx *gorm.DB) error {
	myValue, ok := tx.Get("my_value")
	if ok {
		// Values from Set() are available even in nested association hooks
		_ = myValue
	}
	return nil
}

// InstanceSetGetExample demonstrates using InstanceSet and InstanceGet.
// InstanceSet values are scoped to the current statement only.
func InstanceSetGetExample(db *gorm.DB) error {
	myValue := 456
	result := db.InstanceSet("my_instance_value", myValue).Create(&User{Name: "Jane"})
	return result.Error
}

// BeforeCreateInstance demonstrates the difference between Set and InstanceSet.
// InstanceGet values are only available in the immediate model's hooks,
// not in associated model hooks.
func (u *User) BeforeCreateInstance(tx *gorm.DB) error {
	// This works for InstanceSet values
	myValue, ok := tx.InstanceGet("my_instance_value")
	if ok {
		// myValue is available in User's hooks
		_ = myValue
	}
	return nil
}

// BeforeCreateCreditCardInstance shows InstanceSet limitation with associations.
// Note: CreditCard hooks will NOT have access to InstanceSet values
// because GORM creates a new *Statement for associations.
func (c *CreditCard) BeforeCreateCreditCardInstance(tx *gorm.DB) error {
	myValue, ok := tx.InstanceGet("my_instance_value")
	// ok will be false here - InstanceSet values don't propagate to associations
	if !ok {
		// Expected: InstanceSet values are not available in association hooks
		_ = myValue
	}
	return nil
}

// TableOptionsExample demonstrates setting table options for migrations.
func TableOptionsExample(db *gorm.DB) error {
	// Set table engine option when creating/migrating tables
	result := db.Set("gorm:table_options", "ENGINE=InnoDB").AutoMigrate(&User{})
	return result
}

// TableCommentExample demonstrates setting table comment.
func TableCommentExample(db *gorm.DB) error {
	// Set table comment
	result := db.Set("gorm:table_options", "COMMENT='User accounts table'").AutoMigrate(&User{})
	return result
}

// MultipleSettingsExample demonstrates chaining multiple Set calls.
func MultipleSettingsExample(db *gorm.DB) error {
	result := db.
		Set("audit_user", "admin").
		Set("audit_action", "create").
		Create(&User{Name: "Created by admin"})
	return result.Error
}

// ConditionalLogicInHook demonstrates using settings for conditional logic.
type Order struct {
	gorm.Model
	Total     float64
	SkipAudit bool `gorm:"-"` // Non-persisted field
}

// BeforeCreate demonstrates conditional logic based on settings.
func (o *Order) BeforeCreate(tx *gorm.DB) error {
	skipAudit, ok := tx.Get("skip_audit")
	if ok && skipAudit.(bool) {
		// Skip audit logging when explicitly requested
		return nil
	}
	// Perform audit logging...
	return nil
}

// SkipAuditCreate demonstrates conditionally skipping functionality.
func SkipAuditCreate(db *gorm.DB, order *Order) error {
	result := db.Set("skip_audit", true).Create(order)
	return result.Error
}
