package booking

import "github.com/gin-gonic/gin"

// RegisterRoutes mendaftarkan semua route untuk fitur booking.
func RegisterRoutes(router *gin.Engine, controller *Controller) {
	bookings := router.Group("/api/v1/bookings")
	{
		bookings.POST("", controller.CreateBooking)
		bookings.GET("", controller.ListBookings)
		bookings.GET("/summary", controller.BookingSummary)

		bookings.PUT("/:id", controller.UpdateBooking)
		bookings.DELETE("/:id", controller.DeleteBooking)
		bookings.PATCH("/:id/status", controller.UpdateBookingStatus)
	}
}
