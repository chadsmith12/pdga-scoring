package scoring

import (
	"context"

	"github.com/chadsmith12/pdga-scoring/internal/repository"
	"github.com/chadsmith12/pdga-scoring/pkgs/pdga"
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

func (service *TournamentService) InsertTournament(ctx context.Context, tournamentInfo pdga.TournamentData, tournamentId int) (repository.Tournament, error) {
    numberRounds := int(tournamentInfo.Rounds)
    tournamentInsert := repository.CreateTournamentParams {
        Name: tournamentInfo.Name,
        ExternalID: int64(tournamentId),
        NumberOfRounds: int32(numberRounds),
    }

    tournament, err := service.repository.CreateTournament(ctx, tournamentInsert)

    return tournament, err
}
