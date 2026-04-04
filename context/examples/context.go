
package examples

import (
	"context"
	"fmt"
	"net/http"
	"time"

	"gorm.io/gorm"
)

// User model for context examples
type User struct {
	gorm.Model
	Name string
	Role string
}

// BasicContext demonstrates passing a context to a single query.
func BasicContext(db *gorm.DB) error {
	ctx, cancel := context.WithTimeout(context.Background(), 100*time.Millisecond)
	defer cancel()

	var user User
	// This will return a timeout error if the query takes too long
	return db.WithContext(ctx).First(&user, 1).Error
}

// ContextInTransaction demonstrates using context within a transaction.
func ContextInTransaction(db *gorm.DB) error {
	ctx := context.WithValue(context.Background(), "request_id", "some-unique-id")

	return db.WithContext(ctx).Transaction(func(tx *gorm.DB) error {
		// The transaction and all operations within it will carry the context.
		if err := tx.Create(&User{Name: "transaction-user"}).Error; err != nil {
			return err
		}
		return nil
	})
}

// BeforeCreate is a GORM hook that uses the context.
func (u *User) BeforeCreate(tx *gorm.DB) (err error) {
	ctx := tx.Statement.Context
	if requestID := ctx.Value("request_id"); requestID != nil {
		fmt.Printf("BeforeCreate hook: processing request %v\n", requestID)
	}
	return
}

// ChiMiddleware is an example of a middleware for a web server like Chi.
func ChiMiddleware(db *gorm.DB) func(http.Handler) http.Handler {
	return func(next http.Handler) http.Handler {
		return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
			timeoutContext, cancel := context.WithTimeout(r.Context(), 2*time.Second)
			defer cancel()

			// Add the database connection to the context
			ctx := context.WithValue(timeoutContext, "DB", db.WithContext(timeoutContext))
			next.ServeHTTP(w, r.WithContext(ctx))
		})
	}
}
