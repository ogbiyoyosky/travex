package main

import (
	"encoding/json"
	"fmt"
	"math/rand"
	"os"
	"path"
	"runtime"

	connection "github.com/ogbiyoyosky/travex/db"
	"github.com/ogbiyoyosky/travex/models"
)

type SeedData struct {
	Users         []models.User
	Locations     []models.Location
	LocationTypes []models.Location
}

func main() {
	fmt.Println("--------------------------------")
	RunSeedData()
}

func RunSeedData() {
	connection.Connect()
	var seed *SeedData
	_, filename, _, _ := runtime.Caller(1)

	filepath := path.Join(path.Dir(filename), "location.json")

	file, err := os.OpenFile(filepath, os.O_RDONLY, 0777)

	if err != nil {
		fmt.Println("ERRR READ", err)
	}
	defer file.Close()

	fmt.Println("--------------------------------")

	if err := json.NewDecoder(file).Decode(&seed); err != nil {

		fmt.Println("--------------------------------")
		locationTypeIds := make([]string, 0)
		userIds := make([]string, 0)

		fmt.Println(userIds)
		fmt.Println("--------------------------------")

		for i := 0; i < len(seed.LocationTypes); i++ {

			locationType := models.LocationType{Name: seed.LocationTypes[i].Name}
			connection.DB.Create(&locationType)
			locationTypeIds = append(locationTypeIds, locationType.Id)
		}

		for i := 0; i < len(seed.Users); i++ {
			fmt.Println("__USERS", seed.Users)
			hashPassword := models.HashPassword("Miracle123")
			user := models.User{First_name: seed.Users[i].First_name, Last_name: seed.Users[i].Last_name, Email: seed.Users[i].Email, Password: hashPassword, Role: "customer"}
			connection.DB.Create(&user)
			userIds = append(userIds, user.Id)

		}

		for i := 0; i < len(seed.Locations); i++ {
			randomIndexLocationType := rand.Intn(len(locationTypeIds))

			randomIndexUserId := rand.Intn(len(userIds))
			location := models.Location{Name: seed.Locations[i].Name, Location_type_id: locationTypeIds[randomIndexLocationType], Description: seed.Locations[i].Description, UserId: userIds[randomIndexUserId], Image: seed.Locations[i].Image, Address: seed.Locations[i].Address}
			connection.DB.Create(&location)
		}
	}

}
