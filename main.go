package main

import (
	"github.com/gin-gonic/gin"
	"net/http"
	"tusk/config"
	"tusk/controllers"
	"tusk/models"
)

func main() {
	// Database
	db := config.DatabaseConnection()
	db.AutoMigrate(&models.User{}, &models.Task{})
	config.CreateOwnerAccount(db)

	// Controller
	userController := controllers.UserController{DB: db}

	// Router
	router := gin.Default()
	router.GET("/", func(c *gin.Context) {
		c.JSON(http.StatusOK, "Hello World")
	})

	router.POST("/users/login", userController.Login)
	router.POST("/users", userController.CreateAccount)
	router.Static("/attachments", "./attachments")
	// Start server
	router.Run("192.168.110.216:8080")
}
