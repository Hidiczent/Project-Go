package usecase

import (
	"myapp/internal/user/model"
	"myapp/internal/user/repository"
)

type UserUsecase interface {
	GetAll() ([]model.User, error)
	GetByID(id int64) (model.User, error)
	Create(user model.User) error
	Update(user model.User) error
	Delete(id int64) error
}

type userUsecase struct {
	repo repository.UserRepository
}

func NewUserUsecase(repo repository.UserRepository) UserUsecase {
	return &userUsecase{repo: repo}
}
func (u *userUsecase) GetAll() ([]model.User, error)        { return u.repo.GetAll() }
func (u *userUsecase) GetByID(id int64) (model.User, error) { return u.repo.GetByID(id) }
func (u *userUsecase) Create(user model.User) error         { return u.repo.Create(user) }
func (u *userUsecase) Update(user model.User) error         { return u.repo.Update(user) }
func (u *userUsecase) Delete(id int64) error                { return u.repo.Delete(id) }
