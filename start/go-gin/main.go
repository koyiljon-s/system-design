package main

import (
	"github.com/gin-gonic/gin"
	"gorm.io/driver/sqlite"
	"gorm.io/gorm"
)

var db *gorm.DB

func main() {
	var err error
	db, err = gorm.Open(sqlite.Open("test.db"), &gorm.Config{})
	if err != nil {
		panic("failed to connect database")
	}

	db.AutoMigrate(&User{})

	r := gin.Default()
	r.POST("/users", createUser)
	r.GET("/users/:id", getUser)
	r.GET("/users", getAllUsers) // add this for getting all users

	r.Run(":8080")
}

func getAllUsers(c *gin.Context) {
	var users []User
	db.Find(&users)
	c.JSON(200, users)
}

func getUser(c *gin.Context) {
	id := c.Param("id")
	var user User
	
	if result := db.First(&user, id); result.Error != nil {
		c.JSON(404, gin.H{"error": "User not found"})
		return
	}
	
	c.JSON(200, user)
}

func createUser(c *gin.Context) {
	var user User
	if err := c.ShouldBindJSON(&user); err != nil {
		c.JSON(400, gin.H{"error": err.Error()})
		return
	}
	
	if result := db.Create(&user); result.Error != nil {
		c.JSON(500, gin.H{"error": result.Error.Error()})
		return
	}
	
	c.JSON(201, user)
}