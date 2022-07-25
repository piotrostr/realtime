package rest

import (
	"github.com/gin-gonic/gin"
	_ "github.com/joho/godotenv/autoload"
	database "github.com/piotrostr/realtime/db"
)

func GetRouter() *gin.Engine {
	gin.SetMode(gin.ReleaseMode)
	router := gin.Default()
	db := database.DB{}

	router.POST("/create", func(c *gin.Context) {
		var user database.User
		err := c.BindJSON(&user)
		if err != nil {
			c.JSON(400, gin.H{"error": err.Error()})
			return
		}
		db.Create(user)
		c.JSON(200, gin.H{"message": "success"})
	})

	router.GET("/read", func(c *gin.Context) {
		var user database.User
		db.ReadOne()
		c.JSON(200, user)
	})

	router.PUT("/update", func(c *gin.Context) {
	})

	router.DELETE("/delete", func(c *gin.Context) {})

	return router
}
