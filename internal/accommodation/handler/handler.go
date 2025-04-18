package handler

import (
	"encoding/json"
	"myapp/internal/accommodation/model"
	"myapp/internal/accommodation/usecase"

	"net/http"
	"strconv"

	"github.com/gorilla/mux"
)

type AccommodationHandler struct {
	Usecase usecase.AccommodationUsecase
}

func NewAccommodationHandler(u usecase.AccommodationUsecase) *AccommodationHandler {
	return &AccommodationHandler{Usecase: u}
}

func (h *AccommodationHandler) GetAll(w http.ResponseWriter, r *http.Request) {
	list, err := h.Usecase.GetAll()
	if err != nil {
		http.Error(w, err.Error(), 500)
		return
	}
	json.NewEncoder(w).Encode(list)
}

func (h *AccommodationHandler) GetByID(w http.ResponseWriter, r *http.Request) {
	id, _ := strconv.ParseInt(mux.Vars(r)["id"], 10, 64)
	data, err := h.Usecase.GetByID(id)
	if err != nil {
		http.Error(w, err.Error(), 404)
		return
	}
	json.NewEncoder(w).Encode(data)
}

func (h *AccommodationHandler) Create(w http.ResponseWriter, r *http.Request) {
	var a model.Accommodation
	json.NewDecoder(r.Body).Decode(&a)
	if err := h.Usecase.Create(a); err != nil {
		http.Error(w, err.Error(), 500)
		return
	}
	w.WriteHeader(http.StatusCreated)
}

func (h *AccommodationHandler) Update(w http.ResponseWriter, r *http.Request) {
	var a model.Accommodation
	json.NewDecoder(r.Body).Decode(&a)
	if err := h.Usecase.Update(a); err != nil {
		http.Error(w, err.Error(), 500)
	}
}

func (h *AccommodationHandler) Delete(w http.ResponseWriter, r *http.Request) {
	id, _ := strconv.ParseInt(mux.Vars(r)["id"], 10, 64)
	if err := h.Usecase.Delete(id); err != nil {
		http.Error(w, err.Error(), 500)
	}
}
