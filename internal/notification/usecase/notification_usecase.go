package usecase

import (
	"myapp/internal/notification/model"
	"myapp/internal/notification/repository"
)

type NotificationUseCase interface {
	GetAll() ([]model.Notification, error)
	GetByID(id int) (*model.Notification, error)
	Create(n model.Notification) error
	Update(n model.Notification) error
	Delete(id int) error
}

type notificationUsecase struct {
	repo repository.NotificationRepository
}

func NewNotificationUseCase(repo repository.NotificationRepository) NotificationUseCase {
	return &notificationUsecase{repo}
}

func (u *notificationUsecase) GetAll() ([]model.Notification, error) {
	return u.repo.GetAll()
}

func (u *notificationUsecase) GetByID(id int) (*model.Notification, error) {
	return u.repo.GetByID(id)
}

func (u *notificationUsecase) Create(n model.Notification) error {
	return u.repo.Create(n)
}

func (u *notificationUsecase) Update(n model.Notification) error {
	return u.repo.Update(n)
}

func (u *notificationUsecase) Delete(id int) error {
	return u.repo.Delete(id)
}
