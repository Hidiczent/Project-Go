package main

import (
	"database/sql"
	"log"
	"net/http"

	user "myapp/internal/user"

	_ "github.com/go-sql-driver/mysql"
)

func main() {
	log.Println("🚀 Starting API server...")

	dsn := "jimmy:admin123@tcp(172.20.10.2:3306)/flutterprojecttt?parseTime=true"
	log.Println("🔌 Connecting to MySQL:", dsn)

	db, err := sql.Open("mysql", dsn)
	if err != nil {
		log.Fatal("❌ Failed to connect to database:", err)
	}
	defer db.Close()

	if err := db.Ping(); err != nil {
		log.Fatal("❌ Database not responding:", err)
	}
	log.Println("✅ Connected to MySQL database")

	// Init router from user module
	r := user.InitRouter(db)

	// ✅ Wrap with CORS middleware
	handler := corsMiddleware(r)

	log.Println("🌐 Server running at http://0.0.0.0:5000")
	log.Fatal(http.ListenAndServe("0.0.0.0:5000", handler))
}

// ✅ Global CORS middleware
func corsMiddleware(next http.Handler) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		w.Header().Set("Access-Control-Allow-Origin", "*")
		w.Header().Set("Access-Control-Allow-Headers", "Content-Type")
		w.Header().Set("Access-Control-Allow-Methods", "GET, POST, PUT, DELETE, OPTIONS")
		
		if r.Method == http.MethodOptions {
			w.WriteHeader(http.StatusOK)
			return
		}

		next.ServeHTTP(w, r)
	})
}
