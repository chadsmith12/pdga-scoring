package main

import (
	"errors"
	"fmt"
	"log"
	"path/filepath"

	"github.com/chadsmith12/pdga-scoring/internal/config"
	"github.com/golang-migrate/migrate/v4"
	_ "github.com/golang-migrate/migrate/v4/database/postgres"
	_ "github.com/golang-migrate/migrate/v4/source/file"
	"github.com/joho/godotenv"
)

func main() {
    err := godotenv.Load()
    if err != nil {
        log.Fatal(err)
    }
    
    dbConfig, err := config.LoadDatabase()
    if err != nil {
        log.Fatal(err)
    }

    currentPath, err := filepath.Abs("db/migrations")
    if err != nil  {
	log.Fatal(err)
    }
    sourcePath := fmt.Sprintf("file://%s", currentPath)
    fmt.Printf("%s\n", sourcePath)
    migrator, err := migrate.New(sourcePath, dbConfig.String())
    if err != nil {
	log.Fatal(err)
    }
    
    if err := migrator.Up(); err != nil && !errors.Is(err, migrate.ErrNoChange) {
        log.Fatal(err)
    }

    fmt.Println("Migrated Database")
}
