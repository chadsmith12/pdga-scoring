// Code generated by sqlc. DO NOT EDIT.
// versions:
//   sqlc v1.27.0
// source: tournaments.sql

package repository

import (
	"context"
)

const createTournament = `-- name: CreateTournament :one
insert into tournaments (external_id, name, number_of_rounds)
values ($1, $2, $3)
returning id, external_id, name, number_of_rounds
`

type CreateTournamentParams struct {
	ExternalID     int64
	Name           string
	NumberOfRounds int32
}

func (q *Queries) CreateTournament(ctx context.Context, arg CreateTournamentParams) (Tournament, error) {
	row := q.db.QueryRow(ctx, createTournament, arg.ExternalID, arg.Name, arg.NumberOfRounds)
	var i Tournament
	err := row.Scan(
		&i.ID,
		&i.ExternalID,
		&i.Name,
		&i.NumberOfRounds,
	)
	return i, err
}