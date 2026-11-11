-- 001b_init.sql — Part 2: relationships, artifacts, conversations, jobs, evaluation schema, RLS.

create table file_imports (
  id uuid primary key default gen_random_uuid(),
  analysis_id uuid not null references analyses(id) on delete cascade,
  from_file_id uuid not null references analysis_files(id) on delete cascade,
  to_file_id uuid references analysis_files(id) on delete set null,
  raw_specifier text not null,
  kind import_kind not null default 'unknown',
  created_at timestamptz not null default now(),
  unique (analysis_id, from_file_id, raw_specifier)
);

create table analysis_artifacts (
  id uuid primary key default gen_random_uuid(),
  analysis_id uuid not null references analyses(id) on delete cascade,
  kind artifact_kind not null,
  format text not null default 'markdown',
  content text,
  storage_path text,
  model text,
  created_at timestamptz not null default now()
);

create table conversations (
  id uuid primary key default gen_random_uuid(),
  analysis_id uuid not null references analyses(id) on delete cascade,
  workspace_id uuid not null references workspaces(id) on delete cascade,
  kind conversation_kind not null default 'qa',
  last_message_at timestamptz not null default now(),
  created_at timestamptz not null default now()
);

create table messages (
  id uuid primary key default gen_random_uuid(),
  conversation_id uuid not null references conversations(id) on delete cascade,
  role message_role not null,
  content text not null,
  citations jsonb,
  facts jsonb,
  hypotheses jsonb,
  unknowns jsonb,
  model text,
  input_tokens integer,
  output_tokens integer,
  latency_ms integer,
  created_at timestamptz not null default now()
);
create index messages_conversation_idx on messages(conversation_id);

create table exports (
  id uuid primary key default gen_random_uuid(),
  analysis_id uuid not null references analyses(id) on delete cascade,
  workspace_id uuid not null references workspaces(id) on delete cascade,
  storage_path text not null,
  byte_size integer not null default 0,
  expires_at timestamptz not null,
  created_at timestamptz not null default now()
);

create table job_events (
  id bigserial primary key,
  analysis_id uuid not null references analyses(id) on delete cascade,
  step text not null,
  event job_event_kind not null,
  detail jsonb,
  created_at timestamptz not null default now()
);

create table jobs (
  id uuid primary key default gen_random_uuid(),
  analysis_id uuid not null references analyses(id) on delete cascade,
  step text not null,
  status job_status not null default 'pending',
  attempts integer not null default 0,
  max_attempts integer not null default 3,
  run_at timestamptz not null default now(),
  payload jsonb,
  locked_by text,
  locked_at timestamptz,
  last_error text,
  created_at timestamptz not null default now(),
  unique (analysis_id, step)
);
create index jobs_claim_idx on jobs(status, run_at);

-- ---------- evaluation harness (separate schema) ----------
create schema eval;
create table eval.eval_benchmarks (
  id uuid primary key default gen_random_uuid(),
  repo_full_name text not null,
  pinned_commit_sha char(40),
  language text,
  notes text,
  created_at timestamptz not null default now()
);
create table eval.eval_questions (
  id uuid primary key default gen_random_uuid(),
  benchmark_id uuid references eval.eval_benchmarks(id) on delete cascade,
  category text not null check (category in ('identifier','conceptual','holistic','error_investigation','trap')),
  question text not null,
  expected jsonb,
  created_at timestamptz not null default now()
);
create table eval.eval_runs (
  id uuid primary key default gen_random_uuid(),
  started_at timestamptz not null default now(),
  finished_at timestamptz,
  config jsonb
);
create table eval.eval_results (
  id uuid primary key default gen_random_uuid(),
  run_id uuid not null references eval.eval_runs(id) on delete cascade,
  question_id uuid not null references eval.eval_questions(id) on delete cascade,
  answer text,
  citations jsonb,
  metrics jsonb,
  created_at timestamptz not null default now()
);

-- ---------- RLS: deny-by-default (server connects with a privileged role) ----------
alter table workspaces enable row level security;
alter table repositories enable row level security;
alter table analyses enable row level security;
alter table analysis_files enable row level security;
alter table code_symbols enable row level security;
alter table code_chunks enable row level security;
alter table file_imports enable row level security;
alter table analysis_artifacts enable row level security;
alter table conversations enable row level security;
alter table messages enable row level security;
alter table exports enable row level security;
alter table job_events enable row level security;
alter table jobs enable row level security;
