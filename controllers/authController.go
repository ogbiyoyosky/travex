package controller

import (
	"net/http"
	"time"

	"github.com/dgrijalva/jwt-go"
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
		First_name: data.FirstName,
		Last_name:  data.LastName,
		Email:      data.Email,
		Password:   hashPassword,
		Role:       "customer",
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

	//generate JWT token
	claims := jwt.NewWithClaims(jwt.SigningMethodHS256, jwt.StandardClaims{
		Issuer:    user.Id,
		ExpiresAt: time.Now().Add(time.Hour * 24).Unix(),
	})

	token, error := claims.SignedString([]byte("secret"))

	if error != nil {
		c.Status(http.StatusInternalServerError)
		return c.JSON(fiber.Map{
			"status":  false,
			"message": "internal server error",
		})
	}

	c.Status(http.StatusCreated)

	return c.JSON(fiber.Map{
		"status":  true,
		"message": "successfully created a user",
		"data": fiber.Map{
			"token": token,
			"user":  user,
		},
	})

}

func Login(c *fiber.Ctx) error {
	var data dto.RegisterDto

	if err := c.BodyParser(&data); err != nil {
		return err
	}

	var user models.User
	connection.DB.Where("email = ?", data.Email).First(&user)

	if user.Id == "" {
		c.Status(http.StatusBadRequest)
		return c.JSON(fiber.Map{
			"status":  false,
			"message": "invalid credentials",
		})
	}

	err := models.ComparePassword(user.Password, []byte(data.Password))

	if err != nil {
		c.Status(http.StatusBadRequest)
		return c.JSON(fiber.Map{
			"status":  false,
			"message": "invalid credentials",
		})
	}

	//generate JWT token
	claims := jwt.NewWithClaims(jwt.SigningMethodHS256, jwt.StandardClaims{
		Issuer:    user.Id,
		ExpiresAt: time.Now().Add(time.Hour * 24).Unix(),
	})

	token, error := claims.SignedString([]byte("secret"))

	if error != nil {
		c.Status(http.StatusInternalServerError)
		return c.JSON(fiber.Map{
			"status":  false,
			"message": "internal server error",
		})
	}

	// cookie := fiber.Cookie{
	// 	Name:     "jwt",
	// 	Value:    token,
	// 	Expires:  time.Now().Add(time.Hour * 24),
	// 	HTTPOnly: true,
	// }

	// c.Cookie(&cookie)

	return c.JSON(fiber.Map{
		"status":      true,
		"message":     "successfully logged in",
		"accessToken": token,
	})
}

func Logout(c *fiber.Ctx) error {

	cookie := fiber.Cookie{
		Name:     "jwt",
		Value:    "",
		Expires:  time.Now().Add(-time.Hour),
		HTTPOnly: true,
	}

	c.Cookie(&cookie)

	return c.JSON(fiber.Map{
		"message": "successfully logged out",
	})
}
