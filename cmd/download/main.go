package main

import (
	"context"
	"log"
	"log/slog"

	"github.com/chadsmith12/pdga-scoring/internal/database"
	"github.com/chadsmith12/pdga-scoring/internal/extractor"
	"github.com/chadsmith12/pdga-scoring/internal/repository"
	"github.com/chadsmith12/pdga-scoring/pkgs/pdga"
	"github.com/joho/godotenv"
)

var tournamentIds = []int {77772, 82419, 77773, 84345, 77098, 77774, 78651, 77775, 77758, 77759, 77091, 77760, 77761, 77762, 77099, 79049, 77763, 78193, 77764, 77765, 78647, 77766, 78666, 78194, 78271, 78654, 78195, 77768, 78196, 77750, 78197, 78655, 77769, 77771, 78646, 71315}

var tempids = []int {77761, 77762, 77099, 79049, 77763, 78193, 77764, 77765, 78647, 77766, 78666, 78194, 78271, 78654, 78195, 77768, 78196, 77750, 78197, 78655, 77769, 77771, 78646, 71315}

var insertIds = []int {77771, 78646, 71315}
func main() {
    godotenv.Load()
    db, err := database.Connect(context.Background())
    if err != nil {
        log.Fatal(err)
    }
    
    repo := repository.New(db)
    client := pdga.NewClient()
    service := extractor.NewTournamentExtractor(repo, client, slog.Default(), 3, 6)

    service.Extract(context.Background(), tournamentIds)
}
