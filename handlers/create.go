package handlers

import (
	"encoding/json"
	"fmt"
	"net/http"

	"go-careers/models"
	"go-careers/repository"
)

type CreateCareersHandler struct {
	repo *repository.OccupationRepository
}

func NewCreateCareersHandler(repo *repository.OccupationRepository) *CreateCareersHandler {
	return &CreateCareersHandler{repo: repo}
}

func (h *CreateCareersHandler) CreateBatch(w http.ResponseWriter, r *http.Request) {
	var occupations []models.Occupation

	if err := json.NewDecoder(r.Body).Decode(&occupations); err != nil {
		http.Error(w, fmt.Sprintf("Invalid JSON: %s", err.Error()), http.StatusBadRequest)
		return
	}

	if len(occupations) == 0 {
		http.Error(w, "Empty array: at least one occupation is required", http.StatusBadRequest)
		return
	}

	// Validate each occupation
	for i, occ := range occupations {
		if err := occ.Validate(); err != nil {
			http.Error(w, fmt.Sprintf("Validation error at index %d: %s", i, err.Error()), http.StatusBadRequest)
			return
		}
	}

	// Insert batch
	if err := h.repo.CreateBatch(occupations); err != nil {
		http.Error(w, fmt.Sprintf("Database error: %s", err.Error()), http.StatusInternalServerError)
		return
	}

	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(http.StatusCreated)
	json.NewEncoder(w).Encode(map[string]interface{}{
		"message": "Successfully created occupations",
		"count":   len(occupations),
	})
}
