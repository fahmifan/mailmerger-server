-- +migrate Up
CREATE TABLE
  IF NOT EXISTS files (
    id varchar(255) primary key,
    file_name text not null,
    folder text not null,
    created_at timestamp not null default now(),
    updated_at timestamp not null default now(),
    deleted_at timestamp
  );

CREATE TABLE
  IF NOT EXISTS campaigns (
    id varchar(255) not null primary key,
    file_id varchar(255) REFERENCES files (id),
    template_id varchar(255) null,
    "name" text not null,
    body text not null,
    "subject" text not null,
    created_at timestamp not null default now(),
    updated_at timestamp not null default now(),
    deleted_at timestamp
  );

CREATE TABLE
  IF NOT EXISTS templates (
    id varchar(255) primary key,
    "name" varchar(255) not null, 
    html text not null,
    created_at timestamp not null default now(),
    updated_at timestamp not null default now(),
    deleted_at timestamp
  );

CREATE TABLE IF NOT EXISTS events (
    id varchar(255) primary key,
    campaign_id varchar(255) not null references campaigns(id),
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
