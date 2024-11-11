package fantasy_test

import (
	"testing"

	"github.com/chadsmith12/pdga-scoring/internal/fantasy"
)

type testRun struct {
    name    string
    config  fantasy.ScoringConfig
    results fantasy.Results
    expected float64
}


var testConfig = fantasy.ScoringConfig {
    EventWinner: 5,
    Podiums: 3,
    Top10s: 1,
    HotRound: 0.5,
    RoundBirdies: fantasy.TimesConfig{Length: 3, Score: 0.25},
    EaglesOrBetter: fantasy.TimesConfig{Length: 1, Score: 0.25},
    Bogeys: fantasy.TimesConfig{Length: 3, Score: -0.25},
    DoubleOrWorse: fantasy.TimesConfig{Length: 1, Score: -0.25},
}

func TestCalculateTeamEaglesBetter(t *testing.T) {
    testRuns := []testRun{
        {
            name:   "Calculates even eagles length",
            config: fantasy.ScoringConfig{EaglesOrBetter: fantasy.TimesConfig{Length: 1, Score: 0.25}},
            results: fantasy.Results{RoundEaglesBetter: []map[int64]int{
                    {1: 1},
            }},
            expected: 0.25,
        },
        {
            name:   "Calculates double eagles length",
            config: fantasy.ScoringConfig{EaglesOrBetter: fantasy.TimesConfig{Length: 1, Score: 0.25}},
            results: fantasy.Results{RoundEaglesBetter: []map[int64]int{
                    {1: 2},
            }},
            expected: 0.50,
        },
        {
            name:   "Calculates triple eagles length",
            config: fantasy.ScoringConfig{EaglesOrBetter: fantasy.TimesConfig{Length: 1, Score: 0.25}},
            results: fantasy.Results{RoundEaglesBetter: []map[int64]int{
                    {1: 3},
            }},
            expected: 0.75,
        },
        {
            name:   "Calculates no eagles length",
            config: fantasy.ScoringConfig{EaglesOrBetter: fantasy.TimesConfig{Length: 1, Score: 0.25}},
            results: fantasy.Results{RoundEaglesBetter: []map[int64]int{}},
            expected: 0.0,
        },
        {
            name:   "Won't calculate with no player",
            config: fantasy.ScoringConfig{EaglesOrBetter: fantasy.TimesConfig{Length: 3, Score: 0.25}},
            results: fantasy.Results{RoundEaglesBetter: []map[int64]int{
                    {2: 2},
            }},
            expected: 0.0,
        },
    }
    team := fantasy.CurrentTeam {
        Players: []int64{1},
    }

    for _, test := range testRuns {
        t.Run(test.name, func(t *testing.T) {
            actual := team.CalculateTeamBirdiesBetter(test.config.EaglesOrBetter, test.results)
            if actual != test.expected {
                t.Fatalf("actual was %2f; expected %2f", actual, test.expected)
            }
        })
    }
}

func TestCalculateTeamBirdies(t *testing.T) {
    testRuns := []testRun{
        {
            name:   "Calculates even birdies length",
            config: fantasy.ScoringConfig{RoundBirdies: fantasy.TimesConfig{Length: 3, Score: 0.25}},
            results: fantasy.Results{RoundBirdies: []map[int64]int{
                    {1: 3},
            }},
            expected: 0.25,
        },
        {
            name:   "Calculates ueven birdies length",
            config: fantasy.ScoringConfig{RoundBirdies: fantasy.TimesConfig{Length: 3, Score: 0.25}},
            results: fantasy.Results{RoundBirdies: []map[int64]int{
                    {1: 4},
            }},
            expected: 0.25,
        },
        {
            name:   "Calculates double birdies length",
            config: fantasy.ScoringConfig{RoundBirdies: fantasy.TimesConfig{Length: 3, Score: 0.25}},
            results: fantasy.Results{RoundBirdies: []map[int64]int{
                    {1: 6},
            }},
            expected: 0.50,
        },
        {
            name:   "Calculates triple birdies length",
            config: fantasy.ScoringConfig{RoundBirdies: fantasy.TimesConfig{Length: 3, Score: 0.25}},
            results: fantasy.Results{RoundBirdies: []map[int64]int{
                    {1: 9},
            }},
            expected: 0.75,
        },
        {
            name:   "Calculates no birdies length",
            config: fantasy.ScoringConfig{RoundBirdies: fantasy.TimesConfig{Length: 3, Score: 0.25}},
            results: fantasy.Results{RoundBirdies: []map[int64]int{}},
            expected: 0.0,
        },
        {
            name:   "Calculates less than value birdies length",
            config: fantasy.ScoringConfig{RoundBirdies: fantasy.TimesConfig{Length: 3, Score: 0.25}},
            results: fantasy.Results{RoundBirdies: []map[int64]int{
                    {1: 2},
            }},
            expected: 0.0,
        },
        {
            name:   "Won't calculate with no player",
            config: fantasy.ScoringConfig{RoundBirdies: fantasy.TimesConfig{Length: 3, Score: 0.25}},
            results: fantasy.Results{RoundBirdies: []map[int64]int{
                    {2: 2},
            }},
            expected: 0.0,
        },
    }
    team := fantasy.CurrentTeam {
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

func TestCalculateRoundScore(t *testing.T) {
    testRuns := []testRun{
        {
            name:   "Calculates single player round",
            config: testConfig,
            results: fantasy.Results{
                RoundEaglesBetter: []map[int64]int{},
                RoundBirdies: []map[int64]int{{1: 10}},
                RoundBogeys: []map[int64]int{{1: 2}},
                RoundDoubleWorse: []map[int64]int{{1: 1}},
                HotRounds: map[int][]int64{1: {1}},
            },
            expected: 1,
        },
        {
            name:   "Calculates multiple player round",
            config: testConfig,
            results: fantasy.Results{
                RoundEaglesBetter: []map[int64]int{{2: 1}},
                RoundBirdies: []map[int64]int{{1: 10, 2: 9}},
                RoundBogeys: []map[int64]int{{1: 2}, {2: 2}},
                RoundDoubleWorse: []map[int64]int{{1: 1}, {2: 1}},
                HotRounds: map[int][]int64{1: {1}},
            },
            expected: 1.75,
        },
    }
    team := fantasy.CurrentTeam {
        Players: []int64{1,2},
    }

    for _, test := range testRuns {
        t.Run(test.name, func(t *testing.T) {
            actual := team.ScoreTeam(test.config, test.results)
            if actual != test.expected {
                t.Fatalf("actual was %2f; expected %2f", actual, test.expected)
            }
        })
    }
}

func TestCalculateWinner(t *testing.T) {
    testRuns := []testRun{
        {
            name:   "Calculates winner",
            config: testConfig,
            results: fantasy.Results{
                MpoWinner: 3,
            },
            expected: 5,
        },
        {
            name:   "Calculates no winner",
            config: testConfig,
            results: fantasy.Results{
                MpoWinner: 4,
            },
            expected: 0,
        },
    }
    team := fantasy.CurrentTeam {
        Players: []int64{3},
    }

    for _, test := range testRuns {
        t.Run(test.name, func(t *testing.T) {
            actual := team.ScoreTeam(test.config, test.results)
            if actual != test.expected {
                t.Fatalf("actual was %2f; expected %2f", actual, test.expected)
            }
        })
    }
}

func TestEventScores(t *testing.T) {
    testRuns := []testRun{
        {
            name:   "Calculates podiums",
            config: testConfig,
            results: fantasy.Results{
                MpoWinner: 45,
                Podiums: []int64{1, 2},
            },
            expected: 6,
        },
        {
            name:   "Calculates winner and podiums",
            config: testConfig,
            results: fantasy.Results{
                MpoWinner: 3,
                Podiums: []int64{1, 2},
            },
            expected: 11,
        },
        {
            name:   "Calculates winner and podiums",
            config: testConfig,
            results: fantasy.Results{
                FpoWinner: 3,
                Podiums: []int64{1, 2},
            },
            expected: 11,
        },
        {
            name:   "Calculates winner, podiums, 1 top 10",
            config: testConfig,
            results: fantasy.Results{
                MpoWinner: 3,
                Podiums: []int64{1, 2},
                Top10s: []int64{4,5},
            },
            expected: 12,
        },
        {
            name:   "Calculates winner, podiums, 1 top 10, hot round",
            config: testConfig,
            results: fantasy.Results{
                MpoWinner: 3,
                Podiums: []int64{1, 2},
                Top10s: []int64{4,5},
                HotRounds: map[int][]int64{1: {6}},
            },
            expected: 12.5,
        },
        {
            name:   "Calculates winner, podiums, 1 top 10, multiple hot round",
            config: testConfig,
            results: fantasy.Results{
                MpoWinner: 3,
                Podiums: []int64{1, 2},
                Top10s: []int64{4,5},
                HotRounds: map[int][]int64{1: {6}, 2: {1,2}},
            },
            expected: 13.5,
        },
        {
            name:   "Calculates winner, podiums, 1 top 10, multiple hot round, birdies",
            config: testConfig,
            results: fantasy.Results{
                MpoWinner: 3,
                Podiums: []int64{1, 2},
                Top10s: []int64{4,5},
                HotRounds: map[int][]int64{1: {6}, 2: {1,2}},
                RoundBirdies: []map[int64]int{{1: 3}},
            },
            expected: 13.75,
        },
        {
            name:   "Calculates winner, podiums, 1 top 10, multiple hot round, birdies, eagles",
            config: testConfig,
            results: fantasy.Results{
                MpoWinner: 3,
                Podiums: []int64{1, 2},
                Top10s: []int64{4,5},
                HotRounds: map[int][]int64{1: {6}, 2: {1,2}},
                RoundBirdies: []map[int64]int{{1: 3}, {7: 3}},
                RoundEaglesBetter: []map[int64]int{{2: 1}, {7: 2}},
            },
            expected: 14,
        },
        {
            name:   "Calculates winner, podiums, 1 top 10, multiple hot round, birdies, eagles, bogeys",
            config: testConfig,
            results: fantasy.Results{
                MpoWinner: 3,
                Podiums: []int64{1, 2},
                Top10s: []int64{4,5},
                HotRounds: map[int][]int64{1: {6}, 2: {1,2}},
                RoundBirdies: []map[int64]int{{1: 3}, {7: 3}},
                RoundEaglesBetter: []map[int64]int{{2: 1}, {7: 2}},
                RoundBogeys: []map[int64]int{{3: 3}},
            },
            expected: 13.75,
        },
        {
            name:   "Calculates winner, podiums, 1 top 10, multiple hot round, birdies, eagles, bogeys, doubles",
            config: testConfig,
            results: fantasy.Results{
                MpoWinner: 3,
                Podiums: []int64{1, 2},
                Top10s: []int64{4,5},
                HotRounds: map[int][]int64{1: {6}, 2: {1,2}},
                RoundBirdies: []map[int64]int{{1: 3}, {7: 3}},
                RoundEaglesBetter: []map[int64]int{{2: 1}, {7: 2}},
                RoundBogeys: []map[int64]int{{3: 3}},
                RoundDoubleWorse: []map[int64]int{{3: 2}, {5: 2}},
            },
            expected: 13.25,
        },
        {
            name:   "Calculates podiums, 1 top 10, multiple hot round, birdies, eagles, bogeys, doubles",
            config: testConfig,
            results: fantasy.Results{
                MpoWinner: 10,
                Podiums: []int64{1, 2},
                Top10s: []int64{4,5},
                HotRounds: map[int][]int64{1: {6}, 2: {1,2}},
                RoundBirdies: []map[int64]int{{1: 3}, {7: 3}},
                RoundEaglesBetter: []map[int64]int{{2: 1}, {7: 2}},
                RoundBogeys: []map[int64]int{{3: 3}},
                RoundDoubleWorse: []map[int64]int{{3: 2}, {5: 2}},
            },
            expected: 8.25,
        },
    }
    team := fantasy.CurrentTeam {
        Players: []int64{1,2,3,4,6},
    }

    for _, test := range testRuns {
        t.Run(test.name, func(t *testing.T) {
            actual := team.ScoreTeam(test.config, test.results)
            if actual != test.expected {
                t.Fatalf("actual was %2f; expected %2f", actual, test.expected)
            }
        })
    }
}
