package repositories

import "github.com/ad3n/golang-testable/models"

type (
	CustomerRepository interface {
		Find(Id int) (*models.Customer, error)
		Saves(customers ...*models.Customer) error
	}

	AccountRepository interface {
		Find(Id int) (*models.Account, error)
		Saves(accounts ...*models.Account) error
	}
)
