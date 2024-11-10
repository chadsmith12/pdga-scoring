package main

import (
	"flag"
	"fmt"
	"log"
	"os"

	"github.com/chadsmith12/pdga-scoring/internal/fantasy"
)

func main() {
    teamsPath := flag.String("teams", "", "Path to a json file that defines the teams to use for this simulation")
    configPath := flag.String("config", "", "Path to the json file that defines the scoring config for the simulation")
    //tournaments := flag.String("tournaments", "2024.json", "Path to the json file that lists the tournaments to use for the simulation") 

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

    teams, err := fantasy.LoadTeams(teamsFile)
    if err != nil {
        log.Fatal(err)
    }
    // scoreConfig, err := scoring.ParseConfig(configFile)
    // if err != nil {
    //     log.Fatal(err)
    // }

    fmt.Printf("Number Teams: %d\n", len(teams))
}
