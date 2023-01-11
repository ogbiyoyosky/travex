package models

import (
	"time"

	"gorm.io/gorm"
)

type Location struct {
	Id               string    `gorm:"column:id"`
	Name             string    `gorm:"column:name"`
	Image            string    `gorm:"column:image"`
	Address          string    `gorm:"column:address"`
	Location_type_id string    `gorm:"column:location_type_id"`
	Description      string    `gorm:"column:description"`
	UserId           string    `gorm:"column:user_id"`
	Created_at       time.Time `gorm:"column:created_at"`
	Updated_at       time.Time `gorm:"column:updated_at"`
	IsApproved       bool      `gorm:"Column:is_approved"`
	IsApprovedAt     time.Time `gorm:"Column:is_approved_at"`
	User             User
	LocationType     LocationType
	Reviews          []Review `gorm:"foreignKey:location_id;references:id"`
}

func (u *Location) BeforeCreate(tx *gorm.DB) (err error) {
	u.Created_at = time.Now()
	u.Updated_at = time.Now()
	return nil
}

type LocationType struct {
	Id         string    `gorm:"column:id"`
	Name       string    `gorm:"column:name"`
	Created_at time.Time `gorm:"column:created_at"`
	Updated_at time.Time `gorm:"column:updated_at"`
}

func (u *LocationType) BeforeCreate(tx *gorm.DB) (err error) {
	u.Created_at = time.Now()
	u.Updated_at = time.Now()
	return nil
}
