package pdga_test

import (
	"testing"

	"github.com/chadsmith12/pdga-scoring/pkgs/pdga"
)

func TestCanExtractPlayers(t *testing.T) {
    data := loadTestFile(t, "pool_test.json")
    roundInfo, _ := pdga.UnmarshalTournamentRoundData(data)
    allRounds := []pdga.TournamentRoundResponse{ roundInfo }
    fullTournament := pdga.FullTournamentRound(allRounds)

    players := fullTournament.Players()
    if len(players) != 208 {
        t.Fatalf("numbers of players was %d, expected 208", len(players))
    }
}
