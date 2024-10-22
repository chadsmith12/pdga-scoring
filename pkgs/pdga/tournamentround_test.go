package pdga_test

import (
	"os"
	"testing"

	"github.com/chadsmith12/pdga-scoring/pkgs/pdga"
)

func loadTestFile(t *testing.T) []byte {
    data, err := os.ReadFile("tournament_round_test.json")
    if err != nil {
        t.Fatalf("failed to get tournament round data: %v", err)
    }

    return data
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
