package models

import (
	"final_project/helpers"

	"github.com/asaskevich/govalidator"
	"gorm.io/gorm"
)

type User struct {
	gorm.Model
	Username     string        `gorm:"not null;uniqueIndex" json:"username" form:"username" valid:"required~Your username is required"`
	Email        string        `gorm:"not null;uniqueIndex" json:"email" form:"email" valid:"required~Your email is required,email~Invalid email format"`
	Password     string        `gorm:"not null" json:"password" form:"password" valid:"required~Your password is required,minstringlength(6)~Password has to have a minimum length of 6 characters"`
	Age          uint          `gorm:"not null;gte=9" json:"age" form:"age" valid:"required~Your age is required"`
	Photos       []Photo       `gorm:"constraint:OnUpdate:CASCADE,OnDelete:SET NULL;" json:"photos"`
	SocialMedias []SocialMedia `gorm:"constraint:OnUpdate:CASCADE,OnDelete:SET NULL;" json:"social_media"`
	Comments     []Comment     `gorm:"constraint:OnUpdate:CASCADE,OnDelete:SET NULL;" json:"comments"`
}

func (u *User) BeforeCreate(tx *gorm.DB) error {
	if _, err := govalidator.ValidateStruct(u); err != nil {
		return err
	}

	u.Password = helpers.HashPass(u.Password)
	return nil
}
