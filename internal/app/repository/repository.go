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
	ID int
	Title string
	Description string
	ImagePath string
}

func (r *Repository) GetStars() ([]Star, error) {
	stars := []Star{
		{
			ID: 1,
			Title: "Gliese 667",
			Description: "Красный карлик спектрального класса M1.5V",
			ImagePath: "http://localhost:9000/exocalc/gliese-667.png",
		},
		{
			ID: 2,
			Title: "TRAPPIST-1",
			Description: "Ультрахолодный красный карлик спектрального класса M8V",
			ImagePath: "http://localhost:9000/exocalc/trappist-1.png",
		},
		{
			ID: 3,
			Title: "Ross 128",
			Description: "Красный карлик спектрального класса M4V",
			ImagePath: "http://localhost:9000/exocalc/ross-128.png",
		},
		{
			ID: 4,
			Title: "Proxima Centauri",
			Description: "Красный карлик спектрального класса M5.5V",
			ImagePath: "http://localhost:9000/exocalc/proxima.png",
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
