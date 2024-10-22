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
    roundData, err := pdgaClient.FetchTournamentRound(ctx, 77774, 2, pdga.Mpo)
    if err != nil {
        log.Fatal(err)
    }

    fmt.Printf("Tournament Name: %s\n", roundData.Data.Layouts[0].Name)
    
    for _, standing := range roundData.Top10() {
        fmt.Printf("%d: %s %d\n", standing.RunningPlace, standing.Name, standing.TotalScore())
    }
}
