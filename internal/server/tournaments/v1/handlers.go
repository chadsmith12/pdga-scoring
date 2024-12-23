package v1

import (
	"context"
	"log/slog"
	"net/http"

	"github.com/chadsmith12/pdga-scoring/internal/repository"
	"github.com/chadsmith12/pdga-scoring/internal/server/results"
	"github.com/chadsmith12/pdga-scoring/pkgs/pulse"
	"github.com/jackc/pgx/v5/pgxpool"
)

type handlers struct {
    repo repository.Queries
    logger *slog.Logger
}

func NewHandlers(pool *pgxpool.Pool, logger *slog.Logger) *handlers {
    return &handlers{repo: *repository.New(pool), logger: logger}
}

func (h *handlers) TournamentRoutes(group *pulse.Group) {
    group.Get("/tournaments", h.listTournaments)
}

func (h *handlers) listTournaments(req *http.Request) pulse.PuleHttpWriter {
	tournaments, err := h.repo.GetAllTournaments(context.Background())
	if err != nil {
		return pulse.ErrorJson(500, err)	
	}

	return results.List(tournaments)
}
