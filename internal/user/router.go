package user

import (
	"database/sql"
	"log"
	"myapp/internal/user/handler"
	"myapp/internal/user/repository"
	"myapp/internal/user/usecase"
	"net/http"

	"github.com/gorilla/mux"
)

func InitRouter(db *sql.DB) *mux.Router {
	r := mux.NewRouter()

	repo := repository.NewUserRepository(db)
	uc := usecase.NewUserUsecase(repo)
	h := handler.NewUserHandler(uc)
		// User CRUD

	r.HandleFunc("/users", h.GetAll).Methods("GET")
	r.HandleFunc("/users/{id}", h.GetByID).Methods("GET")
	r.HandleFunc("/users/register", func(w http.ResponseWriter, r *http.Request) {
		log.Println("🔥 Router matched /users/register [POST]")
		h.Create(w, r)
	}).Methods("POST")

	r.HandleFunc("/users/{id}", h.Update).Methods("PUT")
	r.HandleFunc("/users/{id}", h.Delete).Methods("DELETE")

// Auth & User-specific updates
	r.HandleFunc("/login", h.Login).Methods("POST")
	r.HandleFunc("/users/{id}/profile", h.UpdateProfile).Methods("PUT")
	r.HandleFunc("/users/{id}/email", h.UpdateEmail).Methods("PUT")
	r.HandleFunc("/users/{id}/password", h.UpdatePassword).Methods("PUT")
	return r
}
