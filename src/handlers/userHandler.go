package handlers

import (
	"github.com/NazeemNato/tuto/src/database"
	"github.com/NazeemNato/tuto/src/models"
	"github.com/gofiber/fiber/v2"
)

func Ambassador(c *fiber.Ctx) error {
	var user []models.User

	database.DB.Where("is_ambassador=true").Find(&user)

	return c.JSON(user)
}