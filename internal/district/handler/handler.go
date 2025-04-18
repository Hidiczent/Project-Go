package handler

import (
	"encoding/json"
	"net/http"
	"strconv"

	"github.com/gorilla/mux"
	"myapp/internal/district/model"
	"myapp/internal/district/usecase"
)

type DistrictHandler struct {
	Usecase usecase.DistrictUsecase
}

func NewDistrictHandler(u usecase.DistrictUsecase) *DistrictHandler {
	return &DistrictHandler{Usecase: u}
}

func (h *DistrictHandler) GetAll(w http.ResponseWriter, r *http.Request) {
	list, err := h.Usecase.GetAll()
	if err != nil {
		http.Error(w, err.Error(), 500)
		return
	}
	json.NewEncoder(w).Encode(list)
}

func (h *DistrictHandler) GetByID(w http.ResponseWriter, r *http.Request) {
	id, _ := strconv.ParseInt(mux.Vars(r)["id"], 10, 64)
	data, err := h.Usecase.GetByID(id)
	if err != nil {
		http.Error(w, err.Error(), 404)
		return
	}
	json.NewEncoder(w).Encode(data)
}

func (h *DistrictHandler) Create(w http.ResponseWriter, r *http.Request) {
	var d model.District
	json.NewDecoder(r.Body).Decode(&d)
	if err := h.Usecase.Create(d); err != nil {
		http.Error(w, err.Error(), 500)
		return
	}
	w.WriteHeader(http.StatusCreated)
}

func (h *DistrictHandler) Update(w http.ResponseWriter, r *http.Request) {
	var d model.District
	json.NewDecoder(r.Body).Decode(&d)
	if err := h.Usecase.Update(d); err != nil {
		http.Error(w, err.Error(), 500)
	}
}

func (h *DistrictHandler) Delete(w http.ResponseWriter, r *http.Request) {
	id, _ := strconv.ParseInt(mux.Vars(r)["id"], 10, 64)
	if err := h.Usecase.Delete(id); err != nil {
		http.Error(w, err.Error(), 500)
	}
}
