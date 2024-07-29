package handlers

import (
	"errors"
	"net/http"
	"strconv"
	"todowithgin/models"
	"todowithgin/service"

	"github.com/gin-gonic/gin"
)

func (h *TaskHandler) GetTaskHandler(ctx *gin.Context) {
	tid := ctx.Param("id")
	id, err := strconv.Atoi(tid)
	if err != nil {
		ctx.JSON(http.StatusBadRequest, gin.H{
			"message": "Invalid id param",
		})
		return
	}
	var task models.Task
	task, err = h.Service.GetTask(id)
	if err != nil {
		if errors.Is(err, service.ErrorNoRows) {
			ctx.JSON(http.StatusNotFound, gin.H{
				"message": err.Error(),
			})

		} else {
			ctx.JSON(http.StatusInternalServerError, gin.H{
				"message": err.Error(),
			})
		}
		return

	}
	ctx.JSON(http.StatusOK, gin.H{
		"data": task,
	})

}
