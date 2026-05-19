package auth

import (
	"github.com/gin-gonic/gin"
)

// RegisterRoutes mendaftarkan semua route untuk fitur auth.
func RegisterRoutes(router *gin.Engine, controller *Controller) {
	authGroup := router.Group("/api/v1/auth")
	{
		authGroup.POST("/register", controller.Register)
		authGroup.POST("/login", controller.Login)
	}
}
