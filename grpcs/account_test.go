package grpcs

import (
	"context"
	"errors"
	"testing"

	"github.com/ad3n/golang-testable/mocks"
	"github.com/ad3n/golang-testable/models"
	"github.com/ad3n/golang-testable/protos"
	"github.com/ad3n/golang-testable/services"
	"github.com/stretchr/testify/assert"
)

func TestAccountBalanceAccountNotFound(t *testing.T) {
	accountRepository := mocks.AccountRepository{}
	accountRepository.On("Find", 123).Return(&models.Account{}, errors.New("account not found")).Once()

	customerRepository := mocks.CustomerRepository{}

	accountService := services.Account{Repository: &accountRepository, Customer: &customerRepository}

	rpc := Account{Service: &accountService}

	request := protos.BalanceRequest{}
	request.AccountNumber = 123

	_, err := rpc.Balance(context.Background(), &request)

	assert.Equal(t, "account not found", err.Error())
}

func TestAccountBalanceCustomerNotFound(t *testing.T) {
	account := models.Account{}
	account.ID = 123
	account.CustomerID = 321
	account.Balance = 10000

	accountRepository := mocks.AccountRepository{}
	accountRepository.On("Find", account.ID).Return(&account, nil).Once()

	customerRepository := mocks.CustomerRepository{}
	customerRepository.On("Find", account.CustomerID).Return(&models.Customer{}, errors.New("customer not found")).Once()

	accountService := services.Account{Repository: &accountRepository, Customer: &customerRepository}

	rpc := Account{Service: &accountService}

	request := protos.BalanceRequest{}
	request.AccountNumber = 123

	_, err := rpc.Balance(context.Background(), &request)

	assert.Equal(t, "customer not found", err.Error())
}

func TestAccountBalanceSuccess(t *testing.T) {
	account := models.Account{}
	account.ID = 123
	account.CustomerID = 321
	account.Balance = 10000

	customer := models.Customer{}
	customer.ID = account.CustomerID
	customer.Name = "John Doe"

	accountRepository := mocks.AccountRepository{}
	accountRepository.On("Find", account.ID).Return(&account, nil).Once()

	customerRepository := mocks.CustomerRepository{}
	customerRepository.On("Find", account.CustomerID).Return(&customer, nil).Once()

	accountService := services.Account{Repository: &accountRepository, Customer: &customerRepository}

	rpc := Account{Service: &accountService}

	request := protos.BalanceRequest{}
	request.AccountNumber = int32(account.ID)

	response, err := rpc.Balance(context.Background(), &request)

	assert.Equal(t, nil, err)
	assert.Equal(t, request.AccountNumber, response.AccountNumber)
}
