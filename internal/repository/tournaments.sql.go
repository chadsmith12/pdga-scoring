// Code generated by sqlc. DO NOT EDIT.
// versions:
//   sqlc v1.27.0
// source: tournaments.sql

package repository

import (
	"context"
	"time"

	"github.com/jackc/pgx/v5/pgtype"
)

const createTournament = `-- name: CreateTournament :one
insert into tournaments (external_id, name, start_date, end_date, tier, location, country)
values ($1, $2, $3, $4, $5, $6, $7)
returning external_id, name, start_date, end_date, tier, location, country
`

type CreateTournamentParams struct {
	ExternalID int64
	Name       string
	StartDate  time.Time
	EndDate    time.Time
	Tier       pgtype.Text
	Location   pgtype.Text
	Country    pgtype.Text
}

func (q *Queries) CreateTournament(ctx context.Context, arg CreateTournamentParams) (Tournament, error) {
	row := q.db.QueryRow(ctx, createTournament,
		arg.ExternalID,
		arg.Name,
		arg.StartDate,
		arg.EndDate,
		arg.Tier,
		arg.Location,
		arg.Country,
	)
	var i Tournament
	err := row.Scan(
		&i.ExternalID,
		&i.Name,
		&i.StartDate,
		&i.EndDate,
		&i.Tier,
		&i.Location,
		&i.Country,
	)
	return i, err
}
