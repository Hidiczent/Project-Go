package usecase

import (
	"myapp/internal/district/model"
	"myapp/internal/district/repository"
)

type DistrictUsecase interface {
	GetAll() ([]model.District, error)
	GetByID(int64) (model.District, error)
	Create(model.District) error
	Update(model.District) error
	Delete(int64) error
}

type districtUsecase struct {
	repo repository.DistrictRepository
}

func NewDistrictUsecase(r repository.DistrictRepository) DistrictUsecase {
	return &districtUsecase{repo: r}
}

func (u *districtUsecase) GetAll() ([]model.District, error) {
	return u.repo.GetAll()
}

func (u *districtUsecase) GetByID(id int64) (model.District, error) {
	return u.repo.GetByID(id)
}

func (u *districtUsecase) Create(d model.District) error {
	return u.repo.Create(d)
}

func (u *districtUsecase) Update(d model.District) error {
	return u.repo.Update(d)
}

func (u *districtUsecase) Delete(id int64) error {
	return u.repo.Delete(id)
}
