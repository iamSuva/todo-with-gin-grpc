package handlers

import (
	"errors"
	"net/http"
	"time"
	"todowithgin/models"
	"todowithgin/service"
	"todowithgin/utils"

	"github.com/gin-gonic/gin"
)


func (h *TaskHandler) CreateTaskHandler(ctx *gin.Context) {
	var newTask models.Task
	err := ctx.BindJSON(&newTask)
	if err != nil {
		ctx.JSON(http.StatusBadRequest, gin.H{
			"message": "Invalid input",
		})
		return
	}
	validationerr,err:=utils.ValidateTaskHandler(newTask)
	
	if err!=nil{
		ctx.JSON(http.StatusBadRequest,validationerr)
		return
	}

	currTime := time.Now().UTC()
	newTask.CreatedAt_UTC = currTime
	newTask.UpdatedAt_UTC = currTime
	userId, exist := ctx.Get("userId")
	if !exist {
		ctx.JSON(http.StatusUnauthorized, gin.H{
			"message": "user not authorized",
		})
		return
	}
	useridInt, _:= userId.(int)
	
	newTask.UserId = useridInt
	err = h.Service.CreateTask(newTask)
	if err != nil {
		if errors.Is(err, service.ErrorUniqueTitle) {
			ctx.JSON(http.StatusConflict, gin.H{
				"message": err.Error(),
			})
			return
		} else {
			ctx.JSON(http.StatusInternalServerError, gin.H{
				"message": err.Error(),
			})

			return
		}
	}
	ctx.JSON(http.StatusOK, gin.H{
		"message": "Task is added",
	})

}
