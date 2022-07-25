package rest

import (
	"fmt"

	"github.com/gin-gonic/gin"
	_ "github.com/joho/godotenv/autoload"
	database "github.com/piotrostr/realtime/db"
)

func GetRouter() *gin.Engine {
	gin.SetMode(gin.ReleaseMode)
	router := gin.Default()

	db := database.DB{}
	db.Init()

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
		userByName, _, err := db.ReadOne(name)
		fmt.Println(err)
		if err != nil {
			c.JSON(404, gin.H{"error": err.Error()})
			return
		}
		c.JSON(200, userByName)
	})

	// returns metadata associated with the delete
	router.DELETE("/delete", func(c *gin.Context) {
		var user database.User
		err := c.BindJSON(&user)
		fmt.Printf("%+v\n", user)
		if err != nil {
			c.JSON(400, gin.H{"error": err.Error()})
			return
		}

		meta := db.Delete(user)
		if meta == nil {
			c.JSON(404, gin.H{"error": "user not found"})
			return
		}

		c.JSON(204, meta)
	})

	return router
}
