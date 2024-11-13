package fantasy_test

import (
	"testing"

	"github.com/chadsmith12/pdga-scoring/internal/fantasy"
)

type teamTestRun struct {
	name       string
	team       fantasy.Team
	mpoPlayers []int64
	fpoPlayers []int64
	expected   fantasy.CurrentTeam
}

func TestFantasyTeams(t *testing.T) {
	testRuns := []teamTestRun{
		{
			name: "Uses Default team",
			team: fantasy.Team{
				Name: "Test",
				Team: fantasy.FantasyPlayers{
					Players: []int64{1, 2, 3, 4, 5},
					Bench:   []int64{6, 7},
				},
			},
			mpoPlayers: []int64{1, 2, 3},
			fpoPlayers: []int64{4, 5},
			expected: fantasy.CurrentTeam{
				MpoPlayers: []int64{1, 2, 3},
				FpoPlayers: []int64{4, 5},
				Players:    []int64{1, 2, 3, 4, 5},
			},
		},
		{
			name: "Uses the MPo Bench with player not playing",
			team: fantasy.Team{
				Name: "Test",
				Team: fantasy.FantasyPlayers{
					Players: []int64{1, 2, 3, 4, 5},
					Bench:   []int64{6, 7},
				},
			},
			mpoPlayers: []int64{1, 2, 8},
			fpoPlayers: []int64{4, 5},
			expected: fantasy.CurrentTeam{
				MpoPlayers: []int64{1, 2, 7},
				FpoPlayers: []int64{4, 5},
				Players:    []int64{1, 2, 4, 5, 7},
			},
		},
		{
			name: "Uses the Fpo Bench with player not playing",
			team: fantasy.Team{
				Name: "Test",
				Team: fantasy.FantasyPlayers{
					Players: []int64{1, 2, 3, 4, 5},
					Bench:   []int64{6, 7},
				},
			},
			mpoPlayers: []int64{1, 2, 3},
			fpoPlayers: []int64{6, 5},
			expected: fantasy.CurrentTeam{
				MpoPlayers: []int64{1, 2, 3},
				FpoPlayers: []int64{5, 6},
				Players:    []int64{1, 2, 3, 5, 6},
			},
		},
	}

	for _, testRun := range testRuns {
		t.Run(testRun.name, func(t *testing.T) {
			actual := testRun.team.CreateTeam(testRun.mpoPlayers, testRun.fpoPlayers)
			assertSlicesEqual(t, "Mpo Players", testRun.expected.MpoPlayers, actual.MpoPlayers)
			assertSlicesEqual(t, "Fpo Players", testRun.expected.FpoPlayers, actual.FpoPlayers)
			assertSlicesEqual(t, "Players", testRun.expected.Players, actual.Players)
		})
	}
}

func assertSlicesEqual(t *testing.T, title string, expected []int64, actual []int64) {
	if len(expected) != len(actual) {
		t.Fatalf("%s - slices are wrong lengths. expected: %d, got %d", title, len(expected), len(actual))
	}

	for i, actual := range actual {
		if expected[i] != actual {
			t.Fatalf("%s - actual[%d] is wrong. expected: %d, actual: %d", title, i, expected[i], actual)
		}
	}
}
