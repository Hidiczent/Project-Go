package handler

import (
	"encoding/json"
	"github.com/golang-jwt/jwt/v5"
	"github.com/gorilla/mux"
	"golang.org/x/crypto/bcrypt"
	"log"
	"myapp/internal/user/model"
	"myapp/internal/user/usecase"
	"net/http"
	"strconv"
	"strings"
	"time"
)

var jwtSecret = []byte("MySuperSecretKey")

type UserHandler struct {
	Usecase usecase.UserUsecase
}
type LoginRequest struct {
	Email    string `json:"email"`
	Password string `json:"password"`
}

func NewUserHandler(u usecase.UserUsecase) *UserHandler {
	return &UserHandler{Usecase: u}
}
func (h *UserHandler) GetAll(w http.ResponseWriter, r *http.Request) {
	users, err := h.Usecase.GetAll()
	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}

	// ✅ Optional: แปลงข้อมูลก่อนแสดง (ถ้าอยากซ่อนรหัสผ่าน)
	for i := range users {
		users[i].Password = "" // ซ่อนรหัสผ่านจากการส่งกลับ
	}

	w.Header().Set("Content-Type", "application/json")
	json.NewEncoder(w).Encode(users)
}

func (h *UserHandler) GetByID(w http.ResponseWriter, r *http.Request) {
	idStr := mux.Vars(r)["id"]
	id, err := strconv.ParseInt(idStr, 10, 64)
	if err != nil {
		http.Error(w, "Invalid ID format", http.StatusBadRequest)
		return
	}

	user, err := h.Usecase.GetByID(id)
	if err != nil {
		http.Error(w, err.Error(), http.StatusNotFound)
		return
	}

	user.Password = "" // ✅ ซ่อนรหัสผ่านใน response

	w.Header().Set("Content-Type", "application/json")
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

	// ✅ เช็กข้อมูลจำเป็น
	if user.FirstName == "" || user.Email == "" || user.Password == "" {
		http.Error(w, "first_name, email, and password are required", http.StatusBadRequest)
		return
	}

	// ✅ ตรวจสอบว่า email ซ้ำหรือไม่
	if _, err := h.Usecase.GetByEmail(user.Email); err == nil {
		http.Error(w, "Email is already in use", http.StatusConflict) // 409
		return
	}

	// ✅ เข้ารหัส password
	hashedPassword, err := bcrypt.GenerateFromPassword([]byte(user.Password), bcrypt.DefaultCost)
	if err != nil {
		log.Println("❌ Failed to hash password:", err)
		http.Error(w, "Internal Server Error", http.StatusInternalServerError)
		return
	}
	user.Password = string(hashedPassword)

	log.Printf("📝 Creating user: FirstName=%s, Email=%s\n", user.FirstName, user.Email)

	// ✅ สร้างผู้ใช้
	if err := h.Usecase.Create(user); err != nil {
		log.Println("❌ Failed to create user:", err)
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}

	// ✅ ส่งกลับข้อความสำเร็จ
	w.WriteHeader(http.StatusCreated)
	w.Header().Set("Content-Type", "application/json")
	json.NewEncoder(w).Encode(map[string]interface{}{
		"message": "✅ User created successfully",
		"email":   user.Email,
		"name":    user.FirstName,
	})
}

// ✅ Update User profile
func (h *UserHandler) Update(w http.ResponseWriter, r *http.Request) {
	idStr := mux.Vars(r)["id"]
	id, err := strconv.ParseInt(idStr, 10, 64)
	if err != nil {
		http.Error(w, "Invalid user ID", http.StatusBadRequest)
		return
	}

	var user model.User
	if err := json.NewDecoder(r.Body).Decode(&user); err != nil {
		http.Error(w, "Invalid JSON", http.StatusBadRequest)
		return
	}
	user.ID = id // set user ID จาก URL

	if err := h.Usecase.Update(user); err != nil {
		http.Error(w, "Failed to update user: "+err.Error(), http.StatusInternalServerError)
		return
	}

	w.WriteHeader(http.StatusOK)
	json.NewEncoder(w).Encode(map[string]string{
		"message": "User updated successfully",
	})
}

func (h *UserHandler) Delete(w http.ResponseWriter, r *http.Request) {
	idStr := mux.Vars(r)["id"]
	id, _ := strconv.ParseInt(idStr, 10, 64)
	if err := h.Usecase.Delete(id); err != nil {
		http.Error(w, err.Error(), 500)
	}
	w.WriteHeader(http.StatusOK)
}
func (h *UserHandler) Login(w http.ResponseWriter, r *http.Request) {
	var req LoginRequest
	log.Println("📥 Login request for email:", req.Email)

	if err := json.NewDecoder(r.Body).Decode(&req); err != nil {
		http.Error(w, "Invalid request", http.StatusBadRequest)
		return
	}
	log.Println("📥 Login request for email:", req.Email)

	user, err := h.Usecase.GetByEmail(strings.TrimSpace(req.Email))
	if err != nil {
		log.Println("❌ Error fetching user by email:", err)
		http.Error(w, "User not found", http.StatusUnauthorized)
		return
	}

	err = bcrypt.CompareHashAndPassword([]byte(user.Password), []byte(req.Password))
	if err != nil {
		http.Error(w, "Incorrect password", http.StatusUnauthorized)
		return
	}

	// ✅ สร้าง JWT token
	token := jwt.NewWithClaims(jwt.SigningMethodHS256, jwt.MapClaims{
		"user_id": user.ID,
		"email":   user.Email,
		"role":    user.Role,
		"exp":     time.Now().Add(time.Hour * 24).Unix(),
	})

	tokenString, err := token.SignedString(jwtSecret)
	if err != nil {
		http.Error(w, "Could not create token", http.StatusInternalServerError)
		return
	}

	// ✅ ส่งกลับ token + ข้อมูล user
	w.Header().Set("Content-Type", "application/json")
	json.NewEncoder(w).Encode(map[string]string{
		"token": tokenString,
		"name":  user.FirstName,
		"email": user.Email,
	})

	log.Println("✅ Login successful for:", user.Email)
}

// ✅ 1. Update User Profile Handler (/users/{id}/profile)
func (h *UserHandler) UpdateProfile(w http.ResponseWriter, r *http.Request) {
	idStr := mux.Vars(r)["id"]
	log.Println("🧩 mux.Vars:", mux.Vars(r))
	log.Println("🧩 idStr received:", idStr)
	id, err := strconv.ParseInt(idStr, 10, 64)
	if err != nil {
		http.Error(w, "Invalid user ID", http.StatusBadRequest)
		return
	}

	var req struct {
		FirstName   string  `json:"first_name"`
		LastName    string  `json:"lastname"`
		PhoneNumber *int64  `json:"phone_number"`
		Photo       *string `json:"photo"`
	}

	if err := json.NewDecoder(r.Body).Decode(&req); err != nil {
		http.Error(w, "Invalid JSON", http.StatusBadRequest)
		return
	}

	user, err := h.Usecase.GetByID(id)
	if err != nil {
		http.Error(w, "User not found", http.StatusNotFound)
		return
	}

	user.FirstName = req.FirstName
	user.LastName = req.LastName
	user.PhoneNumber = req.PhoneNumber
	user.Photo = req.Photo

	if err := h.Usecase.Update(user); err != nil {
		http.Error(w, "Failed to update profile", http.StatusInternalServerError)
		return
	}

	json.NewEncoder(w).Encode(map[string]string{"message": "Profile updated successfully"})
}

// Update Email
func (h *UserHandler) UpdateEmail(w http.ResponseWriter, r *http.Request) {
	idStr := mux.Vars(r)["id"]

	id, err := strconv.ParseInt(idStr, 10, 64)
	if err != nil {
		http.Error(w, "Invalid user ID", http.StatusBadRequest)
		return
	}

	var req struct {
		Email string `json:"email"`
	}

	if err := json.NewDecoder(r.Body).Decode(&req); err != nil || req.Email == "" {
		http.Error(w, "Invalid email", http.StatusBadRequest)
		return
	}

	user, err := h.Usecase.GetByID(id)
	if err != nil {
		http.Error(w, "User not found", http.StatusNotFound)
		return
	}

	user.Email = strings.TrimSpace(req.Email)
	email := strings.TrimSpace(req.Email)

	if err := h.Usecase.UpdateEmail(id, email); err != nil {
		if err.Error() == "email is already in use" {
			http.Error(w, "Email is already in use", http.StatusConflict)
			return
		}
		http.Error(w, "Failed to update email", http.StatusInternalServerError)
		return
	}
}

// ✅ 3. Update Password (/users/{id}/password)
func (h *UserHandler) UpdatePassword(w http.ResponseWriter, r *http.Request) {
	idStr := mux.Vars(r)["id"]
	id, err := strconv.ParseInt(idStr, 10, 64)
	if err != nil {
		http.Error(w, "Invalid user ID", http.StatusBadRequest)
		return
	}

	var req struct {
		OldPassword     string `json:"old_password"`
		NewPassword     string `json:"new_password"`
		ConfirmPassword string `json:"confirm_password"` // ✅ เพิ่มตรงนี้
	}

	if err := json.NewDecoder(r.Body).Decode(&req); err != nil {
		http.Error(w, "Invalid JSON", http.StatusBadRequest)
		return
	}

	if req.NewPassword != req.ConfirmPassword {
		http.Error(w, "New password and confirm password do not match", http.StatusBadRequest)
		return
	}

	user, err := h.Usecase.GetByID(id)
	if err != nil {
		http.Error(w, "User not found", http.StatusNotFound)
		return
	}

	if err := bcrypt.CompareHashAndPassword([]byte(user.Password), []byte(req.OldPassword)); err != nil {
		http.Error(w, "Old password is incorrect", http.StatusUnauthorized)
		return
	}

	hashedPassword, err := bcrypt.GenerateFromPassword([]byte(req.NewPassword), bcrypt.DefaultCost)
	if err != nil {
		http.Error(w, "Failed to hash new password", http.StatusInternalServerError)
		return
	}

	if err := h.Usecase.UpdatePassword(id, string(hashedPassword)); err != nil {
		http.Error(w, "Failed to update password", http.StatusInternalServerError)
		return
	}

	json.NewEncoder(w).Encode(map[string]string{
		"message": "Password updated successfully",
	})
}
