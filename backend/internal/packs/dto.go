package paket

type CreatePaketRequest struct {
	NamaPaket        string `json:"nama_paket" binding:"required,min=3,max=150"`
	Deskripsi        string `json:"deskripsi"`
	Harga            int64  `json:"harga" binding:"required,min=0"`
	KuotaTotal       int    `json:"kuota_total" binding:"required,min=1"`
	TanggalBerangkat string `json:"tanggal_berangkat" binding:"required"`
}

type UpdatePaketRequest struct {
	NamaPaket        string `json:"nama_paket" binding:"required,min=3,max=150"`
	Deskripsi        string `json:"deskripsi"`
	Harga            int64  `json:"harga" binding:"required,min=0"`
	KuotaTotal       int    `json:"kuota_total" binding:"required,min=1"`
	TanggalBerangkat string `json:"tanggal_berangkat" binding:"required"`
	IsActive         bool   `json:"is_active"`
}
