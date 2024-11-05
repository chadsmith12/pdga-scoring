package scoring

import (

	"github.com/chadsmith12/pdga-scoring/internal/repository"
	"github.com/jackc/pgx/v5/pgxpool"
)

type TournamentService struct {
    repository *repository.Queries
}

func NewTournamentService(conn *pgxpool.Pool) *TournamentService {
    return &TournamentService {
        repository: repository.New(conn),
    }
}

