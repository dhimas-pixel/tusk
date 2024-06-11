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
	taskController := controllers.TaskController{DB: db}

	// Router
	router := gin.Default()
	router.GET("/", func(c *gin.Context) {
		c.JSON(http.StatusOK, "Hello World")
	})

	// User
	router.POST("/users/login", userController.Login)
	router.POST("/users", userController.CreateAccount)
	router.DELETE("/users/:id", userController.DeleteAccount)
	router.GET("/users/employee", userController.GetEmployee)
	router.Static("/attachments", "./attachments")

	// Task
	router.POST("/tasks", taskController.CreateTask)
	router.DELETE("/tasks/:id", taskController.DeleteTask)
	router.PATCH("/tasks/:id/submit", taskController.SubmitTask)
	router.PATCH("/tasks/:id/reject", taskController.RejectTask)
	router.PATCH("/tasks/:id/fix", taskController.FixTask)
	router.PATCH("/tasks/:id/approve", taskController.ApproveTask)
	router.GET("/tasks/:id", taskController.FindById)

	// Start server
	router.Run("192.168.110.216:8080")
}
