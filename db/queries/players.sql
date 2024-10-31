-- name: GetPlayers :many
select * from players
order by first_name;

-- name: CreateManyPlayers :batchexec
insert into players (pdga_number, first_name, last_name, name, division, city, state_prov, country)
values ($1, $2, $3, $4, $5, $6, $7, $8)
on conflict (pdga_number) do nothing;
