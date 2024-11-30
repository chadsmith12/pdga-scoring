package fantasy_test

import (
	"testing"

	"github.com/chadsmith12/pdga-scoring/pkgs/fantasy"
	"gotest.tools/assert"
)

type testRun struct {
	name     string
	config   fantasy.ScoringConfig
	results  fantasy.Results
	expected float64
}

var testConfig = fantasy.ScoringConfig{
	EventWinner:    5,
	Podiums:        3,
	Top10s:         1,
	HotRound:       0.5,
	RoundBirdies:   fantasy.TimesConfig{Length: 3, Score: 0.25},
	EaglesOrBetter: fantasy.TimesConfig{Length: 1, Score: 0.25},
	Bogeys:         fantasy.TimesConfig{Length: 3, Score: -0.25},
	DoubleOrWorse:  fantasy.TimesConfig{Length: 1, Score: -0.25},
}

func setupTestResults() fantasy.Results {
	return fantasy.Results{
		MpoWinner: 1,
		FpoWinner: 0,
		Podiums:   []int64{1, 2},
		Top10s:    []int64{1, 2, 3},
		HotRounds: map[int][]int64{
			1: {1},
			2: {2},
		},
		RoundResults: []fantasy.RoundResult{
			fantasy.NewRoundResult(1,
				map[int64]int{1: 3, 2: 0},
				map[int64]int{1: 1, 2: 0},
				map[int64]int{1: 2, 2: 3},
				map[int64]int{1: 1, 2: 0},
			),
			fantasy.NewRoundResult(2,
				map[int64]int{1: 0, 2: 3},
				map[int64]int{1: 0, 2: 1},
				map[int64]int{1: 0, 2: 1},
				map[int64]int{1: 0, 2: 0},
			),
		},
	}
}

func TestCalculateTeamBirdies(t *testing.T) {
    results := setupTestResults()
    teams := fantasy.CurrentTeam{ Players: []int64{1, 2} }
    points := teams.CalculateTeamBirdies(testConfig.RoundBirdies, results)
    assert.Equal(t, 0.5, points)
}

func TestCalculateTeamEaglesOrBetter(t *testing.T) {
    results := setupTestResults()
    teams := fantasy.CurrentTeam{Players: []int64{1, 2}}
    points := teams.CalculateTeamBirdiesBetter(testConfig.EaglesOrBetter, results)
    assert.Equal(t, 0.5, points)
}

func TestCalculateTeamBogeys(t *testing.T) {
    results := setupTestResults()
    teams := fantasy.CurrentTeam{Players: []int64{1, 2}}
    points := teams.CalculateTeamBogeys(testConfig.Bogeys, results)
    assert.Equal(t, -0.25, points)
}

func TestCalculateTeamDoubleOrWorse(t *testing.T) {
    results := setupTestResults()
    teams := fantasy.CurrentTeam{Players: []int64{1, 2}}
    points := teams.CalculateTeamBogeyWorse(testConfig.DoubleOrWorse, results)
    assert.Equal(t, -0.25, points)
}

func TestCalculateRoundResultsScore(t *testing.T) {
    results := setupTestResults()
    teams := fantasy.CurrentTeam{ Players: []int64{1, 2}}
    total := teams.CalculateTeamBirdies(testConfig.RoundBirdies, results) +
        teams.CalculateTeamBirdiesBetter(testConfig.EaglesOrBetter, results) +
        teams.CalculateTeamBogeys(testConfig.Bogeys, results) +
        teams.CalculateTeamBogeyWorse(testConfig.DoubleOrWorse, results)

    assert.Equal(t, 0.50, total)
}

func TestCalculateWinner(t *testing.T) {
    results := setupTestResults()
    teams := fantasy.CurrentTeam{Players: []int64{1, 2}}
    winner := teams.HasWinner(results)
    assert.Equal(t, true, winner)
}

func TestScoreTeam(t *testing.T) {
    results := setupTestResults()
    teams := fantasy.CurrentTeam{Players: []int64{1, 2}}
    total := teams.ScoreTeam(testConfig, results)
    assert.Equal(t, 14.50, total)
}

