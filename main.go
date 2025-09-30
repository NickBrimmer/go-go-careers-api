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
	"go-careers/server/handlers"
)

var db *sql.DB

func initDB() {
	dbHost := getEnv("DB_HOST", "localhost")
	dbPort := getEnv("DB_PORT", "3306")
	dbUser := getEnv("DB_USER", "root")
	dbPassword := getEnv("DB_PASSWORD", "rootpassword")
	dbName := getEnv("DB_NAME", "go_careers")

	dsn := fmt.Sprintf("%s:%s@tcp(%s:%s)/%s?parseTime=true",
		dbUser, dbPassword, dbHost, dbPort, dbName)

	var err error
	db, err = sql.Open("mysql", dsn)
	if err != nil {
		log.Fatal("Error connecting to database:", err)
	}

	if err = db.Ping(); err != nil {
		log.Fatal("Error pinging database:", err)
	}

	log.Println("Connected to database successfully")
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
	initDB()
	defer db.Close()

	r := mux.NewRouter()

	r.HandleFunc("/health", healthCheck).Methods("GET")
	r.HandleFunc("/occupations", handlers.GetOccupations(db)).Methods("GET")
	r.HandleFunc("/occupations/{id}", handlers.GetOccupation(db)).Methods("GET")

	port := getEnv("PORT", "5000")
	log.Printf("Server starting on port %s", port)
	log.Fatal(http.ListenAndServe(":"+port, r))
}
