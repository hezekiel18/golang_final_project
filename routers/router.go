package routers

import (
	"final_project/controllers"
	"final_project/middlewares"

	"github.com/gin-gonic/gin"
)

func StartApp() *gin.Engine {
	r := gin.Default()

	userRouter := r.Group("/users")
	{
		userRouter.POST("/register", controllers.UserRegister)
		userRouter.POST("/login", controllers.UserLogin)
	}

	photoRouter := r.Group("/photo")
	{
		photoRouter.Use(middlewares.Authentication())
		photoRouter.GET("/all", controllers.GetAllPhotos)
		photoRouter.GET("/:photoId", middlewares.PhotoAuthorization(), controllers.GetOnePhoto)
		photoRouter.POST("/upload", controllers.UploadPhoto)
		photoRouter.PUT("/:photoId", middlewares.PhotoAuthorization(), controllers.UpdatePhoto)
		photoRouter.DELETE("/:photoId", middlewares.PhotoAuthorization(), controllers.DeletePhoto)
	}

	commentRouter := r.Group("/comment")
	{
		commentRouter.Use(middlewares.Authentication())
		commentRouter.GET("/all", controllers.GetAllComments)
		commentRouter.GET("/:commentId", middlewares.CommentAuthorization(), controllers.GetOneComment)
		commentRouter.POST("/post", controllers.PostComment)
		commentRouter.PUT("/:commentId", middlewares.CommentAuthorization(), controllers.UpdateComment)
		commentRouter.DELETE("/:commentId", middlewares.CommentAuthorization(), controllers.DeleteComment)
	}

	r.Static("/images", "./images")

	return r
}
