package controller

import (
	"fmt"
	"net/http"

	"github.com/gofiber/fiber/v2"
	connection "github.com/ogbiyoyosky/travex/db"
	"github.com/ogbiyoyosky/travex/dto"
	"github.com/ogbiyoyosky/travex/models"
)

func AddReview(c *fiber.Ctx) error {
	var data dto.ReviewDto
	userObj := c.Locals("user").(models.User)
	var locationId = c.Params("locationId")

	if err := c.BodyParser(&data); err != nil {
		return err
	}

	var review models.Review

	var location models.Location

	connection.DB.Where("id = ?", locationId).First(&location)

	if location.Id == "" {
		c.Status(http.StatusBadRequest)
		return c.JSON(fiber.Map{
			"status":  false,
			"message": "Location does not exist",
		})
	}

	if location.UserId == userObj.Id {
		c.Status(http.StatusBadRequest)
		return c.JSON(fiber.Map{
			"status":  false,
			"message": "you can't add a review for your location you created",
		})
	}

	connection.DB.Where("author_id = ? AND location_id = ?", userObj.Id, locationId).First(&review)

	fmt.Println("data.Rating", data.Rating/10)

	if review.Author_id != "" {
		c.Status(http.StatusBadRequest)
		return c.JSON(fiber.Map{
			"status":  false,
			"message": "You already given a review for this location",
		})
	}

	if data.Comment != "" {
		review = models.Review{
			Location_id: locationId,
			Rating:      float32(data.Rating) / 10,
			Author_id:   userObj.Id,
			Text:        data.Comment,
			IsApproved:  false,
		}

		connection.DB.Create(&review)

	} else {
		review = models.Review{
			Location_id: locationId,
			Rating:      float32(data.Rating) / 10,
			Author_id:   userObj.Id,
			IsApproved:  true,
		}

		connection.DB.Omit("text").Save(&review)

	}

	return c.JSON(fiber.Map{
		"status":  true,
		"message": "You added a review for this location",
	})

}

func ApproveReview(c *fiber.Ctx) error {
	userObj := c.Locals("user").(models.User)
	var locationId = c.Params("locationId")

	var location models.Location

	var review models.Review

	var data dto.ApproveReviewDto
	if err := c.BodyParser(&data); err != nil {
		return err
	}

	connection.DB.Where("id = ?", locationId).Preload("Author").First(&location)

	if location.Id == "" {
		c.Status(http.StatusBadRequest)
		return c.JSON(fiber.Map{
			"status":  false,
			"message": "Location does not exist",
		})
	}

	if location.UserId != userObj.Id {
		c.Status(http.StatusForbidden)
		return c.JSON(fiber.Map{
			"status":  false,
			"message": "You don't have the permission to carry out this operation",
		})
	}

	connection.DB.Where("id = ? AND location_id = ? ", data.ReviewId, locationId).First(&review)

	if review.Id == "" {
		c.Status(http.StatusBadRequest)
		return c.JSON(fiber.Map{
			"status":  false,
			"message": "Review does not exist",
		})
	}

	if data.IsApproved == 1 {
		review.IsApproved = true
	} else {
		review.IsApproved = false
	}

	connection.DB.Omit("location_id, author_id", "text").Save(&review)

	return c.JSON(fiber.Map{
		"status":  false,
		"message": "Successfully updated Review",
	})
}
