package models

import (
	"time"

	"gorm.io/gorm"
)

type Comment struct {
	Id           string    `gorm:"Column:id"`
	Location_id  string    `gorm:"Column:location_id"`
	Author_id    string    `gorm:"column:author_id"`
	Parent_id    string    `gorm:"column:parent_id"`
	Text         string    `gorm:"Column:text"`
	Review_id    string    `gorm:"column:review_id"`
	IsApproved   bool      `gorm:"Column:is_approved"`
	IsApprovedBy string    `gorm:"Column:is_approved_by"`
	IsApprovedAt time.Time `gorm:"Column:is_approved_at"`
	Created_at   time.Time `gorm:"Column:created_at"`
	Updated_at   time.Time `gorm:"Column:updated_at"`
	Replies      []Comment

	Author User
}

func (u *Comment) BeforeCreate(tx *gorm.DB) (err error) {
	u.Created_at = time.Now()
	u.Updated_at = time.Now()
	return nil
}
