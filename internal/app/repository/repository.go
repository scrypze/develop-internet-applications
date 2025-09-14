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
}

func (r *Repository) GetStars() ([]Star, error) {
	stars := []Star{
		{
			ID: 1,
			Title: "Gliese 667",
		},
		{
			ID: 2,
			Title: "TRAPPIST-1",
		},
		{
			ID: 3,
			Title: "Ross 128",
		},
		{
			ID: 4,
			Title: "Proxima Centauri",
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
