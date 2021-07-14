package services

import (
	"errors"
	"testing"

	"github.com/ad3n/golang-testable/mocks"
	"github.com/ad3n/golang-testable/models"

	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/mock"
)

var accountRepository mocks.AccountRepository
var customerRepository mocks.CustomerRepository

func setUp() {
	accountRepository = mocks.AccountRepository{}
	customerRepository = mocks.CustomerRepository{}
}

func TestAccountBalanceAccountNotFound(t *testing.T) {
	setUp()

	accountRepository.On("Find", 123).Return(&models.Account{}, errors.New("account not found")).Once()

	service := Account{Repository: &accountRepository, Customer: &customerRepository}

	res, err := service.Balance(123)

	accountRepository.AssertExpectations(t)
	customerRepository.AssertExpectations(t)

	assert.Equal(t, "account not found", err.Error())
	assert.Equal(t, 0, res.AccountNumber)
	assert.Equal(t, "", res.CustomerName)
	assert.Equal(t, 0.0, res.Balance)
}

func TestAccountBalanceCustomerNotFound(t *testing.T) {
	setUp()

	account := models.Account{}
	account.ID = 123
	account.CustomerID = 321
	account.Balance = 10000

	accountRepository.On("Find", account.ID).Return(&account, nil).Once()
	customerRepository.On("Find", account.CustomerID).Return(&models.Customer{}, errors.New("customer not found")).Once()

	service := Account{Repository: &accountRepository, Customer: &customerRepository}

	res, err := service.Balance(123)

	accountRepository.AssertExpectations(t)
	customerRepository.AssertExpectations(t)

	assert.Equal(t, "customer not found", err.Error())
	assert.Equal(t, 0, res.AccountNumber)
	assert.Equal(t, "", res.CustomerName)
	assert.Equal(t, 0.0, res.Balance)
}

func TestAccountBalanceSuccess(t *testing.T) {
	setUp()

	account := models.Account{}
	account.ID = 123
	account.CustomerID = 321
	account.Balance = 10000

	customer := models.Customer{}
	customer.ID = account.CustomerID
	customer.Name = "John Doe"

	accountRepository.On("Find", account.ID).Return(&account, nil).Once()
	customerRepository.On("Find", account.CustomerID).Return(&customer, nil).Once()

	service := Account{Repository: &accountRepository, Customer: &customerRepository}

	res, err := service.Balance(123)

	accountRepository.AssertExpectations(t)
	customerRepository.AssertExpectations(t)

	assert.Equal(t, nil, err)
	assert.Equal(t, account.ID, res.AccountNumber)
	assert.Equal(t, customer.Name, res.CustomerName)
	assert.Equal(t, account.Balance, res.Balance)
}

func TestAccountTransferSenderAccountNotFound(t *testing.T) {
	setUp()

	from := 123
	to := 321

	accountRepository.On("Find", from).Return(&models.Account{}, errors.New("account not found")).Once()

	service := Account{Repository: &accountRepository, Customer: &customerRepository}

	err := service.Transfer(from, to, float64(100000))

	accountRepository.AssertExpectations(t)
	customerRepository.AssertExpectations(t)

	assert.Equal(t, "sender account not found", err.Error())
}

func TestAccountTransferInsufficientBalance(t *testing.T) {
	setUp()

	from := 123
	to := 321

	accountRepository.On("Find", from).Return(&models.Account{}, nil).Once()

	service := Account{Repository: &accountRepository, Customer: &customerRepository}

	err := service.Transfer(from, to, float64(100000))

	accountRepository.AssertExpectations(t)
	customerRepository.AssertExpectations(t)

	assert.Equal(t, "insufficient balance", err.Error())
}

func TestAccountTransferReceiverNotFound(t *testing.T) {
	setUp()

	sender := models.Account{}
	sender.ID = 123
	sender.Balance = 100000.0

	accountRepository.On("Find", mock.MatchedBy(func(accountNumber int) bool {
		return accountNumber == sender.ID
	})).Return(&sender, nil).Once()
	accountRepository.On("Find", mock.MatchedBy(func(accountNumber int) bool {
		return accountNumber == 321
	})).Return(&models.Account{}, errors.New("account not found")).Once()

	service := Account{Repository: &accountRepository, Customer: &customerRepository}

	err := service.Transfer(sender.ID, 321, sender.Balance)

	accountRepository.AssertExpectations(t)
	customerRepository.AssertExpectations(t)

	assert.Equal(t, "receiver account not found", err.Error())
}

func TestAccountTransferBalanceEqualTransferAmount(t *testing.T) {
	setUp()

	sender := models.Account{}
	sender.ID = 123
	sender.Balance = 100000.0

	receiver := models.Account{}
	receiver.ID = 321

	accountRepository.On("Find", mock.MatchedBy(func(accountNumber int) bool {
		return accountNumber == sender.ID
	})).Return(&sender, nil).Once()
	accountRepository.On("Find", mock.MatchedBy(func(accountNumber int) bool {
		return accountNumber == receiver.ID
	})).Return(&receiver, nil).Once()

	accountRepository.On("Saves", mock.Anything, mock.Anything).Return(nil).Once()

	service := Account{Repository: &accountRepository, Customer: &customerRepository}

	err := service.Transfer(sender.ID, receiver.ID, sender.Balance)

	accountRepository.AssertExpectations(t)
	customerRepository.AssertExpectations(t)

	assert.Equal(t, nil, err)
	assert.Equal(t, 0.0, sender.Balance)
}

func TestAccountTransferBalanceMoreThanTransferAmount(t *testing.T) {
	sender := models.Account{}
	sender.ID = 123
	sender.Balance = 100001.0

	receiver := models.Account{}
	receiver.ID = 321

	accountRepository.On("Find", mock.MatchedBy(func(accountNumber int) bool {
		return accountNumber == sender.ID
	})).Return(&sender, nil).Once()
	accountRepository.On("Find", mock.MatchedBy(func(accountNumber int) bool {
		return accountNumber == receiver.ID
	})).Return(&receiver, nil).Once()
	accountRepository.On("Saves", mock.Anything, mock.Anything).Return(nil).Once()

	service := Account{Repository: &accountRepository, Customer: &customerRepository}

	err := service.Transfer(sender.ID, receiver.ID, 100000.0)

	accountRepository.AssertExpectations(t)
	customerRepository.AssertExpectations(t)

	assert.Equal(t, nil, err)
	assert.Equal(t, 1.0, sender.Balance)
}
