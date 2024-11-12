package fantasy

import (
	"encoding/json"
	"io"
	"slices"
)

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
    return team.calculateRoundScore(config, results.RoundBirdies)
}

func (team CurrentTeam) CalculateTeamBirdiesBetter(config TimesConfig, results Results) float64 {
    // fmt.Printf("Number Results for birdies or better: %d\n", len(results.RoundEaglesBetter))
    // for _, data := range results.RoundEaglesBetter {
    //     fmt.Printf("Number of players in the eagles or better category: %d\n", len(data))
    // }
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
