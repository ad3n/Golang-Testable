package services

import (
	"errors"

	"github.com/ad3n/golang-testable/repositories"
	"github.com/ad3n/golang-testable/views"
)

type Account struct {
	Repository repositories.AccountRepository
	Customer   repositories.CustomerRepository
}

func (s *Account) Balance(accountNumber int) (views.BalanceResponse, error) {
	account, err := s.Repository.Find(accountNumber)
	if err != nil {
		return views.BalanceResponse{}, err
	}

	customer, err := s.Customer.Find(account.CustomerID)
	if err != nil {
		return views.BalanceResponse{}, err
	}

	view := views.BalanceResponse{}
	view.AccountNumber = account.ID
	view.CustomerName = customer.Name
	view.Balance = account.Balance

	return view, nil
}

func (s *Account) Transfer(fromAccountNumber int, toAccountNumber int, amount float64) error {
	from, err := s.Repository.Find(fromAccountNumber)
	if err != nil {
		return errors.New("sender account not found")
	}

	if from.Balance < amount {
		return errors.New("insufficient balance")
	}

	to, err := s.Repository.Find(toAccountNumber)
	if err != nil {
		return errors.New("receiver account not found")
	}

	to.Balance = to.Balance + amount
	from.Balance = from.Balance - amount

	s.Repository.Saves(from, to)

	return nil
}
