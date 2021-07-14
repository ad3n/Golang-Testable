package main

import (
	"bytes"
	"encoding/json"
	"fmt"
	"io/ioutil"
	"net/http"
	"testing"

	"github.com/ad3n/golang-testable/configs"
	"github.com/ad3n/golang-testable/models"
	"github.com/ad3n/golang-testable/routes"
	"github.com/ad3n/golang-testable/views"
	"github.com/joho/godotenv"

	"github.com/gofiber/fiber/v2"
	"github.com/gofiber/fiber/v2/utils"
)

var app *fiber.App
var URL string

func setUp() {
	godotenv.Load()

	configs.Load()
	configs.Db.AutoMigrate(
		models.Customer{},
		models.Account{},
	)

	URL = fmt.Sprintf("http://localhost:%d", configs.Env.AppPort)
	app = fiber.New()

	(routes.Account{}).RegisterRoute(app)

	go app.Listen(fmt.Sprintf(":%d", configs.Env.AppPort))
}

func TestAccountBalanceAccountNotNumber(t *testing.T) {
	setUp()

	request, err := http.NewRequest(fiber.MethodGet, fmt.Sprintf("%s/account/abc", URL), nil)

	utils.AssertEqual(t, nil, err)

	client := &http.Client{}
	response, err := client.Do(request)

	utils.AssertEqual(t, nil, err)
	utils.AssertEqual(t, http.StatusBadRequest, response.StatusCode)
}

func TestAccountBalanceAccountNotFound(t *testing.T) {
	setUp()

	request, err := http.NewRequest(fiber.MethodGet, fmt.Sprintf("%s/account/123", URL), nil)

	utils.AssertEqual(t, nil, err)

	client := &http.Client{}
	response, err := client.Do(request)

	utils.AssertEqual(t, nil, err)

	body, err := ioutil.ReadAll(response.Body)

	utils.AssertEqual(t, nil, err)

	views := map[string]string{}
	err = json.Unmarshal(body, &views)

	utils.AssertEqual(t, nil, err)
	utils.AssertEqual(t, http.StatusNotFound, response.StatusCode)
	utils.AssertEqual(t, "account not found", views["message"])
}

func TestAccountBalanceSuccess(t *testing.T) {
	setUp()

	request, err := http.NewRequest(fiber.MethodGet, fmt.Sprintf("%s/account/555001", URL), nil)

	utils.AssertEqual(t, nil, err)

	client := &http.Client{}
	response, err := client.Do(request)

	utils.AssertEqual(t, nil, err)

	body, err := ioutil.ReadAll(response.Body)

	utils.AssertEqual(t, nil, err)

	views := views.BalanceResponse{}
	err = json.Unmarshal(body, &views)

	utils.AssertEqual(t, nil, err)
	utils.AssertEqual(t, http.StatusOK, response.StatusCode)
	utils.AssertEqual(t, 555001, views.AccountNumber)
}

func TestAccountTransferAccountNotNumber(t *testing.T) {
	setUp()

	request, err := http.NewRequest(fiber.MethodPost, fmt.Sprintf("%s/account/abc/transfer", URL), nil)

	utils.AssertEqual(t, nil, err)

	client := &http.Client{}
	response, err := client.Do(request)

	utils.AssertEqual(t, nil, err)
	utils.AssertEqual(t, http.StatusBadRequest, response.StatusCode)
}

func TestAccountTransferNotValidPayload(t *testing.T) {
	setUp()

	data := map[string]interface{}{
		"amount": 100,
	}

	body, err := json.Marshal(data)

	utils.AssertEqual(t, nil, err)

	request, err := http.NewRequest(fiber.MethodPost, fmt.Sprintf("%s/account/555001/transfer", URL), bytes.NewReader(body))
	request.Header.Add("content-type", "application/json")

	utils.AssertEqual(t, nil, err)

	client := &http.Client{}
	response, err := client.Do(request)

	utils.AssertEqual(t, nil, err)
	utils.AssertEqual(t, http.StatusBadRequest, response.StatusCode)
}

func TestAccountTransferZeroAmount(t *testing.T) {
	setUp()

	data := map[string]interface{}{
		"to_account_number": 123,
		"amount":            0,
	}

	body, err := json.Marshal(data)

	utils.AssertEqual(t, nil, err)

	request, err := http.NewRequest(fiber.MethodPost, fmt.Sprintf("%s/account/555001/transfer", URL), bytes.NewReader(body))
	request.Header.Add("content-type", "application/json")

	utils.AssertEqual(t, nil, err)

	client := &http.Client{}
	response, err := client.Do(request)

	utils.AssertEqual(t, nil, err)
	utils.AssertEqual(t, http.StatusBadRequest, response.StatusCode)
}

func TestAccountTransferReciverNotFound(t *testing.T) {
	setUp()

	data := map[string]interface{}{
		"to_account_number": 123,
		"amount":            100,
	}

	body, err := json.Marshal(data)

	utils.AssertEqual(t, nil, err)

	request, err := http.NewRequest(fiber.MethodPost, fmt.Sprintf("%s/account/555001/transfer", URL), bytes.NewReader(body))
	request.Header.Add("content-type", "application/json")

	utils.AssertEqual(t, nil, err)

	client := &http.Client{}
	response, err := client.Do(request)

	utils.AssertEqual(t, nil, err)
	utils.AssertEqual(t, http.StatusNotFound, response.StatusCode)
}

func TestAccountTransferSuccess(t *testing.T) {
	setUp()

	data := map[string]interface{}{
		"to_account_number": 555002,
		"amount":            100,
	}

	body, err := json.Marshal(data)

	utils.AssertEqual(t, nil, err)

	request, err := http.NewRequest(fiber.MethodPost, fmt.Sprintf("%s/account/555001/transfer", URL), bytes.NewReader(body))
	request.Header.Add("content-type", "application/json")

	utils.AssertEqual(t, nil, err)

	client := &http.Client{}
	response, err := client.Do(request)

	utils.AssertEqual(t, nil, err)
	utils.AssertEqual(t, http.StatusCreated, response.StatusCode)
}
