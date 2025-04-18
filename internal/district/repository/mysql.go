package repository

import (
	"database/sql"
	"myapp/internal/district/model"
)

type DistrictRepository interface {
	GetAll() ([]model.District, error)
	GetByID(id int64) (model.District, error)
	Create(model.District) error
	Update(model.District) error
	Delete(id int64) error
}

type districtRepo struct {
	db *sql.DB
}

func NewDistrictRepository(db *sql.DB) DistrictRepository {
	return &districtRepo{db: db}
}

func (r *districtRepo) GetAll() ([]model.District, error) {
	rows, err := r.db.Query("SELECT district_id, name, province_id FROM district")
	if err != nil {
		return nil, err
	}
	defer rows.Close()

	var list []model.District
	for rows.Next() {
		var d model.District
		err := rows.Scan(&d.ID, &d.Name, &d.ProvinceID)
		if err != nil {
			return nil, err
		}
		list = append(list, d)
	}
	return list, nil
}

func (r *districtRepo) GetByID(id int64) (model.District, error) {
	var d model.District
	err := r.db.QueryRow("SELECT district_id, name, province_id FROM district WHERE district_id=?", id).
		Scan(&d.ID, &d.Name, &d.ProvinceID)
	return d, err
}

func (r *districtRepo) Create(d model.District) error {
	_, err := r.db.Exec("INSERT INTO district (name, province_id) VALUES (?, ?)", d.Name, d.ProvinceID)
	return err
}

func (r *districtRepo) Update(d model.District) error {
	_, err := r.db.Exec("UPDATE district SET name=?, province_id=? WHERE district_id=?", d.Name, d.ProvinceID, d.ID)
	return err
}

func (r *districtRepo) Delete(id int64) error {
	_, err := r.db.Exec("DELETE FROM district WHERE district_id=?", id)
	return err
}
