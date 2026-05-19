package booking

import (
	"fmt"
	"strings"
	"time"

	"backend-travelku/internal/packs"
	"backend-travelku/pkg/logger"

	"github.com/google/uuid"
	"gorm.io/gorm"
)

// Service adalah core business logic untuk booking.
type Service interface {
	Create(req *CreateRequest) (*Booking, error)
	List(query FilterQuery) ([]Booking, int64, error)
	Update(id string, req *UpdateRequest) (*Booking, error)
	Delete(id string) error
	UpdateStatus(id string, next BookingStatus) (*Booking, error)
	Summary(query FilterQuery) (int64, int64, error)
}

type service struct {
	repo      Repository
	paketRepo paket.Repository
}

func NewService(repo Repository, paketRepo paket.Repository) Service {
	return &service{repo: repo, paketRepo: paketRepo}
}

func (s *service) Create(req *CreateRequest) (*Booking, error) {
	if strings.TrimSpace(req.NamaPemesan) == "" {
		return nil, fmt.Errorf("nama pemesan wajib diisi")
	}
	if strings.TrimSpace(req.Kontak) == "" {
		return nil, fmt.Errorf("kontak wajib diisi")
	}
	if req.JumlahPeserta < 1 {
		return nil, fmt.Errorf("jumlah peserta minimal 1")
	}

	p, err := s.paketRepo.GetByID(req.PaketID)
	if err != nil {
		return nil, fmt.Errorf("paket wisata tidak valid atau tidak ditemukan")
	}

	if !p.IsActive {
		return nil, fmt.Errorf("paket wisata sudah tidak aktif")
	}

	if p.Harga < 0 {
		return nil, fmt.Errorf("harga paket tidak boleh negatif")
	}

	now := time.Now().In(time.Local)
	today := time.Date(now.Year(), now.Month(), now.Day(), 0, 0, 0, 0, now.Location())
	if p.TanggalBerangkat.Before(today) {
		return nil, fmt.Errorf("tanggal keberangkatan paket sudah lewat (tidak boleh di masa lalu)")
	}

	if req.JumlahPeserta > p.SisaKuota {
		return nil, fmt.Errorf("kuota paket tidak mencukupi (sisa kuota: %d)", p.SisaKuota)
	}

	b := &Booking{
		ID:               uuid.NewString(),
		NamaPemesan:      strings.TrimSpace(req.NamaPemesan),
		Kontak:           strings.TrimSpace(req.Kontak),
		PaketID:          p.ID,
		PaketWisata:      p.NamaPaket,
		TanggalBerangkat: p.TanggalBerangkat,
		JumlahPeserta:    req.JumlahPeserta,
		HargaPerOrang:    p.Harga,
		Status:           StatusMenunggu,
		Catatan:          req.Catatan,
	}

	// Kurangi kuota paket
	p.SisaKuota -= req.JumlahPeserta
	if err := s.paketRepo.Update(p); err != nil {
		logger.Error("Failed to update paket kuota: " + err.Error())
		return nil, fmt.Errorf("gagal memperbarui kuota paket")
	}

	if err := s.repo.Create(b); err != nil {
		// Rollback kuota
		p.SisaKuota += req.JumlahPeserta
		_ = s.paketRepo.Update(p)
		logger.Error("Failed to create booking: " + err.Error())
		return nil, fmt.Errorf("failed to create booking")
	}
	return b, nil
}

func (s *service) List(query FilterQuery) ([]Booking, int64, error) {
	filter, err := parseFilter(query)
	if err != nil {
		return nil, 0, err
	}
	total, err := s.repo.Count(filter)
	if err != nil {
		logger.Error("Failed to count bookings: " + err.Error())
		return nil, 0, fmt.Errorf("failed to count bookings")
	}
	bookings, err := s.repo.List(filter)
	if err != nil {
		logger.Error("Failed to list bookings: " + err.Error())
		return nil, 0, fmt.Errorf("failed to list bookings")
	}
	return bookings, total, nil
}

func (s *service) Update(id string, req *UpdateRequest) (*Booking, error) {
	if strings.TrimSpace(req.NamaPemesan) == "" {
		return nil, fmt.Errorf("nama pemesan wajib diisi")
	}
	if strings.TrimSpace(req.Kontak) == "" {
		return nil, fmt.Errorf("kontak wajib diisi")
	}
	if req.JumlahPeserta < 1 {
		return nil, fmt.Errorf("jumlah peserta minimal 1")
	}

	b, err := s.repo.GetByID(id)
	if err != nil {
		if err == gorm.ErrRecordNotFound {
			return nil, fmt.Errorf("booking not found")
		}
		return nil, fmt.Errorf("failed to get booking")
	}

	if b.Status.IsTerminal() {
		return nil, fmt.Errorf("booking tidak bisa diubah karena status sudah final")
	}

	// Cek apakah paket berubah atau jumlah peserta berubah
	if b.PaketID == req.PaketID {
		diff := req.JumlahPeserta - b.JumlahPeserta
		if diff != 0 {
			p, err := s.paketRepo.GetByID(b.PaketID)
			if err != nil {
				return nil, fmt.Errorf("paket wisata tidak ditemukan")
			}
			if diff > p.SisaKuota {
				return nil, fmt.Errorf("kuota paket tidak mencukupi (sisa kuota: %d)", p.SisaKuota)
			}
			p.SisaKuota -= diff
			if err := s.paketRepo.Update(p); err != nil {
				return nil, fmt.Errorf("gagal memperbarui kuota paket")
			}
			b.HargaPerOrang = p.Harga
			b.TanggalBerangkat = p.TanggalBerangkat
			b.PaketWisata = p.NamaPaket
		}
	} else {
		// Paket berubah!
		// 1. Kembalikan kuota paket lama
		oldP, err := s.paketRepo.GetByID(b.PaketID)
		if err == nil {
			oldP.SisaKuota += b.JumlahPeserta
			_ = s.paketRepo.Update(oldP)
		}

		// 2. Ambil paket baru dan potong kuota
		newP, err := s.paketRepo.GetByID(req.PaketID)
		if err != nil {
			return nil, fmt.Errorf("paket wisata baru tidak ditemukan")
		}
		if !newP.IsActive {
			return nil, fmt.Errorf("paket wisata baru sudah tidak aktif")
		}
		if newP.Harga < 0 {
			return nil, fmt.Errorf("harga paket baru tidak boleh negatif")
		}
		now := time.Now().In(time.Local)
		today := time.Date(now.Year(), now.Month(), now.Day(), 0, 0, 0, 0, now.Location())
		if newP.TanggalBerangkat.Before(today) {
			return nil, fmt.Errorf("tanggal keberangkatan paket baru sudah lewat (tidak boleh di masa lalu)")
		}
		if req.JumlahPeserta > newP.SisaKuota {
			return nil, fmt.Errorf("kuota paket baru tidak mencukupi (sisa kuota: %d)", newP.SisaKuota)
		}
		newP.SisaKuota -= req.JumlahPeserta
		if err := s.paketRepo.Update(newP); err != nil {
			return nil, fmt.Errorf("gagal memperbarui kuota paket baru")
		}

		b.PaketID = newP.ID
		b.PaketWisata = newP.NamaPaket
		b.TanggalBerangkat = newP.TanggalBerangkat
		b.HargaPerOrang = newP.Harga
	}

	b.NamaPemesan = strings.TrimSpace(req.NamaPemesan)
	b.Kontak = strings.TrimSpace(req.Kontak)
	b.JumlahPeserta = req.JumlahPeserta
	b.Catatan = req.Catatan

	if err := s.repo.Update(b); err != nil {
		logger.Error("Failed to update booking: " + err.Error())
		return nil, fmt.Errorf("failed to update booking")
	}

	return b, nil
}

func (s *service) Delete(id string) error {
	b, err := s.repo.GetByID(id)
	if err != nil {
		if err == gorm.ErrRecordNotFound {
			return fmt.Errorf("booking not found")
		}
		return fmt.Errorf("failed to get booking")
	}

	if b.Status.IsTerminal() {
		return fmt.Errorf("booking tidak bisa dihapus karena status sudah final")
	}

	// Kembalikan kuota
	if p, err := s.paketRepo.GetByID(b.PaketID); err == nil {
		p.SisaKuota += b.JumlahPeserta
		_ = s.paketRepo.Update(p)
	}

	if err := s.repo.Delete(id); err != nil {
		logger.Error("Failed to delete booking: " + err.Error())
		return fmt.Errorf("failed to delete booking")
	}
	return nil
}

func (s *service) UpdateStatus(id string, next BookingStatus) (*Booking, error) {
	b, err := s.repo.GetByID(id)
	if err != nil {
		if err == gorm.ErrRecordNotFound {
			return nil, fmt.Errorf("booking not found")
		}
		return nil, fmt.Errorf("failed to get booking")
	}

	if b.Status.IsTerminal() {
		return nil, fmt.Errorf("booking tidak bisa diubah karena status sudah final")
	}

	if !isValidTransition(b.Status, next) {
		return nil, fmt.Errorf("transisi status tidak valid")
	}

	// Jika dibatalkan, kembalikan kuota paket
	if next == StatusDibatalkan {
		if p, err := s.paketRepo.GetByID(b.PaketID); err == nil {
			p.SisaKuota += b.JumlahPeserta
			_ = s.paketRepo.Update(p)
		}
	}

	b.Status = next
	if err := s.repo.Update(b); err != nil {
		logger.Error("Failed to update booking status: " + err.Error())
		return nil, fmt.Errorf("failed to update booking status")
	}
	return b, nil
}

func (s *service) Summary(query FilterQuery) (int64, int64, error) {
	filter, err := parseFilter(query)
	if err != nil {
		return 0, 0, err
	}
	count, err := s.repo.Count(filter)
	if err != nil {
		logger.Error("Failed to count bookings: " + err.Error())
		return 0, 0, fmt.Errorf("failed to calculate summary")
	}
	revenue, err := s.repo.SumRevenue(filter)
	if err != nil {
		logger.Error("Failed to sum revenue: " + err.Error())
		return 0, 0, fmt.Errorf("failed to calculate summary")
	}
	return count, revenue, nil
}

func parseDateYYYYMMDD(s string) (time.Time, error) {
	t, err := time.ParseInLocation("2006-01-02", strings.TrimSpace(s), time.Local)
	if err != nil {
		return time.Time{}, err
	}
	return time.Date(t.Year(), t.Month(), t.Day(), 0, 0, 0, 0, time.Local), nil
}

func parseFilter(q FilterQuery) (Filter, error) {
	var filter Filter

	if strings.TrimSpace(q.Status) != "" {
		st := BookingStatus(strings.TrimSpace(q.Status))
		switch st {
		case StatusMenunggu, StatusDikonfirmasi, StatusSelesai, StatusDibatalkan:
			filter.Status = &st
		default:
			return Filter{}, fmt.Errorf("status tidak valid")
		}
	}

	if strings.TrimSpace(q.Paket) != "" {
		p := strings.TrimSpace(q.Paket)
		filter.Paket = &p
	}

	if strings.TrimSpace(q.DateFrom) != "" {
		df, err := parseDateYYYYMMDD(q.DateFrom)
		if err != nil {
			return Filter{}, fmt.Errorf("format date_from harus YYYY-MM-DD")
		}
		filter.DateFrom = &df
	}

	if strings.TrimSpace(q.DateTo) != "" {
		dt, err := parseDateYYYYMMDD(q.DateTo)
		if err != nil {
			return Filter{}, fmt.Errorf("format date_to harus YYYY-MM-DD")
		}
		filter.DateTo = &dt
	}

	if filter.DateFrom != nil && filter.DateTo != nil && filter.DateFrom.After(*filter.DateTo) {
		return Filter{}, fmt.Errorf("date_from tidak boleh setelah date_to")
	}

	if q.Limit > 0 {
		filter.Limit = q.Limit
		if q.Page > 0 {
			filter.Offset = (q.Page - 1) * q.Limit
		}
	}

	return filter, nil
}

func isValidTransition(current, next BookingStatus) bool {
	switch current {
	case StatusMenunggu:
		return next == StatusDikonfirmasi || next == StatusDibatalkan
	case StatusDikonfirmasi:
		return next == StatusSelesai || next == StatusDibatalkan
	default:
		return false
	}
}

