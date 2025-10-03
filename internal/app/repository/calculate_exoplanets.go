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