package dto

import (
	"fmt"

	"github.com/go-playground/validator/v10"
	"github.com/gofiber/fiber/v2"
)

type ReviewDto struct {
	Rating  uint16 `json:"rating" validate:"required,min=1,max=10"`
	Comment string `json:"comment" validate:"required"`
}

func ReviewValidator(c *fiber.Ctx) error {
	var errors []*ErrorResponse
	body := new(ReviewDto)
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
