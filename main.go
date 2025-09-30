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
)

var db *sql.DB

type Occupation struct {
	ID             string `json:"id"`
	SocID          string `json:"soc_id"`
	SocTitle       string `json:"soc_title"`
	Title          string `json:"title"`
	SingularTitle  string `json:"singular_title"`
	Description    string `json:"description"`
	TypicalEdLevel string `json:"typical_ed_level"`
}

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

func getOccupations(w http.ResponseWriter, r *http.Request) {
	rows, err := db.Query("SELECT id, soc_id, soc_title, title, singular_title, description, typical_ed_level FROM occupations LIMIT 10")
	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}
	defer rows.Close()

	var occupations []Occupation
	for rows.Next() {
		var occ Occupation
		if err := rows.Scan(&occ.ID, &occ.SocID, &occ.SocTitle, &occ.Title, &occ.SingularTitle, &occ.Description, &occ.TypicalEdLevel); err != nil {
			http.Error(w, err.Error(), http.StatusInternalServerError)
			return
		}
		occupations = append(occupations, occ)
	}

	w.Header().Set("Content-Type", "application/json")
	json.NewEncoder(w).Encode(occupations)
}

func getOccupation(w http.ResponseWriter, r *http.Request) {
	vars := mux.Vars(r)
	id := vars["id"]

	var occ Occupation
	err := db.QueryRow("SELECT id, soc_id, soc_title, title, singular_title, description, typical_ed_level FROM occupations WHERE id = ?", id).
		Scan(&occ.ID, &occ.SocID, &occ.SocTitle, &occ.Title, &occ.SingularTitle, &occ.Description, &occ.TypicalEdLevel)

	if err == sql.ErrNoRows {
		http.Error(w, "Occupation not found", http.StatusNotFound)
		return
	} else if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}

	w.Header().Set("Content-Type", "application/json")
	json.NewEncoder(w).Encode(occ)
}

func main() {
	initDB()
	defer db.Close()

	r := mux.NewRouter()

	r.HandleFunc("/health", healthCheck).Methods("GET")
	r.HandleFunc("/occupations", getOccupations).Methods("GET")
	r.HandleFunc("/occupations/{id}", getOccupation).Methods("GET")

	port := getEnv("PORT", "5000")
	log.Printf("Server starting on port %s", port)
	log.Fatal(http.ListenAndServe(":"+port, r))
}
