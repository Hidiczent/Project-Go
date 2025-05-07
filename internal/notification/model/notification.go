package model

type Notification struct {
	NotificationID     int    `json:"notification_id"`
	StatusNotification string `json:"status_notification"`
	OrderID            *int   `json:"order_id,omitempty"`
}
