package profile

import (
	"net/http"

	"github.com/mrthoabby/portfolio-api/internal/common"
)

type Handler struct {
	service *Service
}

func NewHandler(service *Service) *Handler {
	return &Handler{service: service}
}

func (h *Handler) GetByID(w http.ResponseWriter, r *http.Request) {
	profileID := r.PathValue("id")
	if profileID == "" {
		common.RespondError(w, http.StatusBadRequest, "BAD_REQUEST", "Profile ID is required", nil)
		return
	}

	// Validate UUID format
	if !common.IsValidUUID(profileID) {
		common.RespondError(w, http.StatusBadRequest, "BAD_REQUEST", "Invalid profile ID format", nil)
		return
	}

	profile, err := h.service.GetByID(r.Context(), profileID)
	if err != nil {
		common.RespondError(w, http.StatusNotFound, "NOT_FOUND", "Profile not found", nil)
		return
	}

	common.RespondJSON(w, http.StatusOK, profile)
}
