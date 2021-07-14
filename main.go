package main

import (
	"fmt"

	"github.com/ad3n/golang-testable/configs"
	"github.com/ad3n/golang-testable/models"
	"github.com/ad3n/golang-testable/routes"

	"github.com/gofiber/fiber/v2"
	"github.com/joho/godotenv"
)

func init() {
	godotenv.Load()

	configs.Load()
	configs.Db.AutoMigrate(
		models.Customer{},
		models.Account{},
	)
}

func main() {
	app := fiber.New()

	(routes.Account{}).RegisterRoute(app)

	app.Listen(fmt.Sprintf(":%d", configs.Env.AppPort))
}
