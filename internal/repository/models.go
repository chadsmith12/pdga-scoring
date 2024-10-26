// Code generated by sqlc. DO NOT EDIT.
// versions:
//   sqlc v1.27.0

package repository

import (
	"github.com/jackc/pgx/v5/pgtype"
)

type HoleScore struct {
	ID                 int64
	PlayerID           int64
	TournamentID       int64
	LayoutID           int64
	RoundNumber        int32
	HoleNumber         int32
	Par                int32
	ScoreRelativeToPar int32
}

type Layout struct {
	ID           int64
	TournamentID int64
	Name         string
	CourseName   string
}

type Player struct {
	ID         int64
	Name       string
	FirstName  string
	LastName   string
	Division   string
	PdgaNumber int64
	City       pgtype.Text
	StateProv  pgtype.Text
	Country    pgtype.Text
}

type Score struct {
	ID           int64
	PlayerID     int64
	TournamentID int64
	LayoutID     int64
	RoundNumber  int32
	Score        int32
}

type Tournament struct {
	ID         int64
	ExternalID int64
	Name       string
	StartDate  pgtype.Date
	EndDate    pgtype.Date
	Tier       pgtype.Text
	Location   pgtype.Text
	Country    pgtype.Text
}
