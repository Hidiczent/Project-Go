package handler

import (
	"encoding/json"
	"log"
	"net/http"
	"strconv"

	"myapp/internal/notification/model"
	"myapp/internal/notification/usecase"

	"github.com/gorilla/mux"
)

type NotificationHandler struct {
	Usecase usecase.NotificationUseCase
}

func NewNotificationHandler(u usecase.NotificationUseCase) *NotificationHandler {
	return &NotificationHandler{u}
}

func (h *NotificationHandler) GetAll(w http.ResponseWriter, r *http.Request) {
	data, err := h.Usecase.GetAll()
	if err != nil {
		http.Error(w, "Error fetching notifications", http.StatusInternalServerError)
		return
	}
	json.NewEncoder(w).Encode(data)
}

func (h *NotificationHandler) GetByID(w http.ResponseWriter, r *http.Request) {
	id, _ := strconv.Atoi(mux.Vars(r)["id"])
	n, err := h.Usecase.GetByID(id)
	if err != nil {
		http.Error(w, "Notification not found", http.StatusNotFound)
		return
	}
	json.NewEncoder(w).Encode(n)
}

func (h *NotificationHandler) Create(w http.ResponseWriter, r *http.Request) {
	var n model.Notification
	if err := json.NewDecoder(r.Body).Decode(&n); err != nil {
		http.Error(w, "Invalid JSON", http.StatusBadRequest)
		return
	}

	if err := h.Usecase.Create(n); err != nil {
		log.Println("❌ Create failed:", err) // ✅ เพิ่ม log error จริงตรงนี้
		http.Error(w, "Create failed", http.StatusBadRequest)
		return
	}

	w.WriteHeader(http.StatusCreated)
}

func (h *NotificationHandler) Update(w http.ResponseWriter, r *http.Request) {
	var n model.Notification
	json.NewDecoder(r.Body).Decode(&n)
	if err := h.Usecase.Update(n); err != nil {
		http.Error(w, "Update failed", http.StatusBadRequest)
		return
	}
	w.WriteHeader(http.StatusOK)
}

func (h *NotificationHandler) Delete(w http.ResponseWriter, r *http.Request) {
	id, _ := strconv.Atoi(mux.Vars(r)["id"])
	if err := h.Usecase.Delete(id); err != nil {
		http.Error(w, "Delete failed", http.StatusInternalServerError)
		return
	}
	w.WriteHeader(http.StatusOK)
}
