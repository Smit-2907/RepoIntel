# Architecture & Data Design — v1.1
## AI Codebase Intelligence & Engineering Assistant
**Design document | 10 November 2026 | Backend revised (Go + Python)**

> **Status:** v1.1 incorporates the approved backend revision (handoff v1.2): Go HTTP API + Go pipeline worker + Python parser service on one host; Cloudflare reduced to DNS + static frontend hosting; Supabase Storage buckets (CDN-backed) serve files. Implementation is approved (handoff v1.1). All schema sketches remain design specifications; migrations are created during implementation.

---

## 1. System overview

### 1.1 Components

| # | Component | Platform | Responsibility |
|---|---|---|---|
| 1 | Web app (UI) | Cloudflare static hosting (Vite SPA) | URL input, progress view, dashboard, Q&A, citations panel, export |
| 2 | HTTP API | Go (single host) | Validation, rate limits, CRUD, query orchestration, capped Gemini tool-calling |
| 3 | Pipeline worker | Go (single host) | Postgres-backed job queue; runs steps 2–8 with retries; progress in analyses.status + job_events |
| 4 | Parser service | Python + py-tree-sitter (Docker) | Deterministic extraction: files, symbols, imports, chunk boundaries, error counts |
| 5 | Job queue | Postgres `jobs` table (Supabase) | Durable step scheduling, leases, retries, backoff (replaces Workflows/Queues) |
| 6 | Database | Supabase Postgres + pgvector | All structured data, vectors, full-text search, RLS |
| 7 | Object storage | Supabase Storage (2 private buckets) | Ephemeral repo archives; Briefing Pack exports |
| 8 | LLM provider | Gemini API | `gemini-3.8-flash` (overview, Q&A, investigation); `gemini-embedding-001` @1536 |
| 9 | Source control | GitHub REST API + codeload | Repo metadata (default branch, head SHA, size) + tarball download |

### 1.2 Flow diagram

```
User → [1] UI (Cloudflare static) → [2] Go API ── create analysis ──► Supabase (analyses: queued; jobs enqueued)
                                    │
                                    └─► [3] Go pipeline worker (claims jobs from the jobs table)
                                          step 2 fetch metadata ──► GitHub REST (default_branch, SHA, size)
                                          step 3 download       ──► codeload tarball → Supabase Storage
                                          step 4 safe extract   ──► file manifest (guards applied)
                                          step 5 parse          ──► [4] Python parser service (py-tree-sitter ×4)
                                          step 6 embed          ──► Gemini → pgvector
                                          step 7 graph+summary  ──► file_imports + 1 Gemini call → Postgres
                                          step 8 finalize       ──► ready / partially_ready; archive deleted ≤1h

User question → [2] Go API → embed query → hybrid RRF (single SQL fn) → capped tool loop
              → answer + backend-validated citations → conversations/messages
```

### 1.3 Design principles carried from approved decisions

- Deterministic facts are extracted **before** any LLM call (Research 1).
- Every analysis is pinned to an immutable commit SHA; all citations refer to that snapshot (Research 6).
- The LLM never sees raw repository instructions as authority; repository text is data (handoff §11).
- All retrieval flows through one auditable SQL function, so evaluation ablations are configuration, not forks (Research 3/8).
- Costs are bounded by design: capped tool rounds, chunk limits, monthly budget (handoff §6.2).

---

## 2. Canonical pipeline (URL → queryable snapshot)

| Step | Name | Owner | Input | Output | Retries | Timeout |
|---|---|---|---|---|---|---|
| 0 | Validate URL | Go API | `repo_url` | normalized `owner/repo` or 422 | — | 2 s |
| 1 | Create analysis | Go API | normalized repo | `analyses` row (queued) + job enqueued | — | 5 s |
| 2 | Fetch metadata | Go worker | owner/repo | `default_branch`, head `commit_sha`, `size_kb`, primary language | 3× | 30 s |
| 3 | Download archive | Go worker | commit SHA | tarball stored in `repo-archives` bucket | 3× | 120 s |
| 4 | Safe extract | Go worker | tarball | file manifest after guards (limits, exclusions) | 1× | 120 s |
| 5 | Parse | Go worker → Python parser | file manifest | files, symbols, imports, chunks, per-file error counts | 2× | 10 min |
| 6 | Embed | Go worker | chunks | 1536-dim vectors in `code_chunks` (pgvector) | 3× per batch | 15 min |
| 7 | Graph + overview | Go worker | imports + manifest | `file_imports` edges + architecture summary artifact (1 Gemini call) | 2× | 120 s |
| 8 | Finalize | Go worker | all | status `ready` / `partially_ready`; archive deletion scheduled | 3× | 30 s |

**Status transitions (single source of truth):**
`queued → fetching → downloading → extracting → parsing → embedding → summarizing → ready`
Any step failure after retries → `failed` (with step + reason recorded). If >60% of source files are unparseable → `partially_ready` with warning.

---

## 3. Pipeline step contracts (detailed)

### Step 0 — URL validation
- Accept only `https://github.com/{owner}/{repo}` (optional trailing `/` or `.git`); **reject any other host** (SSRF guard).
- Owner: 1–39 chars `[A-Za-z0-9-]`; repo: 1–100 chars `[A-Za-z0-9._-]`.
- Reject URLs containing credentials, query strings, or fragments.
- Normalize `owner/repo` to lowercase for identity; keep the original string for display.
- Client-side pre-check is UX only; server-side validation is authoritative.

### Step 1 — Create analysis
- **Idempotency:** if a `ready` analysis exists for the same repo + commit SHA and is <30 days old, return it (HTTP 200) instead of re-indexing.
- Otherwise create the row (status `queued`), after checking: rate limit (5/hour/IP), concurrency (≤3 running), and monthly budget.

### Step 2 — Fetch metadata
- `GET /repos/{owner}/{repo}` with the operator PAT → `default_branch`, `size`, primary language, `archived`.
- `GET /repos/{owner}/{repo}/commits/{default_branch}` → head commit SHA.
- Reject if the repo is empty or `size` > 50 MB (friendly message). Archived repos are allowed but noted.

### Step 3 — Download archive
- `GET https://codeload.github.com/{owner}/{repo}/tar.gz/{sha}` (allow-listed host), streamed with byte caps (compressed 50 MB / expanded 250 MB).
- Stored at `analyses/{analysis_id}/source.tar.gz` in the `repo-archives` bucket.

### Step 4 — Safe extract (guards)
- Reject absolute paths, `..` traversal, symlinks escaping the root, and hardlinks.
- Excluded: `node_modules/`, `vendor/`, `dist/`, `build/`, `.git/`, `target/`, `__pycache__/`, `*.min.js`, lockfiles, binaries (extension + NUL-byte sniff), files > 512 KB.
- Enforce ≤3,000 indexed files; record skipped counts and reasons in the analysis metadata.

### Step 5 — Parse (parser service contract)
**Input:** `{ analysis_id, files: [{ path, content_ref }] }` (signed URLs or mounted directory).
**Output (JSON):**
- `files[]`: path, language, size_bytes, line_count, parse_error_count, is_indexed
- `symbols[]`: path, name, kind (function | class | method | interface | type | const), start_line, end_line, parent_name?, exported?
- `imports[]`: from_path, raw_specifier, kind (internal | external | unknown), resolved_path?
- `chunks[]`: path, symbol_ref?, start_line, end_line, content_hash, token_estimate
- `meta`: parser_version, query_version, duration_ms
**Determinism:** same input + pinned versions ⇒ identical output.

### Step 6 — Embed
- Batches of ≤100 chunks; `gemini-embedding-001`; task type `RETRIEVAL_DOCUMENT` at index time, `RETRIEVAL_QUERY` for questions; 1536 dimensions.
- Batch API during indexing; per-batch retries; permanent failures go to the dead-letter queue and are reported in analysis health.
- Chunk text = symbol-bounded excerpt with a context header (path + symbol name + line range).

### Step 7 — Graph + overview
- Edges: resolve `raw_specifier` → in-repo path using language-specific rules; external imports marked `external`.
- Architecture summary: **one** Gemini call with a structured prompt; output validated; stored as an artifact.

### Step 8 — Finalize
- Compute counts (files, symbols, chunks, parse-error rate); set final status; delete the archive (or schedule deletion ≤1 h); emit a completion event.

---

## 4. Query path (Q&A with citations)

1. Validate question (non-empty, ≤2,000 chars).
2. Embed the question (`RETRIEVAL_QUERY`, 1536 dims).
3. **Hybrid retrieval:** `hybrid_search(analysis_id, query_text, query_embedding, …)` → RRF of full-text (`tsvector`) + vector (HNSW), top 12 chunks.
4. **Capped tool loop** (max 4 rounds, max 2 tool calls per round):
   - `hybrid_search(query, k)` — re-search with refined terms
   - `get_symbol(analysis_id, name)` — exact symbol lookup
   - `get_file_overview(path)` — file structure + symbol list
   - `get_import_neighbors(path, direction)` — dependency context
5. Compose answer under a strict output contract: `{ answer_markdown, citations[], facts[], hypotheses[], unknowns[] }`.
6. **Backend citation validation:** every citation must resolve to an indexed file with a valid line range; invalid citations are dropped and flagged; if no valid citations remain → the answer is downgraded to "insufficient evidence".
7. Persist conversation + message (question, answer, citations JSONB, model, token usage, latency).

---

## 5. Error / stack-trace investigation path

1. Accept stack-trace text (≤20,000 chars).
2. Parse frames: file paths, line numbers, function names (per-language patterns).
3. Map frames to the index: exact path → basename + fuzzy match → symbol name.
4. Retrieve surrounding chunks + import neighbours (2 hops).
5. LLM produces: likely cause, evidence chain (citations), hypotheses, next investigation steps.
6. Same citation validation and persistence rules as Q&A.

---

## 6. Briefing Pack export path

1. Assemble from the index: repo identity + SHA, overview, module map, top symbols, entry points, import-graph summary, key file excerpts, parse-health notes.
2. Generate `BRIEFING.md` + `SKILL.md` (the skill instructs consumer AIs to treat the pack as **data**, cite file:line, and abstain when evidence is absent).
3. Store in the `exports` bucket at `analyses/{analysis_id}/{export_id}/`; write metadata; return a signed URL (1-hour expiry).
4. Exports expire after 7 days and are deleted by the retention job.

---

## 7. Data dictionary (part 1 — core entities)

> Field types are Postgres-oriented design specifications. `enum` = Postgres enum type (or text + CHECK). All tables carry `created_at timestamptz NOT NULL DEFAULT now()` unless stated.

### 7.1 `workspaces` — anonymous session identity
| Field | Type | Constraints | Description |
|---|---|---|---|
| id | uuid | PK | Workspace identifier |
| token_hash | text | NOT NULL, UNIQUE | Hash of the anonymous bearer token (never the raw token) |
| last_seen_at | timestamptz | NOT NULL | Sliding activity timestamp for retention |
| created_at | timestamptz | NOT NULL | Creation time |

*No PII. Workspaces expire after 30 days of inactivity.*

### 7.2 `repositories` — canonical repo identity
| Field | Type | Constraints | Description |
|---|---|---|---|
| id | uuid | PK | Repository identifier |
| provider | text | NOT NULL DEFAULT 'github' | Future-proofing (GitHub only in MVP) |
| owner | text | NOT NULL, lowercase | GitHub owner/org |
| name | text | NOT NULL, lowercase | Repository name |
| full_name | text | NOT NULL, UNIQUE | `owner/name` normalized key |
| html_url | text | NOT NULL | Display URL |
| created_at | timestamptz | NOT NULL | First seen |

### 7.3 `analyses` — one indexed snapshot of a repository
| Field | Type | Constraints | Description |
|---|---|---|---|
| id | uuid | PK | Analysis identifier (used in all child tables and storage paths) |
| repository_id | uuid | FK → repositories, NOT NULL | Parent repository |
| workspace_id | uuid | FK → workspaces, NOT NULL | Who requested it |
| default_branch | text | NOT NULL | Branch recorded at fetch time |
| commit_sha | char(40) | NOT NULL, `^[0-9a-f]{40}$` | Immutable snapshot identity |
| status | enum | NOT NULL | `queued, fetching, downloading, extracting, parsing, embedding, summarizing, ready, partially_ready, failed` |
| status_detail | text | NULL | Human-readable current step detail |
| size_kb | integer | NULL, ≥0 | GitHub `size` at fetch time |
| file_count | integer | NOT NULL DEFAULT 0 | Total files seen |
| indexed_file_count | integer | NOT NULL DEFAULT 0 | Files that entered the index |
| symbol_count | integer | NOT NULL DEFAULT 0 | Extracted symbols |
| chunk_count | integer | NOT NULL DEFAULT 0 | Embedded chunks |
| parse_error_rate | numeric(5,4) | NULL | Share of unparseable source files (0–1) |
| parser_version | text | NULL | Pinned Tree-sitter/query version |
| query_version | text | NULL | Query-file bundle version |
| embedding_model | text | NULL | e.g. `gemini-embedding-001@1536` |
| llm_model | text | NULL | e.g. `gemini-3.8-flash` |
| skipped_files | jsonb | NULL | `{reason: count}` map for exclusions |
| error_step | text | NULL | Which step failed |
| error_message | text | NULL | Sanitized failure reason |
| started_at | timestamptz | NULL | Pipeline start |
| finished_at | timestamptz | NULL | Terminal state time |

**Idempotency:** partial unique index on `(repository_id, commit_sha)` where `status IN ('ready','partially_ready')` — one good snapshot per commit.

### 7.4 `analysis_files` — every file seen in the snapshot
| Field | Type | Constraints | Description |
|---|---|---|---|
| id | uuid | PK | File record id |
| analysis_id | uuid | FK → analyses, NOT NULL | Snapshot |
| path | text | NOT NULL, no leading `/`, no `..` | Repo-relative path |
| language | enum | NOT NULL | `python, go, javascript, typescript, other` |
| size_bytes | integer | NOT NULL, ≥0 | File size |
| line_count | integer | NULL, ≥0 | Lines (null for non-indexed) |
| parse_error_count | integer | NOT NULL DEFAULT 0 | ERROR nodes reported by the parser |
| is_indexed | boolean | NOT NULL DEFAULT false | Whether it entered symbols/chunks |
| exclusion_reason | text | NULL | `binary, too_large, vendored, generated, lockfile, …` |
| content_hash | text | NULL | sha256 of content (change detection, dedupe) |
| UNIQUE | — | `(analysis_id, path)` | One row per file per snapshot |

### 7.5 `code_symbols` — deterministic symbol extraction
| Field | Type | Constraints | Description |
|---|---|---|---|
| id | uuid | PK | Symbol id |
| analysis_id | uuid | FK → analyses, NOT NULL | Snapshot |
| file_id | uuid | FK → analysis_files, NOT NULL | Owning file |
| name | text | NOT NULL, 1–200 chars | Symbol name |
| qualified_name | text | NULL | e.g. `ClassName.method` |
| kind | enum | NOT NULL | `function, class, method, interface, type, const, other` |
| parent_symbol_id | uuid | FK → code_symbols, NULL | Containment (method in class) |
| start_line | integer | NOT NULL, ≥1 | First line (1-based) |
| end_line | integer | NOT NULL, ≥ start_line | Last line |
| is_exported | boolean | NOT NULL DEFAULT false | Exported/public per language rules |

**Indexes:** `(analysis_id, name)` for lookups; `(file_id)` for file views; `(analysis_id, kind)` for overviews.

### 7.6 `code_chunks` — embedded retrieval units
| Field | Type | Constraints | Description |
|---|---|---|---|
| id | uuid | PK | Chunk id |
| analysis_id | uuid | FK → analyses, NOT NULL | Snapshot |
| file_id | uuid | FK → analysis_files, NOT NULL | Owning file |
| symbol_id | uuid | FK → code_symbols, NULL | Symbol this chunk is bounded by (null for file-level chunks) |
| start_line | integer | NOT NULL, ≥1 | First line |
| end_line | integer | NOT NULL, ≥ start_line | Last line |
| content | text | NOT NULL, ≤32,000 chars | Chunk text with context header (path + symbol + lines) |
| content_hash | text | NOT NULL | sha256 for dedupe/change detection |
| token_estimate | integer | NOT NULL, ≥1 | Approximate token count |
| embedding | vector(1536) | NULL until embedded | pgvector embedding |
| tsv | tsvector | GENERATED | Full-text search vector over content |
| UNIQUE | — | `(analysis_id, file_id, start_line, end_line)` | No duplicate ranges |

**Indexes:** HNSW on `embedding` (cosine), GIN on `tsv`. Max 15,000 chunks per analysis (enforced at parse time).

### 7.7 `file_imports` — dependency graph edges
| Field | Type | Constraints | Description |
|---|---|---|---|
| id | uuid | PK | Edge id |
| analysis_id | uuid | FK → analyses, NOT NULL | Snapshot |
| from_file_id | uuid | FK → analysis_files, NOT NULL | Importing file |
| to_file_id | uuid | FK → analysis_files, NULL | Resolved in-repo target (null if external) |
| raw_specifier | text | NOT NULL, ≤500 chars | Import string as written |
| kind | enum | NOT NULL | `internal, external, unknown` |
| UNIQUE | — | `(analysis_id, from_file_id, raw_specifier)` | One edge per specifier per file |

### 7.8 `analysis_artifacts` — generated documents
| Field | Type | Constraints | Description |
|---|---|---|---|
| id | uuid | PK | Artifact id |
| analysis_id | uuid | FK → analyses, NOT NULL | Snapshot |
| kind | enum | NOT NULL | `architecture_summary, briefing_pack, parse_health, eval_report` |
| format | text | NOT NULL | `markdown` or `json` |
| content | text | NULL | Inline content for summaries |
| storage_path | text | NULL | For large artifacts stored in buckets |
| model | text | NULL | Model used, if LLM-generated |
| created_at | timestamptz | NOT NULL | Creation time |

### 7.9 `conversations` — one thread per Q&A or investigation
| Field | Type | Constraints | Description |
|---|---|---|---|
| id | uuid | PK | Conversation id |
| analysis_id | uuid | FK → analyses, NOT NULL | Snapshot context |
| workspace_id | uuid | FK → workspaces, NOT NULL | Owner |
| kind | enum | NOT NULL | `qa` or `investigation` |
| last_message_at | timestamptz | NOT NULL | For retention/listing |

### 7.10 `messages` — questions and answers with evidence
| Field | Type | Constraints | Description |
|---|---|---|---|
| id | uuid | PK | Message id |
| conversation_id | uuid | FK → conversations, NOT NULL | Thread |
| role | enum | NOT NULL | `user` or `assistant` |
| content | text | NOT NULL, ≤100,000 chars | Markdown answer or question |
| citations | jsonb | NULL | `[{file_path, start_line, end_line, symbol?, kind: fact\|hypothesis}]` — backend-validated |
| facts | jsonb | NULL | Explicit facts list |
| hypotheses | jsonb | NULL | Explicit hypotheses list |
| unknowns | jsonb | NULL | Insufficient-evidence notes |
| model | text | NULL | Model used |
| input_tokens | integer | NULL | Usage |
| output_tokens | integer | NULL | Usage |
| latency_ms | integer | NULL | Response time |

### 7.11 `exports` — Briefing Pack records
| Field | Type | Constraints | Description |
|---|---|---|---|
| id | uuid | PK | Export id |
| analysis_id | uuid | FK → analyses, NOT NULL | Snapshot |
| workspace_id | uuid | FK → workspaces, NOT NULL | Owner |
| storage_path | text | NOT NULL | `analyses/{analysis_id}/{export_id}/` |
| byte_size | integer | NOT NULL, ≥0 | Total size |
| expires_at | timestamptz | NOT NULL | created_at + 7 days |
| created_at | timestamptz | NOT NULL | Creation time |

### 7.12 Evaluation harness tables (internal deliverable)
| Table | Key fields | Purpose |
|---|---|---|
| `eval_benchmarks` | id, repo_full_name, pinned_commit_sha, language, notes | The 3 frozen benchmark repos (SHA pinned for reproducibility) |
| `eval_questions` | id, benchmark_id, category (`identifier\|conceptual\|holistic\|error_investigation\|trap`), question, expected jsonb (files/lines/symbols) | The 48-question bank |
| `eval_runs` | id, started_at, finished_at, config jsonb (retrieval mode, model, prompt version) | One harness execution |
| `eval_results` | id, run_id, question_id, answer, citations jsonb, metrics jsonb (citation_valid, correctness, abstained, retrieved_expected, latency_ms, cost_usd) | Per-question measurements → comparison report |

### 7.13 `job_events` — pipeline observability
| Field | Type | Constraints | Description |
|---|---|---|---|
| id | bigserial | PK | Event id |
| analysis_id | uuid | FK → analyses, NOT NULL | Snapshot |
| step | text | NOT NULL | Pipeline step name |
| event | enum | NOT NULL | `started, succeeded, retried, failed` |
| detail | jsonb | NULL | Retry count, error codes, durations |
| created_at | timestamptz | NOT NULL | Event time |

### 7.14 `jobs` — Postgres-backed pipeline queue
| Field | Type | Constraints | Description |
|---|---|---|---|
| id | uuid | PK | Job id |
| analysis_id | uuid | FK → analyses, NOT NULL | Owning snapshot |
| step | text | NOT NULL | Pipeline step name |
| status | enum | NOT NULL | `pending, running, succeeded, failed` |
| attempts | integer | NOT NULL DEFAULT 0 | Attempts so far |
| max_attempts | integer | NOT NULL | Retry budget per step |
| run_at | timestamptz | NOT NULL | Earliest execution time (backoff) |
| locked_by | text | NULL | Worker id holding the lease |
| locked_at | timestamptz | NULL | Lease timestamp (expiry 60 s) |
| last_error | text | NULL | Sanitized error |
| created_at | timestamptz | NOT NULL | Creation time |

**Indexes:** `(status, run_at)` for claim queries; unique `(analysis_id, step)` for idempotency.

**Retention hooks:** `workspaces` (30 days inactive) → cascades; `conversations`/`messages` (30 days); `exports` (7 days); archives in storage (≤1 h after success, 24 h after failure).

---

## 8. Database schema design (design only — no migrations)

### 8.1 Entity relationships

```
workspaces 1──* analyses *──1 repositories
analyses   1──* analysis_files 1──* code_symbols
analysis_files 1──* code_chunks *──0..1 code_symbols
analyses   1──* file_imports  (from_file_id, to_file_id → analysis_files)
analyses   1──* analysis_artifacts
analyses   1──* conversations 1──* messages
analyses   1──* exports
analyses   1──* job_events

eval_benchmarks 1──* eval_questions
eval_runs      1──* eval_results *──1 eval_questions
```

### 8.2 Constraint rules

- **Cascade:** every child FK cascades from `analyses` — deleting an analysis removes its entire index in one operation (retention job uses this).
- **CHECKs:** `start_line ≥ 1`, `end_line ≥ start_line`, `commit_sha ~ '^[0-9a-f]{40}$'`, `parse_error_rate BETWEEN 0 AND 1`, `byte_size ≥ 0`.
- **Idempotency:** partial unique index `analyses(repository_id, commit_sha) WHERE status IN ('ready','partially_ready')`.
- **Dedupe:** `analysis_files(analysis_id, path)` unique; `code_chunks(analysis_id, file_id, start_line, end_line)` unique; `file_imports(analysis_id, from_file_id, raw_specifier)` unique.
- **Limits enforced at write time:** ≤15,000 chunks, ≤3,000 indexed files per analysis (reject/trim with recorded reason, never silent).

### 8.3 The single search function (design specification)

```
hybrid_search(
  p_analysis_id      uuid,
  p_query_text       text,
  p_query_embedding  vector(1536),
  p_match_count      int   DEFAULT 12,
  p_full_text_weight float DEFAULT 1.0,
  p_semantic_weight  float DEFAULT 1.0,
  p_rrf_k            int   DEFAULT 50
) RETURNS TABLE(chunk_id, file_path, symbol_name, start_line, end_line,
                content, ft_rank, vec_rank, rrf_score)
```

- **CTE 1 — full-text:** `websearch_to_tsquery(p_query_text)` over `tsv`, ranked with `ts_rank_cd`.
- **CTE 2 — vector:** `ORDER BY embedding <=> p_query_embedding` directly on the distance (index-friendly), `LIMIT p_match_count × 2`.
- **Fusion:** reciprocal rank fusion `Σ weight / (rrf_k + rank)`, grouped by chunk id.
- **Scope:** always filtered by `p_analysis_id` (isolation enforced inside the function).
- **Ablations (evaluation):** vector-only = `full_text_weight 0`; hybrid = `1,1`; hybrid+graph = hybrid + post-retrieval import-neighbour expansion. No code forks needed.

### 8.4 RLS policy design

- **Deny-by-default:** RLS enabled on every table; **no policies** granted to `anon` or `authenticated`. Only the service role (used exclusively by the API Worker server-side) can read/write.
- **Rationale:** the browser never talks to Supabase directly — all access flows through the Worker, which enforces workspace scoping in application logic and rate limits. RLS is therefore defense-in-depth against key leakage and accidental public exposure.
- **Storage:** both buckets private, no anon policies; the Worker issues signed URLs (1-hour TTL) after ownership checks.
- **Secrets:** `service_role`/secret key and all provider keys live only in Cloudflare secrets (never in client bundles, never in the repo).

### 8.5 Storage bucket layout

| Bucket | Path pattern | Contents | Lifecycle |
|---|---|---|---|
| `repo-archives` (private) | `analyses/{analysis_id}/source.tar.gz` | Original tarball | Deleted ≤1 h after success; 24 h after failure |
| `exports` (private) | `analyses/{analysis_id}/{export_id}/BRIEFING.md`, `SKILL.md` | Briefing Pack | Deleted after 7 days |

### 8.5.1 Environment separation
- Supabase: separate **dev** and **prod** projects (no production data in dev).
- Cloudflare: separate dev/prod Workers + secrets; Containers and Queues namespaced per environment.

### 8.6 Migration policy (implementation phase only)
- Migrations are created **only after implementation approval**, generated with the Supabase CLI, reviewed before applying, and applied dev → prod in order. This document defines the target state; it does not create migrations.

---

## 9. API contracts (design)

### 9.1 Conventions
- Base path: `/api`. JSON only, UTF-8. Max request body 100 KB (investigations ≤25 KB).
- **Auth:** `Authorization: Bearer <workspace_token>` — anonymous token issued on first visit, stored **hashed** in `workspaces`.
- **Error envelope:** `{ "error": { "code": "...", "message": "...", "details": {...} } }`
  Codes: `INVALID_URL, REPO_NOT_FOUND, REPO_TOO_LARGE, RATE_LIMITED, BUDGET_EXCEEDED, ANALYSIS_NOT_READY, ANALYSIS_FAILED, NOT_FOUND, INVALID_INPUT, UPSTREAM_ERROR`
- **Rate-limit headers:** `X-RateLimit-Limit`, `X-RateLimit-Remaining`, `X-RateLimit-Reset`.
- **Idempotency:** `POST /analyses` returns the existing ready analysis (HTTP 200, `deduplicated: true`) when the same repo + SHA is already indexed.

### 9.2 Endpoints

| Method | Path | Request | Response (success) | Notes |
|---|---|---|---|---|
| POST | `/analyses` | `{ repo_url }` | 202 `{ analysis_id, status: "queued", deduplicated: false }` | 200 + `deduplicated: true` on reuse |
| GET | `/analyses/{id}` | — | `{ id, repository: { full_name }, default_branch, commit_sha, status, status_detail, counts: { files, indexed_files, symbols, chunks }, parse_error_rate, started_at, finished_at }` | Poll target for progress UI |
| GET | `/analyses/{id}/overview` | — | `{ summary_markdown, languages[], entry_points[], frameworks[], stats }` | From architecture artifact |
| GET | `/analyses/{id}/files` | `?path=&limit=&offset=` | `{ files: [{ path, language, size_bytes, line_count, is_indexed, exclusion_reason }] }` | Tree/list view |
| GET | `/analyses/{id}/files/{path}` | — | `{ path, language, symbols[], imports[] }` | File overview (no full content for non-indexed files) |
| GET | `/analyses/{id}/symbols` | `?q=&kind=&limit=` | `{ symbols: [{ name, kind, path, start_line, end_line }] }` | Symbol search |
| GET | `/analyses/{id}/graph` | `?depth=1` | `{ nodes: [{ path, language }], edges: [{ from, to, kind, specifier }] }` | Import-graph view |
| POST | `/analyses/{id}/questions` | `{ question }` | `{ message_id, answer_markdown, citations[], facts[], hypotheses[], unknowns[], usage: { input_tokens, output_tokens } }` | Citations backend-validated |
| POST | `/analyses/{id}/investigations` | `{ stack_trace }` | Same shape as question answer | Error-investigation mode |
| POST | `/analyses/{id}/exports` | `{}` | `{ export_id, download_url, expires_at }` | Signed URL, 1-hour TTL |
| GET | `/exports/{id}/download` | — | 302 → signed URL | Re-checks ownership + expiry |
| GET | `/health` | — | `{ status: "ok", version }` | Ops probe |

### 9.3 Status polling contract
`GET /analyses/{id}` returns the current `status` plus `status_detail`; the UI polls every 2 s while non-terminal, stops on `ready | partially_ready | failed`, and displays the step timeline from `job_events` when failed.

---

## 10. Job runner contracts (replaces Workflows)

The Go pipeline worker processes jobs from a Postgres-backed queue. Creating an analysis enqueues the first step; each successful step enqueues the next.

**Worker loop:** claim a due job with `FOR UPDATE SKIP LOCKED` → mark `running` + lease (60 s heartbeat) → execute → on success mark `succeeded` and enqueue the next step → on failure increment `attempts`, set backoff `run_at`, retry until `max_attempts` → mark the analysis `failed`.

**Crash recovery:** expired leases (`locked_at` older than 60 s) are reclaimed on the next tick; steps are idempotent, keyed on `(analysis_id, step)`.

| Step | Input | Output | Retries | Timeout | Failure effect |
|---|---|---|---|---|---|
| `fetch_metadata` | `{ repo_full_name }` | `{ default_branch, commit_sha, size_kb, language, archived }` | 3 | 30 s | `failed` (step: fetching) |
| `download_archive` | `{ repo_full_name, commit_sha, analysis_id }` | `{ storage_path, compressed_bytes }` | 3 | 120 s | `failed` (downloading) |
| `extract` | `{ analysis_id, storage_path }` | `{ manifest, file_count, skipped_counts }` | 1 | 120 s | `failed` (extracting) |
| `parse` | `{ analysis_id, manifest }` → Python parser service | `{ files, symbols, imports, chunks, parse_error_rate }` | 2 | 600 s | `failed` (parsing) |
| `embed` | `{ analysis_id }` | `{ embedded_count, failed_batches }` | 3 per batch | 900 s | `failed` (embedding) if any batch permanently fails |
| `graph_and_summary` | `{ analysis_id }` | `{ edge_count, artifact_id }` | 2 | 120 s | `failed` (summarizing) |
| `finalize` | `{ analysis_id }` | `{ status }` | 3 | 30 s | `failed` (finalizing) |

**Partial success rule:** if parsing completes but >60% of source files are unparseable, the pipeline still finalizes — status becomes `partially_ready` with a warning stored in `status_detail` and the `parse_health` artifact.

---

## 11. Validation rules (Zod) — design specification

> Illustrative design sketches only — **not implementation files**. **Validation split (v1.1):** Zod on the React frontend (UX only); Go struct validation is authoritative for the API; Pydantic is authoritative for the parser service. The contract shapes below are shared across all three.

### 11.1 Shared primitives

| Name | Rule |
|---|---|
| `RepoUrl` | string, ≤300 chars, must match `^https://github\.com/[A-Za-z0-9-]{1,39}/[A-Za-z0-9._-]{1,100}(\.git)?/?$` |
| `RepoFullName` | string, `^[a-z0-9-]{1,39}/[a-z0-9._-]{1,100}$` (normalized lowercase) |
| `CommitSha` | string, `^[0-9a-f]{40}$` |
| `RepoPath` | string, ≤500 chars, no leading `/`, no `..` segments, no `\` |
| `LineRange` | `{ start_line: int ≥1, end_line: int ≥ start_line }` |
| `Language` | enum: `python, go, javascript, typescript, other` |
| `AnalysisStatus` | enum: `queued, fetching, downloading, extracting, parsing, embedding, summarizing, ready, partially_ready, failed` |
| `SymbolKind` | enum: `function, class, method, interface, type, const, other` |

### 11.2 Request schemas

| Schema | Shape |
|---|---|
| `createAnalysisRequest` | `{ repo_url: RepoUrl }` |
| `questionRequest` | `{ question: string, 1–2,000 chars, trimmed }` |
| `investigationRequest` | `{ stack_trace: string, 1–20,000 chars }` |
| `exportRequest` | `{}` (strict empty object) |
| `pagination` | `{ limit: int 1–100 (default 50), offset: int ≥0 }` |

### 11.3 Internal contracts (between components)

**`parseResult`** (Container → Workflow):
- `files[]`: `{ path: RepoPath, language: Language, size_bytes: int ≥0, line_count: int ≥0, parse_error_count: int ≥0, is_indexed: boolean, exclusion_reason?: string }`
- `symbols[]`: `{ path, name: 1–200 chars, kind: SymbolKind, start_line, end_line, parent_name?: string, is_exported: boolean }`
- `imports[]`: `{ from_path: RepoPath, raw_specifier: 1–500 chars, kind: internal|external|unknown, resolved_path?: RepoPath }`
- `chunks[]`: `{ path, symbol_name?: string, start_line, end_line, content_hash: sha256, token_estimate: int ≥1 }`
- `meta`: `{ parser_version, query_version, duration_ms }`
- **Cross-field refinements:** `end_line ≥ start_line` everywhere; chunk ranges must fall within the owning file's `line_count`; `≤3,000` files and `≤15,000` chunks (reject the batch, never silently truncate).

**`embeddingBatch`** (Workflow → Queue): `{ analysis_id: uuid, items: [{ chunk_id: uuid, text: 1–32,000 chars }] }`, max 100 items.

**`citation`**: `{ file_path: RepoPath, start_line, end_line, symbol?: ≤200 chars, kind: fact | hypothesis }` — **must validate against the index** (file exists, is indexed, range within bounds) before display.

**`llmAnswer`**: `{ answer_markdown: 1–50,000 chars, citations: ≤12, facts: ≤20, hypotheses: ≤20, unknowns: ≤20 }` — every citation must reference a chunk that was actually retrieved; otherwise it is dropped and the answer may be downgraded to "insufficient evidence".

### 11.4 Validation policy
- Reject unknown fields; reject over-limit payloads before any processing.
- Cross-entity validation (citations, line ranges) happens after generation and degrades the answer rather than failing the request — honesty over confidence.

---

## 12. Cross-cutting design

### 12.1 Security
- **Untrusted content:** repository text is always wrapped in explicit data delimiters with a system rule that content is data, never instructions; no repo-derived tool calls; all model output validated.
- **SSRF:** strict allow-list (`github.com`, `codeload.github.com`); redirects to other hosts rejected.
- **No execution:** the pipeline never runs repository code, install scripts, or build tooling — archives are only read and parsed.
- **Secrets:** Supabase secret key, Gemini key, and GitHub PAT live only in the host environment (Docker secrets / env vars on the server); never in client bundles, never in the repository. Cloudflare holds only the static frontend.
- **PII:** none collected; anonymous workspace tokens stored hashed.

### 12.2 Observability
- `job_events` gives a per-analysis step timeline (started/succeeded/retried/failed with durations).
- Every message stores token usage and latency → cost tracking per conversation and per month.
- `/health` endpoint; budget counter with 50%/80% alerts and a hard stop at 100%.

### 12.3 Cost controls (enforced in design)
- ≤15,000 chunks per analysis; embedding batches of ≤100; **one** overview LLM call per analysis; ≤4 tool rounds and ≤2 tools/round per question; 50 questions/session/day; $20/month hard budget.

### 12.4 Idempotency & reproducibility
- Every analysis is pinned to a commit SHA; `parser_version`, `query_version`, `embedding_model`, and `llm_model` are recorded per analysis.
- Re-analysis of the same SHA reuses the existing ready analysis (200, `deduplicated: true`); a new analysis row is created only when the SHA or versions change.

---

## 13. Environment configuration (names only)

| Variable | Purpose |
|---|---|
| `SUPABASE_URL` | Supabase project URL |
| `SUPABASE_SECRET_KEY` | Server-side key (bypasses RLS; Worker only) |
| `GEMINI_API_KEY` | Gemini API access (operator-owned) |
| `GITHUB_PAT` | GitHub metadata access (operator-owned) |
| `APP_BASE_URL` | Canonical app URL (links, signed-URL callbacks) |
| `RATE_LIMIT_ANALYSES_PER_HOUR` | Default 5 |
| `RATE_LIMIT_QUESTIONS_PER_DAY` | Default 50 |
| `MONTHLY_BUDGET_USD` | Default 20 |
| `ENVIRONMENT` | `dev` or `prod` |

All values live in the host environment (Docker secrets or env vars) per environment. None are ever exposed to the browser.

---

## 14. Open design questions (small, non-blocking)

| # | Question | Recommendation |
|---|---|---|
| 1 | How is the anonymous workspace token delivered? | **httpOnly signed cookie** (safer than localStorage; no XSS token theft) |
| 2 | Briefing Pack packaging: two loose files or one archive? | **Single `.zip`** containing `BRIEFING.md` + `SKILL.md` (one download, clear instructions) |
| 3 | Where does the retention job run? | **Go worker cron** (internal scheduler) invoking a cleanup routine |
| 4 | Evaluation tables location | **Same Supabase project, separate `eval` schema** (isolated from app data) |

These do not block design approval; they are recorded for the implementation plan.

---

## 15. Traceability (design → approved decisions)

| Design area | Traces to |
|---|---|
| Layered hybrid pipeline, capped tool loop, deterministic-first | Research 1 decisions (v0.6) |
| Tree-sitter engine, query files, partial-tree policy, import graph | Research 2 decisions (v0.6) |
| Relational + JSONB, pgvector, single RRF function, buckets, RLS | Research 3 decisions (v0.7) |
| Go job runner + Python parser topology (replaces Workflows/Queues/Container) | Research 4 revised (v1.2) |
| Supabase buckets as CDN for archives/exports; Cloudflare DNS + static hosting | Backend revision (v1.2) |
| `gemini-3.8-flash`, `gemini-embedding-001@1536`, paid tier, operator key | Research 5 decisions (v0.7) |
| PAT + codeload + default_branch + SHA pin + submodules/LFS exclusion | Research 6 decisions (v0.7) |
| Journey, states, citations panel, fact/hypothesis labels | Research 7 decisions (v0.7) |
| Evaluation tables, ablations via function parameters, 48-question bank | Research 8 decisions (v0.7) |
| Limits, retention, budgets, benchmark/demo repos, frozen features | Scope freeze §6.2 (v0.8) |

---

## 16. Scope statement

This document is **design only**. It contains no implementation files, no migrations, and no deployment configuration. Implementation begins only after Smit's explicit approval, per the standing implementation gate in the project handoff.

**Status (v1.1):** implementation approved (handoff v1.1) and the backend revised (handoff v1.2). Next: build the vertical slice per this design; the 48-question evaluation bank is approved and serves as the regression suite.