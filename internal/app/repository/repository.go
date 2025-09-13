package repository

import "fmt"

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