package repository

import (
	"database/sql"
	"encoding/json"
	"fmt"
	"time"

	"go-careers/cache"
	"go-careers/models"
)

type OccupationRepository struct {
	db    *sql.DB
	cache *cache.RedisCache
}

func NewOccupationRepository(db *sql.DB, redisCache *cache.RedisCache) *OccupationRepository {
	return &OccupationRepository{
		db:    db,
		cache: redisCache,
	}
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
	// Try cache first
	cacheKey := fmt.Sprintf("occupation:%s", id)
	var occ models.Occupation
	if r.cache != nil {
		if err := r.cache.Get(cacheKey, &occ); err == nil {
			return &occ, nil
		}
	}

	// Cache miss - query database
	query := "SELECT id, soc_id, soc_title, title, singular_title, description, typical_ed_level FROM occupations WHERE id = ?"
	err := r.db.QueryRow(query, id).Scan(&occ.ID, &occ.SocID, &occ.SocTitle, &occ.Title, &occ.SingularTitle, &occ.Description, &occ.TypicalEdLevel)

	if err == sql.ErrNoRows {
		return nil, nil
	} else if err != nil {
		return nil, err
	}

	// Store in cache (1 hour TTL)
	if r.cache != nil {
		r.cache.Set(cacheKey, occ, time.Hour)
	}

	return &occ, nil
}

func (r *OccupationRepository) Search(searchTerm string) ([]models.Occupation, error) {
	// Try cache first
	cacheKey := fmt.Sprintf("search:%s", searchTerm)
	var occupations []models.Occupation
	if r.cache != nil {
		if err := r.cache.Get(cacheKey, &occupations); err == nil {
			return occupations, nil
		}
	}

	// Cache miss - query database
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

	occupations = []models.Occupation{}
	for rows.Next() {
		var occ models.Occupation
		if err := rows.Scan(&occ.ID, &occ.SocID, &occ.SocTitle, &occ.Title, &occ.SingularTitle, &occ.Description, &occ.TypicalEdLevel); err != nil {
			return nil, err
		}
		occupations = append(occupations, occ)
	}

	// Store in cache (15 minutes TTL for searches)
	if r.cache != nil {
		r.cache.Set(cacheKey, occupations, 15*time.Minute)
	}

	return occupations, nil
}

func (r *OccupationRepository) CreateBatch(occupations []models.Occupation) error {
	tx, err := r.db.Begin()
	if err != nil {
		return err
	}
	defer tx.Rollback()

	stmt, err := tx.Prepare("INSERT INTO occupations (id, soc_id, soc_title, title, singular_title, description, typical_ed_level) VALUES (?, ?, ?, ?, ?, ?, ?)")
	if err != nil {
		return err
	}
	defer stmt.Close()

	for _, occ := range occupations {
		_, err := stmt.Exec(occ.ID, occ.SocID, occ.SocTitle, occ.Title, occ.SingularTitle, occ.Description, occ.TypicalEdLevel)
		if err != nil {
			return err
		}
	}

	return tx.Commit()
}

func (r *OccupationRepository) GetSimilar(id string) ([]models.Occupation, error) {
	// Try cache first
	cacheKey := fmt.Sprintf("similar:%s", id)
	var occupations []models.Occupation
	if r.cache != nil {
		if err := r.cache.Get(cacheKey, &occupations); err == nil {
			return occupations, nil
		}
	}

	// Cache miss - query database
	// First, get the data JSON for the occupation
	var dataJSON string
	err := r.db.QueryRow("SELECT data FROM occupations WHERE id = ?", id).Scan(&dataJSON)
	if err == sql.ErrNoRows {
		return []models.Occupation{}, nil
	} else if err != nil {
		return nil, err
	}

	// Parse JSON to extract similarOccs array
	var data struct {
		SimilarOccs []string `json:"similarOccs"`
	}
	if err := json.Unmarshal([]byte(dataJSON), &data); err != nil {
		return nil, err
	}

	if len(data.SimilarOccs) == 0 {
		return []models.Occupation{}, nil
	}

	// Build query with placeholders for each similar occupation ID
	query := "SELECT id, soc_id, soc_title, title, singular_title, description, typical_ed_level FROM occupations WHERE id IN ("
	args := make([]interface{}, len(data.SimilarOccs))
	for i, similarID := range data.SimilarOccs {
		if i > 0 {
			query += ","
		}
		query += "?"
		args[i] = similarID
	}
	query += ")"

	rows, err := r.db.Query(query, args...)
	if err != nil {
		return nil, err
	}
	defer rows.Close()

	// Initialize as empty slice so it returns [] instead of null when empty
	occupations = []models.Occupation{}
	for rows.Next() {
		var occ models.Occupation
		if err := rows.Scan(&occ.ID, &occ.SocID, &occ.SocTitle, &occ.Title, &occ.SingularTitle, &occ.Description, &occ.TypicalEdLevel); err != nil {
			return nil, err
		}
		occupations = append(occupations, occ)
	}

	// Store in cache (1 hour TTL)
	if r.cache != nil {
		r.cache.Set(cacheKey, occupations, time.Hour)
	}

	return occupations, nil
}
