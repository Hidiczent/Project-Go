package router

import (
	"database/sql"
	"net/http"

	"github.com/gorilla/mux"

	// accommodation
	accHandler "myapp/internal/accommodation/handler"
	accRepo "myapp/internal/accommodation/repository"
	accUsecase "myapp/internal/accommodation/usecase"

	// district
	districtHandler "myapp/internal/district/handler"
	districtRepo "myapp/internal/district/repository"
	districtUsecase "myapp/internal/district/usecase"
)

func InitRouter(db *sql.DB) *mux.Router {
	r := mux.NewRouter()

	// ✅ Accommodation routes
	accRepository := accRepo.NewAccommodationRepository(db)
	accUC := accUsecase.NewAccommodationUsecase(accRepository)
	accH := accHandler.NewAccommodationHandler(accUC)

	r.HandleFunc("/accommodations", accH.GetAll).Methods("GET")
	r.HandleFunc("/accommodations/{id}", accH.GetByID).Methods("GET")
	r.HandleFunc("/accommodations", accH.Create).Methods("POST")
	r.HandleFunc("/accommodations", accH.Update).Methods("PUT")
	r.HandleFunc("/accommodations/{id}", accH.Delete).Methods("DELETE")

	// ✅ District routes
	dRepo := districtRepo.NewDistrictRepository(db)
	dUC := districtUsecase.NewDistrictUsecase(dRepo)
	dH := districtHandler.NewDistrictHandler(dUC)

	r.HandleFunc("/districts", dH.GetAll).Methods("GET")
	r.HandleFunc("/districts/{id}", dH.GetByID).Methods("GET")
	r.HandleFunc("/districts", dH.Create).Methods("POST")
	r.HandleFunc("/districts", dH.Update).Methods("PUT")
	r.HandleFunc("/districts/{id}", dH.Delete).Methods("DELETE")

	// ✅ Test route
	r.HandleFunc("/ping", func(w http.ResponseWriter, r *http.Request) {
		w.Write([]byte("pong"))
	}).Methods("GET")

	return r
}
