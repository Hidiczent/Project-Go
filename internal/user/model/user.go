package model

import "time"

type User struct {
	ID          int64     `json:"id"`
	FirstName   string    `json:"first_name"`
	LastName    string    `json:"lastname,omitempty"`
	Password    string    `json:"password"`
	PhoneNumber *int64    `json:"phone_number,omitempty"`
	Email       string    `json:"email"`
	Photo       *string   `json:"photo,omitempty"`
	CreatedAt   *time.Time `json:"created_at,omitempty"`
	UpdatedAt   *time.Time `json:"updated_at,omitempty"`
	Role        string    `json:"role,omitempty"`
}
