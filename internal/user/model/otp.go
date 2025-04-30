package model

import "time"

type OTP struct {
	Email     string
	Code      string
	ExpiresAt time.Time
	Verified  bool
}
