package booking

import (
	"strings"
	"time"

	"gorm.io/gorm"
)

type Filter struct {
	Status   *BookingStatus
	Paket    *string
	DateFrom *time.Time
	DateTo   *time.Time
	Limit    int
	Offset   int
}

// Repository adalah data access layer untuk booking.
type Repository interface {
	Create(b *Booking) error
	GetByID(id string) (*Booking, error)
	List(filter Filter) ([]Booking, error)
	Update(b *Booking) error
	Delete(id string) error
	Count(filter Filter) (int64, error)
	SumRevenue(filter Filter) (int64, error)
}

type repository struct {
	db *gorm.DB
}

func NewRepository(db *gorm.DB) Repository {
	return &repository{db: db}
}

func (r *repository) Create(b *Booking) error {
	return r.db.Create(b).Error
}

func (r *repository) GetByID(id string) (*Booking, error) {
	var booking Booking
	if err := r.db.Where("id = ?", id).First(&booking).Error; err != nil {
		return nil, err
	}
	return &booking, nil
}

func (r *repository) List(filter Filter) ([]Booking, error) {
	var bookings []Booking
	q := applyFilter(r.db.Model(&Booking{}), filter)
	if filter.Limit > 0 {
		q = q.Limit(filter.Limit).Offset(filter.Offset)
	}
	if err := q.Order("created_at desc").Find(&bookings).Error; err != nil {
		return nil, err
	}
	return bookings, nil
}

func (r *repository) Update(b *Booking) error {
	return r.db.Save(b).Error
}

func (r *repository) Delete(id string) error {
	return r.db.Where("id = ?", id).Delete(&Booking{}).Error
}

func (r *repository) Count(filter Filter) (int64, error) {
	var count int64
	q := applyFilter(r.db.Model(&Booking{}), filter)
	if err := q.Count(&count).Error; err != nil {
		return 0, err
	}
	return count, nil
}

func (r *repository) SumRevenue(filter Filter) (int64, error) {
	// Revenue hanya dari DIKONFIRMASI dan SELESAI
	q := applyFilter(r.db.Model(&Booking{}), filter).Where("status IN ?", []BookingStatus{StatusDikonfirmasi, StatusSelesai})

	var total int64
	// PostgreSQL: SUM(jumlah_peserta * harga_per_orang)
	row := q.Select("COALESCE(SUM(jumlah_peserta * harga_per_orang), 0)").Row()
	if err := row.Scan(&total); err != nil {
		return 0, err
	}
	return total, nil
}

func applyFilter(q *gorm.DB, filter Filter) *gorm.DB {
	if filter.Status != nil {
		q = q.Where("status = ?", *filter.Status)
	}
	if filter.Paket != nil && *filter.Paket != "" {
		q = q.Where("LOWER(paket_wisata) LIKE ?", "%"+strings.ToLower(*filter.Paket)+"%")
	}
	if filter.DateFrom != nil {
		q = q.Where("tanggal_berangkat >= ?", *filter.DateFrom)
	}
	if filter.DateTo != nil {
		q = q.Where("tanggal_berangkat <= ?", *filter.DateTo)
	}
	return q
}
