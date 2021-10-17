package routes

import (
	"github.com/NazeemNato/tuto/src/handlers"
	"github.com/NazeemNato/tuto/src/middlewares"
	"github.com/gofiber/fiber/v2"
)

func Setup(app *fiber.App) {
	api := app.Group("api")

	// admin
	admin := api.Group("admin")
	admin.Post("register", handlers.Register)
	admin.Post("login", handlers.Login)
	adminAuthenticaed := admin.Use(middlewares.IsAuthenticated)
	adminAuthenticaed.Get("user", handlers.User)
	adminAuthenticaed.Put("user/info", handlers.UpdateInfo)
	adminAuthenticaed.Put("user/password", handlers.UpdatePassword)
	adminAuthenticaed.Post("logout", handlers.Logout)
}
