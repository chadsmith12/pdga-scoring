package fantasy

import (
	"encoding/json"
	"io"
	"slices"
	"strconv"
)

var ScoringHeaders = []string{"Event Winner", "Podiums", "Top 10s", "Hot Rounds", "Birdies", "Eagles Or Better", "Bogeys", "Double or Worse", "Total Score"}

type ScoreExtractor func(round RoundResult, playerId int64) int
var Birdes ScoreExtractor = func(round RoundResult, playerId int64) int {
    return round.BirdiesForPlayer(playerId)
}
var EaglesOrBetter ScoreExtractor = func(round RoundResult, playerId int64) int {
    return round.EaglesBetterForPlayer(playerId)
}
var Bogeys ScoreExtractor = func(round RoundResult, playerId int64) int {
    return round.BogeysForPlayer(playerId)
}
var DoubleBogeysOrWorse ScoreExtractor = func(round RoundResult, playerId int64) int {
    return round.DoublesWorseForPlayer(playerId)
}

type ScoringConfig struct {
    EventWinner float64 `json:"EventWinner"`
    Podiums float64 `json:"Podiums"`
    Top10s float64 `json:"Top10s"`
    HotRound float64 `json:"HotRound"`
    RoundBirdies TimesConfig `json:"RoundBirdies"`
    EaglesOrBetter TimesConfig `json:"EaglesOrBetter"`
    Bogeys TimesConfig `json:"Bogeys"`
    DoubleOrWorse TimesConfig `json:"DoubleOrWorse"`
}

type TournamentScoring struct {
    EventWinner float64 
    Podiums float64
    Top10s float64
    HotRound float64
    RoundBirdies float64 
    EaglesOrBetter float64
    Bogeys float64
    DoubleOrWorse float64 
    TotalScore float64
}

func (t TournamentScoring) Strings() []string {
   return []string{
        strconv.FormatFloat(t.EventWinner, 'f', -1, 64),
        strconv.FormatFloat(t.Podiums, 'f', -1, 64),
        strconv.FormatFloat(t.Top10s, 'f', -1, 64),
        strconv.FormatFloat(t.HotRound, 'f', -1, 64),
        strconv.FormatFloat(t.RoundBirdies, 'f', -1, 64),
        strconv.FormatFloat(t.EaglesOrBetter, 'f', -1, 64),
        strconv.FormatFloat(t.Bogeys, 'f', -1, 64),
        strconv.FormatFloat(t.DoubleOrWorse, 'f', -1, 64),
        strconv.FormatFloat(t.TotalScore, 'f', -1, 64),
    } 
}

type TimesConfig struct {
    Length int `json:"Length"`
    Score float64 `json:"Score"`
}

func UnmarshalConfig(data []byte) (ScoringConfig, error) {
    var config ScoringConfig
    err := json.Unmarshal(data, &config)
    return config, err
}

func (c *ScoringConfig) MarshalConfig() ([]byte, error) {
    return json.Marshal(c)
}

func ParseConfig(reader io.Reader) (ScoringConfig, error) {
    data, err := io.ReadAll(reader)
    if err != nil {
        return ScoringConfig{}, err
    }

    return UnmarshalConfig(data)
}

func (team CurrentTeam) ScoreTournament(config ScoringConfig, results Results) TournamentScoring {
    total := 0.0
    tournamentScoring := TournamentScoring{}
    if team.HasWinner(results) {
        tournamentScoring.EventWinner = config.EventWinner
        total += tournamentScoring.EventWinner
    }
    tournamentScoring.Podiums = config.Podiums * float64(team.NumberPodiums(results))
    total += tournamentScoring.Podiums
    tournamentScoring.Top10s = config.Top10s * float64(team.NumberTop10s(results))
    total += tournamentScoring.Top10s
    tournamentScoring.HotRound = config.HotRound * float64(team.NumberHotRounds(results))
    total += tournamentScoring.HotRound
    tournamentScoring.RoundBirdies = team.CalculateTeamBirdies(config.RoundBirdies, results)
    total += tournamentScoring.RoundBirdies
    tournamentScoring.EaglesOrBetter = team.CalculateTeamBirdiesBetter(config.EaglesOrBetter, results)
    total += tournamentScoring.EaglesOrBetter
    tournamentScoring.Bogeys = team.CalculateTeamBogeys(config.Bogeys, results)
    total += tournamentScoring.Bogeys
    tournamentScoring.DoubleOrWorse = team.CalculateTeamBogeyWorse(config.DoubleOrWorse, results)
    total += tournamentScoring.DoubleOrWorse
    
    tournamentScoring.TotalScore = total
    return tournamentScoring;
}

func (team CurrentTeam) ScoreTeam(config ScoringConfig, results Results) float64 {
    total := 0.0
    if team.HasWinner(results) {
        total += config.EventWinner
    }
    total += config.Podiums * float64(team.NumberPodiums(results))
    total += config.Top10s * float64(team.NumberTop10s(results))
    total += config.HotRound * float64(team.NumberHotRounds(results))
    total += team.CalculateTeamBirdies(config.RoundBirdies, results)
    total += team.CalculateTeamBirdiesBetter(config.EaglesOrBetter, results)
    total += team.CalculateTeamBogeys(config.Bogeys, results)
    total += team.CalculateTeamBogeyWorse(config.DoubleOrWorse, results)

    return total;
}

func (team CurrentTeam) HasWinner(results Results) bool {
    hasMpoWinner := slices.Contains(team.Players, results.MpoWinner)
    if hasMpoWinner {
        return true
    }

    return slices.Contains(team.Players, results.FpoWinner)
}

func (team CurrentTeam) NumberPodiums(results Results) int {
    numberPodiums := 0
    for _, player := range team.Players {
        if slices.Contains(results.Podiums, player) {
            numberPodiums++
        }
    }

    return numberPodiums
}

func (team CurrentTeam) NumberTop10s(results Results) int {
    numberPodiums := 0
    for _, player := range team.Players {
        if slices.Contains(results.Top10s, player) {
            numberPodiums++
        }
    }

    return numberPodiums
}

func (team CurrentTeam) NumberHotRounds(results Results) int {
    numberHotRounds := 0
    for round := range results.HotRounds {
        hotRoundPlayers := results.HotRounds[round]
        for _, player := range team.Players {
            if slices.Contains(hotRoundPlayers, player) {
                numberHotRounds++
            }
        } 
    }

    return numberHotRounds
}

func (team CurrentTeam) CalculateTeamBirdies(config TimesConfig, results Results) float64 {
    return team.calculateRoundScore(config, results.RoundResults, Birdes)
}

func (team CurrentTeam) CalculateTeamBirdiesBetter(config TimesConfig, results Results) float64 {
    return team.calculateRoundScore(config, results.RoundResults, EaglesOrBetter)
}

func (team CurrentTeam) CalculateTeamBogeys(config TimesConfig, results Results) float64 {
    return team.calculateRoundScore(config, results.RoundResults, Bogeys)
}

func (team CurrentTeam) CalculateTeamBogeyWorse(config TimesConfig, results Results) float64 {
    return team.calculateRoundScore(config, results.RoundResults, DoubleBogeysOrWorse)
}

func (team CurrentTeam) calculateRoundScore(config TimesConfig, rounds []RoundResult, scoreExtractor ScoreExtractor) float64 {
    points :=  0.0
    for _, player := range team.Players {
        for _, round := range rounds {
            playerScore := scoreExtractor(round, player)
            multiplier := playerScore / config.Length
            points += (float64(multiplier) * config.Score)
        }
    }

    return points
}
