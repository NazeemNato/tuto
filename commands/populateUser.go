package main

import (
	"github.com/NazeemNato/tuto/src/database"
	"github.com/NazeemNato/tuto/src/models"
	"github.com/bxcodec/faker/v3"
)

func main() {
	database.Connnect()
	
	for i := 0; i < 50; i++ {
		ambassador := models.User{
			Firstname: faker.FirstName(),
			Lastname: faker.LastName(),
			Email: faker.Email(),
			IsAmbassador: true,
		}
		ambassador.SetPassword("1234")
		database.DB.Create(&ambassador)
	}
}
