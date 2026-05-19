package booking

import (
	"net/http"

	"backend-travelku/internal/helper"
	"backend-travelku/pkg/logger"

	"github.com/gin-gonic/gin"
)

// Controller adalah HTTP handler untuk fitur booking.
type Controller struct {
	service Service
}

func NewController(service Service) *Controller {
	return &Controller{service: service}
}

// CreateBooking godoc
// @Summary     Tambah booking
// @Description Tambah pemesanan baru (status otomatis MENUNGGU)
// @Tags        Bookings
// @Accept      json
// @Produce     json
// @Param       request body     CreateRequest true "Booking payload"
// @Success     201     {object} helper.SuccessResponse
// @Failure     400     {object} helper.ErrorResponse
// @Failure     500     {object} helper.ErrorResponse
// @Router      /bookings [post]
func (h *Controller) CreateBooking(c *gin.Context) {
	var req CreateRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		logger.Warn("Invalid create booking request: " + err.Error())
		helper.ErrorJSON(c, http.StatusBadRequest, "Invalid request", err.Error())
		return
	}

	b, err := h.service.Create(&req)
	if err != nil {
		helper.ErrorJSON(c, http.StatusBadRequest, "Failed to create booking", err.Error())
		return
	}

	helper.SuccessJSON(c, http.StatusCreated, "Booking created successfully", ToResponse(b))
}

// ListBookings godoc
// @Summary     List bookings
// @Description Lihat daftar booking (terbaru di atas) dengan filter dan pagenasi optional
// @Tags        Bookings
// @Produce     json
// @Param       status       query    string false "MENUNGGU|DIKONFIRMASI|SELESAI|DIBATALKAN"
// @Param       paket_wisata query    string false "Nama paket"
// @Param       date_from    query    string false "YYYY-MM-DD"
// @Param       date_to      query    string false "YYYY-MM-DD"
// @Param       page         query    int    false "Nomor halaman (default 1)"
// @Param       limit        query    int    false "Jumlah data per halaman (default 10)"
// @Success     200 {object} helper.SuccessResponse
// @Failure     400 {object} helper.ErrorResponse
// @Failure     500 {object} helper.ErrorResponse
// @Router      /bookings [get]
func (h *Controller) ListBookings(c *gin.Context) {
	var q FilterQuery
	_ = c.ShouldBindQuery(&q)

	bookings, total, err := h.service.List(q)
	if err != nil {
		helper.ErrorJSON(c, http.StatusBadRequest, "Failed to list bookings", err.Error())
		return
	}

	responses := make([]Response, len(bookings))
	for i := range bookings {
		responses[i] = ToResponse(&bookings[i])
	}

	payload := ListResponse{Total: total, Count: len(responses), Bookings: responses}
	helper.SuccessJSON(c, http.StatusOK, "Bookings retrieved successfully", payload)
}

// UpdateBooking godoc
// @Summary     Edit booking
// @Description Edit pemesanan yang ada (tidak termasuk ubah status)
// @Tags        Bookings
// @Accept      json
// @Produce     json
// @Param       id      path     string        true "Booking ID"
// @Param       request body     UpdateRequest true "Booking update payload"
// @Success     200 {object} helper.SuccessResponse
// @Failure     400 {object} helper.ErrorResponse
// @Failure     404 {object} helper.ErrorResponse
// @Failure     500 {object} helper.ErrorResponse
// @Router      /bookings/{id} [put]
func (h *Controller) UpdateBooking(c *gin.Context) {
	id := c.Param("id")

	var req UpdateRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		logger.Warn("Invalid update booking request: " + err.Error())
		helper.ErrorJSON(c, http.StatusBadRequest, "Invalid request", err.Error())
		return
	}

	b, err := h.service.Update(id, &req)
	if err != nil {
		status := http.StatusBadRequest
		if err.Error() == "booking not found" {
			status = http.StatusNotFound
		}
		helper.ErrorJSON(c, status, "Failed to update booking", err.Error())
		return
	}

	helper.SuccessJSON(c, http.StatusOK, "Booking updated successfully", ToResponse(b))
}

// DeleteBooking godoc
// @Summary     Hapus booking
// @Description Hapus pemesanan yang ada (tidak bisa jika status final)
// @Tags        Bookings
// @Produce     json
// @Param       id path     string true "Booking ID"
// @Success     200 {object} helper.SuccessResponse
// @Failure     400 {object} helper.ErrorResponse
// @Failure     404 {object} helper.ErrorResponse
// @Failure     500 {object} helper.ErrorResponse
// @Router      /bookings/{id} [delete]
func (h *Controller) DeleteBooking(c *gin.Context) {
	id := c.Param("id")

	if err := h.service.Delete(id); err != nil {
		status := http.StatusBadRequest
		if err.Error() == "booking not found" {
			status = http.StatusNotFound
		}
		helper.ErrorJSON(c, status, "Failed to delete booking", err.Error())
		return
	}

	helper.SuccessEmptyJSON(c, http.StatusOK, "Booking deleted successfully")
}

// UpdateBookingStatus godoc
// @Summary     Ubah status booking
// @Description Ubah status mengikuti alur: MENUNGGU→(DIKONFIRMASI|DIBATALKAN), DIKONFIRMASI→(SELESAI|DIBATALKAN)
// @Tags        Bookings
// @Accept      json
// @Produce     json
// @Param       id      path     string              true "Booking ID"
// @Param       request body     UpdateStatusRequest true "Status payload"
// @Success     200 {object} helper.SuccessResponse
// @Failure     400 {object} helper.ErrorResponse
// @Failure     404 {object} helper.ErrorResponse
// @Failure     500 {object} helper.ErrorResponse
// @Router      /bookings/{id}/status [patch]
func (h *Controller) UpdateBookingStatus(c *gin.Context) {
	id := c.Param("id")

	var req UpdateStatusRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		logger.Warn("Invalid update booking status request: " + err.Error())
		helper.ErrorJSON(c, http.StatusBadRequest, "Invalid request", err.Error())
		return
	}

	b, err := h.service.UpdateStatus(id, req.Status)
	if err != nil {
		status := http.StatusBadRequest
		if err.Error() == "booking not found" {
			status = http.StatusNotFound
		}
		helper.ErrorJSON(c, status, "Failed to update booking status", err.Error())
		return
	}

	helper.SuccessJSON(c, http.StatusOK, "Booking status updated successfully", ToResponse(b))
}

// BookingSummary godoc
// @Summary     Ringkasan booking
// @Description Jumlah booking & total estimasi pendapatan (status DIKONFIRMASI+SELESAI), mengikuti filter aktif
// @Tags        Bookings
// @Produce     json
// @Param       status       query    string false "MENUNGGU|DIKONFIRMASI|SELESAI|DIBATALKAN"
// @Param       paket_wisata query    string false "Nama paket"
// @Param       date_from    query    string false "YYYY-MM-DD"
// @Param       date_to      query    string false "YYYY-MM-DD"
// @Success     200 {object} helper.SuccessResponse
// @Failure     400 {object} helper.ErrorResponse
// @Failure     500 {object} helper.ErrorResponse
// @Router      /bookings/summary [get]
func (h *Controller) BookingSummary(c *gin.Context) {
	var q FilterQuery
	_ = c.ShouldBindQuery(&q)

	count, revenue, err := h.service.Summary(q)
	if err != nil {
		helper.ErrorJSON(c, http.StatusBadRequest, "Failed to get summary", err.Error())
		return
	}

	payload := SummaryResponse{JumlahBooking: int(count), TotalRevenue: revenue}
	helper.SuccessJSON(c, http.StatusOK, "Summary retrieved successfully", payload)
}
