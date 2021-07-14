package models

type Customer struct {
	ID   int    `gorm:"column:customer_number"`
	Name string `gorm:"column:customer_name"`
}
