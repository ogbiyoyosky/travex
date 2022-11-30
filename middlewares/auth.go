package middleware

import (
	"fmt"
	"net/http"

	"github.com/dgrijalva/jwt-go"
	"github.com/gofiber/fiber/v2"
	connection "github.com/ogbiyoyosky/travex/db"
	"github.com/ogbiyoyosky/travex/models"
)

func ValidateJwt(c *fiber.Ctx) error {
	cookie := c.Cookies("jwt")

	token, error := jwt.ParseWithClaims(cookie, &jwt.StandardClaims{}, func(token *jwt.Token) (interface{}, error) {
		return []byte("secret"), nil
	})

	if error != nil {
		c.Status(http.StatusUnauthorized)
		return c.JSON(fiber.Map{
			"message": "unauthorized",
			"status":  false,
		})
	}

	claims := token.Claims.(*jwt.StandardClaims)

	var user models.User

	connection.DB.Where("id = ?", claims.Issuer).First(&user)

	if user.Id == "" {
		c.Status(http.StatusUnauthorized)
		return c.JSON(fiber.Map{
			"message": "unauthorized",
			"status":  false,
		})
	}

	c.Locals("user", user)

	fmt.Println("user", c.Locals("user"))

	return c.Next()
}
