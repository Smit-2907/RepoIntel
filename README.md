# AI Codebase Intelligence & Engineering Assistant

Paste a public GitHub URL → deterministic Tree-sitter index → evidence-grounded Q&A with file:line citations.

**Stack (approved):** Go HTTP API + Go pipeline worker · Python parser service (py-tree-sitter) · React + Vite + Tailwind · Supabase Postgres + pgvector + Storage (CDN-backed) · Gemini 3.8 Flash + embedding-001 · Cloudflare for DNS/static hosting.

**Docs:** all planning + design docs live in [`docs/`](./docs) — start with `AI_Codebase_Intelligence_Project_Handoff_v1.2.md` and `Architecture_and_Data_Design_v1.1.md`.

## Quick start (local, no API keys required)

```bash
cp .env.example .env          # optional: add GEMINI_API_KEY and GITHUB_PAT
docker compose up --build     # db + parser + api + worker
pnpm install
pnpm dev:web                  # web UI on http://localhost:5173
```

Open the UI, paste a repository URL (e.g. `https://github.com/pallets/flask`), watch the pipeline states, then ask a question and inspect the citations.

Without `GEMINI_API_KEY` the app runs in **stub mode**: deterministic local embeddings (hash-based, dev only) and an extractive answer built from the top retrieved chunks — the full pipeline, retrieval, and citation validation still run.

## Services

| Service | Port | Notes |
|---|---|---|
| `db` | 54322 | `pgvector/pgvector:pg16`; migrations auto-applied on first start |
| `parser` | 8090 | Python + FastAPI + py-tree-sitter (Python, Go, JS, TS) |
| `api` | 8080 | Go HTTP API (`/api/*`, `/health`) |
| `worker` | — | Go pipeline worker (jobs table: fetch → download → extract → parse → embed → graph+summary → finalize) |
| `web` | 5173 | Vite dev server (Cloudflare static hosting in production) |

## Repository layout

```
apps/api/        Go: HTTP API + pipeline worker + job runner
apps/parser/     Python: py-tree-sitter parsing service
apps/web/        React + Vite + Tailwind UI
packages/core/   Shared TypeScript types + Zod schemas (frontend validation)
supabase/        SQL migrations (schema + hybrid_search function)
docs/            Handoff, architecture design, evaluation bank
```

## API (vertical slice)

```
POST /api/analyses                 { "repo_url": "https://github.com/owner/repo" }
GET  /api/analyses/{id}            status + counts (poll while running)
POST /api/analyses/{id}/questions  { "question": "..." }  → answer + validated citations
GET  /health
```

## Design guardrails (from the approved design)

- Deterministic facts before any LLM call; repository content is **data, never instructions**.
- Every citation is validated against the index (file exists, indexed, line range in bounds).
- Analyses are pinned to the default-branch head commit SHA; parser/model versions recorded.
- Limits: 50 MB repo, 3,000 files, 512 KB/file, 15,000 chunks, 10-minute analysis cap.
- No repository code is ever executed — archives are read and parsed only.

## Tests

```bash
cd apps/api && go test ./...
cd apps/parser && pytest -q
cd apps/web && pnpm typecheck
```
