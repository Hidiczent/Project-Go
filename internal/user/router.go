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
	r.HandleFunc("/users", h.GetAll).Methods("GET")
	r.HandleFunc("/users/{id}", h.GetByID).Methods("GET")
	r.HandleFunc("/users", func(w http.ResponseWriter, r *http.Request) {
		log.Println("🔥 Router matched /users POST")
		h.Create(w, r)
	}).Methods("POST")

	r.HandleFunc("/users", h.Update).Methods("PUT")
	r.HandleFunc("/users/{id}", h.Delete).Methods("DELETE")

	return r
}
