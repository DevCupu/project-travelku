package user

import (
	"backend-travelku/internal/middleware"

	"github.com/gin-gonic/gin"
)

// RegisterRoutes mendaftarkan semua route untuk fitur user.
func RegisterRoutes(router *gin.Engine, controller *Controller) {
	users := router.Group("/api/v1/users")
	{
		// Protected routes (perlu authentication)
		protected := users.Use(middleware.AuthMiddleware())
		{
			protected.GET("", controller.GetAllUsers)
			protected.GET("/:id", controller.GetUser)
			protected.PUT("/:id", controller.UpdateProfile)                   // update profile (name, email, phone)
			protected.POST("/:id/change-password", controller.ChangePassword) // change password
			protected.DELETE("/:id", controller.DeleteUser)
		}
	}
}
