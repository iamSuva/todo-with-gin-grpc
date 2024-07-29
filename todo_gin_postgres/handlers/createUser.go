package handlers

import (
	"errors"
	"fmt"
	"net/http"
	"todowithgin/models"
	"todowithgin/service"
	"todowithgin/utils"

	"github.com/gin-gonic/gin"
)

func (h *TaskHandler) CreateUserHandler(ctx *gin.Context) {
	var user models.User
	err := ctx.BindJSON(&user)
	if err != nil {
		ctx.JSON(http.StatusBadRequest, gin.H{
			"message": "Invalid input",
		})
		return
	}

	validationerr, err := utils.ValidateUserHandler(user)
	if err != nil {
		ctx.JSON(http.StatusBadRequest, validationerr)
		return;
	}
	fmt.Println(user)
	hash, _ := utils.HashedPassword(user.Password)
	user.Password = hash

	id, err := h.Service.SignUpUser(user)
	if err != nil {
		if errors.Is(err, service.ErrorUniqueUsername) {
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
		"message": "user created",
		"userid":  id,
	})

}
