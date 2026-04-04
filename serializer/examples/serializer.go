
package examples

import (
	"context"
	"database/sql/driver"
	"encoding/json"
	"errors"
	"fmt"
	"reflect"

	"gorm.io/gorm"
	"gorm.io/gorm/schema"
)

// Preferences is a struct that will be serialized to JSON.	ype Preferences struct {
	Theme  string `json:"theme"`
	Notify bool   `json:"notify"`
}

// SerializerUser model demonstrates using serializers.
type SerializerUser struct {
	gorm.Model
	Name string

	// This will be stored as a JSON string in the database.
	Preferences Preferences `gorm:"serializer:json"`

	// This uses a custom type that implements the serializer interface.
	EncryptedData EncryptedString
}

// EncryptedString is a custom type that encrypts/decrypts its value when stored/retrieved.	ype EncryptedString string

// Scan implements the gorm.Scanner interface.
func (es *EncryptedString) Scan(value interface{}) error {
	b, ok := value.([]byte)
	if !ok {
		return errors.New(fmt.Sprint("Failed to unmarshal EncryptedString value:", value))
	}

	// Simple "decryption" for example purposes (in reality, use a proper crypto library)
	*es = EncryptedString(string(b)[8:]) // Remove "ENCRYPTED:"
	return nil
}

// Value implements the driver.Valuer interface.
func (es EncryptedString) Value() (driver.Value, error) {
	// Simple "encryption"
	return "ENCRYPTED:" + string(es), nil
}

// To make EncryptedString work with GORM v2, it also needs to implement schema.SerializerInterface.
func (es *EncryptedString) ScanGORM(ctx context.Context, field *schema.Field, dst reflect.Value, dbValue interface{}) error {
    return es.Scan(dbValue)
}

func (es EncryptedString) ValueGORM(ctx context.Context, field *schema.Field, dst reflect.Value, fieldValue interface{}) (interface{}, error) {
    return es.Value()
}

// UseJSONSerializer demonstrates storing a struct as a JSON blob.
func UseJSONSerializer(db *gorm.DB) error {
	user := SerializerUser{
		Name: "json_user",
		Preferences: Preferences{
			Theme:  "dark",
			Notify: true,
		},
	}

	if err := db.Create(&user).Error; err != nil {
		return err
	}

	var foundUser SerializerUser
	db.First(&foundUser, user.ID)

	fmt.Printf("Found user with theme: %s\n", foundUser.Preferences.Theme) // Should be "dark"
	return nil
}

// UseCustomSerializer demonstrates a custom type that handles its own serialization.
func UseCustomSerializer(db *gorm.DB) error {
	user := SerializerUser{
		Name:          "encrypted_user",
		EncryptedData: "my_secret_data",
	}

	if err := db.Create(&user).Error; err != nil {
		return err
	}

	var foundUser SerializerUser
	db.First(&foundUser, user.ID)

	fmt.Printf("Decrypted data: %s\n", foundUser.EncryptedData) // Should be "my_secret_data"

	// Check the raw data in the database
	var rawResult map[string]interface{}
	db.Model(&SerializerUser{}).Where("id = ?", user.ID).First(&rawResult)
	fmt.Printf("Raw encrypted data: %s\n", rawResult["encrypted_data"]) // Should be "ENCRYPTED:my_secret_data"

	return nil
}
