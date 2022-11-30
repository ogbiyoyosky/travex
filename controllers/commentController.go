package controller

import (
	"fmt"
	"net/http"
	"time"

	"github.com/gofiber/fiber/v2"
	connection "github.com/ogbiyoyosky/travex/db"
	"github.com/ogbiyoyosky/travex/dto"
	"github.com/ogbiyoyosky/travex/models"
)

func AddComment(c *fiber.Ctx) error {
	var data dto.CommentDto
	userObj := c.Locals("user").(models.User)
	var locationId = c.Params("locationId")
	var location models.Location

	if err := c.BodyParser(&data); err != nil {
		return err
	}
	if userObj.Role != "admin" {
		return c.JSON(fiber.Map{
			"status":  false,
			"message": "UnAuthorized",
		})
	}

	var comment models.Comment

	connection.DB.Model(&models.Location{
		Id: locationId,
	}).First(&location)

	if location.Id == "" {
		c.Status(http.StatusBadRequest)
		return c.JSON(fiber.Map{
			"status":  false,
			"message": "Location does not exist",
		})
	}

	comment = models.Comment{
		Location_id: locationId,
		Text:        data.Text,
		Author_id:   userObj.Id,
	}

	fmt.Println("locationId", locationId)
	fmt.Println("userObj.Id", userObj.Id)

	connection.DB.Omit("is_approved_at", "is_approved_by").Create(&comment)

	return c.JSON(fiber.Map{
		"status":  true,
		"message": "Successfully added a comment",
	})

}

func ApproveComment(c *fiber.Ctx) error {
	userObj := c.Locals("user").(models.User)
	var locationId = c.Params("locationId")
	var commentId = c.Params("commentId")
	var location models.Location

	var comment models.Comment

	connection.DB.Where("id = ?", locationId).Preload("Author").First(&location)

	if location.Id == "" {
		c.Status(http.StatusBadRequest)
		return c.JSON(fiber.Map{
			"status":  false,
			"message": "Location does not exist",
		})
	}

	connection.DB.Where("id = ? AND location_id = ?", commentId, locationId).Preload("Author").First(&comment)

	if comment.Id == "" {
		c.Status(http.StatusBadRequest)
		return c.JSON(fiber.Map{
			"status":  false,
			"message": "comment does not exist",
		})
	}

	if location.Id == "" {
		c.Status(http.StatusBadRequest)
		return c.JSON(fiber.Map{
			"status":  false,
			"message": "Location does not exist",
		})
	}

	if comment.IsApproved {
		return c.JSON(fiber.Map{
			"status":  false,
			"message": "Comment Already Approved",
		})
	}

	fmt.Println(" here", userObj.Id)
	fmt.Println(" userObj.Id", userObj.Id)

	comment.IsApproved = true
	comment.IsApprovedAt = time.Now()
	comment.IsApprovedBy = userObj.Id

	connection.DB.Omit("location_id, author_id", "text").Save(&comment)

	return c.JSON(fiber.Map{
		"status":  false,
		"message": "Successfully Approved Comment",
	})
}
