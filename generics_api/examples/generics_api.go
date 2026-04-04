package examples

import (
	"context"

	"gorm.io/gorm"
	"gorm.io/gorm/clause"
	"gorm.io/hints"
)

// User is a basic GORM model for generics API examples.
type User struct {
	gorm.Model
	Name    string
	Age     int
	Friends []User `gorm:"many2many:user_friends;"`
	Pets    []Pet
}

// Pet represents a user's pet.
type Pet struct {
	gorm.Model
	Name   string
	UserID uint
}

// Company represents a company.
type Company struct {
	gorm.Model
	Name string
}

// BasicCRUD demonstrates basic CRUD operations with generics API.
func BasicCRUD(db *gorm.DB) {
	ctx := context.Background()

	// Create
	gorm.G[User](db).Create(ctx, &User{Name: "Alice", Age: 25})

	// Create in batches
	users := []User{
		{Name: "Bob", Age: 30},
		{Name: "Carol", Age: 28},
	}
	gorm.G[User](db).CreateInBatches(ctx, users, 10)

	// Query - First
	user, err := gorm.G[User](db).Where("name = ?", "Alice").First(ctx)
	if err != nil {
		return
	}
	_ = user

	// Query - Find (multiple)
	foundUsers, err := gorm.G[User](db).Where("age >= ?", 25).Find(ctx)
	if err != nil {
		return
	}
	_ = foundUsers

	// Update single field
	gorm.G[User](db).Where("name = ?", "Alice").Update(ctx, "age", 26)

	// Update multiple fields
	gorm.G[User](db).Where("name = ?", "Alice").Updates(ctx, User{Name: "Alice Updated", Age: 27})

	// Delete
	gorm.G[User](db).Where("name = ?", "Bob").Delete(ctx)
}

// AdvancedOptions demonstrates using clauses and hints.
func AdvancedOptions(db *gorm.DB) {
	ctx := context.Background()

	// OnConflict - do nothing
	gorm.G[User](db, clause.OnConflict{DoNothing: true}).Create(ctx, &User{Name: "Test"})

	// OnConflict - update on conflict
	gorm.G[User](db, clause.OnConflict{
		Columns:   []clause.Column{{Name: "name"}},
		DoUpdates: clause.AssignmentColumns([]string{"age"}),
	}).Create(ctx, &User{Name: "Test", Age: 30})

	// Execution hints
	_, _ = gorm.G[User](db,
		hints.New("MAX_EXECUTION_TIME(100)"),
	).Find(ctx)

	// Get result metadata
	result := gorm.WithResult()
	gorm.G[User](db, result).Create(ctx, &User{Name: "WithResult"})
	_ = result.RowsAffected
}

// JoinsExample demonstrates enhanced joins with generics.
func JoinsExample(db *gorm.DB) {
	ctx := context.Background()

	// Inner join - users who have a company
	_, _ = gorm.G[User](db).Joins(clause.Has("Company"), nil).Find(ctx)

	// Left join with custom filter
	_, _ = gorm.G[User](db).Joins(clause.LeftJoin.Association("Company"), func(jb gorm.JoinBuilder, joinTable, curTable clause.Table) error {
		jb.Where(map[string]any{"name": "ACME"})
		return nil
	}).Find(ctx)
}

// PreloadExample demonstrates enhanced preload with generics.
func PreloadExample(db *gorm.DB) {
	ctx := context.Background()

	// Basic preload with conditions
	_, _ = gorm.G[User](db).Preload("Friends", func(pb gorm.PreloadBuilder) error {
		pb.Where("age > ?", 18)
		return nil
	}).Find(ctx)

	// Nested preload
	_, _ = gorm.G[User](db).Preload("Friends.Pets", nil).Find(ctx)

	// Preload with ordering and per-record limit
	_, _ = gorm.G[User](db).Preload("Friends", func(pb gorm.PreloadBuilder) error {
		pb.Select("id", "name").Order("age desc")
		return nil
	}).Preload("Friends.Pets", func(pb gorm.PreloadBuilder) error {
		pb.LimitPerRecord(2)
		return nil
	}).Find(ctx)
}

// RawSQLExample demonstrates raw SQL with generics.
func RawSQLExample(db *gorm.DB) {
	ctx := context.Background()

	// Raw query
	_, _ = gorm.G[User](db).Raw("SELECT * FROM users WHERE name = ?", "Alice").Find(ctx)

	// Raw query with primitive type result
	_, _ = gorm.G[int](db).Raw("SELECT COUNT(*) FROM users").Find(ctx)
}
