package handlers

import (
	"github.com/JonasLindermayr/FileBeam/handlers/controller"
	"github.com/gin-gonic/gin"
)

func GetUserHandler(c *gin.Context) {
	controller.GetUser(c)
}

func VerifyUserOTP(g *gin.Context) {
	controller.VerifyOTP(g)
}

func LogoutUserHandler(g *gin.Context) {
	controller.LogoutUser(g)
}
