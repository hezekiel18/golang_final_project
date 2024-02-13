package models

import (
	"github.com/asaskevich/govalidator"
	"gorm.io/gorm"
)

type Comment struct {
	gorm.Model
	UserId  uint   `json:"user_id" form:"user_id"`
	PhotoId uint   `json:"photo_id" form:"photo_id"`
	Message string `gorm:"not null" json:"message" form:"message" valid:"required~Comment of the photo is required"`
	User    *User
	Photo   *Photo
}

func (c *Comment) BeforeCreate(tx *gorm.DB) error {
	if _, err := govalidator.ValidateStruct(c); err != nil {
		return err
	}

	return nil
}
