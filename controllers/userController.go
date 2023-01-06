package controller

import (
	"fmt"

	"github.com/gofiber/fiber/v2"
	connection "github.com/ogbiyoyosky/travex/db"
	"github.com/ogbiyoyosky/travex/models"
)

func GetProfile(c *fiber.Ctx) error {
	fmt.Println("Called profile")
	var user models.User
	userObj := c.Locals("user").(models.User)

	connection.DB.Where("id = ?", userObj.Id).Preload("Locations").First(&user)

	return c.JSON(fiber.Map{
		"status":  true,
		"message": "Successfully fetched profile",
		"data":    user,
	})
}
