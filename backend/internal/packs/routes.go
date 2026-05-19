package paket

import (
	"backend-travelku/internal/middleware"
	"github.com/gin-gonic/gin"
)

func RegisterRoutes(r *gin.Engine, ctrl *Controller) {
	// Publik / User bisa melihat daftar paket aktif
	r.GET("/api/v1/pakets", ctrl.List)
	r.GET("/api/v1/pakets/:id", ctrl.Get)

	// Admin / Staff routes
	protected := r.Group("/api/v1/pakets")
	protected.Use(middleware.AuthMiddleware())
	{
		protected.POST("", ctrl.Create)
		protected.PUT("/:id", ctrl.Update)
		protected.DELETE("/:id", ctrl.Delete)
	}
}
