package handlers

import (
	"database/sql"
	"encoding/json"
	"net/http"

	"github.com/gorilla/mux"
)

func GetOccupation(db *sql.DB) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
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
}
