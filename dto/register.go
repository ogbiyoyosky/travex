package dto

import (
	"fmt"

	"github.com/go-playground/validator/v10"
	"github.com/gofiber/fiber/v2"
)

type RegisterDto struct {
	FirstName string `json:"firstName" validate:"required"`
	LastName  string `json:"lastName" validate:"required"`
	Email     string `json:"email"  validate:"required,email"`
	Password  string `json:"password" validate:"required"`
}

func RegisterValidator(c *fiber.Ctx) error {
	var errors []*ErrorResponse
	body := new(RegisterDto)
	c.BodyParser(&body)

	err := Validator.Struct(body)
	if err != nil {
		for _, err := range err.(validator.ValidationErrors) {
			var el ErrorResponse
			el.Property = err.Field()
			el.Tag = err.Tag()
			el.Value = err.Param()
			el.Message = fmt.Sprintf("Invalid property %s", err.Field())
			errors = append(errors, &el)
		}

		return c.Status(fiber.StatusBadRequest).JSON(errors)
	}

	return c.Next()
}
