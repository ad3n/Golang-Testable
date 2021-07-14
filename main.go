package main

import (
	"fmt"
	"log"
	"net"

	"github.com/ad3n/golang-testable/configs"
	"github.com/ad3n/golang-testable/grpcs"
	"github.com/ad3n/golang-testable/models"
	"github.com/ad3n/golang-testable/protos"
	"github.com/ad3n/golang-testable/routes"

	"github.com/gofiber/fiber/v2"
	"github.com/joho/godotenv"
	"google.golang.org/grpc"
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

	lis, err := net.Listen("tcp", fmt.Sprintf(":%d", configs.Env.GRpcPort))
	if err != nil {
		log.Fatalf("failed to listen: %v", err)
	}

	account := grpcs.Account{}

	grpcServer := grpc.NewServer()
	protos.RegisterAccountServer(grpcServer, &account)

	go grpcServer.Serve(lis)

	app.Listen(fmt.Sprintf(":%d", configs.Env.AppPort))
}
