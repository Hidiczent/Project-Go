package repository

import (
	"database/sql"
	"myapp/internal/notification/model"
)

type NotificationRepository interface {
	GetAll() ([]model.Notification, error)
	GetByID(id int) (*model.Notification, error)
	Create(n model.Notification) error
	Update(n model.Notification) error
	Delete(id int) error
}

type notificationRepo struct {
	db *sql.DB
}

func NewNotificationRepository(db *sql.DB) NotificationRepository {
	return &notificationRepo{db}
}

func (r *notificationRepo) GetAll() ([]model.Notification, error) {
	rows, err := r.db.Query(`SELECT notification_id, status_notification, order_id FROM notification`)
	if err != nil {
		return nil, err
	}
	defer rows.Close()

	var notifications []model.Notification
	for rows.Next() {
		var n model.Notification
		if err := rows.Scan(&n.NotificationID, &n.StatusNotification, &n.OrderID); err != nil {
			return nil, err
		}
		notifications = append(notifications, n)
	}
	return notifications, nil
}

func (r *notificationRepo) GetByID(id int) (*model.Notification, error) {
	var n model.Notification
	err := r.db.QueryRow(`SELECT notification_id, status_notification, order_id FROM notification WHERE notification_id = ?`, id).
		Scan(&n.NotificationID, &n.StatusNotification, &n.OrderID)
	if err != nil {
		return nil, err
	}
	return &n, nil
}

func (r *notificationRepo) Create(n model.Notification) error {
	_, err := r.db.Exec(`INSERT INTO notification (status_notification, order_id) VALUES (?, ?)`, n.StatusNotification, n.OrderID)
	return err
}

func (r *notificationRepo) Update(n model.Notification) error {
	_, err := r.db.Exec(`UPDATE notification SET status_notification = ?, order_id = ? WHERE notification_id = ?`,
		n.StatusNotification, n.OrderID, n.NotificationID)
	return err
}

func (r *notificationRepo) Delete(id int) error {
	_, err := r.db.Exec(`DELETE FROM notification WHERE notification_id = ?`, id)
	return err
}
