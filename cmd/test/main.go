package main

import (
	"fmt"
	"log"

	"github.com/chadsmith12/pdga-scoring/pkgs/pdga"
)

func main() {
    pdgaClient := pdga.NewClient()

    roundData, err := pdgaClient.FetchTournamentRound(77774, 2, pdga.Mpo)
    if err != nil {
        log.Fatal(err)
    }

    fmt.Printf("Tournament Name: %s\n", roundData.Data.Layouts[0].Name)
    
    for _, standing := range roundData.Data.Scores[:10] {
        fmt.Printf("%d: %s\n", standing.RunningPlace, standing.Name)
    }
}
