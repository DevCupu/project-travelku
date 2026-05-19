package booking

import (
	"time"
)

// CreateRequest untuk membuat booking baru.
// Status otomatis MENUNGGU.
type CreateRequest struct {
	NamaPemesan   string  `json:"nama_pemesan" binding:"required,min=3,max=120"`
	Kontak        string  `json:"kontak" binding:"required,min=3,max=120"`
	PaketID       string  `json:"paket_id" binding:"required"`
	JumlahPeserta int     `json:"jumlah_peserta" binding:"required,min=1"`
	Catatan       *string `json:"catatan,omitempty"`
}

// UpdateRequest untuk edit booking.
// Tidak termasuk status.
type UpdateRequest struct {
	NamaPemesan   string  `json:"nama_pemesan" binding:"required,min=3,max=120"`
	Kontak        string  `json:"kontak" binding:"required,min=3,max=120"`
	PaketID       string  `json:"paket_id" binding:"required"`
	JumlahPeserta int     `json:"jumlah_peserta" binding:"required,min=1"`
	Catatan       *string `json:"catatan,omitempty"`
}

// UpdateStatusRequest untuk ubah status booking.
type UpdateStatusRequest struct {
	Status BookingStatus `json:"status" binding:"required"`
}

// FilterQuery untuk filter list & summary.
// date_from/date_to format: YYYY-MM-DD
type FilterQuery struct {
	Status   string `form:"status"`
	Paket    string `form:"paket_wisata"`
	DateFrom string `form:"date_from"`
	DateTo   string `form:"date_to"`
	Page     int    `form:"page"`
	Limit    int    `form:"limit"`
}

type Response struct {
	ID               string        `json:"id"`
	NamaPemesan      string        `json:"nama_pemesan"`
	Kontak           string        `json:"kontak"`
	PaketID          string        `json:"paket_id"`
	PaketWisata      string        `json:"paket_wisata"`
	TanggalBerangkat string        `json:"tanggal_berangkat"`
	JumlahPeserta    int           `json:"jumlah_peserta"`
	HargaPerOrang    int64         `json:"harga_per_orang"`
	Status           BookingStatus `json:"status"`
	Catatan          *string       `json:"catatan,omitempty"`
	CreatedAt        time.Time     `json:"created_at"`
	UpdatedAt        time.Time     `json:"updated_at"`
}

type ListResponse struct {
	Total    int64      `json:"total"`
	Count    int        `json:"count"`
	Bookings []Response `json:"data"`
}

type SummaryResponse struct {
	JumlahBooking int   `json:"jumlah_booking"`
	TotalRevenue  int64 `json:"total_estimasi_pendapatan"`
}

func ToResponse(b *Booking) Response {
	return Response{
		ID:               b.ID,
		NamaPemesan:      b.NamaPemesan,
		Kontak:           b.Kontak,
		PaketID:          b.PaketID,
		PaketWisata:      b.PaketWisata,
		TanggalBerangkat: b.TanggalBerangkat.Format("2006-01-02"),
		JumlahPeserta:    b.JumlahPeserta,
		HargaPerOrang:    b.HargaPerOrang,
		Status:           b.Status,
		Catatan:          b.Catatan,
		CreatedAt:        b.CreatedAt,
		UpdatedAt:        b.UpdatedAt,
	}
}

