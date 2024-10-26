package main

import (
	"log"
	"sync"

	"github.com/chadsmith12/pdga-scoring/pkgs/pdga"
	"github.com/joho/godotenv"
)

var tournamentIds = []int {77772, 82419, 77773, 84345, 77098, 77774, 78651, 77775, 77758, 77759, 77091, 77760, 77761, 77762, 77099, 79049, 77763, 78193, 77764, 77765, 78647, 77766, 78666, 78194, 78271, 78654, 78195, 77768, 78196, 77750, 78197, 78655, 77769, 77771, 78646, 71315}

func main() {
    err := godotenv.Load()
    if err != nil {
        log.Fatal(err)
    }

    wg := sync.WaitGroup{}
    for _, id := range tournamentIds {
        wg.Add(1)
        go func() {
            defer wg.Done()

            if err := pdga.DownloadTournament("tournaments/", id); err != nil {
                log.Printf("failed to download tournament data for id: %d - %v\n", id, err)
            }

            if err := pdga.DownloadRoundData("tournaments/", id, 1); err != nil {
                log.Printf("failed download round data for id: %d - %v\n", id, err)
            }
        }()
    }
    wg.Wait()
}

