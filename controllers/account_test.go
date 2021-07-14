package controllers

import (
	"bytes"
	"encoding/json"
	"errors"
	"io/ioutil"
	"net/http"
	"net/http/httptest"
	"testing"

	"github.com/ad3n/golang-testable/mocks"
	"github.com/ad3n/golang-testable/models"
	"github.com/ad3n/golang-testable/services"
	"github.com/ad3n/golang-testable/views"
	"github.com/stretchr/testify/mock"

	"github.com/gofiber/fiber/v2"
	"github.com/gofiber/fiber/v2/utils"
)

func TestAccountBalanceAccountNotNumber(t *testing.T) {
	accountRepository := mocks.AccountRepository{}
	customerRepository := mocks.CustomerRepository{}

	accountService := services.Account{Repository: &accountRepository, Customer: &customerRepository}

	controller := Account{Service: &accountService}

	app := fiber.New()

	app.Get("/account/:number", controller.Balance)

	resp, err := app.Test(httptest.NewRequest(fiber.MethodGet, "/account/abc", nil))

	utils.AssertEqual(t, nil, err)
	utils.AssertEqual(t, http.StatusBadRequest, resp.StatusCode)
}

func TestAccountBalanceAccountNotFound(t *testing.T) {
	accountRepository := mocks.AccountRepository{}
	accountRepository.On("Find", 123).Return(&models.Account{}, errors.New("account not found")).Once()

	customerRepository := mocks.CustomerRepository{}

	accountService := services.Account{Repository: &accountRepository, Customer: &customerRepository}

	controller := Account{Service: &accountService}

	app := fiber.New()

	app.Get("/account/:number", controller.Balance)

	resp, err := app.Test(httptest.NewRequest(fiber.MethodGet, "/account/123", nil))

	utils.AssertEqual(t, nil, err)

	body, err := ioutil.ReadAll(resp.Body)

	utils.AssertEqual(t, nil, err)

	views := map[string]string{}
	err = json.Unmarshal(body, &views)

	utils.AssertEqual(t, nil, err)
	utils.AssertEqual(t, http.StatusNotFound, resp.StatusCode)
	utils.AssertEqual(t, "account not found", views["message"])
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

	controller := Account{Service: &accountService}

	app := fiber.New()

	app.Get("/account/:number", controller.Balance)

	resp, err := app.Test(httptest.NewRequest(fiber.MethodGet, "/account/123", nil))

	utils.AssertEqual(t, nil, err)

	body, err := ioutil.ReadAll(resp.Body)

	utils.AssertEqual(t, nil, err)

	views := map[string]string{}
	err = json.Unmarshal(body, &views)

	utils.AssertEqual(t, nil, err)
	utils.AssertEqual(t, http.StatusNotFound, resp.StatusCode)
	utils.AssertEqual(t, "customer not found", views["message"])
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

	controller := Account{Service: &accountService}

	app := fiber.New()

	app.Get("/account/:number", controller.Balance)

	resp, err := app.Test(httptest.NewRequest(fiber.MethodGet, "/account/123", nil))

	utils.AssertEqual(t, nil, err)

	body, err := ioutil.ReadAll(resp.Body)

	utils.AssertEqual(t, nil, err)

	views := views.BalanceResponse{}
	err = json.Unmarshal(body, &views)

	utils.AssertEqual(t, nil, err)
	utils.AssertEqual(t, http.StatusOK, resp.StatusCode)
	utils.AssertEqual(t, account.ID, views.AccountNumber)
}

func TestAccountTransferAccountNotNumber(t *testing.T) {
	accountRepository := mocks.AccountRepository{}
	customerRepository := mocks.CustomerRepository{}

	accountService := services.Account{Repository: &accountRepository, Customer: &customerRepository}

	controller := Account{Service: &accountService}

	app := fiber.New()

	app.Post("/account/:number/transfer", controller.Transfer)

	resp, err := app.Test(httptest.NewRequest(fiber.MethodPost, "/account/abc/transfer", nil))

	utils.AssertEqual(t, nil, err)
	utils.AssertEqual(t, http.StatusBadRequest, resp.StatusCode)
}

func TestAccountTransferNotValidPayload(t *testing.T) {
	accountRepository := mocks.AccountRepository{}
	customerRepository := mocks.CustomerRepository{}

	accountService := services.Account{Repository: &accountRepository, Customer: &customerRepository}

	controller := Account{Service: &accountService}

	app := fiber.New()

	app.Post("/account/:number/transfer", controller.Transfer)

	data := map[string]interface{}{
		"amount": 100,
	}

	body, err := json.Marshal(data)

	utils.AssertEqual(t, nil, err)

	req := httptest.NewRequest(fiber.MethodPost, "/account/123/transfer", bytes.NewReader(body))
	req.Header.Add("content-type", "application/json")
	resp, err := app.Test(req)

	utils.AssertEqual(t, nil, err)
	utils.AssertEqual(t, http.StatusBadRequest, resp.StatusCode)
}

func TestAccountTransferZeroAmount(t *testing.T) {
	accountRepository := mocks.AccountRepository{}
	customerRepository := mocks.CustomerRepository{}

	accountService := services.Account{Repository: &accountRepository, Customer: &customerRepository}

	controller := Account{Service: &accountService}

	app := fiber.New()

	app.Post("/account/:number/transfer", controller.Transfer)

	data := map[string]interface{}{
		"to_account_number": 321,
		"amount":            0,
	}

	body, err := json.Marshal(data)

	utils.AssertEqual(t, nil, err)

	req := httptest.NewRequest(fiber.MethodPost, "/account/123/transfer", bytes.NewReader(body))
	req.Header.Add("content-type", "application/json")
	resp, err := app.Test(req)

	utils.AssertEqual(t, nil, err)
	utils.AssertEqual(t, http.StatusBadRequest, resp.StatusCode)
}

func TestAccountTransferReciverNotFound(t *testing.T) {
	sender := models.Account{}
	sender.ID = 123
	sender.Balance = 100001.0

	receiver := models.Account{}
	receiver.ID = 321

	accountRepository := mocks.AccountRepository{}
	accountRepository.On("Find", mock.MatchedBy(func(accountNumber int) bool {
		return accountNumber == sender.ID
	})).Return(&sender, nil).Once()
	accountRepository.On("Find", mock.MatchedBy(func(accountNumber int) bool {
		return accountNumber == receiver.ID
	})).Return(&receiver, errors.New("account not found")).Once()

	customerRepository := mocks.CustomerRepository{}

	accountService := services.Account{Repository: &accountRepository, Customer: &customerRepository}

	controller := Account{Service: &accountService}

	app := fiber.New()

	app.Post("/account/:number/transfer", controller.Transfer)

	data := map[string]interface{}{
		"to_account_number": 321,
		"amount":            100,
	}

	body, err := json.Marshal(data)

	utils.AssertEqual(t, nil, err)

	req := httptest.NewRequest(fiber.MethodPost, "/account/123/transfer", bytes.NewReader(body))
	req.Header.Add("content-type", "application/json")
	resp, err := app.Test(req)

	utils.AssertEqual(t, nil, err)
	utils.AssertEqual(t, http.StatusNotFound, resp.StatusCode)
}

func TestAccountTransferSuccess(t *testing.T) {
	sender := models.Account{}
	sender.ID = 123
	sender.Balance = 100001.0

	receiver := models.Account{}
	receiver.ID = 321

	accountRepository := mocks.AccountRepository{}
	accountRepository.On("Find", mock.MatchedBy(func(accountNumber int) bool {
		return accountNumber == sender.ID
	})).Return(&sender, nil).Once()
	accountRepository.On("Find", mock.MatchedBy(func(accountNumber int) bool {
		return accountNumber == receiver.ID
	})).Return(&receiver, nil).Once()
	accountRepository.On("Saves", mock.Anything, mock.Anything).Return(nil).Once()

	customerRepository := mocks.CustomerRepository{}

	accountService := services.Account{Repository: &accountRepository, Customer: &customerRepository}

	controller := Account{Service: &accountService}

	app := fiber.New()

	app.Post("/account/:number/transfer", controller.Transfer)

	data := map[string]interface{}{
		"to_account_number": 321,
		"amount":            100,
	}

	body, err := json.Marshal(data)

	utils.AssertEqual(t, nil, err)

	req := httptest.NewRequest(fiber.MethodPost, "/account/123/transfer", bytes.NewReader(body))
	req.Header.Add("content-type", "application/json")
	resp, err := app.Test(req)

	utils.AssertEqual(t, nil, err)
	utils.AssertEqual(t, http.StatusCreated, resp.StatusCode)
}
