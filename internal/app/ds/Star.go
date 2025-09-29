package ds

type Star struct {
	ID                      int
	Title                   string
	Description             string
	ImagePath               string
	SpectralType            string
	Temperature             string
	Radius                  string
	Mass                    string
	Luminosity              string
	Metallicity             string
	Age                     string
	Distance                string
	Scientist               string
	Date                    string
	Probability             int
	ProbableNumberOfPlanets float32
	HabitableZone           string
}

type SelectedStars struct {
	ID                 int
	SelectedStarsItems []Star
}
