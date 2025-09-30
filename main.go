package main

import (
	"database/sql"
	"encoding/json"
	"fmt"
	"log"
	"net/http"
	"os"

	_ "github.com/go-sql-driver/mysql"
	"github.com/gorilla/mux"
	"go-careers/handlers"
	"go-careers/middleware"
	"go-careers/repository"
)

func initDB() *sql.DB {
	dbHost := getEnv("DB_HOST", "localhost")
	dbPort := getEnv("DB_PORT", "3306")
	dbUser := getEnv("DB_USER", "root")
	dbPassword := getEnv("DB_PASSWORD", "rootpassword")
	dbName := getEnv("DB_NAME", "go_careers")

	dsn := fmt.Sprintf("%s:%s@tcp(%s:%s)/%s?parseTime=true",
		dbUser, dbPassword, dbHost, dbPort, dbName)

	db, err := sql.Open("mysql", dsn)
	if err != nil {
		log.Fatal("Error connecting to database:", err)
	}

	if err = db.Ping(); err != nil {
		log.Fatal("Error pinging database:", err)
	}

	log.Println("Connected to database successfully")
	return db
}

func getEnv(key, fallback string) string {
	if value := os.Getenv(key); value != "" {
		return value
	}
	return fallback
}

func healthCheck(w http.ResponseWriter, r *http.Request) {
	w.Header().Set("Content-Type", "application/json")
	json.NewEncoder(w).Encode(map[string]string{"status": "ok"})
}


func main() {
	db := initDB()
	defer db.Close()

	// Initialize repository
	occupationRepo := repository.NewOccupationRepository(db)

	// Initialize handlers
	occupationHandler := handlers.NewOccupationHandler(occupationRepo)
	searchHandler := handlers.NewSearchHandler(occupationRepo)
	createHandler := handlers.NewCreateCareersHandler(occupationRepo)

	// Setup routes
	r := mux.NewRouter()
	r.HandleFunc("/health", healthCheck).Methods("GET")
	r.HandleFunc("/search", searchHandler.Search).Methods("GET")
	r.HandleFunc("/occupations", occupationHandler.GetAll).Methods("GET")
	r.HandleFunc("/occupations", createHandler.CreateBatch).Methods("POST")
	r.HandleFunc("/occupations/{id}", occupationHandler.GetByID).Methods("GET")
	r.HandleFunc("/occupations/{id}/similar", occupationHandler.GetSimilar).Methods("GET")

	// Apply security middleware
	rateLimiter := middleware.NewRateLimiter(100) // 100 requests per minute
	handler := middleware.CORS(r)
	handler = middleware.SecurityHeaders(handler)
	handler = middleware.RequestSizeLimit(1048576)(handler) // 1MB limit
	handler = rateLimiter.Limit(handler)

	port := getEnv("PORT", "5000")
	log.Printf("Server starting on port %s", port)
	log.Fatal(http.ListenAndServe(":"+port, handler))
}
