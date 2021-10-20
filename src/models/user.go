package models

import (
	"golang.org/x/crypto/bcrypt"
	"gorm.io/gorm"
)

type User struct {
	Model
	Firstname    string `json:"first_name"`
	Lastname     string `json:"last_name"`
	Email        string `json:"email" gorm:"unique"`
	Password     []byte `json:"-"`
	IsAmbassador bool   `json:"-"`
	Revenue *float64 `json:"revenue,omitempty" gorm:"-"`
}

func (user *User) SetPassword(pwd string) {
	password, _ := bcrypt.GenerateFromPassword([]byte(pwd), 12)
	user.Password = password
}

func (user *User) ComparePassword(pwd string) bool {
	err := bcrypt.CompareHashAndPassword(user.Password, []byte(pwd))
	return err == nil
}

type Admin User


func ( admin *Admin) CalculateRevenue(database *gorm.DB) {
	
}

type Ambassador User

func ( ambassador *Ambassador) CalculateRevenue(database *gorm.DB) {
	var orders []Order
	database.Preload("OrderItems").Find(&orders, &Order{
		UserId: ambassador.Id,
		Complete: true,
	})

	var revenue float64 = 0.0
	for _ , order := range orders {
		for _, orderItem := range order.OrderItems {
			revenue += orderItem.AmbassadorRevenue
		}
	}

	ambassador.Revenue = &revenue
}