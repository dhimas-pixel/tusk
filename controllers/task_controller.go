package controllers

import (
	"fmt"
	"github.com/gin-gonic/gin"
	"gorm.io/gorm"
	"net/http"
	"os"
	"strconv"
	"tusk/models"
)

type TaskController struct {
	DB *gorm.DB
}

func (t *TaskController) CreateTask(c *gin.Context) {

	task := models.Task{}

	errBindJson := c.ShouldBindJSON(&task)
	if errBindJson != nil {
		c.JSON(http.StatusInternalServerError, models.Response{
			StatusCode: http.StatusInternalServerError,
			Message:    "Something went wrong",
			Data:       nil,
		})
		fmt.Println("Error: errBindJson", errBindJson.Error())
		return
	}

	errDB := t.DB.Create(&task).Error
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
		Message:    "Task created successfully",
		Data:       nil,
	})
}

func (t *TaskController) DeleteTask(c *gin.Context) {
	id := c.Param("id")
	task := models.Task{}

	err := t.DB.First(&task, id).Error
	if err != nil {
		c.JSON(http.StatusNotFound, models.Response{
			StatusCode: http.StatusNotFound,
			Message:    "Task not found",
			Data:       nil,
		})
		fmt.Println("Error: err", err.Error())
		return
	}

	errDB := t.DB.Delete(&models.Task{}, id).Error
	if errDB != nil {
		c.JSON(http.StatusInternalServerError, models.Response{
			StatusCode: http.StatusInternalServerError,
			Message:    "Something went wrong",
			Data:       nil,
		})
		fmt.Println("Error: errDB", errDB.Error())
		return
	}

	if task.Attachment != "" {
		os.Remove("attachments/" + task.Attachment)
	}

	c.JSON(http.StatusOK, models.Response{
		StatusCode: http.StatusOK,
		Message:    "Task deleted successfully",
		Data:       nil,
	})
}

func (t *TaskController) SubmitTask(c *gin.Context) {
	id := c.Param("id")
	task := models.Task{}
	submitDate := c.PostForm("submitDate")
	file, errFile := c.FormFile("attachment")

	if errFile != nil {
		c.JSON(http.StatusInternalServerError, models.Response{
			StatusCode: http.StatusInternalServerError,
			Message:    "Something went wrong",
			Data:       nil,
		})
		fmt.Println("Error: errFile", errFile.Error())
		return
	}

	err := t.DB.First(&task, id).Error
	if err != nil {
		c.JSON(http.StatusNotFound, models.Response{
			StatusCode: http.StatusNotFound,
			Message:    "Task not found",
			Data:       nil,
		})
		fmt.Println("Error: err", err.Error())
		return
	}

	// remove old attachment
	attachment := task.Attachment
	fileInfo, _ := os.Stat("attachments/" + attachment)
	if fileInfo != nil {
		os.Remove("attachments/" + attachment)
	}

	// create new attachment
	attachment = file.Filename
	errSave := c.SaveUploadedFile(file, "attachments/"+attachment)
	if errSave != nil {
		c.JSON(http.StatusInternalServerError, models.Response{
			StatusCode: http.StatusInternalServerError,
			Message:    "Something went wrong",
			Data:       nil,
		})
		fmt.Println("Error: errSave", errSave.Error())
		return
	}

	errDB := t.DB.Where("id = ?", id).Updates(models.Task{
		Status:     "Review",
		SubmitDate: submitDate,
		Attachment: attachment,
	}).Error
	if errDB != nil {
		c.JSON(http.StatusInternalServerError, models.Response{
			StatusCode: http.StatusInternalServerError,
			Message:    "Something went wrong",
			Data:       nil,
		})
		fmt.Println("Error: errDB", errDB.Error())
		return
	}

	c.JSON(http.StatusOK, models.Response{
		StatusCode: http.StatusOK,
		Message:    "Task submitted to review successfully",
		Data:       nil,
	})
}

func (t *TaskController) RejectTask(c *gin.Context) {
	id := c.Param("id")
	task := models.Task{}
	rejectedDate := c.PostForm("rejectedDate")
	reason := c.PostForm("reason")

	err := t.DB.First(&task, id).Error
	if err != nil {
		c.JSON(http.StatusNotFound, models.Response{
			StatusCode: http.StatusNotFound,
			Message:    "Task not found",
			Data:       nil,
		})
		fmt.Println("Error: err", err.Error())
		return
	}

	errDB := t.DB.Where("id = ?", id).Updates(models.Task{
		Status:     "Rejected",
		Reason:     reason,
		RejectDate: rejectedDate,
	}).Error
	if errDB != nil {
		c.JSON(http.StatusInternalServerError, models.Response{
			StatusCode: http.StatusInternalServerError,
			Message:    "Something went wrong",
			Data:       nil,
		})
		fmt.Println("Error: errDB", errDB.Error())
		return
	}

	c.JSON(http.StatusOK, models.Response{
		StatusCode: http.StatusOK,
		Message:    "Task review to rejected successfully",
		Data:       nil,
	})
}

func (t *TaskController) FixTask(c *gin.Context) {
	id := c.Param("id")
	revision, errConv := strconv.Atoi(c.PostForm("revision"))
	if errConv != nil {
		c.JSON(http.StatusInternalServerError, models.Response{
			StatusCode: http.StatusInternalServerError,
			Message:    "Something went wrong",
			Data:       nil,
		})
		fmt.Println("Error: errConv", errConv.Error())
		return
	}

	err := t.DB.First(&models.Task{}, id).Error
	if err != nil {
		c.JSON(http.StatusNotFound, models.Response{
			StatusCode: http.StatusNotFound,
			Message:    "Task not found",
			Data:       nil,
		})
		fmt.Println("Error: err", err.Error())
		return
	}

	errDB := t.DB.Where("id = ?", id).Updates(models.Task{
		Status:   "Queue",
		Revision: int8(revision),
	}).Error
	if errDB != nil {
		c.JSON(http.StatusInternalServerError, models.Response{
			StatusCode: http.StatusInternalServerError,
			Message:    "Something went wrong",
			Data:       nil,
		})
		fmt.Println("Error: errDB", errDB.Error())
		return
	}

	c.JSON(http.StatusOK, models.Response{
		StatusCode: http.StatusOK,
		Message:    "Task fix to queue successfully",
		Data:       nil,
	})
}

func (t *TaskController) ApproveTask(c *gin.Context) {
	id := c.Param("id")
	approvedDate := c.PostForm("approvedDate")

	err := t.DB.First(&models.Task{}, id).Error
	if err != nil {
		c.JSON(http.StatusNotFound, models.Response{
			StatusCode: http.StatusNotFound,
			Message:    "Task not found",
			Data:       nil,
		})
		fmt.Println("Error: err", err.Error())
		return
	}

	errDB := t.DB.Where("id = ?", id).Updates(models.Task{
		Status:       "Approved",
		ApprovedDate: approvedDate,
	}).Error
	if errDB != nil {
		c.JSON(http.StatusInternalServerError, models.Response{
			StatusCode: http.StatusInternalServerError,
			Message:    "Something went wrong",
			Data:       nil,
		})
		fmt.Println("Error: errDB", errDB.Error())
		return
	}

	c.JSON(http.StatusOK, models.Response{
		StatusCode: http.StatusOK,
		Message:    "Task approved successfully",
		Data:       nil,
	})
}

func (t *TaskController) FindById(c *gin.Context) {

	task := models.Task{}
	id := c.Param("id")

	err := t.DB.First(&models.Task{}, id).Error
	if err != nil {
		c.JSON(http.StatusNotFound, models.Response{
			StatusCode: http.StatusNotFound,
			Message:    "Task not found",
			Data:       nil,
		})
		fmt.Println("Error: err", err.Error())
		return
	}

	errDB := t.DB.Preload("User", func(db *gorm.DB) *gorm.DB {
		return db.Select("id, name, email, role")
	}).Find(&task, id).Error
	if errDB != nil {
		c.JSON(http.StatusInternalServerError, models.Response{
			StatusCode: http.StatusInternalServerError,
			Message:    "Something went wrong",
			Data:       nil,
		})
		return
	}

	c.JSON(http.StatusOK, models.Response{
		StatusCode: http.StatusOK,
		Message:    "Task found successfully",
		Data:       task,
	})
}
