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

type CommentResponse struct {
	Id      int    `json:"id"`
	PhotoId int    `json:"photo_id"`
	Message string `json:"message"`
}

func PostComment(c *gin.Context) {
	db := database.GetDB()
	contentType := helpers.GetContentType(c)
	userData := c.MustGet("userData").(jwt.MapClaims)
	Comment := models.Comment{}

	if contentType == appJson {
		c.ShouldBindJSON(&Comment)
	} else {
		c.ShouldBind(&Comment)
	}

	Comment.UserId = uint(userData["id"].(float64))

	if err := db.Debug().Create(&Comment).Error; err != nil {
		c.JSON(http.StatusBadRequest, gin.H{
			"error":   "Bad Request",
			"message": err.Error(),
		})
		return
	}

	c.JSON(http.StatusCreated, gin.H{
		"id": Comment.ID,
	})
}

func UpdateComment(c *gin.Context) {
	db := database.GetDB()
	contentType := helpers.GetContentType(c)
	userData := c.MustGet("userData").(jwt.MapClaims)
	commentID, err := strconv.Atoi(c.Param("commentId"))
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{
			"error":   "Bad Request",
			"message": "Invalid comment ID",
		})
		return
	}

	var updatedComment models.Comment

	if contentType == appJson {
		c.ShouldBindJSON(&updatedComment)
	} else {
		c.ShouldBind(&updatedComment)
	}

	// Find the comment by ID
	var existingComment models.Comment
	if err := db.Where("id = ? AND user_id = ?", commentID, uint(userData["id"].(float64))).First(&existingComment).Error; err != nil {
		c.JSON(http.StatusNotFound, gin.H{
			"error":   "Not Found",
			"message": "Comment not found",
		})
		return
	}

	// Update the existing comment
	if err := db.Model(&existingComment).Updates(updatedComment).Error; err != nil {
		c.JSON(http.StatusBadRequest, gin.H{
			"error":   "Bad Request",
			"message": err.Error(),
		})
		return
	}

	c.JSON(http.StatusOK, gin.H{
		"message": "Comment updated successfully",
	})
}

func DeleteComment(c *gin.Context) {
	db := database.GetDB()
	userData := c.MustGet("userData").(jwt.MapClaims)
	commentID, err := strconv.Atoi(c.Param("commentId"))
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{
			"error":   "Bad Request",
			"message": "Invalid comment ID",
		})
		return
	}

	// Find the comment by ID
	var comment models.Comment
	if err := db.Where("id = ? AND user_id = ?", commentID, uint(userData["id"].(float64))).First(&comment).Error; err != nil {
		c.JSON(http.StatusNotFound, gin.H{
			"error":   "Not Found",
			"message": "Comment not found",
		})
		return
	}

	// Delete the comment
	if err := db.Delete(&comment).Error; err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{
			"error":   "Internal Server Error",
			"message": "Failed to delete the comment",
		})
		return
	}

	c.JSON(http.StatusOK, gin.H{
		"message": "Comment deleted successfully",
	})
}

func GetOneComment(c *gin.Context) {
	db := database.GetDB()
	userData := c.MustGet("userData").(jwt.MapClaims)
	commentID, err := strconv.Atoi(c.Param("commentId"))
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{
			"error":   "Bad Request",
			"message": "Invalid comment ID",
		})
		return
	}

	// Find the comment by ID
	var comment models.Comment
	if err := db.Where("id = ? AND user_id = ?", commentID, uint(userData["id"].(float64))).First(&comment).Error; err != nil {
		c.JSON(http.StatusNotFound, gin.H{
			"error":   "Not Found",
			"message": "Comment not found",
		})
		return
	}

	c.JSON(http.StatusOK, gin.H{
		"comment": CommentResponse{
			Id:      int(comment.ID),
			PhotoId: int(comment.PhotoId),
			Message: comment.Message,
		},
	})
}

func GetAllComments(c *gin.Context) {
	db := database.GetDB()
	userData := c.MustGet("userData").(jwt.MapClaims)

	var comments []models.Comment

	// Retrieve all comments belonging to the authenticated user
	userID := uint(userData["id"].(float64))
	if err := db.Where("user_id = ?", userID).Find(&comments).Error; err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{
			"error":   "Internal Server Error",
			"message": "Failed to fetch comments",
		})
		return
	}

	// If no comments are found, return an empty response
	if len(comments) == 0 {
		c.JSON(http.StatusOK, gin.H{
			"message":  "No comments found",
			"comments": []interface{}{},
		})
		return
	}

	var responses []CommentResponse
	for _, value := range comments {
		responses = append(responses, CommentResponse{
			Id:      int(value.ID),
			PhotoId: int(value.PhotoId),
			Message: value.Message,
		})
	}

	// Return the list of comments
	c.JSON(http.StatusOK, gin.H{
		"comments": responses,
	})
}
