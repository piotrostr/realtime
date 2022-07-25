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

	// returns User obj
	router.POST("/create", func(c *gin.Context) {
		var user database.User
		err := c.BindJSON(&user)
		if err != nil {
			c.JSON(400, gin.H{"error": err.Error()})
			return
		}
		createdUser := db.Create(user)
		c.JSON(201, createdUser)
	})

	// returns []User obj
	router.GET("/read", func(c *gin.Context) {
		users := db.ReadAll()
		c.JSON(200, users)
	})

	// returns metadata associated with the update
	router.PUT("/update", func(c *gin.Context) {
		var user database.User
		err := c.BindJSON(&user)
		if err != nil {
			c.JSON(400, gin.H{"error": err.Error()})
			return
		}
		meta := db.Update(user)
		c.JSON(204, meta)
	})

	// returns User obj
	router.GET("/read/:name", func(c *gin.Context) {
		name := c.Param("name")
		var user database.User
		userByName, _, err := db.ReadOne(name)
		if err != nil {
			c.JSON(400, gin.H{"error": err.Error()})
			return
		}
		if userByName == nil {
			c.JSON(404, gin.H{"error": "user not found"})
			return
		}
		c.JSON(200, user)
	})

	// returns metadata associated with the delete
	router.DELETE("/delete", func(c *gin.Context) {
		var user database.User
		err := c.BindJSON(&user)
		if err != nil {
			c.JSON(400, gin.H{"error": err.Error()})
			return
		}
		meta := db.Delete(user.Name)
		c.JSON(204, meta)
	})

	return router
}
