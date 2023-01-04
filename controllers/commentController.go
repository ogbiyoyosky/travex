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

	commentParentId := data.Parent_id

	if commentParentId != "" {
		connection.DB.Model(&comment).Where("id = ?", commentParentId).First(&comment)

		if comment.Id == "" {

			c.Status(http.StatusNotFound)
			return c.JSON(fiber.Map{
				"status":  false,
				"message": "Can not reply to a comment that does not exist",
			})
		}

		comment = models.Comment{
			Location_id: locationId,
			Text:        data.Text,
			Author_id:   userObj.Id,
			Parent_id:   commentParentId,
		}

		if userObj.Role != "admin" {
			connection.DB.Omit("is_approved_at", "is_approved_by").Create(&comment)

		} else {
			comment.IsApproved = true
			comment.IsApprovedAt = time.Now()
			comment.IsApprovedBy = userObj.Id

			connection.DB.Save(&comment)
		}

	} else {
		comment = models.Comment{
			Location_id: locationId,
			Text:        data.Text,
			Author_id:   userObj.Id,
		}

		if userObj.Role != "admin" {
			connection.DB.Omit("is_approved_at", "is_approved_by", "parent_id").Create(&comment)

		} else {
			comment.IsApproved = true
			comment.IsApprovedAt = time.Now()
			comment.IsApprovedBy = userObj.Id

			connection.DB.Omit("location_id, author_id", "parent_id").Save(&comment)
		}

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

	comment.IsApproved = true
	comment.IsApprovedAt = time.Now()
	comment.IsApprovedBy = userObj.Id

	connection.DB.Omit("location_id, author_id", "text").Save(&comment)

	return c.JSON(fiber.Map{
		"status":  false,
		"message": "Successfully Approved Comment",
	})
}
