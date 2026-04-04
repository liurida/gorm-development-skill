package examples

import (
	"database/sql"
	"time"

	"gorm.io/gorm"
)

// User demonstrates a basic GORM model.
// It includes a gorm.Model to get ID, CreatedAt, UpdatedAt, and DeletedAt fields.
type User struct {
	gorm.Model
	Name         string
	Email        *string `gorm:"unique"`
	Age          uint8
	Birthday     *time.Time
	MemberNumber sql.NullString
	ActivatedAt  sql.NullTime
}
