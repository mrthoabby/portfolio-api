package health

import (
	"context"
	"net/http"
	"time"

	"go.mongodb.org/mongo-driver/mongo"

	"github.com/mrthoabby/portfolio-api/internal/common"
	"github.com/mrthoabby/portfolio-api/internal/version"
)

type Handler struct {
	db *mongo.Database
}

func NewHandler(db *mongo.Database) *Handler {
	return &Handler{db: db}
}

// HealthResponse represents the health check response
type HealthResponse struct {
	Status   string       `json:"status"`
	Database string       `json:"database"`
	Version  version.Info `json:"version"`
}

func (h *Handler) Check(w http.ResponseWriter, r *http.Request) {
	health := HealthResponse{
		Status:  "ok",
		Version: version.Get(),
	}

	// Check database connection
	ctx, cancel := context.WithTimeout(context.Background(), 2*time.Second)
	defer cancel()

	if err := h.db.Client().Ping(ctx, nil); err != nil {
		health.Status = "unhealthy"
		health.Database = "disconnected"
		common.RespondJSON(w, http.StatusServiceUnavailable, health)
		return
	}

	health.Database = "connected"
	common.RespondJSON(w, http.StatusOK, health)
}

