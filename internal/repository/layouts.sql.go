// Code generated by sqlc. DO NOT EDIT.
// versions:
//   sqlc v1.27.0
// source: layouts.sql

package repository

import (
	"context"
)

const getLayoutIdsByTournament = `-- name: GetLayoutIdsByTournament :many
select id from layouts
where tournament_id = $1
`

func (q *Queries) GetLayoutIdsByTournament(ctx context.Context, tournamentID int64) ([]int64, error) {
	rows, err := q.db.Query(ctx, getLayoutIdsByTournament, tournamentID)
	if err != nil {
		return nil, err
	}
	defer rows.Close()
	var items []int64
	for rows.Next() {
		var id int64
		if err := rows.Scan(&id); err != nil {
			return nil, err
		}
		items = append(items, id)
	}
	if err := rows.Err(); err != nil {
		return nil, err
	}
	return items, nil
}
