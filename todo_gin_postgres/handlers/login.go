package handlers

import (
	"errors"
	"net/http"
	"todowithgin/models"
	"todowithgin/service"

	"github.com/gin-gonic/gin"
)

func (h *TaskHandler) LoginHandler(ctx *gin.Context) {
	var user models.User
	err := ctx.BindJSON(&user)
	if err != nil {
		ctx.JSON(http.StatusBadRequest, gin.H{
			"message": "Invalid input",
		})
		return
	}

	user, err = h.Service.LoginUser(user)

	if err != nil {
		if errors.Is(err, service.ErrorNoRows) {
			ctx.JSON(http.StatusNotFound, gin.H{
				"message": "User not found",
			})
			return
		} else if errors.Is(err, service.ErrorPasswordNotMatch) {
			ctx.JSON(http.StatusUnauthorized, gin.H{ //unauthor
				"message": err.Error(),
			})
			return
		}
		ctx.JSON(http.StatusInternalServerError, gin.H{ //internal
			"message": err.Error(),
		})
		return

	}

	token, err := h.GrpcService.GetTokenHandler(user.Username, user.UserId)

	if err != nil {
		ctx.JSON(http.StatusInternalServerError, gin.H{
			"message": err.Error(),
		})
		return
	}
	ctx.JSON(http.StatusOK, gin.H{
		"token": token,
	})
}
