package middlewares

import (
	"strconv"
	"strings"
	"time"

	"github.com/dgrijalva/jwt-go"
	"github.com/gofiber/fiber/v2"
)

const PrivateKey = "buckthorn"

type CliamsWithScope struct {
	jwt.StandardClaims
	Scope string
}

func IsAuthenticated(c *fiber.Ctx) error {
	cookie := c.Cookies("access_token")
	token, err := jwt.ParseWithClaims(cookie, &CliamsWithScope{}, func(t *jwt.Token) (interface{}, error) {
		return []byte(PrivateKey), nil
	})

	if err != nil || !token.Valid {
		return c.Status(fiber.StatusUnauthorized).JSON(fiber.Map{"message": "Invalid access token"})
	}

	payload := token.Claims.(*CliamsWithScope)

	IsAmbassador := strings.Contains(c.Path(), "/api/ambassador")

	if (payload.Scope == "admin" && IsAmbassador) || (payload.Scope == "ambassador" && !IsAmbassador) {
		return c.Status(fiber.StatusUnauthorized).JSON(fiber.Map{"message": "Unauthorized"})

	}

	return c.Next()
}

func GenerateJWT(id uint, scope string) (string, error) {
	var payload = CliamsWithScope{}

	payload.Subject = strconv.Itoa(int(id))
	payload.ExpiresAt = time.Now().Add(time.Hour * 24).Unix()
	payload.Scope = scope

	return jwt.NewWithClaims(jwt.SigningMethodHS256, payload).SignedString([]byte(PrivateKey))
}

func GetUserId(c *fiber.Ctx) (uint, error) {
	cookie := c.Cookies("access_token")
	token, err := jwt.ParseWithClaims(cookie, &CliamsWithScope{}, func(t *jwt.Token) (interface{}, error) {
		return []byte(PrivateKey), nil
	})

	if err != nil {
		return 0, err
	}

	payload := token.Claims.(*CliamsWithScope)
	id, _ := strconv.Atoi(payload.Subject)
	return uint(id), nil
}
