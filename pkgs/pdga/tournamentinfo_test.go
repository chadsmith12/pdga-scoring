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

func TestCalculateNumberRounds(t *testing.T) {
	data := loadTestFile(t, "tournament_info_round_test.json")
	
	tournamentData, err := pdga.UnmarshalTournamentInfo(data)
	if err != nil {
		t.Fatalf("failed to unmarshall tournament data: %v", err)
	}

	numberRounds := tournamentData.Data.NumberRounds()
	if numberRounds != 4 {
		t.Fatalf("expected NumberRounds to be 4; got %d", numberRounds)
	}
}
