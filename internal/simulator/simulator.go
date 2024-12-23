package simulator

import (
	"context"
	"encoding/csv"
	"fmt"
	"log"
	"os"
	"path"
	"slices"
	"sync"

	"github.com/chadsmith12/pdga-scoring/internal/database"
	"github.com/chadsmith12/pdga-scoring/internal/repository"
	"github.com/chadsmith12/pdga-scoring/pkgs/fantasy"
	"github.com/chadsmith12/pdga-scoring/pkgs/pdga"
)

type scoreFilter func(score repository.HoleScore) bool
type roundPlayerGrouping map[int32]map[int64][]repository.HoleScore

var onlyBirdies scoreFilter = func(score repository.HoleScore) bool {
    return score.ScoreRelativeToPar == -1
}

var betterThanBirdies scoreFilter = func(score repository.HoleScore) bool {
    return score.ScoreRelativeToPar <= -2
}

var onlyBogeys scoreFilter = func(score repository.HoleScore) bool {
    return score.ScoreRelativeToPar == 1
}

var worseThanBogey scoreFilter = func(score repository.HoleScore) bool {
    return score.ScoreRelativeToPar >= 2
}

type Simulator struct {
    scoringConfig fantasy.ScoringConfig
    teams fantasy.Teams
    tournaments []int64
    repo repository.Queries
    teamResults map[string]float64
    tournamentScoring map[string][]fantasy.TournamentScoring
    name string
}

func NewSimulator(config fantasy.ScoringConfig, teams fantasy.Teams, tournaments []int64, db repository.DBTX, name string) *Simulator {
    query := repository.New(db)
    teamResults := make(map[string]float64)
    for _, team := range teams {
        fmt.Printf("Team Name is %s\n", team.Name)
        teamResults[team.Name] = 0
    }
    return &Simulator{
    	scoringConfig: config,
    	teams:         teams,
    	tournaments:   tournaments,
        repo: *query,
        teamResults: teamResults,
        tournamentScoring: make(map[string][]fantasy.TournamentScoring),
        name: name,
    }
}

func (sim *Simulator) Run() {
    for _, tournamentId := range sim.tournaments {
        sim.scoreTournament(tournamentId)     
    }

    fmt.Printf("After Simulation: \n")
    for team, result := range sim.teamResults {
        fmt.Printf("%s: %2f\n", team, result)
    }
}

func (sim *Simulator) ExportResults() {
    csvPath := path.Join("results/", sim.name)
    err := os.MkdirAll(csvPath, 0700)
    if err != nil {
        log.Printf("Failed to create the results directory: %v\n", err)
    }
    for name, scoring := range sim.tournamentScoring {
        csvFile, err := os.Create(fmt.Sprintf("%s/%s-%s-results.csv", csvPath, sim.name, name))
        if err != nil {
            log.Printf("Failed to create csv file for export: %v\n", err)
        }
        csvWriter := csv.NewWriter(csvFile)
        csvWriter.Write(fantasy.ScoringHeaders)
        defer csvWriter.Flush()
        for _, score := range scoring {
            if err := csvWriter.Write(score.Strings()); err != nil {
                fmt.Printf("failed to write row: %v\n", err)
            }
        }
    }
}

func (sim *Simulator) scoreTournament(tournamentId int64) {
    tournamentPlayers, err := sim.repo.GetPlayersInTournament(context.Background(), tournamentId)
    if err != nil {
        log.Fatal(err)
    }
    mpoPlayers, fpoPlayers := partitionPlayers(tournamentPlayers)
    currentTeams := make([]fantasy.CurrentTeam, 0, 2)
    for _, team := range sim.teams {
        currentTeams = append(currentTeams, team.CreateTeam(mpoPlayers, fpoPlayers))
    }

    wg := sync.WaitGroup{}
    for _, player := range tournamentPlayers {
        wg.Add(1)
        go func() {
            defer wg.Done()
            team := fantasy.SingleTeam(player.PlayerID, pdga.Division(player.Division))
            playerResults := sim.collectTournamentResults(tournamentId, []fantasy.CurrentTeam{team})
            for _, roundResult := range playerResults.RoundResults {
                playerBatch := sim.repo.InsertFantasyRoundScores(context.Background(), []repository.InsertFantasyRoundScoresParams{
                    {
                        PlayerID: database.BigIntToPgInt8(player.PlayerID),
                        RoundNumber: database.IntToPgInt(int(roundResult.RoundNumber())),
                        TournamentID: database.BigIntToPgInt8(tournamentId),
                        Birdies: database.IntToPgInt(roundResult.BirdiesForPlayer(player.PlayerID)),
                        EaglesOrBetter: database.IntToPgInt(roundResult.EaglesBetterForPlayer(player.PlayerID)),
                        Bogeys: database.IntToPgInt(roundResult.BogeysForPlayer(player.PlayerID)),
                    },
                })
                playerBatch.Exec(func(i int, err error) {
                    if err != nil {
                        fmt.Printf("error inserting %d of fantasy round score: %v\n", i, err)
                    }
                })
            }
            batch := sim.repo.InsertFantasyTournamentScores(context.Background(), []repository.InsertFantasyTournamentScoresParams{
                {
                    PlayerID: database.BigIntToPgInt8(player.PlayerID),
                    TournamentID: database.BigIntToPgInt8(tournamentId),
                    WonEvent: database.BoolToPgBool(team.HasWinner(playerResults)),
                    PodiumFinish: database.BoolToPgBool(team.NumberPodiums(playerResults) > 0),
                    Top10Finish: database.BoolToPgBool(team.NumberTop10s(playerResults) > 0),
                    HotRounds: database.IntToPgInt(playerResults.PlayerHotRounds(player.PlayerID)),
                },
            }) 
            batch.Exec(func(i int, err error) {
                if err != nil {
                    fmt.Printf("error inserting index %d of fantasy scores: %v\n", i, err)
                }
            })
        }()
    }

    tournamentResults := sim.collectTournamentResults(tournamentId, currentTeams)
    for _, currentTeam := range currentTeams {
        tournamentScores := currentTeam.ScoreTournament(sim.scoringConfig, tournamentResults)
        sim.tournamentScoring[currentTeam.Name] = append(sim.tournamentScoring[currentTeam.Name], tournamentScores)
        sim.teamResults[currentTeam.Name] += tournamentScores.TotalScore
    }
}


func (sim *Simulator) collectTournamentResults(tournamentId int64, currentTeams []fantasy.CurrentTeam) fantasy.Results {
    results := fantasy.Results{}
    top10, err := sim.repo.GetTop10ByTournament(context.Background(), tournamentId)
    if err != nil {
        log.Printf("failed to get top 10 for %d tournament. Err = %v\n", tournamentId, err)
    }
    hotRounds, err := sim.repo.GetHotRoundsForTournament(context.Background(), tournamentId)
    if err != nil {
        log.Printf("failed to get hot rounds for %d tournament, Err = %v\n", tournamentId, err)
    }
    teamPlayers := make([]int64, 0, 10)
    for _, team := range currentTeams {
        teamPlayers = append(teamPlayers, team.Players...)
    }
    scores, err := sim.repo.GetPlayersHoleScores(context.Background(), repository.GetPlayersHoleScoresParams{
    	TournamentID: tournamentId,
    	PlayerIds:    teamPlayers,
    })

    roundNumbers := getRoundNumbers(scores)

    results.MpoWinner = getWinner(top10, pdga.Mpo).PlayerID
    results.FpoWinner = getWinner(top10, pdga.Fpo).PlayerID
    results.Podiums = getPodiums(top10)
    results.Top10s = getTop10(top10)
    results.HotRounds = mapHotRounds(hotRounds)

    for _, roundNumber := range roundNumbers {
        birdies := getRoundScoresByFilter(scores, roundNumber, onlyBirdies)        
        eagles := getRoundScoresByFilter(scores, roundNumber, betterThanBirdies)
        bogeys := getRoundScoresByFilter(scores, roundNumber, onlyBogeys)
        doubles := getRoundScoresByFilter(scores, roundNumber, worseThanBogey)

        roundResult := fantasy.NewRoundResult(roundNumber, birdies, eagles, bogeys, doubles)
        results.RoundResults = append(results.RoundResults, roundResult)
    }

    return results
}

func getWinner(results []repository.GetTop10ByTournamentRow, division pdga.Division) repository.GetTop10ByTournamentRow {
    for _, result := range results {
        if result.Division == string(division) && result.Rank == 1 {
            return result
        }
    }

    return repository.GetTop10ByTournamentRow{}
}

func getPodiums(results []repository.GetTop10ByTournamentRow) []int64 {
    podiums := make([]int64, 0, 6)
    for _, result := range results {
        if result.Rank <= 3 {
            podiums = append(podiums, result.PlayerID)
        }
    }

    return podiums
}

func getTop10(results []repository.GetTop10ByTournamentRow) []int64 {
    top10 := make([]int64, 0, 20) 
    for _, result := range results {
        if result.Rank <= 10 {
            top10 = append(top10, result.PlayerID)
        }
    }

    return top10
}

func mapHotRounds(results []repository.GetHotRoundsForTournamentRow) map[int][]int64 {
    hotRounds := make(map[int][]int64, 5)

    for _, result := range results {
        hotRounds[int(result.RoundNumber)] = append(hotRounds[int(result.RoundNumber)], result.PlayerID)
    }

    return hotRounds
}

func getRoundScoresByFilter(scores []repository.HoleScore, roundNumber int32, filter scoreFilter) map[int64]int {
    roundPlayerScores := make(map[int64]int)
    for _, score := range scores {
        if score.RoundNumber == roundNumber && filter(score) {
            roundPlayerScores[score.PlayerID]++
        }
    }

    return roundPlayerScores
}

func getRoundNumbers(results []repository.HoleScore) []int32 {
    roundNumbers := make([]int32, 0, 10)
    for _, result := range results {
        if !slices.Contains(roundNumbers, result.RoundNumber) {
            roundNumbers = append(roundNumbers, result.RoundNumber)
        }
    }

    return roundNumbers
}

func partitionPlayers(players []repository.GetPlayersInTournamentRow) ([]int64, []int64) {
    mpoPlayers := make([]int64, 0, 100)
    fpoPlayers := make([]int64, 0, 50)

    for _, player := range players {
        if player.Division == pdga.Mpo.String() {
            mpoPlayers = append(mpoPlayers, player.PlayerID)
        } else if player.Division == pdga.Fpo.String() {
            fpoPlayers = append(fpoPlayers, player.PlayerID)
        }
    }

    return mpoPlayers, fpoPlayers
}
