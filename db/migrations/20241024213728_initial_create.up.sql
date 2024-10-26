
create table if not exists tournaments (
  id bigint primary key generated always as identity,
  external_id bigint not null,
  name text not null,
  start_date date not null,
  end_date date not null,
  tier text,
  location text,
  country text
);

create table if not exists layouts (
  id bigint primary key generated always as identity,
  tournament_id bigint references tournaments (id) not null,
  name text not null,
  course_name text not null
);

create table if not exists players (
  id bigint primary key generated always as identity,
  name text not null,
  first_name text not null,
  last_name text not null,
  division text not null check (division in ('MPO', 'FPO')),
  pdga_number bigint unique not null,
  city text,
  state_prov text,
  country text
);

create table if not exists scores (
  id bigint primary key generated always as identity,
  player_id bigint references players (id) not null,
  tournament_id bigint references tournaments (id) not null,
  layout_id bigint references layouts (id) not null,
  round_number int not null,
  score int not null
);

create table if not exists hole_scores (
  id bigint primary key generated always as identity,
  player_id bigint references players (id) not null,
  tournament_id bigint references tournaments (id) not null,
  layout_id bigint references layouts (id) not null,
  round_number int not null,
  hole_number int not null,
  par int not null,
  score_relative_to_par int not null
);

-- Create indices after table creation
create index if not exists idx_hole_scores_player_id on hole_scores using btree (player_id);
create index if not exists idx_scores_player_id on scores using btree (player_id);
create index if not exists idx_layouts_tournament_id on layouts using btree (tournament_id);
create index if not exists idx_scores_tournament_id on scores using btree (tournament_id);
create index if not exists idx_scores_layout_id on scores using btree (layout_id);
create index if not exists idx_scores_round_number on scores using btree (round_number);
create index if not exists idx_hole_scores_tournament_id on hole_scores using btree (tournament_id);
create index if not exists idx_hole_scores_layout_id on hole_scores using btree (layout_id);
create index if not exists idx_hole_scores_round_number on hole_scores using btree (round_number);
create index if not exists idx_players_pdga_number on players using btree (pdga_number);

