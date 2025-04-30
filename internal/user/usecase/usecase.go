package usecase

import (
	"errors"
	"myapp/internal/user/model"
	"myapp/internal/user/repository"
)

type UserUsecase interface {
	GetAll() ([]model.User, error)
	GetByID(id int64) (model.User, error)
	Create(user model.User) error
	Update(user model.User) error
	Delete(id int64) error
	GetByEmail(email string) (model.User, error)
	UpdateEmail(id int64, email string) error
	UpdatePassword(id int64, hashedPassword string) error
	IsEmailTaken(email string, excludeID int64) (bool, error)
	UpdateProfilePhoto(id int64, photoPath string) error // ✅ เพิ่ม

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
func (u *userUsecase) GetByEmail(email string) (model.User, error) {
	return u.repo.GetByEmail(email)
}

func (u *userUsecase) UpdateEmail(id int64, email string) error {
	// ตรวจสอบว่า email ซ้ำหรือไม่
	taken, err := u.repo.IsEmailTaken(email, id)
	if err != nil {
		return err
	}
	if taken {
		return errors.New("email is already in use")
	}

	return u.repo.UpdateEmail(id, email)
}

func (u *userUsecase) UpdatePassword(id int64, hashedPassword string) error {
	return u.repo.UpdatePassword(id, hashedPassword)
}

func (u *userUsecase) IsEmailTaken(email string, excludeID int64) (bool, error) {
	return u.repo.IsEmailTaken(email, excludeID)
}

func (u *userUsecase) UpdateProfilePhoto(id int64, photoPath string) error {
	return u.repo.UpdateProfilePhoto(id, photoPath) // ไปเรียกที่ repository ต่อ
}
