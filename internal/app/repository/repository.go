package repository

import (
	"fmt"
	"strings"
)

type Repository struct {
}

func NewRepository() (*Repository, error) {
	return &Repository{}, nil
}

type Star struct {
	ID           int
	Title        string
	Description  string
	ImagePath    string
	SpectralType string
	Temperature  string
	Radius       string
	Mass         string
	Luminosity   string
	Metallicity  string
	Age          string
	Distance     string
}

func (r *Repository) GetStars() ([]Star, error) {
	stars := []Star{
		{
			ID:           1,
			Title:        "Gliese 667",
			Description:  "Красный карлик спектрального класса M1.5V",
			ImagePath:    "http://localhost:9000/exocalc/gliese-667.png",
			SpectralType: "M1.5V",
			Temperature:  "~3700 K",
			Radius:       "~0.42 R☉",
			Mass:         "~0.33 M☉",
			Luminosity:   "~0.014 L☉",
			Metallicity:  "-0.55",
			Age:          "~2-10 млрд лет",
			Distance:     "~6.8 pc (~22 световых лет)",
		},
		{
			ID:           2,
			Title:        "TRAPPIST-1",
			Description:  "Ультрахолодный красный карлик спектрального класса M8V",
			ImagePath:    "http://localhost:9000/exocalc/trappist-1.png",
			SpectralType: "M8V",
			Temperature:  "~2550 K",
			Radius:       "~0.12 R☉",
			Mass:         "~0.09 M☉",
			Luminosity:   "~0.0005 L☉",
			Metallicity:  "0.04",
			Age:          "~7.6 млрд лет",
			Distance:     "~12.1 pc (~39 световых лет)",
		},
		{
			ID:           3,
			Title:        "Ross 128",
			Description:  "Красный карлик спектрального класса M4V",
			ImagePath:    "http://localhost:9000/exocalc/ross-128.png",
			SpectralType: "M4V",
			Temperature:  "~3192 K",
			Radius:       "~0.21 R☉",
			Mass:         "~0.16 M☉",
			Luminosity:   "~0.0036 L☉",
			Metallicity:  "0.00",
			Age:          "~9.45 млрд лет",
			Distance:     "~3.37 pc (~11 световых лет)",
		},
		{
			ID:           4,
			Title:        "Proxima Centauri",
			Description:  "Красный карлик спектрального класса M5.5V",
			ImagePath:    "http://localhost:9000/exocalc/proxima.png",
			SpectralType: "M5.5V",
			Temperature:  "~3042 K",
			Radius:       "~0.15 R☉",
			Mass:         "~0.12 M☉",
			Luminosity:   "~0.0017 L☉",
			Metallicity:  "0.21",
			Age:          "~4.85 млрд лет",
			Distance:     "~1.30 pc (~4.24 световых лет)",
		},
	}

	if len(stars) == 0 {
		return nil, fmt.Errorf("массив пустой")
	}

	return stars, nil
}

func (r *Repository) GetStar(id int) (Star, error) {
	stars, err := r.GetStars()

	if err != nil {
		return Star{}, err
	}

	for _, star := range stars {
		if star.ID == id {
			return star, nil
		}
	}

	return Star{}, fmt.Errorf("звезда не найдена")
}

func (r *Repository) GetStarsByTitle(title string) ([]Star, error) {
	stars, err := r.GetStars()
	if err != nil {
		return []Star{}, err
	}

	var result []Star
	for _, star := range stars {
		if strings.Contains(strings.ToLower(star.Title), strings.ToLower(title)) {
			result = append(result, star)
		}
	}

	return result, nil
}

func (r *Repository) GetCart() ([]Star, int, error) {
	cart := []Star{
		{
			ID:           1,
			Title:        "Gliese 667",
			Description:  "Красный карлик спектрального класса M1.5V",
			ImagePath:    "http://localhost:9000/exocalc/gliese-667.png",
			SpectralType: "M1.5V",
			Temperature:  "~3700 K",
			Radius:       "~0.42 R☉",
			Mass:         "~0.33 M☉",
			Luminosity:   "~0.014 L☉",
			Metallicity:  "-0.55",
			Age:          "~2-10 млрд лет",
			Distance:     "~6.8 pc (~22 световых лет)",
		},
		{
			ID:           2,
			Title:        "TRAPPIST-1",
			Description:  "Ультрахолодный красный карлик спектрального класса M8V",
			ImagePath:    "http://localhost:9000/exocalc/trappist-1.png",
			SpectralType: "M8V",
			Temperature:  "~2550 K",
			Radius:       "~0.12 R☉",
			Mass:         "~0.09 M☉",
			Luminosity:   "~0.0005 L☉",
			Metallicity:  "0.04",
			Age:          "~7.6 млрд лет",
			Distance:     "~12.1 pc (~39 световых лет)",
		},
	}

	if len(cart) == 0 {
		return nil, len(cart), fmt.Errorf("массив пустой")
	}

	return cart, len(cart), nil
}
