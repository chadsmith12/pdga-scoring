package simulator

import (
	"context"
	"fmt"
	"log"
	"slices"

	"github.com/chadsmith12/pdga-scoring/internal/fantasy"
	"github.com/chadsmith12/pdga-scoring/internal/repository"
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
}

func NewSimulator(config fantasy.ScoringConfig, teams fantasy.Teams, tournaments []int64, db repository.DBTX) *Simulator {
    query := repository.New(db)
    return &Simulator{
    	scoringConfig: config,
    	teams:         teams,
    	tournaments:   tournaments,
        repo: *query,
    }
}

func (sim *Simulator) Run() {
    for _, tournamentId := range sim.tournaments {
        sim.scoreTournament(tournamentId)     
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
        currentTeams = append(currentTeams, team.Team.CreateTeam(mpoPlayers, fpoPlayers))
    }
    tournamentResults := sim.collectTournamentResults(tournamentId, currentTeams)
    for _, currentTeam := range currentTeams {
        currentScore := currentTeam.ScoreTeam(sim.scoringConfig, tournamentResults)
        fmt.Printf("For TournamentId %d a team scored: %2f\n", tournamentId, currentScore)
    }

    fmt.Println("---------------------------------")
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

    results.MpoWinner = getWinner(top10, pdga.Mpo).PlayerID
    results.FpoWinner = getWinner(top10, pdga.Fpo).PlayerID
    results.Podiums = getPodiums(top10)
    results.Top10s = getTop10(top10)
    results.HotRounds = mapHotRounds(hotRounds)
    results.RoundBirdies = getRoundScores(scores, onlyBirdies)
    results.RoundEaglesBetter = getRoundScores(scores, betterThanBirdies)
    results.RoundBogeys = getRoundScores(scores, onlyBogeys)
    results.RoundDoubleWorse = getRoundScores(scores, onlyBogeys)

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
        if result.Rank > 1 && result.Rank <= 3 {
            podiums = append(podiums, result.PlayerID)
        }
    }

    return podiums
}

func getTop10(results []repository.GetTop10ByTournamentRow) []int64 {
    top10 := make([]int64, 0, 20) 
    for _, result := range results {
        if result.Rank > 3 && result.Rank <= 10 {
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

func getRoundScores(results []repository.HoleScore, filter scoreFilter) []map[int64]int {
    groupedPlayers := groupBy(results)
    roundNumbers := getRoundNumbers(results)

    roundPlayerScores := make([]map[int64]int, len(roundNumbers))
    for i, round := range roundNumbers {
        roundPlayerScores[i] = make(map[int64]int)
        for player, scores := range groupedPlayers[round] {
            count := 0
            for _, score := range scores {
                if filter(score) {
                    count++
                }
            }
            roundPlayerScores[i][player] = count
        }
    }

    return roundPlayerScores
}

func groupBy(results []repository.HoleScore) roundPlayerGrouping {
    groupedData := make(roundPlayerGrouping)

    for _, result := range results {
        if groupedData[result.RoundNumber] == nil {
            groupedData[result.RoundNumber] = make(map[int64][]repository.HoleScore)
        }
        groupedData[result.RoundNumber][result.PlayerID] = append(groupedData[result.RoundNumber][result.PlayerID], result)
    }

    return groupedData
}

func getRoundNumbers(results []repository.HoleScore) []int32 {
    roundNumbers := make([]int32, 0, 10)
    for _, result := range results {
        if !slices.Contains(roundNumbers, result.HoleNumber) {
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
