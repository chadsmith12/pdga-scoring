package main

import (
	"context"
	"fmt"
	"log"

	"github.com/chadsmith12/pdga-scoring/internal/database"
	"github.com/chadsmith12/pdga-scoring/internal/repository"
	"github.com/chadsmith12/pdga-scoring/internal/scoring"
	"github.com/chadsmith12/pdga-scoring/pkgs/pdga"
	"github.com/jackc/pgx/v5/pgxpool"
	"github.com/joho/godotenv"
)

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
    tourneyInfo, err := pdgaClient.FetchTournamentInfo(context.Background(), 77774)
    if err != nil {
        log.Fatal(err)
    }
    insertedTourney, err := tournamentService.InsertTournament(context.Background(), tourneyInfo.Data, 77774)
    if err != nil {
        log.Fatal(err)
    }

    fmt.Printf("Created Tournament: %s\n", insertedTourney.Name)
    numberFpo, err := insertPlayers(pdgaClient, conn, pdga.Fpo)
    if err != nil {
        log.Fatal(err)
    }

    numberMpo, err := insertPlayers(pdgaClient, conn, pdga.Mpo)
    if err != nil {
        log.Fatal(err)
    }
    fmt.Printf("Inserted a total of %d fpo and %d mpo players\n", numberFpo, numberMpo)
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
