package repository

import "develop-internet-applications/internal/app/ds"

func (r *Repository) GetCalculateExoplanets() (ds.CalculateExoplanets, error) {
	var list ds.CalculateExoplanets

	err := r.db.Find(&list).Error
	if err != nil {
		return ds.CalculateExoplanets{}, err
	}

	return list, nil
}

func (r *Repository) GetCalculateExoplanetsBySelectedStarsID(selectedStarsID int) (map[int]ds.CalculateExoplanets, error) {
	var rows []ds.CalculateExoplanets
	if err := r.db.Where("selected_stars_id = ?", selectedStarsID).Find(&rows).Error; err != nil {
		return nil, err
	}
	result := make(map[int]ds.CalculateExoplanets, len(rows))
	for _, row := range rows {
		result[row.StarID] = row
	}
	return result, nil
}
