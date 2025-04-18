package model

type User struct {
	ID    int64  `json:"id"`    // รหัสผู้ใช้
	Name  string `json:"name"`  // ชื่อผู้ใช้
	Email string `json:"email"` 
	Password string `json:"password"` // อีเมลผู้ใช้
}
