
package examples

import (
	"gorm.io/gorm"
	"gorm.io/hints"
)

// HintUser model for hints examples
type HintUser struct {
	gorm.Model
	Name string `gorm:"index:idx_hint_users_name"`
	Age  int
}

// UseIndexHint demonstrates how to provide an index hint to the database.
func UseIndexHint(db *gorm.DB) error {
	var users []HintUser
	// This query suggests that the database should use the 'idx_hint_users_name' index.
	// The exact SQL generated depends on the database dialect.
	return db.Clauses(hints.UseIndex("idx_hint_users_name")).Where("name = ?", "hint_user").Find(&users).Error
}

// ForceIndexHint demonstrates forcing the database to use a specific index.
func ForceIndexHint(db *gorm.DB) error {
	var users []HintUser
	// This query forces the database to use the 'idx_hint_users_name' index.
	return db.Clauses(hints.ForceIndex("idx_hint_users_name")).Find(&users).Error
}

// OptimizerHint demonstrates passing a general optimizer hint.
func OptimizerHint(db *gorm.DB) error {
	var users []HintUser
	// The hint format is specific to the database (e.g., MySQL, PostgreSQL).
	// This example is generic.
	return db.Clauses(hints.New("SET_VAR(optimizer_switch=\'index_merge=on\')")).Find(&users).Error
}

// CommentHint demonstrates adding a SQL comment to a query.
func CommentHint(db *gorm.DB) error {
	var users []HintUser
	// This adds a comment to the SELECT clause, which can be useful for debugging or tracing.
	return db.Clauses(hints.Comment("select", "query_for_active_users")).Where("age > ?", 18).Find(&users).Error
}
