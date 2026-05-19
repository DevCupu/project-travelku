package paket

import (
	"net/http"
	"strconv"

	"backend-travelku/internal/helper"

	"github.com/gin-gonic/gin"
)

type Controller struct {
	service Service
}

func NewController(service Service) *Controller {
	return &Controller{service: service}
}

// Create godoc
// @Summary     Tambah paket wisata
// @Description Buat paket wisata baru dengan kuota dan harga
// @Tags        Pakets
// @Accept      json
// @Produce     json
// @Security    BearerAuth
// @Param       request body     CreatePaketRequest true "Payload pembuatan paket"
// @Success     201     {object} helper.SuccessResponse
// @Failure     400     {object} helper.ErrorResponse
// @Router      /pakets [post]
func (ctrl *Controller) Create(c *gin.Context) {
	var req CreatePaketRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		helper.ErrorJSON(c, http.StatusBadRequest, "Invalid request", err.Error())
		return
	}

	p, err := ctrl.service.Create(&req)
	if err != nil {
		helper.ErrorJSON(c, http.StatusBadRequest, "Failed to create paket", err.Error())
		return
	}

	helper.SuccessJSON(c, http.StatusCreated, "Paket wisata berhasil dibuat", p)
}

// List godoc
// @Summary     List paket wisata
// @Description Ambil daftar paket wisata, bisa difilter hanya yang aktif
// @Tags        Pakets
// @Produce     json
// @Param       active_only query    bool false "Filter hanya paket aktif (true/false)"
// @Success     200         {object} helper.SuccessResponse
// @Failure     500         {object} helper.ErrorResponse
// @Router      /pakets [get]
func (ctrl *Controller) List(c *gin.Context) {
	activeOnly, _ := strconv.ParseBool(c.Query("active_only"))
	pakets, err := ctrl.service.List(activeOnly)
	if err != nil {
		helper.ErrorJSON(c, http.StatusInternalServerError, "Failed to list pakets", err.Error())
		return
	}

	helper.SuccessJSON(c, http.StatusOK, "Berhasil mengambil daftar paket", pakets)
}

// Get godoc
// @Summary     Detail paket wisata
// @Description Ambil data detail paket wisata berdasarkan ID
// @Tags        Pakets
// @Produce     json
// @Param       id path     string true "Paket ID"
// @Success     200 {object} helper.SuccessResponse
// @Failure     404 {object} helper.ErrorResponse
// @Router      /pakets/{id} [get]
func (ctrl *Controller) Get(c *gin.Context) {
	id := c.Param("id")
	p, err := ctrl.service.GetByID(id)
	if err != nil {
		helper.ErrorJSON(c, http.StatusNotFound, "Paket not found", err.Error())
		return
	}

	helper.SuccessJSON(c, http.StatusOK, "Berhasil mengambil data paket", p)
}

// Update godoc
// @Summary     Edit paket wisata
// @Description Perbarui informasi paket wisata yang ada
// @Tags        Pakets
// @Accept      json
// @Produce     json
// @Security    BearerAuth
// @Param       id      path     string             true "Paket ID"
// @Param       request body     UpdatePaketRequest true "Payload pembaruan paket"
// @Success     200     {object} helper.SuccessResponse
// @Failure     400     {object} helper.ErrorResponse
// @Router      /pakets/{id} [put]
func (ctrl *Controller) Update(c *gin.Context) {
	id := c.Param("id")
	var req UpdatePaketRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		helper.ErrorJSON(c, http.StatusBadRequest, "Invalid request", err.Error())
		return
	}

	p, err := ctrl.service.Update(id, &req)
	if err != nil {
		helper.ErrorJSON(c, http.StatusBadRequest, "Failed to update paket", err.Error())
		return
	}

	helper.SuccessJSON(c, http.StatusOK, "Paket wisata berhasil diperbarui", p)
}

// Delete godoc
// @Summary     Hapus paket wisata
// @Description Hapus paket wisata berdasarkan ID
// @Tags        Pakets
// @Produce     json
// @Security    BearerAuth
// @Param       id path     string true "Paket ID"
// @Success     200 {object} helper.SuccessResponse
// @Failure     400 {object} helper.ErrorResponse
// @Router      /pakets/{id} [delete]
func (ctrl *Controller) Delete(c *gin.Context) {
	id := c.Param("id")
	if err := ctrl.service.Delete(id); err != nil {
		helper.ErrorJSON(c, http.StatusBadRequest, "Failed to delete paket", err.Error())
		return
	}

	helper.SuccessEmptyJSON(c, http.StatusOK, "Paket wisata berhasil dihapus")
}
