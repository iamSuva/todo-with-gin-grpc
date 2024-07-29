package handlers

import (
	"errors"
	"net/http"
	"strconv"
	"todowithgin/models"
	"todowithgin/service"
	"todowithgin/utils"

	"github.com/gin-gonic/gin"
)

func (h *TaskHandler) UpdateTaskHandler(ctx *gin.Context) {
	tid := ctx.Param("id")
	id, err := strconv.Atoi(tid)
	if err != nil {
		ctx.JSON(http.StatusBadRequest, gin.H{
			"message": "Invalid id param",
		})
		return
	}
	var updateTask models.Task
	err = ctx.BindJSON(&updateTask)
	if err != nil {
		ctx.JSON(http.StatusBadRequest, gin.H{
			"message": "Invalid input",
		})
		return
	}
	validationerr, err := utils.ValidateTaskHandler(updateTask)
	if err != nil {
		ctx.JSON(http.StatusBadRequest, validationerr)
		return
	}
	userId, exist := ctx.Get("userId")
	if !exist {
		ctx.JSON(http.StatusUnauthorized, gin.H{
			"message": "user not authorized",
		})
		return
	}
	useridInt, _ := userId.(int)

	err = h.Service.UpdateTask(updateTask, id, useridInt)
	if err != nil {
		if errors.Is(err, service.ErrorUniqueTitle) {
			ctx.JSON(http.StatusConflict, gin.H{
				"message": err.Error(),
			})
			return

		} else if errors.Is(err, service.ErrorNoRows) {
			ctx.JSON(http.StatusUnauthorized, gin.H{
				"message": "unauthorized",
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
		"message": "Task is updated",
	})
}
