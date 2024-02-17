package controllers

import (
	"final_project/database"
	"final_project/models"
	"fmt"
	"net/http"
	"os"
	"path/filepath"
	"strconv"
	"strings"

	"github.com/dgrijalva/jwt-go"
	"github.com/gin-gonic/gin"
	"github.com/google/uuid"
)

type PhotoResponse struct {
	Title   string `json:"title"`
	Caption string `json:"caption"`
	Url     string `json:"url"`
}

func UploadPhoto(c *gin.Context) {
	db := database.GetDB()
	userData := c.MustGet("userData").(jwt.MapClaims)
	photo := models.Photo{}
	photo.Title = c.Request.FormValue("title")
	photo.Caption = c.Request.FormValue("caption")
	photo.UserId = uint(userData["id"].(float64))

	file, err := c.FormFile("image")
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{
			"error":   "Bad Request",
			"message": "Error uploading image",
		})
		return
	}

	// Limit the maximum file size to 1 MB (1 * 1024 * 1024 bytes)
	maxFileSize := int64(1 * 1024 * 1024)
	if file.Size > maxFileSize {
		c.JSON(http.StatusBadRequest, gin.H{
			"error":   "Bad Request",
			"message": "File size exceeds the maximum limit (1 MB)",
		})
		return
	}

	filename := strings.Replace(uuid.New().String(), "-", "", -1)
	fileExt := strings.Split(file.Filename, ".")[1]
	image := fmt.Sprintf("%s.%s", filename, fileExt)

	if err = c.SaveUploadedFile(file, fmt.Sprintf("./images/%s", image)); err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{
			"error":   "Internal Server Error",
			"message": "Error in saving Image",
		})
		return
	}

	photo.Url = fmt.Sprintf("http://localhost:8080/images/%s", image)

	if err := db.Debug().Create(&photo).Error; err != nil {
		c.JSON(http.StatusBadRequest, gin.H{
			"error":   "Bad Request",
			"message": err.Error(),
		})
		return
	}

	c.JSON(http.StatusCreated, gin.H{
		"id":      photo.ID,
		"title":   photo.Title,
		"caption": photo.Caption,
		"url":     photo.Url,
	})
}

func UpdatePhoto(c *gin.Context) {
	db := database.GetDB()
	userData := c.MustGet("userData").(jwt.MapClaims)
	photo := models.Photo{}

	photoId, _ := strconv.Atoi(c.Param("photoId"))
	newTitle := c.Request.FormValue("title")
	newCaption := c.Request.FormValue("caption")
	photo.UserId = uint(userData["id"].(float64))

	file, err := c.FormFile("image")
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{
			"error":   "Bad Request",
			"message": "Error uploading image",
		})
		return
	}

	maxFileSize := int64(1 * 1024 * 1024)
	if file.Size > maxFileSize {
		c.JSON(http.StatusBadRequest, gin.H{
			"error":   "Bad Request",
			"message": "File size exceeds the maximum limit (1 MB)",
		})
		return
	}

	if err := db.First(&photo, photoId).Error; err != nil {
		c.JSON(http.StatusNotFound, gin.H{
			"error":   "Not Found",
			"message": err.Error(),
		})
		return
	}

	if err := db.Model(&photo).Where("id = ?", photoId).Updates(models.Photo{
		Title:   newTitle,
		Caption: newCaption,
	}).Error; err != nil {
		c.JSON(http.StatusBadRequest, gin.H{
			"error":   "Bad Request",
			"message": err.Error(),
		})
		return
	}

	findImageName := strings.Split(photo.Url, "/")

	if err = c.SaveUploadedFile(file, fmt.Sprintf("./images/%s", findImageName[len(findImageName)-1])); err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{
			"error":   "Internal Server Error",
			"message": "Error in saving Image",
		})
		return
	}

	c.JSON(http.StatusCreated, gin.H{
		"id":      photo.ID,
		"title":   photo.Title,
		"caption": photo.Caption,
		"url":     photo.Url,
	})
}

func DeletePhoto(c *gin.Context) {
	db := database.GetDB()
	userData := c.MustGet("userData").(jwt.MapClaims)
	photoId, err := strconv.Atoi(c.Param("photoId"))
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{
			"error":   "Bad Request",
			"message": "Invalid photo ID",
		})
		return
	}

	// Check if the photo exists
	var photo models.Photo
	result := db.Where("id = ? AND user_id = ?", photoId, uint(userData["id"].(float64))).First(&photo)
	if result.Error != nil {
		c.JSON(http.StatusNotFound, gin.H{
			"error":   "Not Found",
			"message": "Photo not found",
		})
		return
	}

	// Delete the photo
	if err := db.Delete(&photo).Error; err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{
			"error":   "Internal Server Error",
			"message": "Failed to delete photo",
		})
		return
	}

	// Delete the associated image file
	imageName := filepath.Base(photo.Url)
	imagePath := fmt.Sprintf("./images/%s", imageName)
	if err := os.Remove(imagePath); err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{
			"error":   "Internal Server Error",
			"message": "Error deleting image file",
		})
	}

	c.JSON(http.StatusOK, gin.H{
		"message": "Photo deleted successfully",
	})
}

func GetOnePhoto(c *gin.Context) {
	db := database.GetDB()
	userData := c.MustGet("userData").(jwt.MapClaims)
	photoId, err := strconv.Atoi(c.Param("photoId"))
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{
			"error":   "Bad Request",
			"message": "Invalid photo ID",
		})
		return
	}

	// Check if the photo exists
	var photo models.Photo
	result := db.Where("id = ? AND user_id = ?", photoId, uint(userData["id"].(float64))).First(&photo)
	if result.Error != nil {
		c.JSON(http.StatusNotFound, gin.H{
			"error":   "Not Found",
			"message": "Photo not found",
		})
		return
	}

	c.JSON(http.StatusOK, gin.H{
		"id":      photo.ID,
		"title":   photo.Title,
		"caption": photo.Caption,
		"url":     photo.Url,
	})
}

func GetAllPhotos(c *gin.Context) {
	db := database.GetDB()
	userData := c.MustGet("userData").(jwt.MapClaims)

	var photos []models.Photo
	userId := uint(userData["id"].(float64))

	// Retrieve all photos belonging to the authenticated user
	if err := db.Where("user_id = ?", userId).Find(&photos).Error; err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{
			"error":   "Internal Server Error",
			"message": "Failed to fetch photos",
		})
		return
	}

	// If no photos are found, return an empty response
	if len(photos) == 0 {
		c.JSON(http.StatusOK, gin.H{
			"message": "No photos found",
			"photos":  []interface{}{},
		})
		return
	}

	var responses []PhotoResponse
	for _, value := range photos {
		responses = append(responses, PhotoResponse{
			Title:   value.Title,
			Caption: value.Caption,
			Url:     value.Url,
		})
	}

	// Return the list of photos
	c.JSON(http.StatusOK, gin.H{
		"photos": responses,
	})
}
