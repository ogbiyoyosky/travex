package routes

import (
	"github.com/gofiber/fiber/v2"
	controller "github.com/ogbiyoyosky/travex/controllers"
	"github.com/ogbiyoyosky/travex/dto"
)

func SetupRoutes(app *fiber.App) {
	app.Post("/api/register", dto.RegisterValidator, controller.Register)
	app.Post("/api/login", dto.RegisterValidator, controller.Register)
}
