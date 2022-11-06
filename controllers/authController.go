package controller

import (
	"net/http"

	"github.com/gofiber/fiber/v2"
	connection "github.com/ogbiyoyosky/travex/db"
	"github.com/ogbiyoyosky/travex/dto"
	"github.com/ogbiyoyosky/travex/models"
)

func Register(c *fiber.Ctx) error {

	var data dto.RegisterDto

	if err := c.BodyParser(&data); err != nil {
		return err
	}

	hashPassword := models.HashPassword(data.Password)
	user := models.User{
		FirstName: data.FirstName,
		LastName:  data.LastName,
		Email:     data.Email,
		Password:  hashPassword,
		Role:      "customer",
	}

	connection.DB.Where("email = ?", data.Email).First(&user)

	if user.Id != "" {
		c.Status(http.StatusBadRequest)
		return c.JSON(fiber.Map{
			"status":  false,
			"message": "user already exists",
		})
	}
	connection.DB.Create(&user)

	c.Status(http.StatusCreated)

	return c.JSON(fiber.Map{
		"status":  true,
		"message": "successfully created a user",
		"data":    user,
	})

}
