// ✅ router.go สำหรับใช้งาน User + OTP (verify email)
package user

import (
	"database/sql"
	"log"
	"net/http"

	"myapp/internal/user/handler"
	"myapp/internal/user/repository"
	"myapp/internal/user/routes/otpRoutes"
	"myapp/internal/user/usecase"

	"github.com/gorilla/mux"
)

func InitRouter(db *sql.DB) *mux.Router {
	r := mux.NewRouter()

	// ✅ Repository & Usecase
	repo := repository.NewUserRepository(db)
	otpRepo := repository.NewOTPRepository(db)

	userUsecase := usecase.NewUserUsecase(repo)
	emailSender := usecase.NewEmailSender()
	otpUsecase := usecase.NewOTPUsecase(otpRepo, emailSender)

	// ✅ Handler พร้อม OTP
	h := handler.NewUserHandler(userUsecase, otpUsecase)

	// ✅ OTP Routes
	otpRoutes.RegisterOtpRoutes(r, db, userUsecase)

	// ✅ User CRUD
	r.HandleFunc("/users", h.GetAll).Methods("GET")
	r.HandleFunc("/users/{id}", h.GetByID).Methods("GET")
	r.HandleFunc("/users/register", func(w http.ResponseWriter, r *http.Request) {
		log.Println("🔥 Router matched /users/register [POST]")
		h.Create(w, r)
	}).Methods("POST")

	r.HandleFunc("/users/{id}", h.Update).Methods("PUT")
	r.HandleFunc("/users/{id}", h.Delete).Methods("DELETE")

	// ✅ Auth & User Updates
	r.HandleFunc("/login", h.Login).Methods("POST")
	r.HandleFunc("/users/{id}/profile", h.UpdateProfile).Methods("PUT")
	r.HandleFunc("/users/{id}/email", h.UpdateEmail).Methods("PUT")
	r.HandleFunc("/users/{id}/password", h.UpdatePassword).Methods("PUT")
	r.HandleFunc("/users/reset-password", h.ResetPassword).Methods("POST")
	r.HandleFunc("/users/reset-password", h.UpdateProfilePhoto).Methods("PUT")

	// r.Path("/users/profile-photo").Methods("PUT").HandlerFunc(h.UpdateProfilePhoto)

	// log.Println("✅ Route /users/profile-photo [PUT] registered")

	// r.Walk(func(route *mux.Route, router *mux.Router, ancestors []*mux.Route) error {
	// 	path, _ := route.GetPathTemplate()
	// 	methods, _ := route.GetMethods()
	// 	log.Printf("🛣️      Route registered: %s [%s]\n", path, methods)
	// 	return nil
	// })

	return r
}
