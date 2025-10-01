package ds

type Star struct {
	ID                      int     `gorm:"primaryKey;autoIncrement"`
	Title                   string  `gorm:"type:varchar(255);not null;index"`
	Description             string  `gorm:"type:text"`
	ImagePath               string  `gorm:"type:varchar(512)"`
	SpectralType            string  `gorm:"type:varchar(64);index"`
	Temperature             string  `gorm:"type:varchar(64)"`
	Radius                  string  `gorm:"type:varchar(64)"`
	Mass                    string  `gorm:"type:varchar(64)"`
	Luminosity              string  `gorm:"type:varchar(64)"`
	Metallicity             string  `gorm:"type:varchar(64)"`
	Age                     string  `gorm:"type:varchar(64)"`
	Distance                string  `gorm:"type:varchar(64)"`
	ProbableNumberOfPlanets float32 `gorm:"type:real"`
	HabitableZone           string  `gorm:"type:varchar(64)"`
}

func (Star) TableName() string { return "stars" }
