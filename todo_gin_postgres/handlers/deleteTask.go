package handlers

import (
	"errors"

	"net/http"
	"strconv"
	"todowithgin/service"
	"github.com/gin-gonic/gin"
)

func (h TaskHandler) DeleteTaskHandler(ctx *gin.Context) {
	tid := ctx.Param("id")
	id, err := strconv.Atoi(tid)
	if err != nil {
		ctx.JSON(http.StatusBadRequest, gin.H{
			"message": "Invalid id param",
		})
		return
	}
	//get user id
	userId, _ := ctx.Get("userId")

	
	useridInt, _ := userId.(int)
	err = h.Service.DeleteTask(id, useridInt)
	if err != nil {
		if errors.Is(err, service.ErrorNoRows) {
			ctx.JSON(http.StatusOK, gin.H{
				"message": "task is deleted",
			})
			return
		}
		ctx.JSON(http.StatusInternalServerError, gin.H{
			"message": err.Error(),
		})
		return
	}
	ctx.JSON(http.StatusOK, gin.H{
		"message": "task is deleted",
	})
}
