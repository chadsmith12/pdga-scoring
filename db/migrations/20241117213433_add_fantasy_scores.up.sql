create table fantasy_round_scores (
  id bigint primary key generated always as identity,
  player_id bigint references players (pdga_number),
  tournament_id bigint references tournaments (external_id),
  round_number int,
  birdies int,
  eagles_or_better int,
  bogeys int,
  double_or_worse_bogeys int,
  unique (player_id, tournament_id, round_number)
);

create index idx_fantasy_round_scores_player_tournament 
on fantasy_round_scores using btree (player_id, tournament_id);

create table fantasy_tournament_scores (
  id bigint primary key generated always as identity,
  player_id bigint references players (pdga_number),
  tournament_id bigint references tournaments (external_id),
  won_event boolean,
  podium_finish boolean,
  top_10_finish boolean,
  hot_rounds int,
  unique (player_id, tournament_id)
);

create index idx_fantasy_tournament_scores_player_tournament
on fantasy_tournament_scores using btree (player_id, tournament_id);
