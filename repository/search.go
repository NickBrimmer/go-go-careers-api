package repository

import (
	"database/sql"
	"go-careers/models"
)

type SearchRepository struct {
	db *sql.DB
}

func NewSearchRepository(db *sql.DB) *SearchRepository {
	return &SearchRepository{db: db}
}

func (r *SearchRepository) SearchOccupations(searchTerm string) ([]models.Occupation, error) {
	query := `
		SELECT id, soc_id, soc_title, title, singular_title, description, typical_ed_level
		FROM occupations
		WHERE title LIKE ? OR soc_title LIKE ?
		LIMIT 50
	`

	searchPattern := "%" + searchTerm + "%"
	rows, err := r.db.Query(query, searchPattern, searchPattern)
	if err != nil {
		return nil, err
	}
	defer rows.Close()

	var occupations []models.Occupation
	for rows.Next() {
		var occ models.Occupation
		if err := rows.Scan(&occ.ID, &occ.SocID, &occ.SocTitle, &occ.Title, &occ.SingularTitle, &occ.Description, &occ.TypicalEdLevel); err != nil {
			return nil, err
		}
		occupations = append(occupations, occ)
	}

	return occupations, nil
}
