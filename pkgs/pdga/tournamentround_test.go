package pdga_test

import (
	"context"
	"testing"
	"time"

	"github.com/chadsmith12/pdga-scoring/pkgs/pdga"
)

func assertSlice[T comparable](t *testing.T, actual, expected []T) {
	for i, item := range actual {
		expectedItem := expected[i]
		if item != expectedItem {
			t.Fatalf("actual[%d] != expected[%d] - expected: %v - wanted: %v", i, i, expectedItem, item)
		}
	}
}

func TestUnmarshalTournamentRound(t *testing.T) {
	data := loadTestFile(t, "test_files/tournament_round_test.json")

	roundData, err := pdga.UnmarshalTournamentRoundData(data)
	if err != nil {
		t.Fatalf("failed to unmarshal tournament round data: %v", err)
	}

	if roundData.Hash == "" {
		t.Fatal("tournament round data hash is empty")
	}
}


func TestUnmarshalPoolTournamentRound(t *testing.T) {
	data := loadTestFile(t, "test_files/pool_test.json")

	roundData, err := pdga.UnmarshalTournamentRoundData(data)
	if err != nil {
		t.Fatalf("failed to unmarshal tournament round data: %v", err)
	}

	if roundData.Hash == "" {
		t.Fatal("tournament round data hash is empty")
	}
	if len(roundData.Data.RoundData) == 0 {
		t.Fatal("tournament round data was empty")
	}
}

func TestFetchTournamentRoundData(t *testing.T) {
	client := pdga.NewClient()

	ctx, cancel := context.WithTimeout(context.Background(), time.Second*2)
	defer cancel()
	tournamentData, err := client.FetchTournamentRound(ctx, 77774, 2, pdga.Mpo)
	if err != nil {
		t.Fatalf("failed to fetch tournament round data: %v", err)
	}

	if tournamentData.Hash == "" {
		t.Fatalf("failed to get the valid hash for tournament round data")
	}
}
