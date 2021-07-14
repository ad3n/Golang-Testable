package models

type Account struct {
	ID         int     `gorm:"column:account_number"`
	CustomerID int     `gorm:"column:customer_number"`
	Balance    float64 `gorm:"column:balance;type:decimal(17,2)"`
}
