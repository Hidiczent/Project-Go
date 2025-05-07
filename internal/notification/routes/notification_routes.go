package routes

import (
	"database/sql"
	"github.com/gorilla/mux"
	"myapp/internal/notification/handler"
	"myapp/internal/notification/repository"
	"myapp/internal/notification/usecase"
)

func RegisterNotificationRoutes(r *mux.Router, db *sql.DB) {
	repo := repository.NewNotificationRepository(db)
	uc := usecase.NewNotificationUseCase(repo)
	h := handler.NewNotificationHandler(uc)

	// ✅ CRUD สำหรับ /notifications
	r.HandleFunc("/notifications", h.GetAll).Methods("GET")
	r.HandleFunc("/notifications/{id:[0-9]+}", h.GetByID).Methods("GET")
	r.HandleFunc("/notifications", h.Create).Methods("POST")
	r.HandleFunc("/notifications", h.Update).Methods("PUT")
	r.HandleFunc("/notifications/{id:[0-9]+}", h.Delete).Methods("DELETE")
}
