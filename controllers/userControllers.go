package controllers

import (
	"final_project/database"
	"final_project/helpers"
	"final_project/models"
	"net/http"

	"github.com/gin-gonic/gin"
)

var (
	appJson = "application/json"
)

func UserRegister(c *gin.Context) {
	db := database.GetDB()
	contentType := helpers.GetContentType(c)
	_, _ = db, contentType
	User := models.User{}

	if contentType == appJson {
		c.ShouldBindJSON(&User)
	} else {
		c.ShouldBind(&User)
	}

	if User.Age < 9 {
		c.JSON(http.StatusBadRequest, gin.H{
			"error":   "Bad Request",
			"message": "Age must be greater than 8",
		})
		return
	}

	if err := db.Debug().Create(&User).Error; err != nil {
		c.JSON(http.StatusBadRequest, gin.H{
			"error":   "Bad Request",
			"message": err.Error(),
		})
		return
	}

	c.JSON(http.StatusCreated, gin.H{
		"id":       User.ID,
		"email":    User.Email,
		"username": User.Username,
	})
}

func UserLogin(c *gin.Context) {
	db := database.GetDB()
	contentType := helpers.GetContentType(c)
	_, _ = db, contentType
	User := models.User{}
	password := ""

	if contentType == appJson {
		c.ShouldBindJSON(&User)
	} else {
		c.ShouldBind(&User)
	}

	password = User.Password

	if err := db.Debug().Where("email = ?", User.Email).Take(&User).Error; err != nil {
		c.JSON(http.StatusUnauthorized, gin.H{
			"error":   "Unauthorized",
			"message": "invalid email/password",
		})
		return
	}

	comparePass := helpers.ComparePass([]byte(User.Password), []byte(password))

	if !comparePass {
		c.JSON(http.StatusUnauthorized, gin.H{
			"error":   "Unauthorized",
			"message": "invalid email/password",
		})
		return
	}

	token := helpers.GenerateToken(User.ID, User.Email)

	c.JSON(http.StatusOK, gin.H{
		"token": token,
	})
}
