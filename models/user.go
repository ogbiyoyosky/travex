package models

import (
	"time"

	"golang.org/x/crypto/bcrypt"
)

type User struct {
	Id        string
	FirstName string `gorm:"column:firstName"`
	LastName  string `gorm:"column:lastName"`
	Email     string
	Password  []byte    `gorm:"password"`
	CreatedAt time.Time `gorm:"column:createdAt"`
	UpdatedAt time.Time `gorm:"column:updatedAt"`
	Role      string    `gorm:"role"`
}

func HashPassword(plainPassword string) []byte {
	password, _ := bcrypt.GenerateFromPassword([]byte(plainPassword), 10)
	return password
}
