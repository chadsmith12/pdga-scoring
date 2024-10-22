package pdga_test

import (
	"context"
	"os"
	"testing"
	"time"

	"github.com/chadsmith12/pdga-scoring/pkgs/pdga"
)

func loadTestFile(t *testing.T) []byte {
	data, err := os.ReadFile("tournament_round_test.json")
	if err != nil {
		t.Fatalf("failed to get tournament round data: %v", err)
	}

	return data
}

func assertSlice[T comparable](t *testing.T, actual, expected []T) {
	for i, item := range actual {
		expectedItem := expected[i]
		if item != expectedItem {
			t.Fatalf("actual[%d] != expected[%d] - expected: %v - wanted: %v", i, i, expectedItem, item)
		}
	}
}

func TestUnmarshalTournamentRound(t *testing.T) {
	data := loadTestFile(t)

	roundData, err := pdga.UnmarshalTournamentRoundData(data)
	if err != nil {
		t.Fatalf("found to unmarshal tournament round data: %v", err)
	}

	if roundData.Hash == "" {
		t.Fatal("tournament round data hash is empty")
	}
}

func TestHoleScoring(t *testing.T) {
	data := loadTestFile(t)
	roundData, err := pdga.UnmarshalTournamentRoundData(data)
	if err != nil {
		t.Fatalf("found to unmarshal tournament round data: %v", err)
	}

	firstScores := roundData.Data.Scores[0]
	expected := []int{3, 3, 4, 4, 4, 4, 3, 3, 4, 2, 3, 4, 4, 2, 3, 2, 3, 3}
	holeScores := firstScores.HoleScoring()

	if len(expected) != len(holeScores) {
		t.Fatal("len(expected) != len(holeScores)")
	}

	assertSlice(t, holeScores, expected)
}

func TestFetchTournamentRoundData(t *testing.T) {
	client := pdga.NewClient()

	ctx, cancel := context.WithTimeout(context.Background(), time.Second*2)
	defer cancel()
	tournamentData, err := client.FetchTournamentRound(ctx, 77774, 2, pdga.Mpo)
	if err != nil {
		t.Fatalf("failed to fetch tournament round data: %v", err)
	}

	if tournamentData.Hash != "c28a4381c04257d7610e9bdb8ac84fd3" {
		t.Fatalf("failed to get the valid hash for tournament round data")
	}
}
