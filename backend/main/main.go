package main

import (
	c "blog/controllers"

	"github.com/gin-gonic/gin"
)

func main() {
	router := gin.Default()
	router.Use(func(c *gin.Context) {
		c.Writer.Header().Set("Access-Control-Allow-Origin", "*")
		c.Writer.Header().Set("Access-Control-Allow-Methods", "GET, POST, PUT, DELETE, OPTIONS")
		c.Writer.Header().Set("Access-Control-Allow-Headers", "Content-Type, Authorization")
		if c.Request.Method == "OPTIONS" {
			c.AbortWithStatus(200)
			return
		}
		c.Next()
	})
	router.GET("/ws", c.CreateSocket)
	router.Static("/images", "../images")
	router.Static("/ppimages", "../pp_images")
	router.POST("/create_user", c.CreateUser)
	router.POST("/get_user", c.GetUser)
	router.POST("/upload_image", c.UploadImage)
	router.POST("/create_blog", c.CreateBlog)
	router.GET("/get_blogs", c.GetAllBlogs)
	router.GET("/single_blog", c.GetSingleBlog)
	router.GET("/get_notif", c.CheckNotif)
	router.GET("/get_pp", c.GetPPImage)
	router.GET("/get_notifications", c.GetNotifications)
	router.GET("/change_notif", c.ChangeNotif)
	router.Run("localhost:8080")

}
