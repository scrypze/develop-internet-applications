package repository

import (
	"develop-internet-applications/internal/app/ds"
	"fmt"
)

func (r *Repository) GetSelectedStars() ([]ds.SelectedStars, error) {
	selected := []ds.SelectedStars{
		{
			ID: 1,
			SelectedStarsItems: []ds.Star{
				{
					ID:                      1,
					Title:                   "Gliese 667",
					Description:             "Красный карлик спектрального класса M1.5V",
					ImagePath:               "http://localhost:9000/exocalc/gliese-667.png",
					SpectralType:            "M1.5V",
					Temperature:             "~3700 K",
					Radius:                  "~0.42 R☉",
					Mass:                    "~0.33 M☉",
					Luminosity:              "~0.014 L☉",
					Metallicity:             "-0.55",
					Age:                     "~2-10 млрд лет",
					Distance:                "~6.8 pc (~22 световых лет)",
					Scientist:               "Михаил Ломоносов",
					Date:                    "12.10.2003",
					ProbableNumberOfPlanets: 2.5,
					HabitableZone:           "0.12-0.24 a.e.",
				},
				{
					ID:                      2,
					Title:                   "TRAPPIST-1",
					Description:             "Ультрахолодный красный карлик спектрального класса M8V",
					ImagePath:               "http://localhost:9000/exocalc/trappist-1.png",
					SpectralType:            "M8V",
					Temperature:             "~2550 K",
					Radius:                  "~0.12 R☉",
					Mass:                    "~0.09 M☉",
					Luminosity:              "~0.0005 L☉",
					Metallicity:             "0.04",
					Age:                     "~7.6 млрд лет",
					Distance:                "~12.1 pc (~39 световых лет)",
					Scientist:               "Григорий Шайн",
					Date:                    "22.10.2010",
					ProbableNumberOfPlanets: 2.5,
					HabitableZone:           "0.12-0.24 a.e.",
				},
				{
					ID:                      3,
					Title:                   "Ross 128",
					Description:             "Красный карлик спектрального класса M4V",
					ImagePath:               "http://localhost:9000/exocalc/ross-128.png",
					SpectralType:            "M4V",
					Temperature:             "~3192 K",
					Radius:                  "~0.21 R☉",
					Mass:                    "~0.16 M☉",
					Luminosity:              "~0.0036 L☉",
					Metallicity:             "0.00",
					Age:                     "~9.45 млрд лет",
					Distance:                "~3.37 pc (~11 световых лет)",
					Scientist:               "Михаил Ломоносов",
					Date:                    "12.10.2003",
					ProbableNumberOfPlanets: 2.5,
					HabitableZone:           "0.12-0.24 a.e.",
				},
			},
		},
	}

	if len(selected) == 0 {
		return nil, fmt.Errorf("список выбранных звёзд пуст")
	}

	return selected, nil
}

func (r *Repository) GetSelectedStarsByID(id int) (ds.SelectedStars, error) {
	lists, err := r.GetSelectedStars()
	if err != nil {
		return ds.SelectedStars{}, err
	}

	for _, list := range lists {
		if list.ID == id {
			return list, nil
		}
	}

	return ds.SelectedStars{}, fmt.Errorf("выбранные звёзды не найдены")
}

func (r *Repository) GetStars() ([]ds.Star, error) {
	var stars []ds.Star

	err := r.db.Find(&stars).Error

	if err != nil {
		return nil, err
	}

	return stars, nil
}

func (r *Repository) GetStar(id int) (ds.Star, error) {
	var star ds.Star

	err := r.db.First(&star, id).Error

	if err != nil {
		return ds.Star{}, err
	}

	return star, nil
}

func (r *Repository) GetStarsByTitle(title string) ([]ds.Star, error) {
	var stars []ds.Star

	err := r.db.Where("title ILIKE ?", "%"+title+"%").Find(&stars).Error

	if err != nil {
		return nil, err
	}
	
	return stars, nil
}
