package handlers

import (
	"github.com/NazeemNato/tuto/src/database"
	"github.com/NazeemNato/tuto/src/models"
	"github.com/gofiber/fiber/v2"
	"github.com/stripe/stripe-go/v72"
	"github.com/stripe/stripe-go/v72/checkout/session"
)

func Orders(c *fiber.Ctx) error {
	var orders []models.Order
	database.DB.Preload("OrderItems").Find(&orders)

	for i, order := range orders {
		orders[i].Name = order.FulName()
		orders[i].Total = order.GetTotal()
	}

	return c.JSON(orders)
}

type CreateOrderRequest struct {
	Code      string
	FirstName string
	LastName  string
	Email     string
	Address   string
	Country   string
	City      string
	Zip       string
	Products  []map[string]int
}

func CreateOrder(c *fiber.Ctx) error {
	var request CreateOrderRequest
	if err := c.BodyParser(&request); err != nil {
		return err
	}

	link := models.Link{Code: request.Code}
	tx := database.DB.Begin()
	tx.Preload("User").First(&link)

	if link.Id == 0 {
		return c.Status(fiber.StatusBadRequest).JSON(fiber.Map{"message": "Invalid link"})
	}
	order := models.Order{
		Code:            link.Code,
		UserId:          link.Id,
		AmbassadorEmail: link.User.Email,
		FirstName:       request.FirstName,
		LastName:        request.LastName,
		Email:           request.Email,
		Address:         request.Address,
		Country:         request.Country,
		City:            request.City,
		Zip:             request.Zip,
	}

	if err := tx.Create(&order).Error; err != nil {
		tx.Rollback()
		return c.Status(fiber.StatusBadRequest).JSON(fiber.Map{"message": err.Error()})
	}

	var lineItems []*stripe.CheckoutSessionLineItemParams

	for _, requestProduct := range request.Products {
		product := models.Product{}
		product.Id = uint(requestProduct["product_id"])
		tx.First(&product)
		total := product.Price * float64(requestProduct["qty"])
		item := models.OrderItem{
			OrderId:           order.Id,
			ProductTitle:      product.Title,
			Price:             product.Price,
			Quantity:          uint(requestProduct["qty"]),
			AmbassadorRevenue: 0.1 * total,
			AdminRevenue:      0.9 * total,
		}

		if err := tx.Create(&item).Error; err != nil {
			tx.Rollback()
			return c.Status(fiber.StatusBadRequest).JSON(fiber.Map{"message": err.Error()})
		}

		lineItems = append(lineItems, &stripe.CheckoutSessionLineItemParams{
			Name:        stripe.String(product.Title),
			Description: stripe.String(product.Description),
			Images:      []*string{stripe.String(product.Image)},
			Amount:      stripe.Int64(100 * int64(product.Price + 10)),
			Currency:    stripe.String("usd"),
			Quantity:    stripe.Int64(int64(requestProduct["qty"])),
		})
	}

	stripe.Key = "sk_test_51F8zwQKpXlWW10jvORV427XSS2DFTjP9Av4A5UxJW7EeNoVRo79NxvpnAGENpRagQBoe8I7dLA14cLst5FV1mAHR00MbZN9jCz"
	params := stripe.CheckoutSessionParams{
		SuccessURL: stripe.String("http://localhost:5000/sucess?source={CHECKOUT_SESSION_ID}"),
		CancelURL: stripe.String("http://localhost:5000/error"),
		PaymentMethodTypes: stripe.StringSlice([]string{"card"}),
		LineItems: lineItems,
	}

	source ,err := session.New(&params)

	if err != nil {
		tx.Rollback()
		return c.Status(fiber.StatusBadRequest).JSON(fiber.Map{"message": err.Error()})
	}

	order.TransactionId = source.ID

	tx.Commit()

	return c.JSON(order)
}
