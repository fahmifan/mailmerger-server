-- +migrate Up
CREATE TABLE
  IF NOT EXISTS files (
    id text primary key,
    file_name text not null,
    folder text not null,
    created_at timestamp not null default now(),
    updated_at timestamp not null default now(),
    deleted_at timestamp
  );

CREATE TABLE
  IF NOT EXISTS campaigns (
    id text not null primary key,
    file_id text REFERENCES files (id),
    "name" text not null,
    created_at timestamp not null default now(),
    updated_at timestamp not null default now(),
    deleted_at timestamp
  );

CREATE TABLE
  IF NOT EXISTS templates (
    id text primary key,
    campaign_id text not null references campaigns(id),
    body text not null,
    "subject" text not null,
    created_at timestamp not null default now(),
    updated_at timestamp not null default now(),
    deleted_at timestamp
  );

CREATE TABLE IF NOT EXISTS events (
    id text primary key,
    campaign_id text not null references campaigns(id),
    detail text not null default '',
    "status" varchar(100) not null,
    created_at timestamp not null default now(),
    updated_at timestamp not null default now(),
    deleted_at timestamp
);

-- +migrate Down
DROP TABLE IF EXISTS events;
DROP TABLE IF EXISTS templates;
DROP TABLE IF EXISTS campaigns;
DROP TABLE IF EXISTS files;
