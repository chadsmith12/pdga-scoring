package main

import (
	"context"
	"encoding/json"
	"flag"
	"io"
	"log"
	"os"

	"github.com/chadsmith12/pdga-scoring/internal/database"
	"github.com/chadsmith12/pdga-scoring/internal/fantasy"
	"github.com/chadsmith12/pdga-scoring/internal/simulator"
	"github.com/joho/godotenv"
)

func main() {
    godotenv.Load()

    db, err := database.Connect(context.Background())
    if err != nil {
        log.Fatal(err)
    }
    teamsPath := flag.String("teams", "", "Path to a json file that defines the teams to use for this simulation")
    configPath := flag.String("config", "", "Path to the json file that defines the scoring config for the simulation")
    tournamentsPath := flag.String("tournaments", "2024.json", "Path to the json file that lists the tournaments to use for the simulation") 

    flag.Parse()

    teamsFile, err := os.Open(*teamsPath)
    if err != nil {
       log.Fatal(err) 
    }
    defer teamsFile.Close()

    configFile, err := os.Open(*configPath)
    if err != nil {
        log.Fatal(err)
    }
    defer configFile.Close()
    
    tournamentsFile, err := os.Open(*tournamentsPath)
    if err != nil {
        log.Fatal(err)
    }
    defer tournamentsFile.Close()

    teams, err := fantasy.LoadTeams(teamsFile)
    if err != nil {
        log.Fatal(err)
    }
    scoreConfig, err := fantasy.ParseConfig(configFile)
    if err != nil {
        log.Fatal(err)
    }
    
    sim := simulator.NewSimulator(scoreConfig, teams, unmarshalTournaments(tournamentsFile), db)
    sim.Run()
    sim.ExportResults()
}

func unmarshalTournaments(r io.Reader) []int64 {
    data, err := io.ReadAll(r)
    if err != nil {
        log.Fatal(err) 
    }
    var tournaments []int64
    err = json.Unmarshal(data, &tournaments)
    if err != nil {
        log.Fatal(err)
    }

    return tournaments
}
