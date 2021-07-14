package controllers

import (
	"fmt"
	"net/http"
	"strconv"

	"github.com/ad3n/golang-testable/services"
	"github.com/ad3n/golang-testable/views"

	"github.com/go-playground/validator/v10"
	"github.com/gofiber/fiber/v2"
)

type Account struct {
	Service *services.Account
}

func (ctl *Account) Balance(c *fiber.Ctx) error {
	accountNumber, err := strconv.Atoi(c.Params("number"))
	if err != nil {
		c.JSON(map[string]string{
			"message": "account number is not number",
		})

		return c.SendStatus(http.StatusBadRequest)
	}

	result, err := ctl.Service.Balance(accountNumber)
	if err != nil {
		c.JSON(map[string]string{
			"message": err.Error(),
		})

		return c.SendStatus(http.StatusNotFound)
	}

	return c.JSON(result)
}

func (ctl *Account) Transfer(c *fiber.Ctx) error {
	accountNumber, err := strconv.Atoi(c.Params("number"))
	if err != nil {
		c.JSON(map[string]string{
			"message": "account number is not number",
		})

		return c.SendStatus(http.StatusBadRequest)
	}

	request := views.TransferRequest{}
	c.BodyParser(&request)

	validate := validator.New()
	err = validate.Struct(request)
	if err != nil {
		messages := []string{}
		for _, err := range err.(validator.ValidationErrors) {
			messages = append(messages, fmt.Sprintf("%s is %s", err.Field(), err.Tag()))
		}

		c.JSON(map[string]interface{}{
			"message": messages,
		})

		return c.SendStatus(http.StatusBadRequest)
	}

	err = ctl.Service.Transfer(accountNumber, request.Receiver, request.Amount)
	if err != nil {
		c.JSON(map[string]string{
			"message": err.Error(),
		})

		return c.SendStatus(http.StatusNotFound)
	}

	c.JSON("")

	return c.SendStatus(http.StatusCreated)
}
