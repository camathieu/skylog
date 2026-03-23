package handlers

import (
	"encoding/json"
	"net/http"
	"strconv"
	"strings"

	"github.com/camathieu/skylog/server/models"
	"github.com/go-chi/chi/v5"
	"gorm.io/gorm"
)

type JumpHandler struct {
	DB *gorm.DB
}

func NewJumpHandler(db *gorm.DB) *JumpHandler {
	return &JumpHandler{DB: db}
}

// respondJSON writes a JSON response with the given status code.
func respondJSON(w http.ResponseWriter, status int, v any) {
	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(status)
	_ = json.NewEncoder(w).Encode(v)
}

// respondError writes a JSON error response.
func respondError(w http.ResponseWriter, status int, msg string) {
	respondJSON(w, status, map[string]string{"error": msg})
}

// ListJumps handles GET /api/jumps
func (h *JumpHandler) ListJumps(w http.ResponseWriter, r *http.Request) {
	page, _ := strconv.Atoi(r.URL.Query().Get("page"))
	if page < 1 {
		page = 1
	}
	perPage, _ := strconv.Atoi(r.URL.Query().Get("per_page"))
	if perPage < 1 || perPage > 100 {
		perPage = 50
	}
	order := r.URL.Query().Get("order")
	if order != "asc" {
		order = "desc"
	}

	var jumps []models.Jump
	var total int64

	query := h.DB.Model(&models.Jump{})
	query.Count(&total)
	result := query.Order("jump_number " + order).
		Offset((page - 1) * perPage).
		Limit(perPage).
		Find(&jumps)

	if result.Error != nil {
		respondError(w, http.StatusInternalServerError, "failed to list jumps")
		return
	}

	respondJSON(w, http.StatusOK, map[string]any{
		"jumps":   jumps,
		"total":   total,
		"page":    page,
		"perPage": perPage,
	})
}

// GetJump handles GET /api/jumps/{id}
func (h *JumpHandler) GetJump(w http.ResponseWriter, r *http.Request) {
	id, err := strconv.Atoi(chi.URLParam(r, "id"))
	if err != nil {
		respondError(w, http.StatusBadRequest, "invalid id")
		return
	}

	var jump models.Jump
	if err := h.DB.First(&jump, id).Error; err != nil {
		if err == gorm.ErrRecordNotFound {
			respondError(w, http.StatusNotFound, "jump not found")
			return
		}
		respondError(w, http.StatusInternalServerError, "failed to get jump")
		return
	}

	respondJSON(w, http.StatusOK, jump)
}

// CreateJump handles POST /api/jumps
func (h *JumpHandler) CreateJump(w http.ResponseWriter, r *http.Request) {
	var jump models.Jump
	if err := json.NewDecoder(r.Body).Decode(&jump); err != nil {
		respondError(w, http.StatusBadRequest, "invalid request body")
		return
	}
	jump.ID = 0 // ensure auto-assign

	if err := h.DB.Create(&jump).Error; err != nil {
		if strings.Contains(err.Error(), "UNIQUE") {
			respondError(w, http.StatusConflict, "A jump with this number already exists")
			return
		}
		respondError(w, http.StatusInternalServerError, "failed to create jump")
		return
	}

	respondJSON(w, http.StatusCreated, jump)
}

// UpdateJump handles PUT /api/jumps/{id}
func (h *JumpHandler) UpdateJump(w http.ResponseWriter, r *http.Request) {
	id, err := strconv.Atoi(chi.URLParam(r, "id"))
	if err != nil {
		respondError(w, http.StatusBadRequest, "invalid id")
		return
	}

	var existing models.Jump
	if err := h.DB.First(&existing, id).Error; err != nil {
		if err == gorm.ErrRecordNotFound {
			respondError(w, http.StatusNotFound, "jump not found")
			return
		}
		respondError(w, http.StatusInternalServerError, "failed to find jump")
		return
	}

	var update models.Jump
	if err := json.NewDecoder(r.Body).Decode(&update); err != nil {
		respondError(w, http.StatusBadRequest, "invalid request body")
		return
	}
	update.ID = existing.ID
	update.CreatedAt = existing.CreatedAt

	if err := h.DB.Save(&update).Error; err != nil {
		if strings.Contains(err.Error(), "UNIQUE") {
			respondError(w, http.StatusConflict, "A jump with this number already exists")
			return
		}
		respondError(w, http.StatusInternalServerError, "failed to update jump")
		return
	}

	respondJSON(w, http.StatusOK, update)
}

// DeleteJump handles DELETE /api/jumps/{id}
func (h *JumpHandler) DeleteJump(w http.ResponseWriter, r *http.Request) {
	id, err := strconv.Atoi(chi.URLParam(r, "id"))
	if err != nil {
		respondError(w, http.StatusBadRequest, "invalid id")
		return
	}

	result := h.DB.Delete(&models.Jump{}, id)
	if result.Error != nil {
		respondError(w, http.StatusInternalServerError, "failed to delete jump")
		return
	}
	if result.RowsAffected == 0 {
		respondError(w, http.StatusNotFound, "jump not found")
		return
	}

	w.WriteHeader(http.StatusNoContent)
}
