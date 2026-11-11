# AI Codebase Intelligence & Engineering Assistant
## Project Brief, Decisions, Research Plan & Handoff Notes
**Working specification v0.4 | Updated: 10 November 2026**

> **Implementation gate:** This project is in product definition, research, and planning. Do not write code, create implementation files, modify a repository, or deploy services until Smit explicitly approves implementation.

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

## 4. Candidate capabilities

- Public GitHub URL import.
- Repository overview: languages, framework clues, structure, configs, entry points.
- Code structure: files, symbols, functions/classes, imports, relationships.
- Architecture intelligence: modules and dependency relationships, potentially visualised.
- Context-aware Q&A with source evidence and line references.
- Error/stack-trace investigation with likely causes and evidence.
- Documentation generation.
- Later: code-quality insights such as possible cycles, duplication, maintainability concerns.

## 5. Proposed release boundaries (not frozen)

**MVP candidates:** import a public repository; analyse one branch/snapshot and record commit SHA; produce repository overview; extract basic file/symbol/import data for the four target languages; generate architecture overview; answer questions with evidence; investigate a supplied error/stack trace.

**Later candidates:** richer docs/onboarding; deeper dependency/execution tracing; code-quality analysis; more languages; branch/commit comparison; private repositories/auth; collaboration/team workspaces; code modifications or PR workflows. Code editing/PR creation is deliberately outside initial scope unless reconsidered.

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

## 7. Clarifications and open interpretations

- **“Main branch only”:** recommended interpretation is the repository's default branch (which might be `main`, `master`, or another name), not a hard-coded branch named `main`. Needs confirmation.
- **Repository limits:** undecided. Gemini credential ownership is separate from compute/storage/network/rate-limit/abuse/model-cost constraints. Reasonable limits will still be needed.
- **Gemini credentials:** clarify whether the app uses the operator's server-side key or asks each user for their own. Never expose provider secrets in browser code.
- **Snapshot identity:** recommended to store repository URL, branch, commit SHA, timestamps, analysis status, and parser/index versions.
- **Retention/privacy:** decide retention and deletion for archives, parsed data, embeddings, and conversations.
- **Supabase layout:** database and buckets are decided; schemas, RLS, bucket policies, environment separation, and retention need design.
- **Cloudflare:** mostly Cloudflare does not necessarily mean heavy parsing must run inside Workers. Background/compute placement remains open.

## 8. Research plan — one topic at a time

1. **AI/code-intelligence architecture:** compare basic RAG, syntax-aware/static analysis, hybrid analysis + retrieval, graph-based intelligence, and agent/tool-calling. Determine a reliable, maintainable combination.
2. **Parsing/language intelligence:** compare Tree-sitter, native parsers/compiler tooling, and LSP for Python, Go, JS, TS; assess symbols, imports, cross-file references, error tolerance, and deployment fit.
3. **Supabase data/storage design:** use Supabase PostgreSQL + Storage buckets. Compare relational metadata, JSONB, pgvector/retrieval choices, bucket layout, access controls, retention, and whether a graph DB is needed.
4. **Cloudflare deployment:** evaluate Workers, Queues, Workflows, and compute options for heavier parsing; map Cloudflare and Supabase responsibilities.
5. **Gemini model selection:** compare current models for coding/reasoning, context, structured output, tool calling, latency, pricing, limits, and data handling.
6. **GitHub ingestion/safety:** compare archive/API/raw access; assess rate limits, large repos, submodules/LFS, generated/binary files, URL validation, SSRF, extraction safety, prompt injection, and reproducibility.
7. **Product workflow:** define URL-to-results journey, progress/status, failure states, evidence display, and first demo scenario.
8. **Evaluation:** create benchmark questions and bug cases; measure citation correctness, answer correctness, abstention, retrieval coverage, latency, and cost.

Each research topic should end with: evidence/sources, options comparison, recommendation, trade-offs, cost/operational implications, and decisions for Smit to approve. Do not silently convert recommendations into confirmed decisions.

## 9. Technical questions to settle

- Canonical pipeline from GitHub URL to queryable snapshot.
- Deterministic facts extracted before LLM calls.
- Representation of files, symbols, imports, chunks, and line ranges.
- Hybrid lexical/vector/graph retrieval strategy.
- Exact file/line citation generation and validation.
- Behaviour when parsers cannot parse or resolve symbols.
- Background job queue, retries, cancellation, and progress reporting.
- What belongs in Supabase tables versus Storage buckets.
- RLS and bucket access policies.
- Provider key storage and usage/cost controls.
- Repository/file/time limits and exclusion rules.
- How to avoid executing untrusted repository code.

## 10. Security and reliability requirements

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

## 11. Conceptual data entities (not a schema)

Do not create tables/migrations until research and approval. Candidate entities:
- User/session or anonymous workspace identity.
- Repository.
- Repository analysis (branch, commit SHA, status, timestamps, parser/index versions).
- Repository file.
- Code symbol.
- Code relationship (imports, calls/references, containment, inheritance where available).
- Code chunk.
- Embedding/index record.
- Analysis artifact (architecture summary, docs, reports).
- Conversation/query (question, answer, evidence, model/usage metadata; subject to retention).

## 12. Prompt-engineering evidence

Maintain a chronological log: goal, AI tool/model, complete prompt, context supplied, output, shortcomings, revised prompt, and final result. Preserve meaningful iterations. Categories: requirements, architecture research, data model, parser/retrieval design, implementation, code review/debugging, evaluation, docs, presentation.

## 13. Work status

**Completed in planning**
- Assessment brief reviewed.
- Category and project concept selected.
- High-level problem, users, principles, and capability candidates defined.
- Public-only GitHub, URL input/no GitHub auth, and four languages decided.
- Gemini provider family selected.
- Supabase database and Storage buckets selected.
- Mostly-Cloudflare direction selected.
- No-implementation-until-approval rule established.

**Pending, recommended order**
1. Research AI/code-intelligence architecture.
2. Research parsing/tooling for the four languages.
3. Research Supabase PostgreSQL + Storage + vector retrieval.
4. Research Cloudflare/background processing topology.
5. Research/select Gemini model and credential handling.
6. Research GitHub ingestion, limits, and safety.
7. Confirm default-branch interpretation and API-key ownership.
8. Freeze MVP scope and demo scenario.
9. Define end-to-end architecture, data dictionary, schemas, and Zod validation.
10. Define evaluation cases and success criteria.
11. Obtain explicit approval before implementation.

## 14. New-chat continuation prompt

> Continue the AI Codebase Intelligence project from this handoff. We are still in planning/research; do not implement. Start with Research 1 only: compare AI/code-intelligence architectures for our requirements. Use current official sources, explain trade-offs, and end with a recommendation plus decisions for me to approve. Then wait for my response before moving to Research 2.

**Working preference:** one research topic at a time, clear explanations, manageable pace, and explicit separation between confirmed decisions, recommendations, and open questions.

## 15. Change log

- v0.1 — Initial concept and category.
- v0.2 — Capability ideas and research topics.
- v0.3 — Public GitHub, URL/no auth, four languages, Gemini family, main-branch intention, mostly Cloudflare.
- v0.4 — 10 Nov 2026: Supabase database and Storage buckets confirmed; one-topic-at-a-time research preference; consolidated handoff.
