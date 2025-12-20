package projects

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

func (h *Handler) GetByProfileID(w http.ResponseWriter, r *http.Request) {
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

	projects, err := h.service.GetByProfileID(r.Context(), profileID)
	if err != nil {
		common.RespondError(w, http.StatusNotFound, "NOT_FOUND", err.Error(), nil)
		return
	}

	response := Response{Projects: projects}
	common.RespondJSON(w, http.StatusOK, response)
}

