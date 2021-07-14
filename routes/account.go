package routes

import (
	"github.com/ad3n/golang-testable/configs"
	"github.com/ad3n/golang-testable/controllers"
	"github.com/ad3n/golang-testable/repositories"
	"github.com/ad3n/golang-testable/services"

	"github.com/gofiber/fiber/v2"
)

type Account struct {
}

func (Account) RegisterRoute(router fiber.Router) {
	accountRepository := repositories.Account{Storage: configs.Db}
	customerRepository := repositories.Customer{Storage: configs.Db}

	accountService := services.Account{Repository: &accountRepository, Customer: &customerRepository}

	account := controllers.Account{Service: &accountService}

	router.Get("/account/:number", account.Balance)
	router.Post("/account/:number/transfer", account.Transfer)
}
