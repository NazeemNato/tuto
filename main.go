package main

import (
	"log"

	"github.com/NazeemNato/tuto/src/database"
	"github.com/NazeemNato/tuto/src/routes"
	"github.com/gofiber/fiber/v2"
	"github.com/gofiber/fiber/v2/middleware/cors"
)

func main() {

	database.Connnect()
	database.AutoMigrate()
	database.SetupRedis()
	
	app := fiber.New()

	app.Use(cors.New(cors.Config{AllowCredentials: true}))

	routes.Setup(app)

	log.Fatal(app.Listen(":8000"))
}
