package main

import (
	"context"
	"fmt"
	"log"
	"time"

	"github.com/chadsmith12/pdga-scoring/pkgs/pdga"
)

func main() {
    pdgaClient := pdga.NewClient()
    ctx, cancel := context.WithTimeout(context.Background(), time.Second * 10)
    defer cancel()
    tournamentData, err := pdgaClient.FetchTournamentInfo(ctx, 77774)
    if err != nil {
        log.Fatal(err)
    }

    tournamentRounds := make([]pdga.TournamentRoundData, 0, tournamentData.Data.Rounds)
    for roundNumber := range tournamentData.Data.Rounds {
        if (roundNumber == 0) { continue } 
        fmt.Printf("Round: %d\n", roundNumber)
        roundData, err := pdgaClient.FetchTournamentRound(ctx, 77774, int(roundNumber + 1), pdga.Mpo)
        if err != nil {
            log.Fatal(err)
        }
        tournamentRounds = append(tournamentRounds, roundData)
    }
    
    roundData := pdga.TournamentRounds(tournamentRounds)
    finals := roundData.FinalStandings()

    for _, player := range finals {
        fmt.Printf("%d: %s\n", player.RunningPlace, player.Name)
    }
}
