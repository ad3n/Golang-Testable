package main

import (
	"github.com/ad3n/golang-testable/configs"
	"github.com/ad3n/golang-testable/models"
	"github.com/ad3n/golang-testable/repositories"

	"github.com/joho/godotenv"
)

func init() {
	godotenv.Load()
	configs.Load()
}

func main() {
	customerRepository := repositories.Customer{Storage: configs.Db}
	accountRepository := repositories.Account{Storage: configs.Db}

	customer1 := models.Customer{}
	customer1.ID = 1001
	customer1.Name = "Bob Martin"

	customer2 := models.Customer{}
	customer2.ID = 1002
	customer2.Name = "Linus Torvalds"

	customerRepository.Saves(&customer1, &customer2)

	account1 := models.Account{}
	account1.ID = 555001
	account1.CustomerID = customer1.ID
	account1.Balance = 10000

	account2 := models.Account{}
	account2.ID = 555002
	account2.CustomerID = customer2.ID
	account2.Balance = 15000

	accountRepository.Saves(&account1, &account2)
}
