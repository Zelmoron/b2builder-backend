package main

import (
	"log"
	"os"

	"github.com/gofiber/fiber/v2"
	"github.com/gofiber/fiber/v2/middleware/cors"
	"github.com/gofiber/fiber/v2/middleware/logger"

	"main/database"
	fbApp "main/firebase"
	"main/handler"
	"main/repository"
	"main/services"
)

func main() {
	fbApp.InitFirebase()

	dsn := os.Getenv("DATABASE_URL")

	db := database.InitDatabase(dsn)
	database.Migrate(db)

	repo := repository.NewRepository(db)
	service := services.NewService(repo)
	handlers := handler.NewHandler(service, repo)

	app := fiber.New()

	app.Use(logger.New())
	app.Use(cors.New())

	SetupRoutes(app, handlers)

	log.Fatal(app.Listen(":8080"))
}
