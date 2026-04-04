
package examples

import (
	"gorm.io/gorm"
)

// Order model for scope examples
type Order struct {
	gorm.Model
	Amount        float64
	Status        string
	PaymentMethod string
}

// AmountGreaterThan is a scope to find orders with an amount greater than a given value.
func AmountGreaterThan(amount float64) func(db *gorm.DB) *gorm.DB {
	return func(db *gorm.DB) *gorm.DB {
		return db.Where("amount > ?", amount)
	}
}

// Status is a scope to filter orders by their status.
func Status(status string) func(db *gorm.DB) *gorm.DB {
	return func(db *gorm.DB) *gorm.DB {
		return db.Where("status = ?", status)
	}
}

// PaidWithCreditCard is a simple scope to find orders paid by credit card.
func PaidWithCreditCard(db *gorm.DB) *gorm.DB {
	return db.Where("payment_method = ?", "credit_card")
}

// Paginate is a scope for paginating results.
func Paginate(page, pageSize int) func(db *gorm.DB) *gorm.DB {
	return func(db *gorm.DB) *gorm.DB {
		if page <= 0 {
			page = 1
		}
		switch {
		case pageSize > 100:
			pageSize = 100
		case pageSize <= 0:
			pageSize = 10
		}
		offset := (page - 1) * pageSize
		return db.Offset(offset).Limit(pageSize)
	}
}

// GetPaidOrdersWithLargeAmount demonstrates combining multiple scopes.
func GetPaidOrdersWithLargeAmount(db *gorm.DB) ([]Order, error) {
	var orders []Order
	err := db.Scopes(AmountGreaterThan(1000), Status("paid"), PaidWithCreditCard).Find(&orders).Error
	return orders, err
}

// GetSecondPageOfOrders demonstrates a pagination scope.
func GetSecondPageOfOrders(db *gorm.DB) ([]Order, error) {
	var orders []Order
	err := db.Scopes(Paginate(2, 25)).Find(&orders).Error
	return orders, err
}
