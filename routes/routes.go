package routes

import (
	"github.com/gofiber/fiber/v2"
	controller "github.com/ogbiyoyosky/travex/controllers"
	"github.com/ogbiyoyosky/travex/dto"
	middleware "github.com/ogbiyoyosky/travex/middlewares"
)

func SetupRoutes(app *fiber.App) {
	app.Post("/api/register", dto.RegisterValidator, controller.Register)
	app.Post("/api/business/register", dto.RegisterValidator, controller.AdminRegister)
	app.Post("/api/login", dto.LoginValidator, controller.Login)
	app.Post("/api/logout", controller.Logout)
	app.Get("/api/locations", controller.GetLocations)
	app.Get("/api/locations/:locationId", middleware.ValidateJwt, controller.GetLocation)
	app.Get("/api/admin/locations/:locationId", middleware.ValidateJwt, controller.GetAdminLocation)
	app.Get("/api/profile/locations", middleware.ValidateJwt, controller.MyLocations)
	app.Post("/api/locations", middleware.ValidateJwt, controller.AddLocation)
	app.Post("/api/locations/:locationId/reviews", dto.ReviewValidator, middleware.ValidateJwt, controller.AddReview)
	app.Post("/api/locations/:locationId/comments", dto.CommentValidator, middleware.ValidateJwt, controller.AddComment)
	app.Post("/api/locations/:locationId/comments/:commentId/approve", middleware.ValidateJwt, controller.ApproveComment)
	app.Get("/api/profile", middleware.ValidateJwt, controller.GetProfile)
}
