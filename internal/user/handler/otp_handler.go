package handler

import (
	"encoding/json"
	"myapp/internal/user/usecase"
	"net/http"
)

type OTPHandler struct {
	Usecase     usecase.OTPUsecase
	UserUsecase usecase.UserUsecase
}

func NewOTPHandler(otpUC usecase.OTPUsecase, userUC usecase.UserUsecase) *OTPHandler {
	return &OTPHandler{
		Usecase:     otpUC,
		UserUsecase: userUC, // ✅ เพิ่มตรงนี้
	}
}

// ✅ [POST] /otp/send - ส่ง OTP ไปยังอีเมล
func (h *OTPHandler) SendOTP(w http.ResponseWriter, r *http.Request) {
	var req struct {
		Email    string            `json:"email"`
		Action   string            `json:"action"`
		Metadata map[string]string `json:"metadata"`
	}

	if err := json.NewDecoder(r.Body).Decode(&req); err != nil || req.Email == "" || req.Action == "" {
		http.Error(w, "Invalid request", http.StatusBadRequest)
		return
	}

	if req.Action == "register" {
		// ✅ สมัครสมาชิก => ต้องยังไม่มี email นี้
		if _, err := h.UserUsecase.GetByEmail(req.Email); err == nil {
			http.Error(w, "This email is already registered", http.StatusConflict)
			return
		}
	} else {
		// ✅ action อื่นๆ => ต้องมี email นี้ในระบบ
		if _, err := h.UserUsecase.GetByEmail(req.Email); err != nil {
			http.Error(w, "This email is not registered", http.StatusNotFound)
			return
		}
	}
	if err := h.Usecase.SendOTPWithMetadata(req.Email, req.Action, req.Metadata); err != nil {
		http.Error(w, "Failed to send OTP", http.StatusInternalServerError)
		return
	}
	// ✅ ส่ง OTP
	if err := h.Usecase.SendOTP(req.Email, req.Action); err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}

	json.NewEncoder(w).Encode(map[string]string{"message": "OTP sent successfully"})
}

// ✅ [POST] /otp/verify - ตรวจสอบ OTP
func (h *OTPHandler) VerifyOTP(w http.ResponseWriter, r *http.Request) {
	var req struct {
		Email  string `json:"email"`
		Otp    string `json:"otp"`
		Action string `json:"action"`
	}
	if err := json.NewDecoder(r.Body).Decode(&req); err != nil || req.Email == "" || req.Otp == "" || req.Action == "" {
		http.Error(w, "Invalid request", http.StatusBadRequest)
		return
	}

	if err := h.Usecase.VerifyOTP(req.Email, req.Otp, req.Action); err != nil {
		http.Error(w, err.Error(), http.StatusUnauthorized)
		return
	}

	json.NewEncoder(w).Encode(map[string]string{"message": "OTP verified successfully"})
}
