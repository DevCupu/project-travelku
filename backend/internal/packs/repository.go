package paket

import (
	"gorm.io/gorm"
)

type Repository interface {
	Create(p *Paket) error
	GetByID(id string) (*Paket, error)
	List(onlyActive bool) ([]Paket, error)
	Update(p *Paket) error
	Delete(id string) error
}

type repository struct {
	db *gorm.DB
}

func NewRepository(db *gorm.DB) Repository {
	return &repository{db: db}
}

func (r *repository) Create(p *Paket) error {
	return r.db.Create(p).Error
}

func (r *repository) GetByID(id string) (*Paket, error) {
	var p Paket
	if err := r.db.Where("id = ?", id).First(&p).Error; err != nil {
		return nil, err
	}
	return &p, nil
}

func (r *repository) List(onlyActive bool) ([]Paket, error) {
	var pakets []Paket
	q := r.db.Model(&Paket{})
	if onlyActive {
		q = q.Where("is_active = ?", true)
	}
	if err := q.Order("tanggal_berangkat asc").Find(&pakets).Error; err != nil {
		return nil, err
	}
	return pakets, nil
}

func (r *repository) Update(p *Paket) error {
	return r.db.Save(p).Error
}

func (r *repository) Delete(id string) error {
	return r.db.Where("id = ?", id).Delete(&Paket{}).Error
}
