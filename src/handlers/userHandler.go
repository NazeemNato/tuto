package handlers

import (
	"context"

	"github.com/NazeemNato/tuto/src/database"
	"github.com/NazeemNato/tuto/src/models"
	"github.com/go-redis/redis/v8"
	"github.com/gofiber/fiber/v2"
)

func Ambassador(c *fiber.Ctx) error {
	var user []models.User

	database.DB.Where("is_ambassador=true").Find(&user)

	return c.JSON(user)
}

func Rankings (c *fiber.Ctx) error {

	rankings , err := database.Cache.ZRangeByScoreWithScores(context.Background(),"rankings", &redis.ZRangeBy{
		Min: "-inf",
		Max: "+inf",
	}).Result()

	if err != nil {
		return err
	}

	result := make(map[string]float64)

	for _, ranking := range rankings {
		result[ranking.Member.(string)] = ranking.Score
	}

	return c.JSON(result)
}