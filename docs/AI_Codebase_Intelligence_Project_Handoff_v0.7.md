# AI Codebase Intelligence & Engineering Assistant
## Project Brief, Decisions, Research Plan & Handoff Notes
**Working specification v0.7 | Updated: 10 November 2026**

> **Implementation gate:** This project is in product definition, research, and planning. Do not write code, create implementation files, modify a repository, or deploy services until Smit explicitly approves implementation. Recommendations must never be treated as confirmed decisions without Smit's explicit approval.

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

## 5. Proposed release boundaries (not frozen)

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
| Deployment | Mostly Cloudflare | Direction confirmed; topology open |
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
| Parsing placement | Cloudflare Container (Node + Tree-sitter) in the indexing pipeline; subject to Research 4 topology confirmation | Confirmed 10 Nov 2026 (v0.6) |
| Native parsers / LSP | Rejected for MVP (multi-runtime sprawl, dependency-install security conflicts, uneven error tolerance) | Confirmed 10 Nov 2026 (v0.6) |
| Supabase data model | Relational tables for core entities; JSONB for flexible extraction details; import graph as relational tables; no graph DB | Confirmed 10 Nov 2026 (v0.7) |
| Supabase retrieval | pgvector embeddings on chunk rows (HNSW index); all search through one hybrid RRF SQL function | Confirmed 10 Nov 2026 (v0.7) |
| Supabase storage | Two private buckets: repo-archives (deleted after indexing) and exports (signed URLs); RLS on all tables; service keys server-side only | Confirmed 10 Nov 2026 (v0.7) |
| Retention direction | Archives deleted post-indexing; documented retention/deletion policy for analyses and conversations (durations at scope freeze) | Confirmed 10 Nov 2026 (v0.7) |
| Cloudflare topology | Workflows orchestrate the analysis pipeline; Queues fan out embedding batches with dead-letter handling; Container (Node + Tree-sitter) parses; Workers Paid plan required | Confirmed 10 Nov 2026 (v0.7) |
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

### 6.1 Approvals log

All five Research 1 decisions, all six Research 2 decisions, and all Research 3–8 decisions were approved by Smit on 10 November 2026 (v0.6–v0.7) and are recorded in the table above. No recommendations are currently awaiting approval.

## 7. Clarifications and open interpretations

- **"Main branch only":** RESOLVED v0.7 — analysis targets the repository's `default_branch` (from the GitHub API), never a hard-coded branch name.
- **Repository limits:** mechanism confirmed v0.7 (size/file caps enforced with friendly rejection); exact numbers to be set at scope freeze.
- **Gemini credentials:** RESOLVED v0.7 — the operator's server-side key, stored in Cloudflare secrets, with per-session quotas and a monthly budget; never exposed in browser code.
- **Snapshot identity:** direction confirmed v0.7 — store repository URL, default branch, head commit SHA, timestamps, analysis status, parser/index versions, and LLM/embedding model versions.
- **Retention/privacy:** direction confirmed v0.7 — archives deleted after indexing; documented retention/deletion policy for analyses and conversations, with final durations at scope freeze.
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

**Pending, recommended order**
1. Freeze MVP scope and demo scenario (final repo limits, retention durations, benchmark repos, and demo repo).
2. Define end-to-end architecture, data dictionary, schemas, and Zod validation.
3. Finalize evaluation cases and success criteria.
4. Obtain explicit approval before implementation.

## 15. New-chat continuation prompt

> Continue the AI Codebase Intelligence project from this handoff (v0.7). We are still in planning; do not implement or write code. All eight research topics are complete and every decision is approved (§6, §9). Next steps, in order: (1) freeze MVP scope and demo scenario — final repo limits, retention durations, benchmark repos, and the demo repo; (2) define the end-to-end architecture, data dictionary, schemas, and Zod validation; (3) finalize evaluation cases. Present each step for my approval separately, then obtain my explicit approval before any implementation.

**Working preference:** one research topic at a time, clear explanations, manageable pace, and explicit separation between confirmed decisions, recommendations, and open questions.

## 16. Change log

- v0.1 — Initial concept and category.
- v0.2 — Capability ideas and research topics.
- v0.3 — Public GitHub, URL/no auth, four languages, Gemini family, main-branch intention, mostly Cloudflare.
- v0.4 — 10 Nov 2026: Supabase database and Storage buckets confirmed; one-topic-at-a-time research preference; consolidated handoff.
- v0.5 — 10 Nov 2026: Research 1 completed (layered hybrid architecture recommended; approval pending); Research log section added; scope additions confirmed (evaluation harness + comparison report deliverable, Codebase Briefing Pack + SKILL.md export in MVP, live MCP server and user-facing comparison toggle deferred, documentation generation demoted to Later); continuation prompt updated for Research 2.
- v0.6 — 10 Nov 2026: Research 2 completed; all five Research 1 and six Research 2 decisions approved by Smit and recorded in §6 (§6.1 converted to an approvals log); Tree-sitter confirmed as the single extraction engine for the four languages with Cloudflare Container placement direction (subject to Research 4); heuristic import/file graph confirmed for MVP cross-file intelligence; two §10 questions answered; pending order and continuation prompt updated for Research 3.
- v0.7 — 10 Nov 2026: Research 3–8 completed in one pass at Smit's request and all decisions approved — Supabase relational+JSONB model with pgvector hybrid search and two private buckets; Cloudflare Workflows+Queues+Container topology; Gemini 3.8 Flash with embedding-001 @1536 on the Paid tier using an operator server-side key; GitHub ingestion via PAT+codeload pinned to default_branch head SHA with submodules/LFS excluded; product workflow and demo scenario defined; evaluation benchmark, metrics, and success bar set; §7 clarifications resolved; most §10 questions answered; pending order reduced to scope freeze → architecture/data dictionary → evaluation finalization → explicit approval.