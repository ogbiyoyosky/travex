package dto

import (
	"fmt"

	"github.com/go-playground/validator/v10"
	"github.com/gofiber/fiber/v2"
)

type CreateLocationDto struct {
	Name         string `json:"name" validate:"required"`
	Image        string `json:"image" validate:"required"`
	Address      string `json:"address"  validate:"required"`
	LocationType string `json:"locationType" validate:"required"`
	Description  string `json:"description" validate:"required"`
}

func CreateLocationValidator(c *fiber.Ctx) error {
	var errors []*ErrorResponse
	body := new(CreateLocationDto)
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
