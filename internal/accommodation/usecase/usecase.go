package usecase

import (
	"myapp/internal/accommodation/model"
	"myapp/internal/accommodation/repository"
)

type AccommodationUsecase interface {
	GetAll() ([]model.Accommodation, error)
	GetByID(int64) (model.Accommodation, error)
	Create(model.Accommodation) error
	Update(model.Accommodation) error
	Delete(int64) error
}

type accommodationUsecase struct {
	repo repository.AccommodationRepository
}

func NewAccommodationUsecase(r repository.AccommodationRepository) AccommodationUsecase {
	return &accommodationUsecase{repo: r}
}

func (u *accommodationUsecase) GetAll() ([]model.Accommodation, error) {
	return u.repo.GetAll()
}

func (u *accommodationUsecase) GetByID(id int64) (model.Accommodation, error) {
	return u.repo.GetByID(id)
}

func (u *accommodationUsecase) Create(a model.Accommodation) error {
	return u.repo.Create(a)
}

func (u *accommodationUsecase) Update(a model.Accommodation) error {
	return u.repo.Update(a)
}

func (u *accommodationUsecase) Delete(id int64) error {
	return u.repo.Delete(id)
}
