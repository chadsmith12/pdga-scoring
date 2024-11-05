package pdga_test

import (
	"context"
	"net/http"
	"net/http/httptest"
	"os"
	"strings"
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
	data := loadTestFile(t, "test_files/tournament_info_test.json")

	tournamentData, err := pdga.UnmarshalTournamentInfo(data)
	if err != nil {
		t.Fatalf("found to unmarshal tournament data: %v", err)
	}

	if tournamentData.Hash == "" {
		t.Fatal("tournament data hash is empty")
	}
}

func TestCalculateNumberRounds(t *testing.T) {
	data := loadTestFile(t, "test_files/tournament_info_round_test.json")
	
	tournamentData, err := pdga.UnmarshalTournamentInfo(data)
	if err != nil {
		t.Fatalf("failed to unmarshall tournament data: %v", err)
	}

	t.Run("should calculate number rounds", func(t *testing.T) {
		numberRounds := tournamentData.Data.NumberRounds(pdga.Fpo)
		if numberRounds != 4 {
			t.Fatalf("expected NumberRounds to be 4; got %d", numberRounds)
		}
	})

	t.Run("should get 0 for division not played", func(t *testing.T) {
		numberRounds := tournamentData.Data.NumberRounds(pdga.Mpo)
		if numberRounds != 0 {
			t.Fatalf("expected NumberRounds to be 0; got %d", numberRounds)
		}
	})

}

func TestGetTournamentInfo(t *testing.T) {
	server := httptest.NewServer(http.HandlerFunc(testTournamentServer(t)))
	defer server.Close()

	client := testTournamentClient(server.URL)
	data, _ := client.FetchTournamentInfo(context.TODO(), 77774)
	if data.Data.TournamentID != "77774" {
		t.Errorf("expected 'TournamentId' to be 77774, got: %s", data.Data.TournamentID)
	}
}

func testTournamentClient(baseUrl string) *pdga.Client {
	return pdga.NewClient(pdga.WithBaseUrl(baseUrl))
}

func testTournamentServer(t *testing.T) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		if !strings.Contains(r.URL.Path, "live_results_fetch_event") {
			t.Errorf("expected to request `tournament.json`, got: %s", r.URL.Path)
		}
		queryValues := r.URL.Query()
		if queryValues.Get("TournID") == "" {
			t.Error("expected query string 'TournID', got empty string")
		}


		w.WriteHeader(http.StatusOK)
		w.Write(loadTestFile(t, "test_files/tournament_info_test.json"))
	}
} 
