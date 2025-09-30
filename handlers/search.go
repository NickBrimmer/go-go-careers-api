package handlers

import (
	"encoding/json"
	"net/http"

	"go-careers/repository"
)

type SearchHandler struct {
	repo *repository.OccupationRepository
}

func NewSearchHandler(repo *repository.OccupationRepository) *SearchHandler {
	return &SearchHandler{repo: repo}
}

func (h *SearchHandler) Search(w http.ResponseWriter, r *http.Request) {
	query := r.URL.Query().Get("q")

	if query == "" {
		http.Error(w, "Missing search query parameter 'q'", http.StatusBadRequest)
		return
	}

	results, err := h.repo.Search(query)
	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}

	w.Header().Set("Content-Type", "application/json")
	json.NewEncoder(w).Encode(results)
}
