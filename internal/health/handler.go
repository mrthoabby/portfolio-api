package health

import (
	"context"
	"net/http"
	"time"

	"github.com/mrthoabby/portfolio-api/internal/common"
	"github.com/mrthoabby/portfolio-api/internal/common/contracts"
	"github.com/mrthoabby/portfolio-api/internal/version"
)

type Handler struct {
	dataSource contracts.DataSource
}

func NewHandler(dataSource contracts.DataSource) *Handler {
	return &Handler{dataSource: dataSource}
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

	// Check data source connection
	ctx, cancel := context.WithTimeout(context.Background(), 2*time.Second)
	defer cancel()

	if err := h.dataSource.Ping(ctx); err != nil {
		health.Status = "unhealthy"
		health.Database = "disconnected"
		common.RespondJSON(w, http.StatusServiceUnavailable, health)
		return
	}

	health.Database = "connected"
	common.RespondJSON(w, http.StatusOK, health)
}
