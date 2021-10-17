package database

import (
	"github.com/NazeemNato/tuto/src/models"
	"gorm.io/driver/mysql"
	"gorm.io/gorm"
)

var DB *gorm.DB

func Connnect() {
	var err error
	DB, err = gorm.Open(mysql.Open("root:root@tcp(db:3306)/ambassador"), &gorm.Config{})

	if err != nil {
		panic(err)
	}
}

func AutoMigrate() {
	DB.AutoMigrate(models.User{}, models.Product{})
}