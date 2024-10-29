package extractor

import (
	"context"
	"log/slog"
	"sync"
	"time"

	"github.com/chadsmith12/pdga-scoring/internal/repository"
	"github.com/chadsmith12/pdga-scoring/pkgs/pdga"
	"github.com/chadsmith12/pdga-scoring/pkgs/utils"
	"github.com/jackc/pgx/v5/pgtype"
	"golang.org/x/time/rate"
)

type rateLimiter struct {
	limiter *rate.Limiter
}

type TournamentExtrator struct {
	store        *repository.Queries
	client       *pdga.Client
	logger       *slog.Logger
	tournamentWg sync.WaitGroup
	limiter *rateLimiter
	workers int
}

func NewTournamentExtractor(store *repository.Queries, client *pdga.Client, logger *slog.Logger, requestLimit float64, workers int) *TournamentExtrator {
	limiter := rate.NewLimiter(rate.Limit(requestLimit), 1)
	return &TournamentExtrator{
		store:        store,
		client:       client,
		logger:       logger,
		tournamentWg: sync.WaitGroup{},
		limiter: &rateLimiter{limiter: limiter},
		workers: workers,
	}
}

func (service *TournamentExtrator) Extract(ctx context.Context, tournamentIds []int) {
	tournamentCh := make(chan int, len(tournamentIds))
	for _, tournamentId := range tournamentIds {
		tournamentCh <- tournamentId
	}
	close(tournamentCh)
	// start all the workers for processing tournaments
	for i := 0; i < service.workers; i++ {
		service.tournamentWg.Add(1)
		go service.tournamentWorker(ctx, tournamentCh)
	}

	service.tournamentWg.Wait()
}

func (service *TournamentExtrator) tournamentWorker(ctx context.Context, tournamentCh <- chan int) {
	defer service.tournamentWg.Done()
	
	for tournamentId := range tournamentCh {
		service.processTournament(ctx, tournamentId)		
	}

	slog.Info("finishing tournament worker")
}

func (service *TournamentExtrator) processTournament(ctx context.Context, id int) {
	service.limiter.limiter.Wait(ctx)
	tourneyInfo, err := service.client.FetchTournamentInfo(ctx, id)
	if err != nil {
		service.logger.Warn("failed to download touranment", slog.Int64("id", int64(id)), slog.Any("error", err))
		return
	}
	model := newCreateTournamentParams(id, tourneyInfo)
	_, err = service.store.CreateTournament(context.Background(), model)
	if err != nil {
		service.logger.Warn("failed to create tournament", slog.Int64("id", int64(id)), slog.Any("error", err))
		return
	}
	tournamentRounds := pdga.FullTournamentRound(service.extractRounds(ctx, tourneyInfo))
	service.insertPlayers(ctx, tournamentRounds.Players())
}

func (service *TournamentExtrator) extractRounds(ctx context.Context, tourneyInfo pdga.TournamentInfo) []pdga.TournamentRoundResponse {
	id, err := tourneyInfo.IdAsInt()
	if err != nil {
		service.logger.Warn("failed to get the tournament id as an integer", slog.String("id", tourneyInfo.Data.TournamentID))
		return []pdga.TournamentRoundResponse{}
	}

	numberRounds := 2 * len(tourneyInfo.Data.RoundsList)
	roundResponses := make([]pdga.TournamentRoundResponse, 0, numberRounds)
	for _, division := range tourneyInfo.Data.Divisions {
		seenLatest := false
		for round, roundInfo := range tourneyInfo.Data.RoundsList {
			if seenLatest {
				break
			}
			if round == division.LatestRound {
				seenLatest = true
			}
			service.limiter.limiter.Wait(ctx)

			roundResponse, err := service.extractRound(int(roundInfo.Number), id, pdga.Division(division.Division))
			if err != nil {
				service.logger.Warn("failed to get tournament round",
					slog.Int("id", id),
					slog.Int("roundNumber", int(roundInfo.Number)),
					slog.Any("err", err))
				continue
			}
			roundResponses = append(roundResponses, roundResponse)
		}
	}

	return roundResponses
}

func (service *TournamentExtrator) insertPlayers(ctx context.Context, players []pdga.RoundPlayer) {
	playersToInsert := utils.MapSlice(players, func(rp pdga.RoundPlayer) repository.CreateManyPlayersParams {
		return newCreatePlayer(rp)
	})

	results := service.store.CreateManyPlayers(ctx, playersToInsert)
	results.Exec(func(i int, err error) {})
//	results.Close()
}

func (service *TournamentExtrator) extractRound(roundNumber, tournamentId int, division pdga.Division) (pdga.TournamentRoundResponse, error) {
	roundData, err := service.client.FetchTournamentRound(context.Background(), tournamentId, roundNumber, division)

	return roundData, err
}

func newCreatePlayer(player pdga.RoundPlayer) repository.CreateManyPlayersParams {
	return repository.CreateManyPlayersParams{
		FirstName:  player.FirstName,
		LastName:   player.LastName,
		Name:       player.Name,
		Division:   string(player.PlayerDivison),
		PdgaNumber: player.PdgaNumber,
		City:       pgtype.Text{String: player.City, Valid: true},
		StateProv:  pgtype.Text{String: player.StateProv, Valid: true},
		Country:    pgtype.Text{String: player.Country, Valid: true},
	}
}

func newCreateTournamentParams(id int, tourneyInfo pdga.TournamentInfo) repository.CreateTournamentParams {
	return repository.CreateTournamentParams{
		ExternalID: int64(id),
		Name:       tourneyInfo.Data.Name,
		StartDate:  toDate(tourneyInfo.Data.StartDate),
		EndDate:    toDate(tourneyInfo.Data.EndDate),
		Tier:       pgtype.Text{String: tourneyInfo.Data.Tier, Valid: true},
		Location:   pgtype.Text{String: tourneyInfo.Data.Location, Valid: true},
		Country:    pgtype.Text{String: tourneyInfo.Data.Country, Valid: true},
	}
}

func toDate(value string) time.Time {
	parsedTime, _ := time.Parse(time.DateOnly, value)
	return parsedTime
}

