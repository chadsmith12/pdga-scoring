package server

import (
	"context"

	"github.com/chadsmith12/pdga-scoring/internal/database"
	tournamentsV1 "github.com/chadsmith12/pdga-scoring/internal/server/tournaments/v1"
	"github.com/chadsmith12/pdga-scoring/pkgs/pulse"
	"github.com/jackc/pgx/v5/pgxpool"
	"github.com/joho/godotenv"
)

type Server struct {
    pulse *pulse.PulseApp
    db *pgxpool.Pool
}

func NewServer() *Server {
	godotenv.Load()
    pulseApp := pulse.Pulse()

    return &Server{pulse: pulseApp}
}

func (s *Server) Start() error {
    db, err := database.Connect(context.Background())
    if err != nil {
		return err
    }
	s.db = db
	s.setupRoutes()

	return s.pulse.Start()
}

func (s *Server) setupRoutes() {
    group := s.pulse.Group("/api")

	tournamentRoutes := tournamentsV1.NewHandlers(s.db, s.pulse.Logger())
	tournamentRoutes.TournamentRoutes(group)
}
