package model

type Star struct {
	ID           int    `gorm:"primaryKey;autoIncrement" json:"id"`
	Title        string `gorm:"type:varchar(255);not null;index" json:"title"`
	Description  string `gorm:"type:text" json:"description"`
	ImagePath    string `gorm:"type:varchar(512)" json:"image_path"`
	SpectralType string `gorm:"type:varchar(64);index" json:"spectral_type"`
	Temperature  string `gorm:"type:varchar(64)" json:"temperature"`
	Radius       string `gorm:"type:varchar(64)" json:"radius"`
	Mass         string `gorm:"type:varchar(64)" json:"mass"`
	Luminosity   string `gorm:"type:varchar(64)" json:"luminosity"`
	Metallicity  string `gorm:"type:varchar(64)" json:"metallicity"`
	Age          string `gorm:"type:varchar(64)" json:"age"`
	Distance     string `gorm:"type:varchar(64)" json:"distance"`
	IsDeleted    string `gorm:"column:is_deleted;not null;default:false" json:"is_deleted"`
}
