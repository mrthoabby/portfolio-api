package questions

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

	var questionReq Request
	if err := json.NewDecoder(r.Body).Decode(&questionReq); err != nil {
		if err.Error() == "http: request body too large" {
			common.RespondError(w, http.StatusRequestEntityTooLarge, "PAYLOAD_TOO_LARGE", "Request body too large", nil)
			return
		}
		common.RespondError(w, http.StatusBadRequest, "BAD_REQUEST", "Invalid request body", nil)
		return
	}

	// Sanitize message input
	questionReq.Message = common.SanitizeString(common.StripHTMLTags(questionReq.Message))

	if err := h.validator.Struct(&questionReq); err != nil {
		common.RespondError(w, http.StatusBadRequest, "VALIDATION_ERROR", "Validation failed", err.Error())
		return
	}

	// Get client IP from context (set by middleware)
	clientIP := common.ClientIPFromContext(r.Context())

	question, err := h.service.Create(r.Context(), profileID, questionReq.Message, clientIP)
	if err != nil {
		common.RespondError(w, http.StatusNotFound, "NOT_FOUND", err.Error(), nil)
		return
	}

	response := Response{
		ID:        question.ID,
		Message:   "Question received successfully",
		CreatedAt: question.CreatedAt,
	}

	common.RespondJSON(w, http.StatusCreated, response)
}
