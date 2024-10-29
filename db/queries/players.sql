-- name: GetPlayers :many
select * from players
order by first_name;

-- name: CreatePlayers :copyfrom
insert into players (first_name, last_name, name, division, pdga_number, city, state_prov, country)
values ($1, $2, $3, $4, $5, $6, $7, $8);

-- name: CreateManyPlayers :batchexec
insert into players (first_name, last_name, name, division, pdga_number, city, state_prov, country)
values (
    sqlc.arg(first_name),
    sqlc.arg(last_name),
    sqlc.arg(name),
    sqlc.arg(division),
    sqlc.arg(pdga_number),
    sqlc.arg(city),
    sqlc.arg(state_prov),
    sqlc.arg(country)
)
on conflict (pdga_number) do nothing;
