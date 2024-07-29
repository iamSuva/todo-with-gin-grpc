package handlers

import (
	"fmt"
	"net/http"
	"strings"
	"github.com/gin-gonic/gin"
)

func (h *TaskHandler) PrivateRouteHandler(ctx *gin.Context) {
	tokenString := ctx.Request.Header.Get("Authorization")
	if tokenString == "" || !strings.HasPrefix(tokenString, "Bearer ") {
		ctx.JSON(http.StatusUnauthorized, gin.H{
			"message": "missing authorization header",
		})
		ctx.Abort()
		return
	}
	token := tokenString[len("Bearer "):]
	auth, err := h.GrpcService.VeriFyTokenHandler(token)

	fmt.Println("auth and err : ", err)
	if err != nil {
		ctx.JSON(http.StatusUnauthorized, gin.H{
			"message": err.Error(),
		})
		ctx.Abort()
		return
	}

	if auth.UserId<=0{
		ctx.Next()

	}else{

		fmt.Println("auth is : ", auth.UserId, auth.Username)
		ctx.Set("username", auth.Username)
		ctx.Set("userId", auth.UserId)
		ctx.Next()
	}
}
