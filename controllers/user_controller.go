package controllers

import (
	"fmt"
	"github.com/gin-gonic/gin"
	"golang.org/x/crypto/bcrypt"
	"gorm.io/gorm"
	"net/http"
	"tusk/models"
)

type UserController struct {
	DB *gorm.DB
}

func (u *UserController) Login(c *gin.Context) {

	user := models.User{}

	errBindJson := c.ShouldBindJSON(&user)
	if errBindJson != nil {
		c.JSON(http.StatusInternalServerError, models.Response{
			StatusCode: http.StatusInternalServerError,
			Message:    "Something went wrong",
			Data:       nil,
		})
		fmt.Println("Error: errBindJson", errBindJson.Error())

		return
	}

	password := user.Password
	errDB := u.DB.Where("email = ?", user.Email).Take(&user).Error
	if errDB != nil {
		c.JSON(http.StatusUnauthorized, models.Response{
			StatusCode: http.StatusUnauthorized,
			Message:    "Email or Password is incorrect",
			Data:       nil,
		})
		return
	}

	errHash := bcrypt.CompareHashAndPassword([]byte(user.Password), []byte(password))
	if errHash != nil {
		c.JSON(http.StatusUnauthorized, models.Response{
			StatusCode: http.StatusUnauthorized,
			Message:    "Email or Password is incorrect",
			Data:       nil,
		})
		return
	}

	c.JSON(http.StatusOK, models.Response{
		StatusCode: http.StatusOK,
		Message:    "Login success",
		Data: gin.H{
			"role":  user.Role,
			"name":  user.Name,
			"email": user.Email,
		},
	})
}

func (u *UserController) CreateAccount(c *gin.Context) {

	user := models.User{}

	errBindJson := c.ShouldBindJSON(&user)
	if errBindJson != nil {
		c.JSON(http.StatusInternalServerError, models.Response{
			StatusCode: http.StatusInternalServerError,
			Message:    "Something went wrong",
			Data:       nil,
		})
		fmt.Println("Error: errBindJson", errBindJson.Error())
		return
	}

	emailExist := u.DB.Where("email = ?", user.Email).First(&user).RowsAffected != 0
	if emailExist {
		c.JSON(http.StatusBadRequest, models.Response{
			StatusCode: http.StatusBadRequest,
			Message:    "Email already exist",
			Data:       nil,
		})
		return
	}

	hashedPasswordBytes, _ := bcrypt.GenerateFromPassword([]byte("123456"), bcrypt.DefaultCost)

	user.Password = string(hashedPasswordBytes)
	user.Role = "Employee"

	errDB := u.DB.Create(&user).Error
	if errDB != nil {
		c.JSON(http.StatusInternalServerError, models.Response{
			StatusCode: http.StatusInternalServerError,
			Message:    "Something went wrong",
			Data:       nil,
		})
		fmt.Println("Error: errDB", errDB.Error())
		return
	}

	c.JSON(http.StatusCreated, models.Response{
		StatusCode: http.StatusCreated,
		Message:    "Account created successfully",
		Data: gin.H{
			"role":  user.Role,
			"name":  user.Name,
			"email": user.Email,
		},
	})
}
