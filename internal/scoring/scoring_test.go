package scoring_test

import (
	"testing"

	"github.com/chadsmith12/pdga-scoring/internal/scoring"
)

type testRun struct {
    name    string
    config  scoring.ScoringConfig
    results scoring.Results
    expected float64
}

func TestCalculateTeamBirdies(t *testing.T) {
    testRuns := []testRun{
        {
            name:   "Calculates even birdies length",
            config: scoring.ScoringConfig{RoundBirdies: scoring.TimesConfig{Length: 3, Score: 0.25}},
            results: scoring.Results{RoundBirdies: []map[int64]int{
                    {1: 3},
            }},
            expected: 0.25,
        },
        {
            name:   "Calculates ueven birdies length",
            config: scoring.ScoringConfig{RoundBirdies: scoring.TimesConfig{Length: 3, Score: 0.25}},
            results: scoring.Results{RoundBirdies: []map[int64]int{
                    {1: 4},
            }},
            expected: 0.25,
        },
        {
            name:   "Calculates double birdies length",
            config: scoring.ScoringConfig{RoundBirdies: scoring.TimesConfig{Length: 3, Score: 0.25}},
            results: scoring.Results{RoundBirdies: []map[int64]int{
                    {1: 6},
            }},
            expected: 0.50,
        },
        {
            name:   "Calculates triple birdies length",
            config: scoring.ScoringConfig{RoundBirdies: scoring.TimesConfig{Length: 3, Score: 0.25}},
            results: scoring.Results{RoundBirdies: []map[int64]int{
                    {1: 9},
            }},
            expected: 0.75,
        },
        {
            name:   "Calculates no birdies length",
            config: scoring.ScoringConfig{RoundBirdies: scoring.TimesConfig{Length: 3, Score: 0.25}},
            results: scoring.Results{RoundBirdies: []map[int64]int{}},
            expected: 0.0,
        },
        {
            name:   "Calculates less than value birdies length",
            config: scoring.ScoringConfig{RoundBirdies: scoring.TimesConfig{Length: 3, Score: 0.25}},
            results: scoring.Results{RoundBirdies: []map[int64]int{
                    {1: 2},
            }},
            expected: 0.0,
        },
        {
            name:   "Won't calculate with no player",
            config: scoring.ScoringConfig{RoundBirdies: scoring.TimesConfig{Length: 3, Score: 0.25}},
            results: scoring.Results{RoundBirdies: []map[int64]int{
                    {2: 2},
            }},
            expected: 0.0,
        },
    }
    team := scoring.CurrentTeam {
        Players: []int64{1},
    }
    
    for _, test := range testRuns {
        t.Run(test.name, func(t *testing.T) {
            actual := team.CalculateTeamBirdies(test.config.RoundBirdies, test.results)
            if actual != test.expected {
                t.Fatalf("actual was %2f; expected %2f", actual, test.expected)
            }
        })
    }
}
