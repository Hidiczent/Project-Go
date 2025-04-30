package usecase

import (
	"errors"
	"fmt"
	"log"
	"myapp/internal/user/repository"
	"time"
)

type OTPUsecase interface {
	SendOTP(email, action string) error
	SendOTPWithMetadata(email, action string, metadata map[string]string) error
	VerifyOTP(email, otp, action string) error
	VerifyAndGetMetadata(email, otp, action string) (map[string]string, error)
}

type otpUsecase struct {
	repo        repository.OTPRepository
	emailSender EmailSender // ✅ Interface สำหรับส่งอีเมล (mock/test ได้)
}

func NewOTPUsecase(repo repository.OTPRepository, sender EmailSender) OTPUsecase {
	return &otpUsecase{repo: repo, emailSender: sender}
}
func (u *otpUsecase) SendOTP(email, action string) error {
	otp := GenerateRandomOTP()                         // ✅ สร้าง OTP 6 หลัก
	expiresAt := time.Now().Add(5 * time.Minute).UTC() // ✅ หมดอายุใน 5 นาที

	err := u.repo.SaveOTP(email, otp, action, expiresAt)
	if err != nil {
		return err
	}
	log.Printf("🕒 Local time now: %s", time.Now().Format(time.RFC3339))
	log.Printf("🕒 UTC time now  : %s", time.Now().UTC().Format(time.RFC3339))

	// ✅ ส่งอีเมลจริง พร้อม action ในเนื้อหา
	return u.emailSender.Send(email, "Your OTP Code", "Your OTP for "+action+" is: "+otp)
}

func (u *otpUsecase) VerifyOTP(email, otp, action string) error {
	valid, err := u.repo.VerifyOTP(email, otp, action)
	if err != nil {
		return err
	}
	if !valid {
		return errors.New("Invalid or expired OTP")
	}
	return u.repo.MarkVerified(email, otp, action)
}

// EmailSender interface

type EmailSender interface {
	Send(to, subject, body string) error
}

// GenerateRandomOTP สร้างเลข 6 หลักแบบสุ่ม
func GenerateRandomOTP() string {
	r := time.Now().UnixNano() % 1000000
	return fmt.Sprintf("%06d", r)
}
func (u *otpUsecase) VerifyAndGetMetadata(email, otp, action string) (map[string]string, error) {
	if err := u.VerifyOTP(email, otp, action); err != nil {
		return nil, err
	}
	return u.repo.GetOTPMetadata(email, action)
}
func (u *otpUsecase) SendOTPWithMetadata(email, action string, metadata map[string]string) error {
	otp := GenerateRandomOTP()
	expiresAt := time.Now().Add(5 * time.Minute).UTC()

	err := u.repo.SaveOTPWithMetadata(email, otp, action, expiresAt, metadata)
	if err != nil {
		return err
	}

	// ส่ง OTP ทาง Email
	return u.emailSender.Send(email, "Your OTP Code", "Your OTP for "+action+" is: "+otp)
}
