package repository

import (
	"database/sql"
	"go-careers/models"
)

type OccupationRepository struct {
	db *sql.DB
}

func NewOccupationRepository(db *sql.DB) *OccupationRepository {
	return &OccupationRepository{db: db}
}

func (r *OccupationRepository) GetAll() ([]models.Occupation, error) {
	query := "SELECT id, soc_id, soc_title, title, singular_title, description, typical_ed_level FROM occupations LIMIT 10"
	rows, err := r.db.Query(query)
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

func (r *OccupationRepository) GetByID(id string) (*models.Occupation, error) {
	query := "SELECT id, soc_id, soc_title, title, singular_title, description, typical_ed_level FROM occupations WHERE id = ?"

	var occ models.Occupation
	err := r.db.QueryRow(query, id).Scan(&occ.ID, &occ.SocID, &occ.SocTitle, &occ.Title, &occ.SingularTitle, &occ.Description, &occ.TypicalEdLevel)

	if err == sql.ErrNoRows {
		return nil, nil
	} else if err != nil {
		return nil, err
	}

	return &occ, nil
}
