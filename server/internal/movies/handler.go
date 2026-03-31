package movies

import (
	"encoding/json"
	"errors"
	"net/http"

	"github.com/andtkach/cinema/internal/utils"
	"github.com/google/uuid"
)

type handler struct {
	svc *Service
}

func NewHandler(svc *Service) *handler {
	return &handler{svc}
}

func (h *handler) ListMovies(w http.ResponseWriter, r *http.Request) {
	ms, err := h.svc.List()
	if err != nil {
		utils.WriteJSON(w, http.StatusInternalServerError, map[string]string{"error": "internal error"})
		return
	}
	if ms == nil {
		ms = []Movie{}
	}
	utils.WriteJSON(w, http.StatusOK, ms)
}

func (h *handler) GetMovie(w http.ResponseWriter, r *http.Request) {
	id, err := uuid.Parse(r.PathValue("movieID"))
	if err != nil {
		utils.WriteJSON(w, http.StatusBadRequest, map[string]string{"error": "invalid movie id"})
		return
	}
	m, err := h.svc.GetByID(id)
	if errors.Is(err, ErrMovieNotFound) {
		utils.WriteJSON(w, http.StatusNotFound, map[string]string{"error": "not found"})
		return
	}
	if err != nil {
		utils.WriteJSON(w, http.StatusInternalServerError, map[string]string{"error": "internal error"})
		return
	}
	utils.WriteJSON(w, http.StatusOK, m)
}

func (h *handler) CreateMovie(w http.ResponseWriter, r *http.Request) {
	var req struct {
		ID    string `json:"id"`
		Title string `json:"title"`
		Rows  int    `json:"rows"`
		Seats int    `json:"seats"`
	}
	if err := json.NewDecoder(r.Body).Decode(&req); err != nil {
		utils.WriteJSON(w, http.StatusBadRequest, map[string]string{"error": "invalid body"})
		return
	}
	if req.Title == "" || req.Rows <= 0 || req.Seats <= 0 {
		utils.WriteJSON(w, http.StatusBadRequest, map[string]string{"error": "title, rows and seats are required"})
		return
	}
	var movieID uuid.UUID
	if req.ID != "" {
		parsedID, parseErr := uuid.Parse(req.ID)
		if parseErr != nil {
			utils.WriteJSON(w, http.StatusBadRequest, map[string]string{"error": "invalid id"})
			return
		}
		movieID = parsedID
	}
	m, err := h.svc.Create(Movie{ID: movieID, Title: req.Title, Rows: req.Rows, Seats: req.Seats})
	if errors.Is(err, ErrMovieIDConflict) {
		utils.WriteJSON(w, http.StatusConflict, map[string]string{"error": "id already exists"})
		return
	}
	if err != nil {
		utils.WriteJSON(w, http.StatusInternalServerError, map[string]string{"error": "internal error"})
		return
	}
	utils.WriteJSON(w, http.StatusCreated, m)
}

func (h *handler) UpdateMovie(w http.ResponseWriter, r *http.Request) {
	id, err := uuid.Parse(r.PathValue("movieID"))
	if err != nil {
		utils.WriteJSON(w, http.StatusBadRequest, map[string]string{"error": "invalid movie id"})
		return
	}
	var req struct {
		Title string `json:"title"`
		Rows  int    `json:"rows"`
		Seats int    `json:"seats"`
	}
	if err := json.NewDecoder(r.Body).Decode(&req); err != nil {
		utils.WriteJSON(w, http.StatusBadRequest, map[string]string{"error": "invalid body"})
		return
	}
	if req.Title == "" || req.Rows <= 0 || req.Seats <= 0 {
		utils.WriteJSON(w, http.StatusBadRequest, map[string]string{"error": "title, rows and seats are required"})
		return
	}
	m, err := h.svc.Update(Movie{ID: id, Title: req.Title, Rows: req.Rows, Seats: req.Seats})
	if errors.Is(err, ErrMovieNotFound) {
		utils.WriteJSON(w, http.StatusNotFound, map[string]string{"error": "not found"})
		return
	}
	if err != nil {
		utils.WriteJSON(w, http.StatusInternalServerError, map[string]string{"error": "internal error"})
		return
	}
	utils.WriteJSON(w, http.StatusOK, m)
}

func (h *handler) DeleteMovie(w http.ResponseWriter, r *http.Request) {
	id, err := uuid.Parse(r.PathValue("movieID"))
	if err != nil {
		utils.WriteJSON(w, http.StatusBadRequest, map[string]string{"error": "invalid movie id"})
		return
	}
	err = h.svc.Delete(id)
	if errors.Is(err, ErrMovieNotFound) {
		utils.WriteJSON(w, http.StatusNotFound, map[string]string{"error": "not found"})
		return
	}
	if err != nil {
		utils.WriteJSON(w, http.StatusInternalServerError, map[string]string{"error": "internal error"})
		return
	}
	w.WriteHeader(http.StatusNoContent)
}
