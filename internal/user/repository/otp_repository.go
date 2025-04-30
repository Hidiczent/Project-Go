package repository

import (
	"database/sql"
	"log"
	"time"
)

type OTPRepository interface {
	SaveOTP(email, otp, action string, expiresAt time.Time) error
	SaveOTPWithMetadata(email, otp, action string, expiresAt time.Time, metadata map[string]string) error
	VerifyOTP(email, code, action string) (bool, error)
	MarkVerified(email, code, action string) error
	GetOTPMetadata(email, action string) (map[string]string, error)
}
type otpRepo struct {
	db *sql.DB
}

func NewOTPRepository(db *sql.DB) OTPRepository {
	return &otpRepo{db: db}
}

func (r *otpRepo) SaveOTP(email, otp, action string, expiresAt time.Time) error {
	// ✅ แปลงเป็น UTC เพื่อให้ตรงกับเวลาของ MySQL
	expiresAt = expiresAt.UTC()

	log.Printf("📥 Save OTP: email=%s, code=%s, action=%s, expires_at=%s", email, otp, action, expiresAt.Format(time.RFC3339))
	_, err := r.db.Exec("INSERT INTO otps (email, code, action, expires_at) VALUES (?, ?, ?, ?)", email, otp, action, expiresAt)
	return err
}
func (r *otpRepo) VerifyOTP(email, code, action string) (bool, error) {
	log.Printf("🔍 Verifying OTP: email=%s, code=%s", email, code)
	log.Printf("🧪 Checking OTP - email: %s | code: %s | action: %s", email, code, action)

	var count int
	err := r.db.QueryRow(`
		SELECT COUNT(*) FROM otps
		WHERE email = ? AND code = ? AND action = ? AND verified = 0 AND expires_at > UTC_TIMESTAMP()`,
		email, code, action,
	).Scan(&count)

	if err != nil {
		return false, err
	}

	log.Printf("✅ Found OTP match count: %d", count)
	return count > 0, nil
}

func (r *otpRepo) MarkVerified(email, code, action string) error {
	_, err := r.db.Exec(`
		UPDATE otps SET verified = 1, verified_at = NOW()
		WHERE email = ? AND code = ? AND action = ?`, email, code, action)
	return err
}
func (r *otpRepo) GetOTPMetadata(email, action string) (map[string]string, error) {
	rows, err := r.db.Query(`
		SELECT field_key, field_value
		FROM otp_metadata
		WHERE email = ? AND action = ?`,
		email, action,
	)
	if err != nil {
		return nil, err
	}
	defer rows.Close()

	result := make(map[string]string)
	for rows.Next() {
		var key, value string
		if err := rows.Scan(&key, &value); err != nil {
			return nil, err
		}
		result[key] = value
	}
	return result, nil
}

func (r *otpRepo) SaveOTPWithMetadata(email, otp, action string, expiresAt time.Time, metadata map[string]string) error {
	expiresAt = expiresAt.UTC()

	// 1. Save OTP ปกติ
	_, err := r.db.Exec(`
		INSERT INTO otps (email, code, action, expires_at)
		VALUES (?, ?, ?, ?)`,
		email, otp, action, expiresAt,
	)
	if err != nil {
		return err
	}

	// 2. Save metadata
	for key, value := range metadata {
		_, err := r.db.Exec(`
			INSERT INTO otp_metadata (email, action, field_key, field_value)
			VALUES (?, ?, ?, ?)`,
			email, action, key, value,
		)
		if err != nil {
			return err
		}
	}
	log.Printf("📝 Saving metadata for: %s | %s | %v", email, action, metadata)

	return nil
}
