-- 001a_init.sql — AI Codebase Intelligence & Engineering Assistant
-- Part 1: extensions, enums, core tables. Design: docs/Architecture_and_Data_Design_v1.1.md

create extension if not exists vector;
create extension if not exists pgcrypto;

-- ---------- enums ----------
create type analysis_status as enum ('queued','fetching','downloading','extracting','parsing','embedding','summarizing','ready','partially_ready','failed');
create type language_kind as enum ('python','go','javascript','typescript','other');
create type symbol_kind as enum ('function','class','method','interface','type','const','other');
create type import_kind as enum ('internal','external','unknown');
create type artifact_kind as enum ('architecture_summary','briefing_pack','parse_health','eval_report');
create type conversation_kind as enum ('qa','investigation');
create type message_role as enum ('user','assistant');
create type job_status as enum ('pending','running','succeeded','failed');
create type job_event_kind as enum ('started','succeeded','retried','failed');

-- ---------- core ----------
create table workspaces (
  id uuid primary key default gen_random_uuid(),
  token_hash text not null unique,
  last_seen_at timestamptz not null default now(),
  created_at timestamptz not null default now()
);

create table repositories (
  id uuid primary key default gen_random_uuid(),
  provider text not null default 'github',
  owner text not null,
  name text not null,
  full_name text not null unique,
  html_url text not null,
  created_at timestamptz not null default now()
);

create table analyses (
  id uuid primary key default gen_random_uuid(),
  repository_id uuid not null references repositories(id) on delete cascade,
  workspace_id uuid not null references workspaces(id) on delete cascade,
  default_branch text not null,
  commit_sha char(40) not null check (commit_sha ~ '^[0-9a-f]{40}$'),
  status analysis_status not null default 'queued',
  status_detail text,
  size_kb integer check (size_kb >= 0),
  file_count integer not null default 0,
  indexed_file_count integer not null default 0,
  symbol_count integer not null default 0,
  chunk_count integer not null default 0,
  parse_error_rate numeric(5,4) check (parse_error_rate between 0 and 1),
  parser_version text,
  query_version text,
  embedding_model text,
  llm_model text,
  skipped_files jsonb,
  error_step text,
  error_message text,
  started_at timestamptz,
  finished_at timestamptz,
  created_at timestamptz not null default now()
);
create unique index analyses_snapshot_ready_uniq on analyses(repository_id, commit_sha)
  where status in ('ready','partially_ready');
create index analyses_status_idx on analyses(status);

create table analysis_files (
  id uuid primary key default gen_random_uuid(),
  analysis_id uuid not null references analyses(id) on delete cascade,
  path text not null,
  language language_kind not null default 'other',
  size_bytes integer not null check (size_bytes >= 0),
  line_count integer check (line_count >= 0),
  parse_error_count integer not null default 0,
  is_indexed boolean not null default false,
  exclusion_reason text,
  content_hash text,
  created_at timestamptz not null default now(),
  unique (analysis_id, path)
);

create table code_symbols (
  id uuid primary key default gen_random_uuid(),
  analysis_id uuid not null references analyses(id) on delete cascade,
  file_id uuid not null references analysis_files(id) on delete cascade,
  name text not null,
  qualified_name text,
  kind symbol_kind not null default 'other',
  parent_symbol_id uuid references code_symbols(id) on delete set null,
  start_line integer not null check (start_line >= 1),
  end_line integer not null,
  is_exported boolean not null default false,
  created_at timestamptz not null default now(),
  check (end_line >= start_line)
);
create index code_symbols_name_idx on code_symbols(analysis_id, name);
create index code_symbols_file_idx on code_symbols(file_id);

create table code_chunks (
  id uuid primary key default gen_random_uuid(),
  analysis_id uuid not null references analyses(id) on delete cascade,
  file_id uuid not null references analysis_files(id) on delete cascade,
  symbol_id uuid references code_symbols(id) on delete set null,
  start_line integer not null check (start_line >= 1),
  end_line integer not null,
  content text not null,
  content_hash text not null,
  token_estimate integer not null default 1,
  embedding vector(1536),
  tsv tsvector generated always as (to_tsvector('simple', content)) stored,
  created_at timestamptz not null default now(),
  check (end_line >= start_line),
  unique (analysis_id, file_id, start_line, end_line)
);
create index code_chunks_tsv_idx on code_chunks using gin(tsv);
create index code_chunks_embedding_idx on code_chunks using hnsw (embedding vector_cosine_ops);
