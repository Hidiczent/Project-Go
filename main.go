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

	dsn := "jimmy:admin123@tcp(192.168.243.213:3306)/flutterproject"
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

	r := user.InitRouter(db)
	log.Println("✅ Router initialized")

	log.Println("🌐 Server running at http://localhost:5000")
	http.ListenAndServe(":5000", r)
}
