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
	var reviewId = c.Params("reviewId")
	var location models.Location

	if err := c.BodyParser(&data); err != nil {
		return err
	}

	var comment models.Comment
	var review models.Review

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

	connection.DB.Model(&models.Review{
		Id: reviewId,
	}).First(&review)

	if review.Id == "" {
		c.Status(http.StatusBadRequest)
		return c.JSON(fiber.Map{
			"status":  false,
			"message": "Review does not exist",
		})
	}

	comment = models.Comment{
		Location_id: locationId,
		Text:        data.Text,
		Author_id:   userObj.Id,
		Review_id:   reviewId,
	}

	if userObj.Role != "admin" {
		connection.DB.Omit("is_approved_at", "is_approved_by").Create(&comment)

	} else {
		comment.IsApproved = true
		comment.IsApprovedAt = time.Now()
		comment.IsApprovedBy = userObj.Id

		connection.DB.Omit("location_id, author_id", "parent_id").Save(&comment)
	}

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

	var data dto.ApproveCommentDto

	var comment models.Comment

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

	comment.IsApproved = data.IsApproved

	if data.IsApproved {
		comment.IsApprovedAt = time.Now()
		comment.IsApprovedBy = userObj.Id
	}

	connection.DB.Omit("location_id, author_id", "text").Save(&comment)

	return c.JSON(fiber.Map{
		"status":  false,
		"message": "Successfully Approved Comment",
	})
}
