package extractor

import (
	"context"
	"log/slog"
	"sync"
	"time"

	"github.com/chadsmith12/pdga-scoring/internal/database"
	"github.com/chadsmith12/pdga-scoring/internal/repository"
	"github.com/chadsmith12/pdga-scoring/pkgs/pdga"
	"github.com/chadsmith12/pdga-scoring/pkgs/utils"
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
	dbTourney, err := service.store.CreateTournament(context.Background(), model)
	if err != nil {
		service.logger.Warn("failed to create tournament", slog.Int64("id", int64(id)), slog.Any("error", err))
		return
	}
	tournamentRounds := pdga.FullTournamentRound(service.extractRounds(ctx, tourneyInfo))
	service.processLayouts(ctx, dbTourney.ExternalID, tourneyInfo)
	service.insertPlayers(ctx, tournamentRounds.Players())
	service.insertRoundScores(ctx, dbTourney.ExternalID, tournamentRounds)
}

func (service *TournamentExtrator) insertRoundScores(ctx context.Context, tournamentId int64, fullRound pdga.FullTournamentRound) {
	scores := make([]repository.CreateRoundScoresParams, 0, 18)
	for _, poolRound := range fullRound {
		for _, round := range poolRound.Data.RoundData {
			for _, score := range round.Scores {
				roundScore := repository.CreateRoundScoresParams{
					PlayerID:     score.PDGANum,
					TournamentID: tournamentId,
					LayoutID:     score.LayoutID,
					RoundNumber:  int32(score.Round),
					Score:        int32(score.RoundtoPar),
				}
				scores = append(scores, roundScore)
			}	
		}
	}

	results := service.store.CreateRoundScores(ctx, scores)
	results.Exec(func(i int, err error) {
		if err != nil {
			service.logger.Warn("failed to insert score", slog.Int("index", i), slog.Any("err", err))
		}
	})
}

func (service *TournamentExtrator) processLayouts(ctx context.Context, externalId int64, tourneyInfo pdga.TournamentInfo) {
	layouts := tourneyInfo.Data.Layouts
	dbLayouts := make([]repository.CreateManyLayoutsParams, 0, len(layouts))

	for _, layout := range layouts {
		if layout.CourseID == -1 {
			continue
		}
		dbLayout := repository.CreateManyLayoutsParams{
			ID:           layout.LayoutID,
			TournamentID: externalId,
			Name:         layout.Name,
			CourseName:   layout.CourseName,
			Length:       database.IntToPgInt(int(layout.Length)),
			Units:        database.StringToPgText(layout.Units),
			Holes:        database.IntToPgInt(int(layout.Holes)),
			Par:          database.IntToPgInt(int(layout.Par)),
		}
		dbLayouts = append(dbLayouts, dbLayout)
	}
	results := service.store.CreateManyLayouts(ctx, dbLayouts)
	results.Exec(func(i int, err error) {
		if err != nil {
			service.logger.Warn("failed to insert layout", slog.Int("index", i), slog.Any("err", err))
			return
		}
	})
	results.Close()
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
		City: database.StringToPgText(player.City),
		StateProv: database.StringToPgText(player.StateProv),
		Country: database.StringToPgText(player.Country),
	}
}

func newCreateTournamentParams(id int, tourneyInfo pdga.TournamentInfo) repository.CreateTournamentParams {
	return repository.CreateTournamentParams{
		ExternalID: int64(id),
		Name:       tourneyInfo.Data.Name,
		StartDate:  toDate(tourneyInfo.Data.StartDate),
		EndDate:    toDate(tourneyInfo.Data.EndDate),
		Tier:       database.StringToPgText(tourneyInfo.Data.Tier),
		Location:   database.StringToPgText(tourneyInfo.Data.Location),
		Country:    database.StringToPgText(tourneyInfo.Data.Country),
	}
}

func toDate(value string) time.Time {
	parsedTime, _ := time.Parse(time.DateOnly, value)
	return parsedTime
}

