package controllers

import (
	"final_project/database"
	"final_project/helpers"
	"final_project/models"
	"net/http"
	"strconv"

	"github.com/dgrijalva/jwt-go"
	"github.com/gin-gonic/gin"
)

type SocialMediaResponse struct {
	Name string `json:"name"`
	Url  string `json:"social_media_url"`
}

func AddSocialMedia(c *gin.Context) {
	db := database.GetDB()
	contentType := helpers.GetContentType(c)
	userData := c.MustGet("userData").(jwt.MapClaims)
	SocialMedia := models.SocialMedia{}

	if contentType == appJson {
		c.ShouldBindJSON(&SocialMedia)
	} else {
		c.ShouldBind(&SocialMedia)
	}

	SocialMedia.UserId = uint(userData["id"].(float64))

	if err := db.Debug().Create(&SocialMedia).Error; err != nil {
		c.JSON(http.StatusBadRequest, gin.H{
			"error":   "Bad Request",
			"message": err.Error(),
		})
		return
	}

	c.JSON(http.StatusCreated, gin.H{
		"id": SocialMedia.ID,
	})
}

func UpdateSocialMedia(c *gin.Context) {
	db := database.GetDB()
	contentType := helpers.GetContentType(c)
	userData := c.MustGet("userData").(jwt.MapClaims)
	SocialMediaId, err := strconv.Atoi(c.Param("socialMediaId"))
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{
			"error":   "Bad Request",
			"message": "Invalid Social Media ID",
		})
		return
	}

	var updatedSocialMedia models.SocialMedia

	if contentType == appJson {
		c.ShouldBindJSON(&updatedSocialMedia)
	} else {
		c.ShouldBind(&updatedSocialMedia)
	}

	// Find the social media by ID
	var existingSocialMedia models.SocialMedia
	if err := db.Where("id = ? AND user_id = ?", SocialMediaId, uint(userData["id"].(float64))).First(&existingSocialMedia).Error; err != nil {
		c.JSON(http.StatusNotFound, gin.H{
			"error":   "Not Found",
			"message": "Social Media not found",
		})
		return
	}

	// Update the existing social media
	if err := db.Model(&existingSocialMedia).Updates(updatedSocialMedia).Error; err != nil {
		c.JSON(http.StatusBadRequest, gin.H{
			"error":   "Bad Request",
			"message": err.Error(),
		})
		return
	}

	c.JSON(http.StatusOK, gin.H{
		"message": "Social Media updated successfully",
	})
}

func DeleteSocialMedia(c *gin.Context) {
	db := database.GetDB()
	userData := c.MustGet("userData").(jwt.MapClaims)
	socialMediaId, err := strconv.Atoi(c.Param("socialMediaId"))
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{
			"error":   "Bad Request",
			"message": "Invalid social media ID",
		})
		return
	}

	// Find the social media by ID
	var socialMedia models.SocialMedia
	if err := db.Where("id = ? AND user_id = ?", socialMediaId, uint(userData["id"].(float64))).First(&socialMedia).Error; err != nil {
		c.JSON(http.StatusNotFound, gin.H{
			"error":   "Not Found",
			"message": "Social Media not found",
		})
		return
	}

	// Delete the social media
	if err := db.Delete(&socialMedia).Error; err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{
			"error":   "Internal Server Error",
			"message": "Failed to delete the social media",
		})
		return
	}

	c.JSON(http.StatusOK, gin.H{
		"message": "Social Media deleted successfully",
	})
}

func GetOneSocialMedia(c *gin.Context) {
	db := database.GetDB()
	userData := c.MustGet("userData").(jwt.MapClaims)
	socialMediaId, err := strconv.Atoi(c.Param("socialMediaId"))
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{
			"error":   "Bad Request",
			"message": "Invalid social media ID",
		})
		return
	}

	// Find the social media by ID
	var socialMedia models.SocialMedia
	if err := db.Where("id = ? AND user_id = ?", socialMediaId, uint(userData["id"].(float64))).First(&socialMedia).Error; err != nil {
		c.JSON(http.StatusNotFound, gin.H{
			"error":   "Not Found",
			"message": "Social Media not found",
		})
		return
	}

	c.JSON(http.StatusOK, gin.H{
		"social_media": SocialMediaResponse{
			Name: socialMedia.Name,
			Url:  socialMedia.SocialMediaUrl,
		},
	})
}

func GetAllSocialMedias(c *gin.Context) {
	db := database.GetDB()
	userData := c.MustGet("userData").(jwt.MapClaims)

	var socialMedias []models.SocialMedia

	// Retrieve all social medias belonging to the authenticated user
	userID := uint(userData["id"].(float64))
	if err := db.Where("user_id = ?", userID).Find(&socialMedias).Error; err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{
			"error":   "Internal Server Error",
			"message": "Failed to fetch social media",
		})
		return
	}

	// If no social media are found, return an empty response
	if len(socialMedias) == 0 {
		c.JSON(http.StatusOK, gin.H{
			"message":      "No social medias found",
			"social_media": []interface{}{},
		})
		return
	}

	var responses []SocialMediaResponse
	for _, value := range socialMedias {
		responses = append(responses, SocialMediaResponse{
			Name: value.Name,
			Url:  value.SocialMediaUrl,
		})
	}

	// Return the list of social medias
	c.JSON(http.StatusOK, gin.H{
		"social_medias": responses,
	})
}
