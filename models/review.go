package models

import (
	"time"

	"gorm.io/gorm"
)

type Review struct {
	Id          string    `gorm:"Column:id"`
	Location_id string    `gorm:"Column:location_id"`
	Author_id   string    `gorm:"column:author_id"`
	Rating      float32   `gorm:"Column:rating"`
	Created_at  time.Time `gorm:"Column:created_at"`
	Updated_at  time.Time `gorm:"Column:updated_at"`
	Author      User
}

func (u *Review) BeforeCreate(tx *gorm.DB) (err error) {
	u.Created_at = time.Now()
	u.Updated_at = time.Now()
	return nil
}
