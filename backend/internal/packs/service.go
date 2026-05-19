package paket

import (
	"fmt"
	"strings"
	"time"

	"backend-travelku/pkg/logger"

	"github.com/google/uuid"
	"gorm.io/gorm"
)

type Service interface {
	Create(req *CreatePaketRequest) (*Paket, error)
	List(onlyActive bool) ([]Paket, error)
	GetByID(id string) (*Paket, error)
	Update(id string, req *UpdatePaketRequest) (*Paket, error)
	Delete(id string) error
}

type service struct {
	repo Repository
}

func NewService(repo Repository) Service {
	return &service{repo: repo}
}

func (s *service) Create(req *CreatePaketRequest) (*Paket, error) {
	if strings.TrimSpace(req.NamaPaket) == "" {
		return nil, fmt.Errorf("nama paket wajib diisi")
	}
	if req.Harga < 0 {
		return nil, fmt.Errorf("harga tidak boleh negatif")
	}
	if req.KuotaTotal < 1 {
		return nil, fmt.Errorf("kuota total minimal 1")
	}

	dep, err := parseDate(req.TanggalBerangkat)
	if err != nil {
		return nil, fmt.Errorf("format tanggal_berangkat harus YYYY-MM-DD")
	}

	now := time.Now().In(time.Local)
	today := time.Date(now.Year(), now.Month(), now.Day(), 0, 0, 0, 0, now.Location())
	if dep.Before(today) {
		return nil, fmt.Errorf("tanggal keberangkatan tidak boleh di masa lalu")
	}

	p := &Paket{
		ID:               uuid.NewString(),
		NamaPaket:        strings.TrimSpace(req.NamaPaket),
		Deskripsi:        strings.TrimSpace(req.Deskripsi),
		Harga:            req.Harga,
		KuotaTotal:       req.KuotaTotal,
		SisaKuota:        req.KuotaTotal,
		TanggalBerangkat: dep,
		IsActive:         true,
	}

	if err := s.repo.Create(p); err != nil {
		logger.Error("Failed to create paket: " + err.Error())
		return nil, fmt.Errorf("gagal membuat paket wisata")
	}
	return p, nil
}

func (s *service) List(onlyActive bool) ([]Paket, error) {
	return s.repo.List(onlyActive)
}

func (s *service) GetByID(id string) (*Paket, error) {
	p, err := s.repo.GetByID(id)
	if err != nil {
		if err == gorm.ErrRecordNotFound {
			return nil, fmt.Errorf("paket wisata tidak ditemukan")
		}
		return nil, fmt.Errorf("gagal mengambil data paket")
	}
	return p, nil
}

func (s *service) Update(id string, req *UpdatePaketRequest) (*Paket, error) {
	if strings.TrimSpace(req.NamaPaket) == "" {
		return nil, fmt.Errorf("nama paket wajib diisi")
	}
	if req.Harga < 0 {
		return nil, fmt.Errorf("harga tidak boleh negatif")
	}
	if req.KuotaTotal < 1 {
		return nil, fmt.Errorf("kuota total minimal 1")
	}

	p, err := s.repo.GetByID(id)
	if err != nil {
		if err == gorm.ErrRecordNotFound {
			return nil, fmt.Errorf("paket wisata tidak ditemukan")
		}
		return nil, fmt.Errorf("gagal mengambil data paket")
	}

	dep, err := parseDate(req.TanggalBerangkat)
	if err != nil {
		return nil, fmt.Errorf("format tanggal_berangkat harus YYYY-MM-DD")
	}

	now := time.Now().In(time.Local)
	today := time.Date(now.Year(), now.Month(), now.Day(), 0, 0, 0, 0, now.Location())
	if !dep.Equal(p.TanggalBerangkat) && dep.Before(today) {
		return nil, fmt.Errorf("tanggal keberangkatan tidak boleh diubah ke masa lalu")
	}

	// Hitung selisih kuota jika kuota total diubah
	diff := req.KuotaTotal - p.KuotaTotal
	newSisa := p.SisaKuota + diff
	if newSisa < 0 {
		return nil, fmt.Errorf("kuota total baru terlalu kecil, peserta terdaftar melebihi kuota")
	}

	p.NamaPaket = strings.TrimSpace(req.NamaPaket)
	p.Deskripsi = strings.TrimSpace(req.Deskripsi)
	p.Harga = req.Harga
	p.KuotaTotal = req.KuotaTotal
	p.SisaKuota = newSisa
	p.TanggalBerangkat = dep
	p.IsActive = req.IsActive

	if err := s.repo.Update(p); err != nil {
		logger.Error("Failed to update paket: " + err.Error())
		return nil, fmt.Errorf("gagal memperbarui paket wisata")
	}
	return p, nil
}

func (s *service) Delete(id string) error {
	p, err := s.repo.GetByID(id)
	if err != nil {
		if err == gorm.ErrRecordNotFound {
			return fmt.Errorf("paket wisata tidak ditemukan")
		}
		return fmt.Errorf("gagal mengambil data paket")
	}

	if p.KuotaTotal != p.SisaKuota {
		return fmt.Errorf("paket tidak bisa dihapus karena sudah ada peserta terdaftar")
	}

	if err := s.repo.Delete(id); err != nil {
		logger.Error("Failed to delete paket: " + err.Error())
		return fmt.Errorf("gagal menghapus paket wisata")
	}
	return nil
}

func parseDate(s string) (time.Time, error) {
	t, err := time.ParseInLocation("2006-01-02", strings.TrimSpace(s), time.Local)
	if err != nil {
		return time.Time{}, err
	}
	return time.Date(t.Year(), t.Month(), t.Day(), 0, 0, 0, 0, time.Local), nil
}
