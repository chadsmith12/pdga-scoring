package scoring

import (
	"slices"
)

type PositionMap map[int]int64

type ScoringConfig struct {
    EventWinner float64
    Podiums float64
    Top10s float64
    HotRound float64
    RoundBirdies TimesConfig
    EaglesOrBetter TimesConfig
    Bogeys TimesConfig
    DoubleOrWorse TimesConfig
}

type TimesConfig struct {
    Length int
    Score float64
}

type Results struct {
    Winner int64
    Podiums []int64
    Top10s []int64
    HotRounds map[int][]int64
    RoundBirdies []map[int64]int
    RoundEaglesBetter []map[int64]int
    RoundBogeys []map[int64]int
    RoundDoubleWorse []map[int64]int
}

type FantasyTeam struct {
    Players []int64
    Bench []int64
}

type CurrentTeam struct {
    Players []int64
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
    return slices.Contains(team.Players, results.Winner)
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
    return team.calculateRoundScore(config, results.RoundBirdies)
}

func (team CurrentTeam) CalculateTeamBirdiesBetter(config TimesConfig, results Results) float64 {
    return team.calculateRoundScore(config, results.RoundEaglesBetter)
}

func (team CurrentTeam) CalculateTeamBogeys(config TimesConfig, results Results) float64 {
    return team.calculateRoundScore(config, results.RoundBogeys)
}

func (team CurrentTeam) CalculateTeamBogeyWorse(config TimesConfig, results Results) float64 {
    return team.calculateRoundScore(config, results.RoundDoubleWorse)
}

func (team CurrentTeam) calculateRoundScore(config TimesConfig, roundScores []map[int64]int) float64 {
    points :=  0.0
    for _, player := range team.Players {
        for _, round := range roundScores {
            playerScore, ok := round[player]
            if !ok {
                continue
            }
            multiplier := playerScore / config.Length
            points += (float64(multiplier) * config.Score)
        }
    }

    return points
}
