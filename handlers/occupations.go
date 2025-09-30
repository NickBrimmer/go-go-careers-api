package handlers

import (
	"encoding/json"
	"net/http"

	"github.com/gorilla/mux"
	"go-careers/repository"
)

type OccupationHandler struct {
	repo *repository.OccupationRepository
}

func NewOccupationHandler(repo *repository.OccupationRepository) *OccupationHandler {
	return &OccupationHandler{repo: repo}
}

func (h *OccupationHandler) GetAll(w http.ResponseWriter, r *http.Request) {
	occupations, err := h.repo.GetAll()
	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}

	w.Header().Set("Content-Type", "application/json")
	json.NewEncoder(w).Encode(occupations)
}

func (h *OccupationHandler) GetByID(w http.ResponseWriter, r *http.Request) {
	vars := mux.Vars(r)
	id := vars["id"]

	occ, err := h.repo.GetByID(id)
	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}

	if occ == nil {
		http.Error(w, "Occupation not found", http.StatusNotFound)
		return
	}

	w.Header().Set("Content-Type", "application/json")
	json.NewEncoder(w).Encode(occ)
}
