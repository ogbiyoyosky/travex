package dto

import (
	"fmt"

	"github.com/go-playground/validator/v10"
	"github.com/gofiber/fiber/v2"
)

type CommentDto struct {
	Text string `json:"text" validate:"required"`
}

type ApproveCommentDto struct {
	IsApproved int `json:"isApproved" validate:"required"`
}

func CommentValidator(c *fiber.Ctx) error {
	var errors []*ErrorResponse
	body := new(CommentDto)
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

func ApproveCommentValidator(c *fiber.Ctx) error {
	var errors []*ErrorResponse
	body := new(ApproveCommentDto)
	c.BodyParser(&body)

	err := Validator.Struct(body)
	if err != nil {
		for _, err := range err.(validator.ValidationErrors) {

			fmt.Println(err.Param())
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
