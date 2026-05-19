package paket

import "time"

type Paket struct {
	ID               string    `gorm:"primaryKey" json:"id"`
	NamaPaket        string    `gorm:"type:varchar(150);not null" json:"nama_paket"`
	Deskripsi        string    `gorm:"type:text" json:"deskripsi"`
	Harga            int64     `gorm:"not null" json:"harga"`
	KuotaTotal       int       `gorm:"not null" json:"kuota_total"`
	SisaKuota        int       `gorm:"not null" json:"sisa_kuota"`
	TanggalBerangkat time.Time `gorm:"type:date;not null" json:"tanggal_berangkat"`
	IsActive         bool      `gorm:"default:true" json:"is_active"`
	CreatedAt        time.Time `json:"created_at"`
	UpdatedAt        time.Time `json:"updated_at"`
}

func (Paket) TableName() string {
	return "pakets"
}
