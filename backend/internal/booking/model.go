package booking

import "time"

type BookingStatus string

const (
	StatusMenunggu     BookingStatus = "MENUNGGU"
	StatusDikonfirmasi BookingStatus = "DIKONFIRMASI"
	StatusSelesai      BookingStatus = "SELESAI"
	StatusDibatalkan   BookingStatus = "DIBATALKAN"
)

func (s BookingStatus) IsTerminal() bool {
	return s == StatusSelesai || s == StatusDibatalkan
}

type Booking struct {
	ID               string        `gorm:"primaryKey" json:"id"`
	NamaPemesan      string        `json:"nama_pemesan"`
	Kontak           string        `json:"kontak"`
	PaketID          string        `gorm:"type:varchar(100);index" json:"paket_id"`
	PaketWisata      string        `json:"paket_wisata"`
	TanggalBerangkat time.Time     `gorm:"type:date;not null" json:"tanggal_berangkat"`
	JumlahPeserta    int           `json:"jumlah_peserta"`
	HargaPerOrang    int64         `json:"harga_per_orang"`
	Status           BookingStatus `gorm:"type:varchar(20);not null" json:"status"`
	Catatan          *string       `json:"catatan,omitempty"`
	CreatedAt        time.Time     `json:"created_at"`
	UpdatedAt        time.Time     `json:"updated_at"`
}

func (Booking) TableName() string {
	return "bookings"
}

