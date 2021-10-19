package models

import "golang.org/x/crypto/bcrypt"

type User struct {
	Model
	Firstname    string `json:"first_name"`
	Lastname     string `json:"last_name"`
	Email        string `json:"email" gorm:"unique"`
	Password     []byte `json:"-"`
	IsAmbassador bool   `json:"-"`
}

func (user *User) SetPassword(pwd string) {
	password, _ := bcrypt.GenerateFromPassword([]byte(pwd), 12)
	user.Password = password
}

func (user *User) ComparePassword(pwd string) bool {
	err := bcrypt.CompareHashAndPassword(user.Password, []byte(pwd))
	return err == nil
}
