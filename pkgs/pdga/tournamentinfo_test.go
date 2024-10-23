package pdga_test

import (
	"os"
	"testing"

	"github.com/chadsmith12/pdga-scoring/pkgs/pdga"
)

func loadTestFile(t *testing.T, file string) []byte {
	data, err := os.ReadFile(file)
	if err != nil {
		t.Fatalf("failed to get tournament round data: %v", err)
	}

	return data
}

func TestUnmarshallTournamentInfo(t *testing.T) {
	data := loadTestFile(t, "tournament_info_test.json")

	tournamentData, err := pdga.UnmarshalTournamentInfo(data)
	if err != nil {
		t.Fatalf("found to unmarshal tournament data: %v", err)
	}

	if tournamentData.Hash == "" {
		t.Fatal("tournament data hash is empty")
	}
}
