package handler

import (
	"encoding/base64"
	"encoding/json"

	"fmt"
	"io"
	"log"
	"myapp/internal/user/model"
	"myapp/internal/user/usecase"
	"net/http"
	"os"
	"strconv"
	"strings"
	"time"

	"github.com/golang-jwt/jwt/v5"
	"github.com/gorilla/mux"
	"golang.org/x/crypto/bcrypt"
)

var jwtSecret = []byte("MySuperSecretKey")

type UserHandler struct {
	Usecase    usecase.UserUsecase
	OTPUsecase usecase.OTPUsecase // ✅ Inject OTPUsecase
}
type LoginRequest struct {
	Email    string `json:"email"`
	Password string `json:"password"`
}

func NewUserHandler(userUC usecase.UserUsecase, otpUC usecase.OTPUsecase) *UserHandler {
	return &UserHandler{
		Usecase:    userUC,
		OTPUsecase: otpUC,
	}
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

// GetUserByID
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

// CreateUser
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

	// ✅ ส่ง OTP สำหรับ verify_email
	if err := h.OTPUsecase.SendOTP(user.Email, "verify_email"); err != nil {
		log.Printf("⚠️ Failed to send OTP: %v", err)
		// ไม่ return error เพื่อให้ user ยังใช้งานได้แม้ส่ง OTP ไม่สำเร็จ
	}

	// ✅ ส่งกลับข้อความสำเร็จ
	w.WriteHeader(http.StatusCreated)
	w.Header().Set("Content-Type", "application/json")
	json.NewEncoder(w).Encode(map[string]interface{}{
		"message": "✅ User created successfully. Please verify your email.",
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
	log.Println("📥 Received:", req.FirstName, req.LastName, req.PhoneNumber)
	log.Println("📝 Updating lastname:", req.LastName)
	log.Println("📝 Updating phone:", req.PhoneNumber)

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
	log.Printf("📥 Change password for userID: %d", id)
	log.Printf("🔐 Old: %s | New: %s | Confirm: %s", req.OldPassword, req.NewPassword, req.ConfirmPassword)

	json.NewEncoder(w).Encode(map[string]string{
		"message": "Password updated successfully",
	})
}

// ✅ 4.ResetPassword
func (h *UserHandler) ResetPassword(w http.ResponseWriter, r *http.Request) {
	var req struct {
		Email       string `json:"email"`
		OTP         string `json:"otp"`
		NewPassword string `json:"new_password"`
	}

	if err := json.NewDecoder(r.Body).Decode(&req); err != nil {
		http.Error(w, "Invalid request", http.StatusBadRequest)
		return
	}

	// ✅ ตรวจสอบ OTP
	if err := h.OTPUsecase.VerifyOTP(req.Email, req.OTP, "reset_password"); err != nil {
		http.Error(w, "Invalid or expired OTP", http.StatusUnauthorized)
		return
	}

	// ✅ เข้ารหัสรหัสผ่านใหม่
	hashedPassword, err := bcrypt.GenerateFromPassword([]byte(req.NewPassword), bcrypt.DefaultCost)
	if err != nil {
		http.Error(w, "Failed to hash password", http.StatusInternalServerError)
		return
	}

	// ✅ อัปเดตรหัสผ่านใหม่
	user, err := h.Usecase.GetByEmail(req.Email)
	if err != nil {
		http.Error(w, "User not found", http.StatusNotFound)
		return
	}

	if err := h.Usecase.UpdatePassword(user.ID, string(hashedPassword)); err != nil {
		http.Error(w, "Failed to update password", http.StatusInternalServerError)
		return
	}

	json.NewEncoder(w).Encode(map[string]string{
		"message": "Password reset successfully",
	})
}

func (h *UserHandler) UpdateProfilePhoto(w http.ResponseWriter, r *http.Request) {
	defer func() {
		if err := recover(); err != nil {
			log.Printf("🔥 Panic in UpdateProfilePhoto: %v\n", err)
			http.Error(w, "Internal Server Error", http.StatusInternalServerError)
		}
	}()

	log.Println("🟡 Entered UpdateProfilePhoto handler")
	log.Printf("🧩 h.Usecase is nil? = %v", h.Usecase == nil)

	// ✅ ดึง user ID จาก token
	userID, err := getUserIDFromToken(r)
	if err != nil {
		log.Printf("❌ Failed to get user ID from token: %v\n", err)
		http.Error(w, "Unauthorized", http.StatusUnauthorized)
		return
	}
	log.Printf("✅ Extracted user ID: %d\n", userID)

	// ✅ Parse multipart form
	log.Println("🧩 Parsing multipart form...")
	err = r.ParseMultipartForm(10 << 20) // 10 MB
	if err != nil {
		log.Printf("❌ Error parsing form data: %v\n", err)
		http.Error(w, "Error parsing form data", http.StatusBadRequest)
		return
	}

	file, header, err := r.FormFile("photo")
	if err != nil {
		log.Printf("❌ Error reading form file: %v\n", err)
		http.Error(w, "Photo is required", http.StatusBadRequest)
		return
	}
	defer file.Close()

	log.Printf("✅ Received photo: filename=%s, size=%d bytes\n", header.Filename, header.Size)

	// ✅ ตรวจสอบโฟลเดอร์ uploads และสร้างถ้ายังไม่มี
	uploadDir := "uploads"
	log.Printf("🗂️ Checking upload directory: %s", uploadDir)
	if _, err := os.Stat(uploadDir); os.IsNotExist(err) {
		log.Println("📁 Upload directory does not exist, creating...")
		if err := os.MkdirAll(uploadDir, os.ModePerm); err != nil {
			log.Printf("❌ Failed to create upload directory: %v\n", err)
			http.Error(w, "Failed to create uploads folder", http.StatusInternalServerError)
			return
		}
	}

	// ✅ ตั้งชื่อไฟล์ใหม่แบบปลอดภัย
	uploadPath := fmt.Sprintf("%s/user_%d_%s", uploadDir, userID, header.Filename)
	log.Printf("📤 Saving file to: %s\n", uploadPath)

	dst, err := os.Create(uploadPath)
	if err != nil {
		log.Printf("❌ Failed to create file: %v\n", err)
		http.Error(w, "Failed to save image", http.StatusInternalServerError)
		return
	}
	defer dst.Close()

	_, err = io.Copy(dst, file)
	if err != nil {
		log.Printf("❌ Failed to write file: %v\n", err)
		http.Error(w, "Failed to save image", http.StatusInternalServerError)
		return
	}

	// ✅ บันทึก path รูปใน database
	log.Println("💾 Updating photo path in database...")
	if err := h.Usecase.UpdateProfilePhoto(userID, uploadPath); err != nil {
		log.Printf("❌ Failed to update DB: %v\n", err)
		http.Error(w, "Failed to update profile photo in DB", http.StatusInternalServerError)
		return
	}

	// ✅ ตอบกลับเป็น JSON
	log.Println("✅ Profile photo updated successfully")
	w.Header().Set("Content-Type", "application/json")
	err = json.NewEncoder(w).Encode(map[string]string{
		"message": "✅ Profile photo updated successfully",
		"path":    uploadPath,
	})
	if err != nil {
		log.Printf("❌ Failed to encode JSON response: %v\n", err)
	}
}

func getUserIDFromToken(r *http.Request) (int64, error) {
	log.Println("🔐 Entered getUserIDFromToken (No verify)")

	authHeader := r.Header.Get("Authorization")
	log.Printf("🔐 Authorization Header: %s\n", authHeader)
	if authHeader == "" {
		return 0, fmt.Errorf("missing authorization header")
	}

	tokenStr := strings.TrimPrefix(authHeader, "Bearer ")
	log.Printf("🔐 JWT Token string: %s\n", tokenStr)

	parts := strings.Split(tokenStr, ".")
	if len(parts) != 3 {
		return 0, fmt.Errorf("invalid token format")
	}

	payload := parts[1]
	// base64 decode payload
	payloadBytes, err := base64.RawURLEncoding.DecodeString(payload)
	if err != nil {
		log.Printf("❌ Failed to decode payload: %v\n", err)
		return 0, fmt.Errorf("invalid payload encoding")
	}

	var claims map[string]interface{}
	if err := json.Unmarshal(payloadBytes, &claims); err != nil {
		log.Printf("❌ Failed to unmarshal payload: %v\n", err)
		return 0, fmt.Errorf("invalid payload json")
	}

	log.Printf("🧾 Token Claims: %#v\n", claims)

	rawUID := claims["user_id"]
	log.Printf("👉 [DEBUG] claims[user_id] = %#v (type: %T)\n", rawUID, rawUID)

	uidStr := fmt.Sprintf("%v", rawUID)
	uidParsed, err := strconv.ParseInt(uidStr, 10, 64)
	if err != nil {
		log.Printf("❌ Failed to parse user_id: %v", err)
		return 0, fmt.Errorf("invalid user_id format")
	}

	return uidParsed, nil
}

func (h *OTPHandler) ConfirmRegister(w http.ResponseWriter, r *http.Request) {
	var req struct {
		Email  string `json:"email"`
		Otp    string `json:"otp"`
		Action string `json:"action"`
	}

	// ✅ Decode JSON
	if err := json.NewDecoder(r.Body).Decode(&req); err != nil {
		http.Error(w, "Invalid JSON format", http.StatusBadRequest)
		return
	}

	// ✅ Validate input
	if req.Email == "" || req.Otp == "" || req.Action != "register" {
		http.Error(w, "Missing or invalid fields", http.StatusBadRequest)
		return
	}

	// ✅ Verify OTP and fetch metadata
	otpData, err := h.Usecase.VerifyAndGetMetadata(req.Email, req.Otp, req.Action)
	if err != nil {
		http.Error(w, "Invalid or expired OTP", http.StatusUnauthorized)
		return
	}

	// ✅ Check if user already exists
	if _, err := h.UserUsecase.GetByEmail(req.Email); err == nil {
		http.Error(w, "Email is already registered", http.StatusConflict)
		return
	}

	// ✅ Check required metadata fields
	firstName, ok1 := otpData["first_name"]
	password, ok2 := otpData["password"]
	if !ok1 || !ok2 || firstName == "" || password == "" {
		http.Error(w, "Missing required registration data", http.StatusBadRequest)
		return
	}

	// ✅ Hash password
	hashedPassword, err := bcrypt.GenerateFromPassword([]byte(password), bcrypt.DefaultCost)
	if err != nil {
		http.Error(w, "Failed to hash password", http.StatusInternalServerError)
		return
	}

	// ✅ Create user
	user := model.User{
		FirstName: firstName,
		Email:     req.Email,
		Password:  string(hashedPassword),
	}

	if err := h.UserUsecase.Create(user); err != nil {
		http.Error(w, "Failed to create user", http.StatusInternalServerError)
		return
	}

	// ✅ Respond success
	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(http.StatusOK)
	json.NewEncoder(w).Encode(map[string]string{
		"message": "✅ User registered successfully.",
	})
}
