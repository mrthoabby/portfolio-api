package contacts

import (
	"encoding/json"
	"net/http"

	"github.com/go-playground/validator/v10"

	"github.com/mrthoabby/portfolio-api/internal/common"
)

type Handler struct {
	service   *Service
	validator *validator.Validate
}

func NewHandler(service *Service) *Handler {
	return &Handler{
		service:   service,
		validator: validator.New(),
	}
}

func (h *Handler) Create(w http.ResponseWriter, r *http.Request) {
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

	var contactReq Request
	if err := json.NewDecoder(r.Body).Decode(&contactReq); err != nil {
		// Check if body was too large
		if err.Error() == "http: request body too large" {
			common.RespondError(w, http.StatusRequestEntityTooLarge, "PAYLOAD_TOO_LARGE", "Request body too large", nil)
			return
		}
		common.RespondError(w, http.StatusBadRequest, "BAD_REQUEST", "Invalid request body", nil)
		return
	}

	// Sanitize inputs before validation
	sanitized := common.SanitizeContactInput(contactReq.Name, contactReq.Email, contactReq.Message)
	contactReq.Name = sanitized.Name
	contactReq.Email = sanitized.Email
	contactReq.Message = sanitized.Message

	if err := h.validator.Struct(&contactReq); err != nil {
		common.RespondError(w, http.StatusBadRequest, "VALIDATION_ERROR", "Validation failed", err.Error())
		return
	}

	contact, err := h.service.Create(r.Context(), profileID, &contactReq)
	if err != nil {
		common.RespondError(w, http.StatusNotFound, "NOT_FOUND", err.Error(), nil)
		return
	}

	response := Response{
		ID:          contact.ID,
		Message:     "Contact message sent successfully",
		ContactedAt: contact.ContactedAt,
	}

	common.RespondJSON(w, http.StatusCreated, response)
}

