package main

import (
	"context"
	"fmt"
	"log"
	"sync"

	"github.com/chadsmith12/pdga-scoring/internal/database"
	"github.com/chadsmith12/pdga-scoring/internal/repository"
	"github.com/chadsmith12/pdga-scoring/internal/scoring"
	"github.com/chadsmith12/pdga-scoring/pkgs/pdga"
	"github.com/jackc/pgx/v5/pgxpool"
	"github.com/joho/godotenv"
)

var tournamentIds = []int {77772, 82419, 77773, 84345, 77098, 77774, 78651, 77775, 77758, 77759, 77091, 77760, 77761, 77762, 77099, 79049, 77763, 78193, 77764, 77765, 78647, 77766, 78666, 78194, 78271, 78654, 78195, 77768, 78196, 77750, 78197, 78655, 77769, 77771, 78646, 71315}

func main() {
    err := godotenv.Load()
    if err != nil {
        log.Fatal(err)
    }

    conn, err := database.Connect(context.Background())
    if err != nil {
        log.Fatal(err)
    }
    tournamentService := scoring.NewTournamentService(conn)
    pdgaClient := pdga.NewClient()
    wg := sync.WaitGroup{}
    for _, id := range tournamentIds {
        wg.Add(1)
        go func() {
            tourneyInfo, err := pdgaClient.FetchTournamentInfo(context.Background(), id)
            defer wg.Done()
            if err != nil {
                fmt.Printf("Failed to get tournament with id: %d\n", id)
                return
            }
            _, err = tournamentService.InsertTournament(context.Background(), tourneyInfo.Data, id)
            if err != nil {
                fmt.Printf("failed to insert tournament: %s\n", tourneyInfo.Data.Name)
                return
            }
        }()
    }
    wg.Wait()
}

func insertPlayers(client *pdga.Client, conn *pgxpool.Pool, division pdga.Division) (int, error) {
    players, err := pdga.FetchPlayers(client, 77774, division)
    if err != nil {
        return 0, err
    }

    newPlayers := make([]repository.CreatePlayersParams, 0, len(players))
    for _, player := range players {
        newPlayers = append(newPlayers, repository.CreatePlayersParams{
            FirstName: player.FirstName,
            LastName: player.LastName,
            Name: player.Name,
            Division: string(division),
        })
    }
    repo := repository.New(conn)
    numPlayers, err := repo.CreatePlayers(context.Background(), newPlayers)

    return int(numPlayers), err
}
