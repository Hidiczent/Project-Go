package repository

import (
	"database/sql"
	"myapp/internal/accommodation/model"
)

type AccommodationRepository interface {
	GetAll() ([]model.Accommodation, error)
	GetByID(id int64) (model.Accommodation, error)
	Create(model.Accommodation) error
	Update(model.Accommodation) error
	Delete(id int64) error
}

type accommodationRepo struct {
	db *sql.DB
}

func NewAccommodationRepository(db *sql.DB) AccommodationRepository {
	return &accommodationRepo{db: db}
}

func (r *accommodationRepo) GetAll() ([]model.Accommodation, error) {
	rows, err := r.db.Query(`
		SELECT accommodation_id, name, main_image, village_id, about, popular_facilities, latitude, longitude, createdAt, updatedAt 
		FROM accommodation`)
	if err != nil {
		return nil, err
	}
	defer rows.Close()

	var list []model.Accommodation
	for rows.Next() {
		var a model.Accommodation
		err := rows.Scan(&a.ID, &a.Name, &a.MainImage, &a.VillageID, &a.About, &a.PopularFacilities, &a.Latitude, &a.Longitude, &a.CreatedAt, &a.UpdatedAt)
		if err != nil {
			return nil, err
		}
		list = append(list, a)
	}
	return list, nil
}

func (r *accommodationRepo) GetByID(id int64) (model.Accommodation, error) {
	var a model.Accommodation
	err := r.db.QueryRow(`
		SELECT accommodation_id, name, main_image, village_id, about, popular_facilities, latitude, longitude, createdAt, updatedAt 
		FROM accommodation WHERE accommodation_id=?`, id).
		Scan(&a.ID, &a.Name, &a.MainImage, &a.VillageID, &a.About, &a.PopularFacilities, &a.Latitude, &a.Longitude, &a.CreatedAt, &a.UpdatedAt)
	return a, err
}

func (r *accommodationRepo) Create(a model.Accommodation) error {
	_, err := r.db.Exec(`
		INSERT INTO accommodation (name, main_image, village_id, about, popular_facilities, latitude, longitude) 
		VALUES (?, ?, ?, ?, ?, ?, ?)`,
		a.Name, a.MainImage, a.VillageID, a.About, a.PopularFacilities, a.Latitude, a.Longitude)
	return err
}

func (r *accommodationRepo) Update(a model.Accommodation) error {
	_, err := r.db.Exec(`
		UPDATE accommodation SET 
			name=?, main_image=?, village_id=?, about=?, popular_facilities=?, latitude=?, longitude=? 
		WHERE accommodation_id=?`,
		a.Name, a.MainImage, a.VillageID, a.About, a.PopularFacilities, a.Latitude, a.Longitude, a.ID)
	return err
}

func (r *accommodationRepo) Delete(id int64) error {
	_, err := r.db.Exec("DELETE FROM accommodation WHERE accommodation_id=?", id)
	return err
}
