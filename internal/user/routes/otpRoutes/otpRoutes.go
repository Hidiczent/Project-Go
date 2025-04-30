package otpRoutes

import (
	"database/sql"
	"github.com/gorilla/mux"
	otpHandler "myapp/internal/user/handler"
	otpRepo "myapp/internal/user/repository"
	otpUsecase "myapp/internal/user/usecase"
)

func RegisterOtpRoutes(r *mux.Router, db *sql.DB, userUC otpUsecase.UserUsecase) {
	repo := otpRepo.NewOTPRepository(db)
	sender := otpUsecase.NewEmailSender()
	usecase := otpUsecase.NewOTPUsecase(repo, sender)

	handler := otpHandler.NewOTPHandler(usecase, userUC)

	r.HandleFunc("/otp/send", handler.SendOTP).Methods("POST")
	r.HandleFunc("/otp/verify", handler.VerifyOTP).Methods("POST")
	r.HandleFunc("/otp/confirm-register", handler.ConfirmRegister).Methods("POST") // ✅ เพิ่มตรงนี้

}
