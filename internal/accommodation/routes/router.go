package router

import (
	"database/sql"
	"fmt"
	"net/http"

	"github.com/gorilla/mux"

	accHandler "myapp/internal/accommodation/handler"
	accRepo "myapp/internal/accommodation/repository"
	accUsecase "myapp/internal/accommodation/usecase"
)

func InitRouter(db *sql.DB) *mux.Router {
	r := mux.NewRouter()

	accRepository := accRepo.NewAccommodationRepository(db)
	accUC := accUsecase.NewAccommodationUsecase(accRepository)
	accH := accHandler.NewAccommodationHandler(accUC)

	r.HandleFunc("/accommodations", accH.GetAll).Methods("GET")
	r.HandleFunc("/accommodations/{id}", accH.GetByID).Methods("GET")
	r.HandleFunc("/accommodations", accH.Create).Methods("POST")
	r.HandleFunc("/accommodations", accH.Update).Methods("PUT")
	r.HandleFunc("/accommodations/{id}", accH.Delete).Methods("DELETE")
	r.HandleFunc("/ping", func(w http.ResponseWriter, r *http.Request) {
		fmt.Fprint(w, "pong")
	}).Methods("GET")

	return r
}
