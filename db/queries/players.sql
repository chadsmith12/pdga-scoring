-- name: GetPlayers :many
select * from players
order by first_name;

-- name: CreatePlayers :copyfrom
insert into players (first_name, last_name, name, division, pdga_number, city, state_prov, country)
values ($1, $2, $3, $4, $5, $6, $7, $8);
