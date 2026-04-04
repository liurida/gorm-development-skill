package examples

import (
	"database/sql"

	"gorm.io/gorm"
	"gorm.io/gorm/clause"
	"gorm.io/hints"
)

// User is a basic GORM model for SQL builder examples.
type User struct {
	gorm.Model
	Name  string
	Age   int
	Email string
	Money float64
}

// Result is a custom struct for raw query results.
type Result struct {
	ID   int
	Name string
	Age  int
}

// RawQueryExample demonstrates raw SQL queries with Scan.
func RawQueryExample(db *gorm.DB) {
	var result Result
	db.Raw("SELECT id, name, age FROM users WHERE id = ?", 3).Scan(&result)

	var results []Result
	db.Raw("SELECT id, name, age FROM users WHERE age > ?", 18).Scan(&results)

	var totalAge int
	db.Raw("SELECT SUM(age) FROM users WHERE name LIKE ?", "%jin%").Scan(&totalAge)
}

// ExecExample demonstrates raw SQL execution for write operations.
func ExecExample(db *gorm.DB) {
	// Simple exec
	db.Exec("UPDATE users SET name = ? WHERE id = ?", "updated_name", 1)

	// Exec with slice parameter
	db.Exec("DELETE FROM users WHERE id IN ?", []int64{1, 2, 3})

	// Exec with SQL expression
	db.Exec("UPDATE users SET money = ? WHERE name = ?",
		gorm.Expr("money * ? + ?", 1.1, 100),
		"jinzhu")
}

// NamedArgumentExample demonstrates using named arguments.
func NamedArgumentExample(db *gorm.DB) {
	var user User

	// Using sql.Named
	db.Where("name = @name OR email = @name", sql.Named("name", "jinzhu")).Find(&user)

	// Using map
	db.Where("name = @name AND age = @age", map[string]interface{}{
		"name": "jinzhu",
		"age":  18,
	}).Find(&user)

	// Using struct
	type Args struct {
		Name string
		Age  int
	}
	db.Raw("SELECT * FROM users WHERE name = @Name AND age = @Age",
		Args{Name: "jinzhu", Age: 18}).Scan(&user)
}

// DryRunExample demonstrates generating SQL without execution.
func DryRunExample(db *gorm.DB) {
	var user User
	stmt := db.Session(&gorm.Session{DryRun: true}).First(&user, 1).Statement

	// Get generated SQL
	_ = stmt.SQL.String() // SELECT * FROM `users` WHERE `id` = ? ORDER BY `id` LIMIT 1
	_ = stmt.Vars         // []interface{}{1}
}

// ToSQLExample demonstrates ToSQL for debugging.
func ToSQLExample(db *gorm.DB) {
	sql := db.ToSQL(func(tx *gorm.DB) *gorm.DB {
		return tx.Model(&User{}).Where("id = ?", 100).Limit(10).Order("age desc").Find(&[]User{})
	})
	_ = sql // SELECT * FROM "users" WHERE id = 100 AND "users"."deleted_at" IS NULL ORDER BY age desc LIMIT 10
}

// RowExample demonstrates getting results as *sql.Row.
func RowExample(db *gorm.DB) {
	var name string
	var age int

	row := db.Table("users").Where("name = ?", "jinzhu").Select("name", "age").Row()
	row.Scan(&name, &age)

	// With raw SQL
	row = db.Raw("SELECT name, age FROM users WHERE id = ?", 1).Row()
	row.Scan(&name, &age)
}

// RowsExample demonstrates getting results as *sql.Rows.
func RowsExample(db *gorm.DB) {
	rows, err := db.Model(&User{}).Where("age > ?", 18).Select("name", "age").Rows()
	if err != nil {
		return
	}
	defer rows.Close()

	for rows.Next() {
		var name string
		var age int
		rows.Scan(&name, &age)
	}
}

// ScanRowsExample demonstrates scanning rows into structs.
func ScanRowsExample(db *gorm.DB) {
	rows, err := db.Model(&User{}).Where("name LIKE ?", "%jin%").Rows()
	if err != nil {
		return
	}
	defer rows.Close()

	for rows.Next() {
		var user User
		db.ScanRows(rows, &user)
		// Process user
	}
}

// ConnectionExample demonstrates running SQL in the same connection.
func ConnectionExample(db *gorm.DB) {
	db.Connection(func(tx *gorm.DB) error {
		// These run in the same TCP connection (but not a transaction)
		tx.Exec("SET @my_var = ?", "value")
		tx.Raw("SELECT @my_var").Scan(new(string))
		return nil
	})
}

// ClausesExample demonstrates custom clauses.
func ClausesExample(db *gorm.DB) {
	var user User

	// INSERT IGNORE
	db.Clauses(clause.Insert{Modifier: "IGNORE"}).Create(&user)

	// Hints
	db.Clauses(hints.New("MAX_EXECUTION_TIME(1000)")).Find(&user)
}
