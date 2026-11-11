# AI Codebase Intelligence & Engineering Assistant
## Project Brief, Decisions, Research Plan & Handoff Notes
**Working specification v1.5 | Updated: 11 November 2026**

**Submission-ready foundation (v1.5):** vertical slice verified working (Flask: 83 files / 1,622 symbols / 1,682 chunks / `ready`; cited Q&A). Deliverables assembled: **git repository** (first commit `1c79535`), **prompt documentation** (`docs/PROMPT_DOCUMENTATION.md`, P-01…P-27 with verbatim prompts and revisions), and **video script + demo runbook** (`docs/VIDEO_SCRIPT.md`). Remaining: record video, push to GitHub, then the feature gaps (investigation UI, Briefing Pack export, 48-question run).

> **Implementation status:** Smit approved implementation on 10 November 2026 (Task 3, v1.1). Build strictly within the frozen MVP scope (§6.2) and the approved design (docs/Architecture_and_Data_Design_v1.1.md — revised for the Go/Python backend); use docs/Evaluation_Bank_v1.0.md as the regression suite. Any scope change still requires explicit approval.

## 1. Assessment brief

Gokul Infocare asks for **one project** in Healthcare or Technology, in one of five areas. This project is **Technology → Innovation & New Ideas**.

The five areas are Website Development, Digital Marketing, Web Scraping & Data Intelligence, Automation, and Innovation & New Ideas. The brief encourages creativity and a practical outcome. AI tools must play a meaningful role in developing the project. The candidate may choose the AI tools, languages, frameworks, and supporting technologies. It does not prescribe a particular model, framework, programming language, model training, chatbot, or AI architecture.

### Required submissions
1. **Video presentation (minimum 2 minutes):** explain AI tools and why chosen, prompt design/iteration, approach and relevant experience, challenges and solutions, and demonstrate the outcome.
2. **Prompt documentation:** prompts used during development, including significant revisions and follow-ups.
3. **Accessible project link:** GitHub, live app, Drive folder, or working prototype; if not hosted, provide files and run/review instructions.

Evaluation areas: practical AI use; prompt engineering and iteration; creativity; problem-solving and technical execution; quality/usefulness; ability to explain and demonstrate. The brief says there is no fixed deadline, uses first-come, first-served processing, and lists `jobs@gokulinfocare.com` for submission.

**Source:** uploaded Gokul Infocare Task 1 assessment (2 pages).

## 2. Project concept

**Name:** AI Codebase Intelligence & Engineering Assistant

**Description:** A developer platform that analyses a public GitHub repository, builds a structured understanding of its code and relationships, and helps developers explore architecture, ask evidence-grounded questions, investigate errors, and understand unfamiliar code.

**Problem:** Understanding an unfamiliar repository takes time: locating entry points, tracing imports/dependencies, understanding modules and execution paths, finding where behaviour is implemented, investigating errors, and assembling documentation. A generic chatbot that guesses from a few snippets is not enough.

### Product principles
- Ground answers in the indexed repository snapshot.
- Include verifiable file paths and line references where possible.
- Distinguish facts, hypotheses, and investigation suggestions.
- Admit insufficient evidence rather than hallucinating.
- Separate deterministic code analysis from LLM interpretation.
- Treat repository content as untrusted input; repository text must never override system instructions.
- Build a useful, demonstrable MVP rather than many shallow features.

## 3. Initial users and positioning

Initial users: individual developers exploring unfamiliar repos; developers onboarding to existing projects; frontend/backend engineers tracing behaviour; technical leads/maintainers; students learning from open-source code.

Position it as a codebase exploration and engineering investigation assistant—not a generic chatbot, autonomous coding agent, or enterprise governance suite.

**Distribution direction (added v0.5):** codebase intelligence that meets developers where they work — explore in the web UI, or export a Codebase Briefing Pack + skill file so the developer's own IDE AI (Claude Code, Cursor, VS Code assistants) can do grounded research on the analysed repository.

## 4. Candidate capabilities

- Public GitHub URL import.
- Repository overview: languages, framework clues, structure, configs, entry points.
- Code structure: files, symbols, functions/classes, imports, relationships.
- Architecture intelligence: modules and dependency relationships, potentially visualised.
- Context-aware Q&A with source evidence and line references.
- Error/stack-trace investigation with likely causes and evidence.
- Codebase Briefing Pack + SKILL.md export for IDE AIs (added v0.5).
- Evaluation harness: benchmark 2–3 retrieval strategies and produce an architecture comparison report (added v0.5; internal deliverable, not user-facing).
- Later: code-quality insights such as possible cycles, duplication, maintainability concerns; documentation generation (demoted v0.5); live MCP server exposing the index as IDE tools (deferred v0.5).

## 5. Release boundaries (MVP scope frozen 10 Nov 2026 — see §6.2)

> **Scope frozen 10 Nov 2026 (v0.8).** The MVP feature list, repository limits, retention, budgets, benchmark repositories, and demo repository are fixed in §6.2. Any change requires explicit approval.

**MVP candidates:** import a public repository; analyse one branch/snapshot and record commit SHA; produce repository overview; extract basic file/symbol/import data for the four target languages; generate architecture overview; answer questions with evidence; investigate a supplied error/stack trace; export the Codebase Briefing Pack + SKILL.md.

**Project deliverable (internal, not user-facing):** evaluation harness comparing 2–3 retrieval strategies on benchmark repositories; produces a comparison report (citation correctness, answer correctness, abstention, retrieval coverage, latency, cost) that doubles as assessment evidence.

**Later candidates:** documentation generation (demoted from MVP, v0.5); richer docs/onboarding; deeper dependency/execution tracing; code-quality analysis; more languages; branch/commit comparison; private repositories/auth; collaboration/team workspaces; code modifications or PR workflows; live MCP server exposing the index as IDE tools (deferred v0.5); user-facing architecture-comparison toggle (deferred v0.5). Code editing/PR creation is deliberately outside initial scope unless reconsidered.

## 6. Confirmed decisions

| Area | Decision | Status |
|---|---|---|
| Assessment category | Technology → Innovation & New Ideas | Confirmed |
| Project concept | AI Codebase Intelligence & Engineering Assistant | Confirmed |
| Repositories | Public GitHub repositories only initially | Confirmed |
| User flow | Paste GitHub URL; no GitHub OAuth initially | Confirmed |
| Languages | Python, Go, JavaScript, TypeScript | Confirmed |
| AI provider | Google Gemini family initially; exact model open | Provider confirmed |
| Database | Supabase database (PostgreSQL) | Confirmed 10 Nov 2026 |
| Object/file storage | Supabase Storage buckets | Confirmed 10 Nov 2026 |
| Deployment | Mostly Cloudflare | **Superseded v1.2** — Cloudflare = DNS + static frontend hosting; backend runs on a host (Go/Python) |
| Research process | One research topic at a time | Confirmed preference |
| Implementation | Wait for explicit approval | Standing instruction |
| Evaluation report | Evaluation harness comparing 2–3 retrieval strategies; comparison report is a project deliverable | Confirmed 10 Nov 2026 (v0.5) |
| IDE skill export | Codebase Briefing Pack + SKILL.md export is in the MVP | Confirmed 10 Nov 2026 (v0.5) |
| MCP server | Live MCP server deferred to Later | Confirmed 10 Nov 2026 (v0.5) |
| Architecture toggle | User-facing run-all-architectures comparison deferred to Later | Confirmed 10 Nov 2026 (v0.5) |
| Documentation generation | Demoted from MVP to Later | Confirmed 10 Nov 2026 (v0.5) |
| Working architecture | Layered hybrid: deterministic extraction → hybrid retrieval → capped Gemini tool-calling; GraphRAG and free-roaming agents deferred to post-MVP | Confirmed 10 Nov 2026 (v0.6) |
| Retrieval mechanism | Supabase hybrid search: Postgres full-text + pgvector with reciprocal rank fusion | Confirmed 10 Nov 2026 (v0.6) |
| Agent design | Constrained tool-calling: small fixed toolkit (hybrid_search, get_symbol, get_file_overview, get_import_neighbors) with capped rounds | Confirmed 10 Nov 2026 (v0.6) |
| Code graph | Deterministic import/file graph in Postgres + one architecture summary per repo | Confirmed 10 Nov 2026 (v0.6) |
| Parsing engine | Tree-sitter as the single deterministic extraction engine for Python, Go, JavaScript, TypeScript | Confirmed 10 Nov 2026 (v0.6) |
| Parser configuration | Per-language Tree-sitter query files; grammar/parser versions pinned and recorded per analysis | Confirmed 10 Nov 2026 (v0.6) |
| Parse-error policy | Files with syntax errors are still indexed via partial trees; per-file error counts are stored | Confirmed 10 Nov 2026 (v0.6) |
| Cross-file references | Heuristic import/file graph for MVP; precise name resolution (LSP/SCIP) deferred; stack-graphs excluded (archived upstream) | Confirmed 10 Nov 2026 (v0.6) |
| Parsing placement | Cloudflare Container (Node + Tree-sitter) in the indexing pipeline; subject to Research 4 topology confirmation | **Superseded v1.2** — parser runs as a Python service (py-tree-sitter) |
| Native parsers / LSP | Rejected for MVP (multi-runtime sprawl, dependency-install security conflicts, uneven error tolerance) | Confirmed 10 Nov 2026 (v0.6) |
| Supabase data model | Relational tables for core entities; JSONB for flexible extraction details; import graph as relational tables; no graph DB | Confirmed 10 Nov 2026 (v0.7) |
| Supabase retrieval | pgvector embeddings on chunk rows (HNSW index); all search through one hybrid RRF SQL function | Confirmed 10 Nov 2026 (v0.7) |
| Supabase storage | Two private buckets: repo-archives (deleted after indexing) and exports (signed URLs); RLS on all tables; service keys server-side only | Confirmed 10 Nov 2026 (v0.7) |
| Retention direction | Archives deleted post-indexing; documented retention/deletion policy for analyses and conversations (durations at scope freeze) | Confirmed 10 Nov 2026 (v0.7) |
| Cloudflare topology | Workflows orchestrate the analysis pipeline; Queues fan out embedding batches with dead-letter handling; Container (Node + Tree-sitter) parses; Workers Paid plan required | **Superseded v1.2** — replaced by Go job runner + Python parser service (see Backend stack row) |
| Gemini working model | gemini-3.8-flash for all MVP LLM calls (overview, Q&A, error investigation); Gemini 4 Argon deferred as possible later deep mode | Confirmed 10 Nov 2026 (v0.7) |
| Gemini embeddings | gemini-embedding-001 at 1536 dimensions; Batch API during indexing; embedding and LLM model versions recorded per analysis | Confirmed 10 Nov 2026 (v0.7) |
| Gemini tier and privacy | Paid tier only (content not used to improve Google products); free tier never used for user data | Confirmed 10 Nov 2026 (v0.7) |
| Gemini credentials | Operator's server-side API key in Cloudflare secrets; per-session quotas and a monthly model budget; never exposed client-side | Confirmed 10 Nov 2026 (v0.7) |
| GitHub ingestion | REST metadata (operator PAT) + codeload tarball; snapshot pinned to the default branch head commit SHA; default branch = the API default_branch field | Confirmed 10 Nov 2026 (v0.7) |
| GitHub exclusions | Submodules and LFS contents not fetched in MVP (presence recorded); vendored/build/minified/lockfiles/binary excluded from indexing | Confirmed 10 Nov 2026 (v0.7) |
| GitHub safety | Strict URL allow-list (github.com / codeload.github.com); streaming extraction with bomb/traversal/symlink guards; repo size and file limits enforced (numbers at scope freeze) | Confirmed 10 Nov 2026 (v0.7) |
| Product workflow | Journey and job states as proposed (§9 Research 7); citations panel; fact/hypothesis/unknown labeling mandatory in answers | Confirmed 10 Nov 2026 (v0.7) |
| Demo scenario | Analyze a mid-size public repo, evidence-grounded Q&A, seeded stack-trace investigation, briefing-pack export into Claude Code; final repos at scope freeze | Confirmed 10 Nov 2026 (v0.7) |
| Evaluation benchmark | Three repos; four question categories plus trap questions; metrics: citation validity, correctness, abstention, retrieval coverage, latency, cost; three ablation configs | Confirmed 10 Nov 2026 (v0.7) |
| Evaluation success bar | Targets: ≥80% valid citations, ≥70% correct answers, 100% trap-question abstention, cents-level cost per question (finalized at scope freeze) | Confirmed 10 Nov 2026 (v0.7) |
| Repository limits | Max repo 50 MB (API size), extracted 250 MB, 3,000 files, 512 KB/file, 15,000 chunks, 10-minute analysis cap; >60% unparseable → partially ready | Confirmed 10 Nov 2026 (v0.8) |
| Retention durations | Archives deleted ≤1 h after success (24 h after failure); analyses 30 days; conversations 30 days; exports 7 days with 1-hour signed URLs | Confirmed 10 Nov 2026 (v0.8) |
| Usage and cost guards | 5 analyses/hour/IP, 3 concurrent analyses, 50 questions/session/day, $20/month Gemini budget with 50%/80% alerts and hard stop | Confirmed 10 Nov 2026 (v0.8) |
| Benchmark repositories | psf/requests (Python), honojs/hono (TypeScript), gin-gonic/gin (Go, master); backups: expressjs/express, spf13/cobra, koajs/koa | Confirmed 10 Nov 2026 (v0.8) |
| Demo repository | pallets/flask primary (Python, main); expressjs/express backup | Confirmed 10 Nov 2026 (v0.8) |
| Evaluation question bank | 48 questions: 10 identifier + 10 conceptual + 10 holistic + 10 error-investigation + 8 trap questions | Confirmed 10 Nov 2026 (v0.8) |
| Frozen MVP features | The nine-item frozen feature list in §6.2; everything else explicitly out of MVP | Confirmed 10 Nov 2026 (v0.8) |
| Architecture design | End-to-end architecture & data design v1.0 approved (docs/Architecture_and_Data_Design_v1.0.md) | Confirmed 10 Nov 2026 (v0.9) |
| Design resolutions | httpOnly cookie token; single-zip export; cron-Worker retention; separate eval schema; deny-by-default RLS with Worker-only DB access | Confirmed 10 Nov 2026 (v0.9) |
| Evaluation bank | 48-question bank approved (docs/Evaluation_Bank_v1.0.md): 10 identifier, 10 conceptual, 10 holistic, 10 error-investigation, 8 trap; rubric, success bar, run protocol, report shape | Confirmed 10 Nov 2026 (v1.0) |
| Implementation approval | Smit approved implementation 10 Nov 2026; vertical-slice-first build order; strict adherence to frozen scope and approved design | Confirmed 10 Nov 2026 (v1.1) |
| Backend stack | Go HTTP API + Go pipeline worker (Postgres-backed job queue; retries recorded in job_events) + Python parser service (py-tree-sitter) on one host via Docker Compose | Confirmed 10 Nov 2026 (v1.2) |
| Hosting and CDN | Cloudflare = DNS + static frontend hosting; Supabase Storage buckets (CDN-backed) serve archives and exports; Supabase Postgres + pgvector for data | Confirmed 10 Nov 2026 (v1.2) |
| Validation split | Zod on the React frontend (UX only); Go struct validation + Pydantic on the Python parser service (authoritative) | Confirmed 10 Nov 2026 (v1.2) |

### 6.1 Approvals log

All five Research 1 decisions, all six Research 2 decisions, and all Research 3–8 decisions were approved by Smit on 10 November 2026 (v0.6–v0.7), the MVP scope freeze (§6.2) was approved on 10 November 2026 (v0.8), the architecture & data design v1.0 (§6.3) was approved on 10 November 2026 (v0.9), and the 48-question evaluation bank (§6.4) was approved on 10 November 2026 (v1.0). The backend stack was revised on 10 November 2026 (v1.2): Go + Python services on one host, Cloudflare reduced to DNS/static hosting, Supabase buckets (CDN) for files. All are recorded above. No recommendations are currently awaiting approval.

### 6.2 Scope freeze (approved 10 Nov 2026, v0.8)

**Frozen MVP features (nothing else ships in MVP):**
1. Paste public GitHub URL → analysis pinned to the default-branch head SHA.
2. Repository overview (languages, structure, entry points, framework clues).
3. Deterministic index: files, symbols, imports, chunks for the four languages.
4. Architecture summary + import-graph view.
5. Evidence-grounded Q&A with citations panel (facts/hypotheses/unknowns).
6. Error/stack-trace investigation mode.
7. Codebase Briefing Pack + SKILL.md export.
8. Evaluation harness + comparison report (internal deliverable).
9. Full progress and failure states.

**Repository limits:** max repository size 50 MB (GitHub `size` field), max extracted 250 MB, max 3,000 files indexed, max single file 512 KB (skipped files recorded), max 15,000 chunks/embeddings, 10-minute analysis cap, >60% unparseable files → "partially ready" with a prominent warning.

**Retention:** archives deleted within 1 hour after successful indexing (24 hours after failure, for debugging); analyses and index data 30 days from creation; conversations 30 days; Briefing Pack exports 7 days with 1-hour signed URLs; anonymous workspace only — no accounts, no PII.

**Usage and cost guards:** 5 analyses per hour per IP, max 3 concurrent analyses, 50 questions per session per day, $20/month Gemini budget with alerts at 50%/80% and a hard stop at 100%.

**Benchmark repositories (evaluation harness):** psf/requests (Python, `main`, ~13 MB), honojs/hono (TypeScript, `main`, ~10 MB), gin-gonic/gin (Go, `master`, ~10 MB); backups: expressjs/express (JS, `master`), spf13/cobra (Go, `main`), koajs/koa (JS, `master`). All verified live via the GitHub API on 10 Nov 2026 — public, active, non-archived, 2–14 MB.

**Demo repository:** pallets/flask (Python, `main`, ~12 MB, 75k stars) primary; expressjs/express backup.

**Evaluation question bank:** 48 questions — 10 identifier-lookup, 10 conceptual, 10 holistic-architecture, 10 error-investigation, plus 8 trap questions (nonexistent targets → must abstain).

**Evidence note:** 4 of the 8 verified candidates default to `master` (express, gin, koa, chi), validating the default_branch decision from Research 6 — hard-coding `main` would have broken half the candidate set.

**Explicitly out of MVP:** documentation generation, live MCP server, architecture toggle, code-quality analysis, languages beyond the four, private repositories, branch/commit comparison, collaboration, code editing/PRs.

### 6.3 Design document (approved 10 Nov 2026, v0.9)

The approved end-to-end design lives in `docs/Architecture_and_Data_Design_v1.1.md` (revised for the Go/Python backend): system components and flow, canonical pipeline with step contracts, query/investigation/export paths, data dictionary (14 entities incl. `jobs`), schema design (constraints, single `hybrid_search` function, RLS, buckets), API contracts (12 endpoints), job-runner contracts (7 steps), validation split (Zod frontend / Go / Pydantic), cross-cutting design, and traceability to approved decisions.

Four design questions were resolved on approval:
1. Anonymous workspace token delivered via an **httpOnly signed cookie**.
2. Briefing Pack packaged as a **single .zip** (`BRIEFING.md` + `SKILL.md`).
3. Retention job runs as a **scheduled Worker (cron)**.
4. Evaluation tables live in a **separate `eval` schema** in the same Supabase project.

No migrations or code are created by the design; implementation remains gated on explicit approval.

### 6.4 Evaluation bank (approved 10 Nov 2026, v1.0)

The approved 48-question evaluation bank lives in `docs/Evaluation_Bank_v1.0.md`: 10 identifier-lookup, 10 conceptual, 10 holistic-architecture, 10 error-investigation, and 8 trap questions across psf/requests, honojs/hono, and gin-gonic/gin; plus the 0–3 grading rubric, success bar, run protocol, and expected comparison-report shape. Expected evidence is file/symbol-level; line numbers are pinned when benchmark commits are frozen at implementation time.

### 6.5 Backend revision (approved 10 Nov 2026, v1.2)

The Cloudflare-native compute topology (Workers + Workflows + Queues + Container) is **superseded**. The approved backend:

- **Go HTTP API** — all `/api/*` endpoints, validation, rate limits, Gemini orchestration.
- **Go pipeline worker** — the indexing pipeline as a Postgres-backed job queue (jobs table + worker loop; retries recorded in `job_events`), replacing Workflows/Queues.
- **Python parser service** — Tree-sitter via `py-tree-sitter`; exposes an internal HTTP endpoint (the parse contract in the design doc), replacing the Node container.
- **One host** running all services via Docker Compose (e.g. Fly.io/Railway/VPS); local dev identical.
- **Cloudflare** — DNS + static frontend hosting only.
- **Supabase** — Postgres + pgvector + Storage; buckets are CDN-backed and serve archives and exports directly (no extra CDN layer).
- **Validation split** — Zod stays on the React frontend for UX; server-side validation is authoritative in Go (structs/validator) and Pydantic in the parser service.

Unchanged: data dictionary, schema design, `hybrid_search` function, API contracts, retention, budgets, benchmark repos, demo repo, evaluation bank.

## 7. Clarifications and open interpretations

- **"Main branch only":** RESOLVED v0.7 — analysis targets the repository's `default_branch` (from the GitHub API), never a hard-coded branch name.
- **Repository limits:** RESOLVED v0.8 — max repo 50 MB, extracted 250 MB, 3,000 files, 512 KB/file, 15,000 chunks, 10-minute analysis cap (see §6.2).
- **Gemini credentials:** RESOLVED v0.7 — the operator's server-side key, stored in Cloudflare secrets, with per-session quotas and a monthly budget; never exposed in browser code.
- **Snapshot identity:** direction confirmed v0.7 — store repository URL, default branch, head commit SHA, timestamps, analysis status, parser/index versions, and LLM/embedding model versions.
- **Retention/privacy:** RESOLVED v0.8 — archives deleted ≤1 h after success (24 h after failure), analyses and conversations 30 days, exports 7 days, anonymous workspace only (see §6.2).
- **Supabase layout:** direction confirmed v0.7 — relational core + JSONB extras, two private buckets, RLS everywhere; detailed schemas are the next design task.
- **Cloudflare:** RESOLVED v0.7 — Workflows orchestrate the pipeline, Queues fan out embedding work, and parsing runs in a Container; Workers only orchestrate.

## 8. Research plan — one topic at a time

1. **AI/code-intelligence architecture:** compare basic RAG, syntax-aware/static analysis, hybrid analysis + retrieval, graph-based intelligence, and agent/tool-calling. Determine a reliable, maintainable combination. **COMPLETED 10 Nov 2026 — see Research log (§9); all five decisions approved 10 Nov 2026 (v0.6).**
2. **Parsing/language intelligence:** compare Tree-sitter, native parsers/compiler tooling, and LSP for Python, Go, JS, TS; assess symbols, imports, cross-file references, error tolerance, and deployment fit. **COMPLETED 10 Nov 2026 — see Research log (§9); all six decisions approved 10 Nov 2026 (v0.6).**
3. **Supabase data/storage design:** use Supabase PostgreSQL + Storage buckets. Compare relational metadata, JSONB, pgvector/retrieval choices, bucket layout, access controls, retention, and whether a graph DB is needed. **COMPLETED 10 Nov 2026 — see Research log (§9); decisions approved (v0.7).**
4. **Cloudflare deployment:** evaluate Workers, Queues, Workflows, and compute options for heavier parsing; map Cloudflare and Supabase responsibilities. **COMPLETED 10 Nov 2026 — see Research log (§9); decisions approved (v0.7).**
5. **Gemini model selection:** compare current models for coding/reasoning, context, structured output, tool calling, latency, pricing, limits, and data handling. **COMPLETED 10 Nov 2026 — see Research log (§9); decisions approved (v0.7).**
6. **GitHub ingestion/safety:** compare archive/API/raw access; assess rate limits, large repos, submodules/LFS, generated/binary files, URL validation, SSRF, extraction safety, prompt injection, and reproducibility. **COMPLETED 10 Nov 2026 — see Research log (§9); decisions approved (v0.7).**
7. **Product workflow:** define URL-to-results journey, progress/status, failure states, evidence display, briefing-pack export, and first demo scenario. **COMPLETED 10 Nov 2026 — see Research log (§9); decisions approved (v0.7).**
8. **Evaluation:** create benchmark questions and bug cases; measure citation correctness, answer correctness, abstention, retrieval coverage, latency, and cost. Feeds the confirmed evaluation-harness deliverable. **COMPLETED 10 Nov 2026 — see Research log (§9); decisions approved (v0.7).**

Each research topic should end with: evidence/sources, options comparison, recommendation, trade-offs, cost/operational implications, and decisions for Smit to approve. Do not silently convert recommendations into confirmed decisions.

## 9. Research log

### Research 1 — AI/code-intelligence architecture (completed 10 Nov 2026; decisions approved 10 Nov 2026, v0.6)

**Sources consulted (current official):** Tree-sitter documentation; Microsoft GraphRAG documentation; Google Gemini API docs (Interactions API is now the default interface); CodeQL documentation; LSP/LSIF specification; Supabase hybrid search guide; Anthropic context-engineering guidance.

**Options compared:**
- **A. Basic RAG** (fixed chunks → embeddings → vector search): simplest; but fixed-size chunking breaks functions/classes, weak on exact code identifiers, no structural knowledge, approximate citations.
- **B. Syntax-aware/static analysis only** (Tree-sitter/CodeQL/LSP): Tree-sitter is fast, error-tolerant, dependency-free, with upstream parsers for all four target languages; CodeQL requires buildable codebases (unsuitable for arbitrary public repos); LSP servers are stateful and heavy per repo. Alone, cannot answer natural-language how/why questions.
- **C. Hybrid analysis + retrieval:** Tree-sitter extraction → AST/symbol-aware chunking → Postgres full-text + pgvector → Supabase hybrid search (reciprocal rank fusion) → Gemini answers with citations validated against the deterministic index.
- **D. GraphRAG-style LLM knowledge graphs:** strongest on holistic/connect-the-dots questions, but LLM-driven graph extraction is expensive, slow, prompt-tuning heavy, and weakens verifiability. Retain the insight: store a deterministic file/import graph and generate indexed architecture summaries instead.
- **E. Agent/tool-calling:** best for error investigation and follow-ups; but unpredictable cost/latency, larger prompt-injection surface, harder to evaluate. Constrained tool-calling retains most of the value at bounded cost.

**Recommendation:** layered hybrid — deterministic Tree-sitter extraction + Supabase hybrid lexical/vector retrieval + capped Gemini tool-calling with a small fixed toolkit (hybrid_search, get_symbol, get_file_overview, get_import_neighbors), deterministic import/file graph in Postgres, one Gemini call per repo for the architecture overview, and backend validation of file:line citations. Deferred to post-MVP: full GraphRAG, free-roaming agents, LSP/CodeQL precision, graph databases.

**Decisions:** the five items formerly listed here were approved by Smit on 10 November 2026 (v0.6) and are recorded in §6.

**Implications:** more indexing complexity (parser + embedding + graph build) directly informs Research 2 and Research 4; embedding model choice becomes a real decision (feed into Research 3/5); cost profile ≈ N embeddings + 1 overview call per repo and 1–3 Gemini calls per question; citations can be validated, not just asserted.

### Research 2 — Parsing/language intelligence (completed 10 Nov 2026; all six decisions approved v0.6)

**Sources consulted (current official):** Tree-sitter documentation (bindings and parser usage); Python `ast` module docs; Go `go/parser` package docs (go1.27); `@babel/parser` docs; TypeScript Compiler API wiki; Cloudflare Containers docs; Cloudflare Workers WebAssembly docs; GitHub `stack-graphs` repository.

**Options compared:**
- **A. Tree-sitter (single engine):** upstream grammars for all four target languages; official Node/Python/Go bindings; error-tolerant partial trees (core design goal); symbols/imports extracted via the query system; syntactic only (no name resolution); deployable via wasm or a Cloudflare Container.
- **B. Native toolchains:** Python `ast` (exact grammar fidelity, but zero error tolerance and needs a Python runtime); Go `go/parser` (partial ASTs with `ast.Bad` nodes on errors, no dependencies needed, never executes code); `@babel/parser` (JS+TS+JSX with `errorRecovery: true`); TypeScript Compiler API (precise, but cross-file analysis needs installed dependencies — security conflict — and TS 7.1 changes the API). Net: three runtimes, uneven error tolerance, higher ops complexity.
- **C. LSP/LSIF/SCIP:** most precise cross-file references, but live servers are stateful, memory-heavy, and need dependency installs; offline indexers still require toolchains and dependency installation. Rejected for MVP.
- **Cross-file finding:** GitHub's `stack-graphs` (Tree-sitter-driven name resolution) was archived upstream on 9 Sep 2025 — not a viable MVP dependency. MVP cross-file intelligence uses a heuristic import/file graph instead.

**Decisions (approved 10 Nov 2026):** Tree-sitter as the single extraction engine for all four languages; per-language versioned query files with parser versions pinned per analysis; partial-tree indexing of broken files with stored per-file error counts; heuristic import/file graph for cross-file relationships with precise name resolution deferred; parsing in a Cloudflare Container (Node + Tree-sitter), subject to Research 4 topology; native toolchains, LSP/SCIP, and stack-graphs excluded from MVP.

**Implications:** container compute is CPU-bound but batch-amortized (Workers Paid plan required); query files are one-time per-language effort; heuristic import resolution misses dynamic imports and re-exports and must be represented honestly in answers (hypothesis vs. fact); wasm-in-Workers kept as a possible later optimization for small repos.

### Research 3 — Supabase data/storage design (completed 10 Nov 2026; decisions approved v0.7)

**Sources:** Supabase pgvector/vector-columns, RLS, and Storage access-control docs; PostgreSQL JSONB docs; Supabase hybrid-search docs.

**Key findings:** pgvector stores embeddings as typed vector columns (fewer dimensions perform best; HNSW indexes); similarity queries must be wrapped in Postgres functions called via rpc() (the auto-API does not support the operators). RLS policies live in the database and cannot be bypassed except by service keys (never client-side). Storage buckets deny all access by default until policies exist. JSONB is queryable and indexable, but relational columns give integrity and efficient joins.

**Decisions (approved):** relational tables for core entities with JSONB for flexible extraction details; embeddings on chunk rows with a single hybrid RRF search function; import graph as relational tables — no graph database in MVP; two private buckets (repo-archives deleted after indexing; exports via signed URLs); RLS on all tables; archives deleted post-indexing with a documented retention policy for analyses and conversations (durations at scope freeze).

**Implications:** one auditable search function simplifies the evaluation ablations (vector-only vs hybrid vs hybrid+graph = config flags); the chunk embedding dimension (1536) must match the chosen embedding model; any embedding-model change requires a full re-embed, so model versions are recorded per analysis.

### Research 4 — Cloudflare deployment topology (completed 10 Nov 2026; decisions approved v0.7)

**Sources:** Cloudflare Queues, Workflows, and Containers docs; Workers WebAssembly docs (Research 2).

**Key findings:** Workflows provide durable multi-step execution with per-step retries, state persisted for hours or weeks, and no step timeouts — a direct fit for the analysis pipeline. Queues provide guaranteed delivery with batching, retries, and dead-letter queues — a fit for chunk-embedding fan-out and notifications. Containers (Workers Paid plan) run any runtime on demand — the home for Tree-sitter parsing.

**Decisions (approved):** Workflows orchestrate the pipeline (download → parse → embed → graph → summarize; each step retryable with progress persisted); Queues fan out embedding batches with dead-letter handling; parsing runs in a Container as a workflow step; Workers only orchestrate.

**Implications:** Workers Paid plan is required (Containers); durable steps mean a failed embedding batch never loses completed parse work; per-instance observability feeds the UI status states directly.

### Research 5 — Gemini model and credential handling (completed 10 Nov 2026; decisions approved v0.7)

**Sources:** Gemini models, pricing, embeddings, and Gemini 3.8 Flash model docs.

**Key findings:** Gemini 4 Argon is the flagship (1M context; $2/$10 per 1M tokens promotional until 30 Nov 2026, rising to $4/$20; no free tier; thinking tokens and tool declarations billed). Gemini 3.8 Flash (stable, Sept 2026) is engineered for long-horizon software engineering with 1M context, 65K output, function calling, structured outputs, thinking levels, caching, and Batch API (50% cheaper) — at Flash-tier prices. Embeddings: gemini-embedding-001 (text, retrieval task types, up to 3072 dims, truncatable); gemini-embedding-2 is multimodal but incompatible with 001's vector space (a model change requires a full re-embed). On the Paid tier content is not used to improve Google's products; on Free it is.

**Decisions (approved):** gemini-3.8-flash as the single working model for overview generation, Q&A tool-calling, and error investigation (Argon deferred as a possible later deep mode); gemini-embedding-001 at 1536 dimensions with the Batch API during indexing; Paid tier only; operator's server-side key in Cloudflare secrets with per-session quotas and a monthly model budget.

**Implications:** cost per question ≈ 1–3 Flash calls plus one query embedding (cents-level); overview generation can use the Batch API for further savings; model versions are recorded per analysis for reproducibility.

### Research 6 — GitHub ingestion and safety (completed 10 Nov 2026; decisions approved v0.7)

**Sources:** GitHub REST repositories and rate-limit docs; git-scm gitsubmodules docs; GitHub Git LFS docs.

**Key findings:** the repositories API exposes default_branch (resolving the "main branch" clarification); tarball/zipball endpoints serve archives per branch/tag/SHA via codeload; submodules are separate repositories referenced by pointer entries and are not inside archives; Git LFS contents are stored outside the main repo with pointer files in-repo (archive inclusion is repo-owner opt-in); standard rate limits are ~60 requests/hour unauthenticated vs ~5,000/hour with a PAT, plus secondary abuse limits (figures re-verified at implementation).

**Decisions (approved):** metadata via REST with an operator-owned PAT plus tarball via codeload; snapshots pinned to the default branch head commit SHA; default branch = the API's default_branch field; submodules and LFS contents excluded from MVP (presence recorded); strict URL allow-list; streaming extraction with decompression-bomb, path-traversal, and symlink guards; repository size/file limits enforced with friendly rejection; vendored/build/minified/lockfile/binary exclusion per §11.

**Implications:** an operator PAT is required server-side (60 requests/hour is too low for a public service); reproducibility is guaranteed by the SHA pin; archive-only ingestion means we never run git against untrusted repositories, shrinking the attack surface.

### Research 7 — Product workflow (completed 10 Nov 2026; decisions approved v0.7)

**Sources:** none required — design synthesis of Research 1–6 and the confirmed MVP scope.

**Decisions (approved):** journey = paste URL → validation → job states (queued, downloading, parsing, indexing, ready / partially ready / failed with reasons); results dashboard = overview card, file explorer, symbol search, import-graph view, Q&A chat with citations panel, error-investigation mode, briefing-pack export; every answer separates facts, hypotheses, and unknowns and must abstain when evidence is insufficient; demo scenario = analyze a mid-size public repo, evidence-grounded Q&A, seeded stack-trace investigation, and briefing-pack export into Claude Code (final repos at scope freeze).

**Implications:** the workflow states map one-to-one onto the durable workflow steps from Research 4, so progress reporting is native; the citations panel doubles as the trust signal in the video demo.

### Research 8 — Evaluation design (completed 10 Nov 2026; decisions approved v0.7)

**Sources:** evaluation requirements from §8 and the confirmed evaluation-harness deliverable (v0.5).

**Decisions (approved):** benchmark = three repositories across the target languages (small Python, mid JS/TS, mid Go); four question categories (identifier lookup, conceptual, holistic architecture, error investigation with seeded stack traces) plus trap questions about nonexistent things; metrics = backend-validated citation validity, rubric answer correctness, abstention correctness, retrieval coverage, latency, and cost per question; the harness runs the three ablation configs (vector-only, hybrid, hybrid+graph) and emits the comparison report; success bar targets = ≥80% valid citations, ≥70% correct answers, 100% trap-question abstention, cents-level cost per question (finalized at scope freeze).

**Implications:** trap questions operationalize the "abstain rather than hallucinate" principle as a measurable metric; the harness doubles as a regression suite for prompt changes and feeds the prompt-documentation requirement (§13).

## 10. Technical questions to settle

- Canonical pipeline from GitHub URL to queryable snapshot. *(Answered v0.7: REST metadata + codeload tarball → Workflow steps download/parse/embed/graph/summarize → Supabase index.)*
- Deterministic facts extracted before LLM calls. *(Answered v0.6: Tree-sitter extraction — files, symbols, imports, AST-aware chunk boundaries, per-file error counts.)*
- Representation of files, symbols, imports, chunks, and line ranges. *(Answered v0.7: relational tables with JSONB extras; embeddings on chunk rows; line ranges from parse results.)*
- Hybrid lexical/vector/graph retrieval strategy. *(Answered v0.7: one hybrid RRF SQL function + relational import graph.)*
- Exact file/line citation generation and validation. *(Direction confirmed: backend validation against the index; mechanism finalized at design stage.)*
- Behaviour when parsers cannot parse or resolve symbols. *(Answered v0.6: partial trees indexed; per-file error counts recorded.)*
- Background job queue, retries, cancellation, and progress reporting. *(Answered v0.7: Workflows + Queues with dead-letter handling.)*
- What belongs in Supabase tables versus Storage buckets. *(Answered v0.7: index data in tables; archives and exports in buckets.)*
- RLS and bucket access policies. *(Answered v0.7: RLS on all tables; private buckets; signed URLs for exports.)*
- Provider key storage and usage/cost controls. *(Answered v0.7: Cloudflare secrets; per-session quotas; monthly model budget.)*
- Repository/file/time limits and exclusion rules. *(Mechanism answered v0.7: caps enforced with friendly rejection; exact numbers at scope freeze.)*
- How to avoid executing untrusted repository code. *(Answered: never run repo code; Tree-sitter reads only; archive-based ingestion without git execution.)*
- Briefing Pack format, export rules, and how the SKILL.md instructs consumer AIs (added v0.5).

## 11. Security and reliability requirements

- Treat all repo contents as untrusted, including comments, README, and configs.
- Prevent prompt injection from repository content.
- Do not run repo code, install scripts, build scripts, or arbitrary commands in initial analysis.
- Validate GitHub URLs and redirects; prevent SSRF.
- Safely handle archive extraction, path traversal, decompression bombs, huge files, binaries, and symlinks.
- Define exclusion rules for generated output, vendored dependencies, build artefacts, lockfiles, and binaries.
- Avoid exposing accidental secrets found in public source.
- Apply processing limits, throttling, and model budgets.
- Expose clear states: queued, downloading, parsing, indexing, ready, partially ready, failed.
- Tie every analysis to an immutable commit SHA and label stale analyses.
- Ensure citations refer to the exact analysed snapshot and valid line ranges.
- The exported Briefing Pack contains untrusted repository content; its SKILL.md must instruct consumer AIs to treat pack content as data, never as instructions (added v0.5).
- Submodules and Git LFS contents are never fetched in MVP; their presence is recorded in the analysis (added v0.7).

## 12. Conceptual data entities (not a schema)

Do not create tables/migrations until research and approval. Candidate entities:
- User/session or anonymous workspace identity.
- Repository.
- Repository analysis (branch, commit SHA, status, timestamps, parser/index versions).
- Repository file.
- Code symbol.
- Code relationship (imports, calls/references, containment, inheritance where available).
- Code chunk.
- Embedding/index record.
- Analysis artifact (architecture summary, briefing pack, evaluation reports).
- Conversation/query (question, answer, evidence, model/usage metadata; subject to retention).

## 13. Prompt-engineering evidence

Maintain a chronological log: goal, AI tool/model, complete prompt, context supplied, output, shortcomings, revised prompt, and final result. Preserve meaningful iterations. Categories: requirements, architecture research, data model, parser/retrieval design, implementation, code review/debugging, evaluation, docs, presentation. SKILL.md files for the briefing-pack export are prompt artifacts and must be logged with their iterations (added v0.5).

## 14. Work status

**Completed in planning**
- Assessment brief reviewed.
- Category and project concept selected.
- High-level problem, users, principles, and capability candidates defined.
- Public-only GitHub, URL input/no GitHub auth, and four languages decided.
- Gemini provider family selected.
- Supabase database and Storage buckets selected.
- Mostly-Cloudflare direction selected.
- No-implementation-until-approval rule established.
- Research 1 completed; layered hybrid architecture recommended; five decisions approved 10 Nov 2026 (v0.6).
- Research 2 completed; Tree-sitter confirmed as the single parsing engine with Cloudflare Container placement direction; six decisions approved 10 Nov 2026 (v0.6).
- Research 3–8 completed in one pass at Smit's request; all decisions approved 10 Nov 2026 (v0.7): Supabase data model, Cloudflare topology, Gemini 3.8 Flash + embedding-001, GitHub ingestion/safety, product workflow, and evaluation design.
- §7 clarifications resolved: default branch = the API's default_branch; Gemini credentials = the operator's server-side key.
- Scope additions confirmed (v0.5): evaluation harness + architecture comparison report deliverable; Codebase Briefing Pack + SKILL.md export in the MVP; live MCP server and user-facing architecture toggle deferred; documentation generation demoted to Later.
- Step 1 complete (v0.8): MVP scope frozen 10 Nov 2026 — see §6.2 (limits, retention, budgets, benchmark repositories, demo repository, 48-question evaluation bank, nine-item frozen feature list).
- Task 1 complete (v0.9): Architecture & Data Design v1.0 approved — docs/Architecture_and_Data_Design_v1.0.md; four design questions resolved (httpOnly cookie, single-zip export, cron-Worker retention, separate eval schema).
- Task 2 complete (v1.0): 48-question evaluation bank approved — docs/Evaluation_Bank_v1.0.md (10 identifier, 10 conceptual, 10 holistic, 10 error-investigation, 8 trap; rubric, success bar, run protocol, report shape).
- Task 3 complete (v1.1): Smit gave explicit implementation approval on 10 Nov 2026. Implementation phase begins with the vertical slice.
- Backend revised (v1.2): Go API + Go pipeline worker + Python parser service on one host; Cloudflare = DNS/static hosting; Supabase buckets (CDN) serve files; job queue replaces Workflows/Queues; validation split (Zod frontend, Go/Pydantic server).
- Vertical slice scaffolded (v1.3): `apps/api` (Go API + worker, job runner, GitHub client, safe extractor, parser client, Gemini client with stub mode, hybrid search + citation validation), `apps/parser` (FastAPI + py-tree-sitter with per-language queries and chunking), `apps/web` (React + Vite + Tailwind, status polling, Q&A with citations), `packages/core` (Zod schemas), `supabase/migrations` (schema + `hybrid_search`), Docker Compose (db/parser/api/worker). Go build verified (`go build ./...`); Docker/PNPM/Python verification pending on the local machine.
- Vertical slice WORKING (v1.4): first-run issues found and fixed — (1) port 8080 conflict → API on 8081; (2) pnpm `allowBuilds` duplicate key; (3) py-tree-sitter 0.26 uses `QueryCursor` (compat layer added); (4) Gemini key returns 429/503 → resilient fallback to stub embeddings and extractive answers with a provider circuit breaker (20 s timeout); (5) partial unique index blocked finalize for repeat snapshots (migration 003); (6) job lease extended to 10 minutes + retry-budget guard. Verified: Flask → 83 files / 1,622 symbols / 1,682 chunks / 0 parse errors / `ready`; Q&A returns 3 validated citations in ~10 ms (stub fallback).
- Deliverables assembled (v1.5): git repository initialized (first commit `1c79535`, `.env`/`node_modules`/`.build` excluded, `go.sum` committed); `docs/PROMPT_DOCUMENTATION.md` (Deliverable #2 — 27 logged prompts with verbatim text, failures, revisions, techniques, lessons); `docs/VIDEO_SCRIPT.md` (Deliverable #1 — timed 2:30 script, pre-recording checklist, demo runbook, Q&A prep, submission checklist); stub-mode retrieval improved (full-text-only weighting when the provider is unavailable).

**Implementation milestones (in order)**
1. Revise the architecture design to v1.1 for the Go/Python backend (topology, job runner, parser service, validation split).
2. Vertical slice: URL validation → GitHub fetch → safe extract → Tree-sitter parse → local index → one cited Q&A (pallets/flask).
3. Database wiring: Supabase schema + migrations + buckets + RLS; pipeline worker.
4. Remaining frozen MVP features (§6.2 items 2–9), including Briefing Pack export and evaluation harness.
5. Evaluation runs on the 48-question bank; comparison report artifact.
6. Assessment deliverables: video, prompt log, accessible project link.

## 15. New-chat continuation prompt

> Continue the AI Codebase Intelligence project from this handoff (v1.2). Implementation is approved. The backend was revised: Go HTTP API + Go pipeline worker + Python parser service on one host (Docker Compose), Cloudflare for DNS/static hosting, Supabase Postgres + pgvector + Storage buckets (CDN) for data and files. Next: revise the architecture design to v1.1 for this backend, then build the vertical slice (URL → download → parse → store → one cited Q&A on pallets/flask). Stay within the frozen MVP scope (§6.2) and use the 48-question evaluation bank as the regression suite.

**Working preference:** one research topic at a time, clear explanations, manageable pace, and explicit separation between confirmed decisions, recommendations, and open questions.

## 16. Change log

- v0.1 — Initial concept and category.
- v0.2 — Capability ideas and research topics.
- v0.3 — Public GitHub, URL/no auth, four languages, Gemini family, main-branch intention, mostly Cloudflare.
- v0.4 — 10 Nov 2026: Supabase database and Storage buckets confirmed; one-topic-at-a-time research preference; consolidated handoff.
- v0.5 — 10 Nov 2026: Research 1 completed (layered hybrid architecture recommended; approval pending); Research log section added; scope additions confirmed (evaluation harness + comparison report deliverable, Codebase Briefing Pack + SKILL.md export in MVP, live MCP server and user-facing comparison toggle deferred, documentation generation demoted to Later); continuation prompt updated for Research 2.
- v0.6 — 10 Nov 2026: Research 2 completed; all five Research 1 and six Research 2 decisions approved by Smit and recorded in §6 (§6.1 converted to an approvals log); Tree-sitter confirmed as the single extraction engine for the four languages with Cloudflare Container placement direction (subject to Research 4); heuristic import/file graph confirmed for MVP cross-file intelligence; two §10 questions answered; pending order and continuation prompt updated for Research 3.
- v0.7 — 10 Nov 2026: Research 3–8 completed in one pass at Smit's request and all decisions approved — Supabase relational+JSONB model with pgvector hybrid search and two private buckets; Cloudflare Workflows+Queues+Container topology; Gemini 3.8 Flash with embedding-001 @1536 on the Paid tier using an operator server-side key; GitHub ingestion via PAT+codeload pinned to default_branch head SHA with submodules/LFS excluded; product workflow and demo scenario defined; evaluation benchmark, metrics, and success bar set; §7 clarifications resolved; most §10 questions answered; pending order reduced to scope freeze → architecture/data dictionary → evaluation finalization → explicit approval.
- v0.8 — 10 Nov 2026: Step 1 complete — MVP scope frozen (§6.2): repository limits, retention durations, usage/cost guards, benchmark repositories (psf/requests, honojs/hono, gin-gonic/gin), demo repository (pallets/flask with expressjs/express backup), 48-question evaluation bank, and the nine-item frozen MVP feature list; eight candidate repositories verified live via the GitHub API (4 of 8 default to `master`, validating the default_branch decision); §7 repository-limits and retention clarifications resolved; pending order reduced to architecture/data dictionary → evaluation bank → explicit approval.
- v0.9 — 10 Nov 2026: Task 1 complete — Architecture & Data Design v1.0 approved (docs/Architecture_and_Data_Design_v1.0.md): 13-entity data dictionary, schema design with the single hybrid_search function, deny-by-default RLS with Worker-only DB access, API contracts (12 endpoints), workflow contracts (7 steps), Zod validation rules, and traceability; four design questions resolved (httpOnly cookie, single-zip export, cron-Worker retention, separate eval schema); pending reduced to evaluation bank (Task 2) → explicit approval (Task 3).
- v1.0 — 10 Nov 2026: **Planning phase complete.** Task 2 complete — 48-question evaluation bank approved (docs/Evaluation_Bank_v1.0.md): 10 identifier, 10 conceptual, 10 holistic, 10 error-investigation, 8 trap questions across the three benchmark repos, with grading rubric, success bar, run protocol, and report shape; evaluation-bank row and §6.4 added; pending reduced to a single item — Task 3, Smit's explicit implementation approval.
- v1.1 — 10 Nov 2026: Task 3 complete — Smit approved implementation; build begins with the vertical slice (URL → download → parse → store → one cited Q&A on pallets/flask); pending list replaced by implementation milestones; implementation-status banner added.
- v1.2 — 10 Nov 2026: Backend revised on Smit's direction — Go HTTP API + Go pipeline worker (Postgres-backed job queue) + Python parser service (py-tree-sitter) on one host; Cloudflare reduced to DNS + static frontend hosting; Supabase Storage buckets (CDN-backed) serve files; Workflows/Queues/Container decisions marked superseded; validation split recorded (Zod frontend, Go/Pydantic server); §6.5 added; milestones updated to include the design revision to v1.1.
- v1.3 — 10 Nov 2026: Vertical slice scaffolded and Go build verified — repo structure (apps/api, apps/parser, apps/web, packages/core, supabase/migrations), Docker Compose with shared data volume, job runner with leases/retries, GitHub ingestion with guards, Tree-sitter extraction (Python/Go/JS/TS) with partial-tree policy, chunking, embeddings with deterministic stub fallback, hybrid search, citation validation, and the React UI with status polling and Q&A. Remaining verification (docker compose up, pnpm install/build, pytest) documented in the README.
- v1.4 — 10 Nov 2026: **Vertical slice verified working end-to-end.** Fixed first-run issues: API port 8081 (8080 conflict), pnpm allowBuilds config, py-tree-sitter 0.26 QueryCursor compatibility, Gemini 429/503 resilience (stub embeddings + extractive answer fallback + 20 s timeout + provider circuit breaker), duplicate-snapshot finalize failure (migration 003 drops the partial unique index), job lease 10 minutes + retry-budget guard. Flask analysis: 83 files, 1,622 symbols, 1,682 chunks, 0 parse errors, status ready; Q&A returns validated citations (~10 ms stub).
- v1.5 — 11 Nov 2026: **Deliverables assembled (3 of 3 blockers addressed).** Git repository initialized with a clean first commit (`1c79535`); prompt documentation written (`docs/PROMPT_DOCUMENTATION.md`, P-01…P-27, verbatim prompts, failures, revisions, techniques, lessons); video script and demo runbook written (`docs/VIDEO_SCRIPT.md`, timed 2:30 script, checklist, runbook, Q&A prep, submission checklist); stub-mode retrieval now uses full-text-only weighting when the provider is unavailable. Remaining: record the video, push the repository, then the feature gaps (investigation UI, Briefing Pack export, 48-question evaluation run).