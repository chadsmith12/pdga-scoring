-- name: CreateTournament :one
insert into tournaments (external_id, name, start_date, end_date, tier, location, country)
values ($1, $2, $3, $4, $5, $6, $7)
returning *;

