package models

import (
	"time"

	"golang.org/x/crypto/bcrypt"
	"gorm.io/gorm"
)

type User struct {
	Id         string    `gorm:"column:id"`
	First_name string    `gorm:"column:first_name"`
	Last_name  string    `gorm:"column:last_name"`
	Email      string    `gorm:"column:email"`
	Password   []byte    `gorm:"password" json:"-"`
	Created_at time.Time `gorm:"column:created_at"`
	Updated_at time.Time `gorm:"column:updated_at"`
	Role       string    `gorm:"role"`
	Locations  []Location
}

func (u *User) BeforeCreate(tx *gorm.DB) (err error) {
	u.Created_at = time.Now()
	u.Updated_at = time.Now()
	return nil
}

func HashPassword(plainPassword string) []byte {
	password, _ := bcrypt.GenerateFromPassword([]byte(plainPassword), 10)
	return password
}

func ComparePassword(password []byte, plainPassword []byte) error {
	err := bcrypt.CompareHashAndPassword([]byte(password), plainPassword)

	return err
}
