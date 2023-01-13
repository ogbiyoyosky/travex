package controller

import (
	"context"
	"fmt"
	"io"
	"io/ioutil"
	"log"
	"net/http"
	"os"
	"time"

	"github.com/cloudinary/cloudinary-go"
	"github.com/cloudinary/cloudinary-go/api/uploader"
	"github.com/gofiber/fiber/v2"
	connection "github.com/ogbiyoyosky/travex/db"
	"github.com/ogbiyoyosky/travex/dto"
	"github.com/ogbiyoyosky/travex/models"
)

type UploadImage struct {
	log *log.Logger
}

func NewUploadAppImage(log *log.Logger) *UploadImage {
	return &UploadImage{log}
}

func GetLocations(c *fiber.Ctx) error {
	var locations []models.Location

	var search = c.Query("search")

	if search != "" {

		connection.DB.Model(&models.Location{}).Where("locations.name ILIKE ? AND locations.is_approved is ?", "%"+search+"%", true).Or("location_types.name LIKE ?", "%"+search+"%").Joins("JOIN location_types ON location_types.id = locations.location_type_id").Preload("User").Preload("Reviews.Comments", "is_approved = ?", true).Preload("Reviews.Author").Preload("LocationType").Find(&locations)
		c.Status(http.StatusOK)
		fmt.Println("locations", locations)
		return c.JSON(fiber.Map{
			"status":  true,
			"message": "Successfully retrieved locations",
			"data":    locations,
		})
	}

	connection.DB.Model(&models.Location{}).Preload("Reviews.Author").Preload("User").Preload("Reviews", "is_approved = ?", true).Preload("Reviews.Comments", "is_approved = ?", true).Preload("Reviews.Comments.Author").Preload("LocationType").Find(&locations)

	c.Status(http.StatusOK)

	return c.JSON(fiber.Map{
		"status":  true,
		"message": "Successfully retrieved locations",
		"data":    locations,
	})
}

func ApproveLocation(c *fiber.Ctx) error {

	var locationId = c.Params("locationId")

	var location models.Location

	var data dto.ApproveLocationDto

	if err := c.BodyParser(&data); err != nil {
		return err
	}

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

	location.IsApproved = data.IsApproved

	if data.IsApproved {
		location.IsApprovedAt = time.Now()
	}

	connection.DB.Omit("name, image", "address", "location_type_id", "description", "user_id").Save(&location)

	return c.JSON(fiber.Map{
		"status":  true,
		"message": "Successfully approved location",
	})
}

func MyLocations(c *fiber.Ctx) error {
	var locations []models.Location

	userObj := c.Locals("user").(models.User)

	var search = c.Query("search")

	if search != "" {

		connection.DB.Model(&models.Location{}).Where("locations.name ILIKE ?", "%"+search+"%").Or("location_types.name LIKE ?", "%"+search+"%").Joins("JOIN location_types ON location_types.id = locations.location_type_id").Preload("User").Preload("Reviews.Author").Preload("User").Preload("Reviews.Comment").Preload("Reviews.Author").Preload("Reviews.Comments").Preload("Reviews.Comment.Author").Preload("LocationType").Where("user_id = ?", userObj.Id).Find(&locations)
		c.Status(http.StatusOK)
		fmt.Println("locations", locations)
		return c.JSON(fiber.Map{
			"status":  true,
			"message": "Successfully retrieved locations",
			"data":    locations,
		})
	}

	connection.DB.Model(&models.Location{}).Preload("Comments.Author").Preload("Reviews", "is_approved = ?", false).Preload("User").Preload("Comments", "is_approved = ?", true).Preload("LocationType").Preload("Reviews.Author").Where("user_id = ?", userObj.Id).Find(&locations)
	c.Status(http.StatusOK)

	return c.JSON(fiber.Map{
		"status":  true,
		"message": "Successfully retrieved locations",
		"data":    locations,
	})
}

func GetLocation(c *fiber.Ctx) error {
	var location models.Location
	userObj := c.Locals("user").(models.User)
	var locationId = c.Params("locationId")

	fmt.Println(userObj.Role == "admin")
	if userObj.Role == "admin" {
		connection.DB.Model(&models.Location{}).Where("id = ?", locationId).Preload("LocationType").Preload("User").Preload("Reviews.Comments", "is_approved = ?", true).Preload("Reviews.Author").Preload("User").Preload("Reviews.Comments.Author").First(&location)
	} else {

		connection.DB.Model(&models.Location{}).Where("id = ?", locationId).Preload("LocationType").Preload("Reviews", "is_approved = ?", false).Preload("Reviews.Author").Preload("User").Preload("Reviews.Comments", "is_approved = ?", false).Preload("Reviews.Comments.Author").Preload("Reviews.Author").Preload("User").First(&location)
	}
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

func GetAdminLocation(c *fiber.Ctx) error {
	var location models.Location
	//var testLocation models.TestLocation
	var locationId = c.Params("locationId")

	//connection.DB.Model(&models.Location{}).Where("id = ?", locationId).Preload("LocationType").Preload("Comments.Author").Preload("User").Preload("Comments", false).Preload("Reviews.Author").Preload("User").First(&location)
	connection.DB.Raw("SELECT c.text, c.parent_id FROM locations LEFT JOIN comments ON locations.id = comments.location_id LEFT JOIN comments as c ON c.parent_id = comments.id WHERE comments.location_id = ? GROUP BY c.text,c.parent_id", locationId).Scan(&location)

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

func AddLocation(c *fiber.Ctx) error {
	userObj := c.Locals("user").(models.User)

	var location models.Location

	var locationType models.LocationType

	connection.DB.Where("name = ? ", c.FormValue("locationType")).First(&locationType)

	if c.FormValue("locationType") == "" {
		c.Status(http.StatusBadRequest)
		return c.JSON(fiber.Map{
			"status":  false,
			"message": "LocationType does not exist",
		})
	}

	fileheader, err := c.FormFile("image")
	if err != nil {
		panic(err)
	}

	file, err := fileheader.Open()

	if err != nil {
		panic(err)
	}
	defer file.Close()

	buffer, err := io.ReadAll(file)

	if err != nil {
		panic(err)
	}

	fileName, err := UploadAppImage(fileheader.Filename, buffer)

	location = models.Location{
		Name:             c.FormValue("name"),
		Image:            fileName,
		Address:          c.FormValue("address"),
		Location_type_id: locationType.Id,
		Description:      c.FormValue("description"),
		UserId:           userObj.Id,
	}

	connection.DB.Omit("is_approved", "is_approved_at").Create(&location)

	return c.JSON(fiber.Map{
		"status":  true,
		"message": "Successfully created Location",
	})

}

func AddMasterAdminLocation(c *fiber.Ctx) error {
	userObj := c.Locals("user").(models.User)

	if userObj.Role != "master_admin" {
		c.Status(http.StatusBadRequest)
		return c.JSON(fiber.Map{
			"status":  false,
			"message": "No suffcient permission",
		})
	}

	var location models.Location

	var locationType models.LocationType

	connection.DB.Where("name = ? ", c.FormValue("locationType")).First(&locationType)

	if c.FormValue("locationType") == "" {
		c.Status(http.StatusBadRequest)
		return c.JSON(fiber.Map{
			"status":  false,
			"message": "LocationType does not exist",
		})
	}

	fileheader, err := c.FormFile("image")
	if err != nil {
		panic(err)
	}

	file, err := fileheader.Open()

	if err != nil {
		panic(err)
	}
	defer file.Close()

	buffer, err := io.ReadAll(file)

	if err != nil {
		panic(err)
	}

	fileName, _ := UploadAppImage(fileheader.Filename, buffer)

	location = models.Location{
		Name:             c.FormValue("name"),
		Image:            fileName,
		Address:          c.FormValue("address"),
		Location_type_id: locationType.Id,
		Description:      c.FormValue("description"),
		UserId:           c.FormValue("userId"),
	}

	connection.DB.Omit("is_approved", "is_approved_at").Create(&location)

	return c.JSON(fiber.Map{
		"status":  true,
		"message": "Successfully created Location",
	})

}

func UploadAppImage(stringifiedFileName string, fileBuff []byte) (string, error) {
	tempDir, err := ioutil.TempDir("", "uploads")
	if err != nil {
		return "", err
	}

	file, err := ioutil.TempFile(tempDir, stringifiedFileName)
	if err != nil {
		return "", err
	}

	file.Write(fileBuff)

	result, err := Upload(file.Name())

	if err != nil {
		return "", err
	}
	os.RemoveAll(tempDir)

	return result, nil
}

func Upload(filePath string) (string, error) {
	var ctx = context.Background()
	cloudinaryUrl := os.Getenv("CLOUDINARY_URL")
	fmt.Println(cloudinaryUrl)
	cld, _ := cloudinary.NewFromURL(cloudinaryUrl)

	uploadResult, err := cld.Upload.Upload(
		ctx,
		filePath,
		uploader.UploadParams{PublicID: fmt.Sprintf("travex/%v", time.Now())})
	if err != nil {
		fmt.Println(err)
	}

	log.Println(uploadResult.SecureURL)
	return uploadResult.SecureURL, nil
}
