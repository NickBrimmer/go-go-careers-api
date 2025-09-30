package handlers

import (
	"database/sql"
	"encoding/json"
	"net/http"
)

type Occupation struct {
	ID             string `json:"id"`
	SocID          string `json:"soc_id"`
	SocTitle       string `json:"soc_title"`
	Title          string `json:"title"`
	SingularTitle  string `json:"singular_title"`
	Description    string `json:"description"`
	TypicalEdLevel string `json:"typical_ed_level"`
}

func GetOccupations(db *sql.DB) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
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
}
