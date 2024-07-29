package handlers

import (
	"net/http"
	"todowithgin/models"

	"github.com/gin-gonic/gin"
)

func (h *TaskHandler) GetTasksHandler(ctx *gin.Context) {
	var tasks []models.Task
	tasks, err := h.Service.GetTasks()
	if err != nil {
		ctx.JSON(http.StatusInternalServerError, gin.H{
			"message": err.Error(),
		})
		return
	}
	ctx.JSON(http.StatusOK, gin.H{
		"message": "success",
		"data":    tasks,
	})
}
