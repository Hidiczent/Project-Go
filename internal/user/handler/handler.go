package handler

import (
	"encoding/json"
	"log"
	"myapp/internal/user/model"
	"myapp/internal/user/usecase"
	"net/http"
	"strconv"

	"github.com/gorilla/mux"
)

type UserHandler struct {
	Usecase usecase.UserUsecase
}

func NewUserHandler(u usecase.UserUsecase) *UserHandler {
	return &UserHandler{Usecase: u}
}

func (h *UserHandler) GetAll(w http.ResponseWriter, r *http.Request) {
	users, err := h.Usecase.GetAll()
	if err != nil {
		http.Error(w, err.Error(), 500)
		return
	}
	json.NewEncoder(w).Encode(users)
}

func (h *UserHandler) GetByID(w http.ResponseWriter, r *http.Request) {
	idStr := mux.Vars(r)["id"]
	id, _ := strconv.ParseInt(idStr, 10, 64)

	user, err := h.Usecase.GetByID(id)
	if err != nil {
		http.Error(w, err.Error(), 404)
		return
	}
	json.NewEncoder(w).Encode(user)
}
func (h *UserHandler) Create(w http.ResponseWriter, r *http.Request) {
	log.Println("📥 [POST] /users - Create called")

	var user model.User
	err := json.NewDecoder(r.Body).Decode(&user)
	if err != nil {
		log.Println("❌ Failed to decode body:", err)
		http.Error(w, err.Error(), http.StatusBadRequest)
		return
	}

	log.Printf("📝 Creating user: Name=%s, Email=%s\n", user.Name, user.Email)

	if err := h.Usecase.Create(user); err != nil {
		log.Println("❌ Failed to create user:", err)
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}

	log.Println("✅ User created successfully")
	w.WriteHeader(http.StatusCreated)
}

func (h *UserHandler) Update(w http.ResponseWriter, r *http.Request) {
	var user model.User
	json.NewDecoder(r.Body).Decode(&user)
	if err := h.Usecase.Update(user); err != nil {
		http.Error(w, err.Error(), 500)
	}
	w.WriteHeader(http.StatusOK)
}

func (h *UserHandler) Delete(w http.ResponseWriter, r *http.Request) {
	idStr := mux.Vars(r)["id"]
	id, _ := strconv.ParseInt(idStr, 10, 64)
	if err := h.Usecase.Delete(id); err != nil {
		http.Error(w, err.Error(), 500)
	}
	w.WriteHeader(http.StatusOK)
}
