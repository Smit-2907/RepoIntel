# AI Codebase Intelligence & Engineering Assistant
## Project Brief, Decisions, Research Plan & Handoff Notes
**Working specification v0.6 | Updated: 10 November 2026**

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

### 6.1 Approvals log

All five Research 1 architecture decisions and all six Research 2 parsing decisions were approved by Smit on 10 November 2026 (v0.6) and are recorded in the table above. No recommendations are currently awaiting approval.

## 7. Clarifications and open interpretations

- **"Main branch only":** recommended interpretation is the repository's default branch (which might be `main`, `master`, or another name), not a hard-coded branch named `main`. Needs confirmation.
- **Repository limits:** undecided. Gemini credential ownership is separate from compute/storage/network/rate-limit/abuse/model-cost constraints. Reasonable limits will still be needed.
- **Gemini credentials:** clarify whether the app uses the operator's server-side key or asks each user for their own. Never expose provider secrets in browser code.
- **Snapshot identity:** recommended to store repository URL, branch, commit SHA, timestamps, analysis status, and parser/index versions.
- **Retention/privacy:** decide retention and deletion for archives, parsed data, embeddings, and conversations.
- **Supabase layout:** database and buckets are decided; schemas, RLS, bucket policies, environment separation, and retention need design.
- **Cloudflare:** mostly Cloudflare does not necessarily mean heavy parsing must run inside Workers. Background/compute placement remains open.

## 8. Research plan — one topic at a time

1. **AI/code-intelligence architecture:** compare basic RAG, syntax-aware/static analysis, hybrid analysis + retrieval, graph-based intelligence, and agent/tool-calling. Determine a reliable, maintainable combination. **COMPLETED 10 Nov 2026 — see Research log (§9); all five decisions approved 10 Nov 2026 (v0.6).**
2. **Parsing/language intelligence:** compare Tree-sitter, native parsers/compiler tooling, and LSP for Python, Go, JS, TS; assess symbols, imports, cross-file references, error tolerance, and deployment fit. **COMPLETED 10 Nov 2026 — see Research log (§9); all six decisions approved 10 Nov 2026 (v0.6).**
3. **Supabase data/storage design:** use Supabase PostgreSQL + Storage buckets. Compare relational metadata, JSONB, pgvector/retrieval choices, bucket layout, access controls, retention, and whether a graph DB is needed.
4. **Cloudflare deployment:** evaluate Workers, Queues, Workflows, and compute options for heavier parsing; map Cloudflare and Supabase responsibilities.
5. **Gemini model selection:** compare current models for coding/reasoning, context, structured output, tool calling, latency, pricing, limits, and data handling.
6. **GitHub ingestion/safety:** compare archive/API/raw access; assess rate limits, large repos, submodules/LFS, generated/binary files, URL validation, SSRF, extraction safety, prompt injection, and reproducibility.
7. **Product workflow:** define URL-to-results journey, progress/status, failure states, evidence display, briefing-pack export, and first demo scenario.
8. **Evaluation:** create benchmark questions and bug cases; measure citation correctness, answer correctness, abstention, retrieval coverage, latency, and cost. Feeds the confirmed evaluation-harness deliverable.

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

## 10. Technical questions to settle

- Canonical pipeline from GitHub URL to queryable snapshot.
- Deterministic facts extracted before LLM calls. *(Answered v0.6: Tree-sitter extraction — files, symbols, imports, AST-aware chunk boundaries, per-file error counts.)*
- Representation of files, symbols, imports, chunks, and line ranges.
- Hybrid lexical/vector/graph retrieval strategy.
- Exact file/line citation generation and validation.
- Behaviour when parsers cannot parse or resolve symbols. *(Answered v0.6: partial trees indexed; per-file error counts recorded.)*
- Background job queue, retries, cancellation, and progress reporting.
- What belongs in Supabase tables versus Storage buckets.
- RLS and bucket access policies.
- Provider key storage and usage/cost controls.
- Repository/file/time limits and exclusion rules.
- How to avoid executing untrusted repository code.
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
- Scope additions confirmed (v0.5): evaluation harness + architecture comparison report deliverable; Codebase Briefing Pack + SKILL.md export in the MVP; live MCP server and user-facing architecture toggle deferred; documentation generation demoted to Later.

**Pending, recommended order**
1. Research 3: Supabase PostgreSQL + Storage + vector retrieval.
2. Research 4: Cloudflare/background processing topology.
3. Research 5: Gemini model selection and credential handling.
4. Research 6: GitHub ingestion, limits, and safety.
5. Confirm default-branch interpretation and API-key ownership.
6. Freeze MVP scope and demo scenario.
7. Define end-to-end architecture, data dictionary, schemas, and Zod validation.
8. Define evaluation cases and success criteria (includes benchmark design for the evaluation harness).
9. Obtain explicit approval before implementation.

## 15. New-chat continuation prompt

> Continue the AI Codebase Intelligence project from this handoff (v0.6). We are still in planning/research; do not implement or write code. Research 1 and Research 2 are complete and all their decisions are approved (§6). Proceed with Research 3 only: Supabase data/storage design (PostgreSQL schema direction, JSONB vs relational choices, pgvector/hybrid retrieval, bucket layout, RLS and bucket policies, retention, and whether a graph DB is needed). Use current official sources, explain trade-offs, and end with a recommendation plus decisions for me to approve. Then wait for my response before moving to Research 4.

**Working preference:** one research topic at a time, clear explanations, manageable pace, and explicit separation between confirmed decisions, recommendations, and open questions.

## 16. Change log

- v0.1 — Initial concept and category.
- v0.2 — Capability ideas and research topics.
- v0.3 — Public GitHub, URL/no auth, four languages, Gemini family, main-branch intention, mostly Cloudflare.
- v0.4 — 10 Nov 2026: Supabase database and Storage buckets confirmed; one-topic-at-a-time research preference; consolidated handoff.
- v0.5 — 10 Nov 2026: Research 1 completed (layered hybrid architecture recommended; approval pending); Research log section added; scope additions confirmed (evaluation harness + comparison report deliverable, Codebase Briefing Pack + SKILL.md export in MVP, live MCP server and user-facing comparison toggle deferred, documentation generation demoted to Later); continuation prompt updated for Research 2.
- v0.6 — 10 Nov 2026: Research 2 completed; all five Research 1 and six Research 2 decisions approved by Smit and recorded in §6 (§6.1 converted to an approvals log); Tree-sitter confirmed as the single extraction engine for the four languages with Cloudflare Container placement direction (subject to Research 4); heuristic import/file graph confirmed for MVP cross-file intelligence; two §10 questions answered; pending order and continuation prompt updated for Research 3.