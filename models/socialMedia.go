package models

import (
	"github.com/asaskevich/govalidator"
	"gorm.io/gorm"
)

type SocialMedia struct {
	gorm.Model
	Name           string `gorm:"not null" json:"name" form:"name" valid:"required~Name of your social media is required"`
	SocialMediaUrl string `gorm:"not null" json:"social_media_url" form:"social_media_url" valid:"required~Url of your social media is required"`
	UserId         uint   `json:"user_id" form:"user_id"`
	User           *User
}

func (s *SocialMedia) BeforeCreate(tx *gorm.DB) error {
	if _, err := govalidator.ValidateStruct(s); err != nil {
		return err
	}

	return nil
}
