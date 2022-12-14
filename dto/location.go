package dto

import (
	"fmt"
	"mime/multipart"

	"github.com/go-playground/validator/v10"
	"github.com/gofiber/fiber/v2"
)

type CreateLocationDto struct {
	Name         string                `json:"name" validate:"required"`
	Image        *multipart.FileHeader `json:"image" validate:"required"`
	Address      string                `json:"address"  validate:"required"`
	LocationType string                `json:"locationType" validate:"required"`
	Description  string                `json:"description" validate:"required"`
}

type CreateMasterLocationDto struct {
	Name         string                `json:"name" validate:"required"`
	Image        *multipart.FileHeader `json:"image" validate:"required"`
	Address      string                `json:"address"  validate:"required"`
	LocationType string                `json:"locationType" validate:"required"`
	Description  string                `json:"description" validate:"required"`
	UserId       string                `json:"userId" validate:"required"`
}

func CreateMasterValidator(c *fiber.Ctx) error {
	var errors []*ErrorResponse
	body := new(CreateMasterLocationDto)
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
