package main

import (
	"database/sql"
	"fmt"
	_ "github.com/go-sql-driver/mysql"
	"github.com/joho/godotenv"
	"log"
	notification "myapp/internal/notification/routes"
	user "myapp/internal/user/routes"
	"net/http"
)

func main() {
	godotenv.Load()
	// setup IP  & Log API server
	ip := "192.168.80.213"
	log.Println("🚀 Starting API server...")

	// Config    username ,pwd,database
	dsn := fmt.Sprintf("jimmy:admin123@tcp(%s:3306)/flutterprojecttt?parseTime=true", ip)
	// Log Connecting to Database
	log.Println("🔌 Connecting to MySQL:", dsn)

	// show Error when Errors
	db, err := sql.Open("mysql", dsn)
	if err != nil {
		log.Fatal("❌ Failed to connect to database:", err)
	}
	defer db.Close()
	// show Error when Database not responding
	if err := db.Ping(); err != nil {
		log.Fatal("❌ Database not responding:", err)
	}
	log.Println("✅ Connected to MySQL database")

	// Init router from user module
	r := user.InitRouter(db)
	notification.RegisterNotificationRoutes(r, db) // ✅ เพิ่มตรงนี้

	// ✅ Wrap with CORS middleware
	handler := corsMiddleware(r)

	log.Println("🌐 Server running at http://0.0.0.0:5000")
	log.Fatal(http.ListenAndServe("0.0.0.0:5000", handler))

	log.Println("✅ Routes initialized:")
	log.Println(" - PUT /users/profile-photo")

}

func corsMiddleware(next http.Handler) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		defer func() {
			if err := recover(); err != nil {
				log.Printf("🔥 Recovered from panic: %v\n", err)
				http.Error(w, "Internal Server Error", http.StatusInternalServerError)
			}
		}()
		log.Printf("🌐 Incoming request: %s %s, Content-Type: %s", r.Method, r.URL.Path, r.Header.Get("Content-Type"))
		log.Printf("🌐 Incoming request: %s %s", r.Method, r.URL.Path)
		w.Header().Set("Access-Control-Allow-Origin", "*")
		w.Header().Set("Access-Control-Allow-Headers", "Content-Type, Authorization")
		w.Header().Set("Access-Control-Allow-Methods", "GET, POST, PUT, PATCH, DELETE, OPTIONS")

		if r.Method == http.MethodOptions {
			w.WriteHeader(http.StatusOK)
			return
		}

		next.ServeHTTP(w, r)
	})
}

// go func() {
// 	for {
// 		time.Sleep(5 * time.Minute) // ✅ ทุกๆ 5 นาที
// 		_, err := db.Exec("DELETE FROM otps WHERE expires_at < NOW()")
// 		if err != nil {
// 			log.Println("❌ Failed to delete expired OTPs:", err)
// 		} else {
// 			log.Println("🧹 Expired OTPs deleted successfully")
// 		}
// 	}
// }
