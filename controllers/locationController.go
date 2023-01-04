package controller

import (
	"fmt"
	"net/http"

	"github.com/gofiber/fiber/v2"
	connection "github.com/ogbiyoyosky/travex/db"
	"github.com/ogbiyoyosky/travex/models"
)

func GetLocations(c *fiber.Ctx) error {
	var locations []models.Location

	var search = c.Query("search")

	if search != "" {

		connection.DB.Model(&models.Location{}).Where("name ILIKE ?", "%"+search+"%").Preload("User").Preload("Comments.Author").Preload("User").Preload("Comments", "is_approved NOT IN (?)", false).Preload("Reviews.Author").Preload("LocationType").Find(&locations)
		c.Status(http.StatusOK)
		fmt.Println("locations", locations)
		return c.JSON(fiber.Map{
			"status":  true,
			"message": "Successfully retrieved locations",
			"data":    locations,
		})
	}

	connection.DB.Model(&models.Location{}).Preload("Comments.Author").Preload("User").Preload("Comments", "is_approved NOT IN (?)", false).Preload("LocationType").Preload("Reviews.Author").Find(&locations)
	c.Status(http.StatusOK)

	return c.JSON(fiber.Map{
		"status":  true,
		"message": "Successfully retrieved locations",
		"data":    locations,
	})
}

func GetLocation(c *fiber.Ctx) error {
	var location models.Location
	//var testLocation models.TestLocation
	var locationId = c.Params("locationId")

	connection.DB.Model(&models.Location{}).Where("id = ?", locationId).Preload("LocationType").Preload("Comments.Author").Preload("User").Preload("Comments", "is_approved NOT IN (?)", false).Preload("Reviews.Author").Preload("User").First(&location)
	connection.DB.Raw("SELECT c.text, c.parent_id FROM locations LEFT JOIN comments ON locations.id = comments.location_id LEFT JOIN comments as c ON c.parent_id = comments.id WHERE comments.location_id = ? GROUP BY c.text,c.parent_id", locationId).Scan(&location)

	fmt.Println(location)

	if location.Id == "" {
		c.Status(http.StatusBadRequest)
		return c.JSON(fiber.Map{
			"status":  false,
			"message": "Location does not exist",
		})
	}
	c.Status(http.StatusOK)

	return c.JSON(fiber.Map{
		"status":  true,
		"message": "Successfully retrieved location",
		"data":    location,
	})
}
