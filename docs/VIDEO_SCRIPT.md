# Video Presentation Script & Demo Runbook — Deliverable #1
## AI Codebase Intelligence & Engineering Assistant
**Gokul Infocare assessment · Technology → Innovation & New Ideas**
**Target length: 2:30–3:00 (minimum 2:00) | 11 November 2026**

> **Brief requirements covered:** AI tools and why chosen · prompt design/iteration · approach and relevant experience · challenges and solutions · live demonstration of the outcome.

---

## 1. Pre-recording checklist (10 minutes)

| # | Check | Command / action |
|---|---|---|
| 1 | Stack running | `docker compose ps` → db, parser, api, worker all **Up** |
| 2 | API healthy | `curl -s http://localhost:8081/health` → `{"status":"ok"}` |
| 3 | Flask pre-indexed (fast demo) | analysis `521936f8-…` is `ready`; or run `bash scripts/smoke.sh` |
| 4 | UI running | `pnpm dev:web` → http://localhost:5173 (private window) |
| 5 | Clean browser | Close extra tabs; zoom 110%; hide bookmarks bar |
| 6 | Terminal font | 14–16 pt, dark theme, no secrets visible |
| 7 | **Hide your API key** | Do not open `.env` on camera |
| 8 | Screen recorder | 1080p, mic test, record a 10-second sample first |

**Mode choice:** with a working `GEMINI_API_KEY` you get reasoned answers; without one the app runs in **stub mode** (extractive + citations). Both are honest — say which mode you are in on camera.

---

## 2. Timed script (word-for-word)

### 0:00–0:20 — Opening (20 s)
> "Hi, I'm Smit Solanki. This is my submission for the Gokul Infocare assessment, Technology — Innovation and New Ideas: an **AI Codebase Intelligence and Engineering Assistant**. It lets a developer paste a public GitHub URL and then ask questions about unfamiliar code — with answers grounded in exact file and line citations."

**On screen:** UI with the heading and URL input.

### 0:20–0:50 — AI tools and why (30 s)
> "I used AI in three layers. First, as a **development partner**: Cline in VS Code for research, architecture, design and implementation. Second, inside the product: **Gemini 3.8 Flash** for answering, and **Gemini embedding-001** for semantic search. Third, I had an independent **Gemini audit** of the finished codebase.
> The stack is Go for the API and pipeline worker, Python with **Tree-sitter** for deterministic code parsing, React for the UI, and PostgreSQL with **pgvector** for hybrid search."

**On screen:** architecture diagram (docs/Architecture_and_Data_Design_v1.1.md §1.2).

### 0:50–1:20 — Prompt engineering journey (30 s)
> "My prompt work followed a discipline: **evidence first, approval gates, then implementation**. I ran eight research topics one at a time, each ending with a decision for me to approve. Every recommendation was verified against official sources — for example, I discovered that GitHub's stack-graphs project was archived, so I switched to a heuristic import graph.
> For the product prompts, two rules matter most: the model must return **strict JSON**, and repository content is **data, never instructions** — that blocks prompt injection from the code being analysed. And every citation is **machine-validated** against the index before it is shown."

**On screen:** the system prompt in `answer.go`, then a citation card in the UI.

### 1:20–1:50 — Live demo (30 s)
> "Let me show it. I paste the Flask repository — the pipeline pins the default branch commit, downloads the archive, extracts it safely, parses it with Tree-sitter, embeds the chunks, and builds the import graph. You can see the stages: fetching, parsing 83 files, embedding, summarizing — and now **ready**.
> I ask: *'Where is routing handled?'* — and here is the answer, with citations to the exact files and line ranges."

**On screen:** paste URL → status card progresses → ask the question → citations appear.

**Demo tips:**
- Use the pre-indexed analysis for speed; mention that a fresh run takes ~30–60 s.
- If the provider is unavailable, say: "Gemini is rate-limited right now, so it falls back to an extractive answer — still grounded, still cited." That is a feature to show, not hide.

### 1:50–2:15 — Challenges and solutions (25 s)
> "Four real problems I solved: a port conflict that moved the API to 8081; a newer Tree-sitter library that removed the query API I used, fixed with a compatibility layer; a Gemini key with no quota, handled with a circuit breaker and honest fallback; and a database index that blocked repeat analyses, fixed with a migration. Each one is documented with the exact error and fix."

**On screen:** terminal showing the fallback and the migration (optional).

### 2:15–2:35 — Impact and close (20 s)
> "The result is a working vertical slice: deterministic facts before any AI call, citations you can verify, and graceful behaviour when the AI provider is down. Next milestones are the Briefing Pack export for IDE assistants, stack-trace investigation, and the 48-question evaluation harness. Thank you."

**On screen:** summary slide or README.

---

## 3. If something breaks during recording

| Problem | Instant fix on camera |
|---|---|
| UI blank / extension error | Switch to the terminal demo: `bash scripts/smoke.sh` (same flow, same citations) |
| Gemini quota error | Say "provider fallback" — answers become extractive and labelled |
| Analysis slow | Use the pre-indexed Flask analysis ID; mention re-runs take ~30–60 s |
| Docker not running | `docker compose up -d` (60 s); keep a pre-recorded backup clip |

**Always record a 30-second backup clip of the working demo** before the main take.

---

## 4. Demo runbook (fish-safe commands)

Open a terminal in the project:

```fish
cd "/home/smitsolanki/Smit Solanki/Personal/career/Gokul"
docker compose ps
curl -s http://localhost:8081/health
```

**Pre-index Flask (do this before recording):**

```fish
bash scripts/smoke.sh
```

The script prints the `analysis_id` and the final answer — keep that ID handy.

**Demo via API (if the UI is unavailable):**

```fish
curl -s -X POST http://localhost:8081/api/analyses -H 'Content-Type: application/json' -d '{"repo_url":"https://github.com/pallets/flask"}'
```

```fish
set id (curl -s -X POST http://localhost:8081/api/analyses -H 'Content-Type: application/json' -d '{"repo_url":"https://github.com/pallets/flask"}' | python3 -c 'import sys,json;print(json.load(sys.stdin)["analysis_id"])')
echo $id
curl -s http://localhost:8081/api/analyses/$id
```

```fish
curl -s -X POST http://localhost:8081/api/analyses/$id/questions -H 'Content-Type: application/json' -d '{"question":"Where is routing handled?"}' | python3 -m json.tool
```

**Start the UI:**

```fish
pnpm dev:web
```

Open http://localhost:5173 in a **private window** (extensions off).

---

## 5. Recording tips

1. **Script on a second monitor or phone** — do not read from the screen being recorded.
2. **Speak to the "why", not the "what"** — evaluators know what a pipeline is; explain *why* deterministic-first, *why* citations validated, *why* Go + Python.
3. **Show one failure honestly** (the fallback), then show it recovering — it demonstrates engineering maturity.
4. **Keep the cursor still** during narration; move it only when pointing at something.
5. **Do not show**: `.env`, API keys, tokens, personal paths in the browser.
6. **End card**: repository link + one-line description + contact.

---

## 6. Likely questions & prepared answers

| Question | Answer (30 seconds) |
|---|---|
| "Why not just use a chatbot with the repo pasted in?" | Context limits and hallucination. We index deterministically, retrieve precisely, and validate every citation against the snapshot. |
| "Why Tree-sitter?" | One engine for four languages, error-tolerant on broken files, never executes repository code. |
| "Why Go and Python?" | Go for a robust API and job runner; Python because `py-tree-sitter` is the simplest, most reliable Tree-sitter integration. |
| "How do you stop prompt injection from repo content?" | Content is fenced as data with an explicit system rule; output is schema-validated; no repo-derived tool calls. |
| "What happens when the LLM is down or out of quota?" | Circuit breaker → deterministic stub embeddings + extractive answers, clearly labelled in `unknowns`. The pipeline never silently fails. |
| "How do you measure quality?" | A 48-question bank (identifier, conceptual, holistic, error, traps) with citation validity, correctness, abstention, latency, and cost metrics. |
| "What's not built yet?" | Briefing Pack export, investigation UI, evaluation run, retention cron — all scoped and documented. |

---

## 7. Submission checklist

- [ ] Video recorded (≥ 2:00, ideally 2:30–3:00), audio clear, demo visible
- [ ] Video uploaded (Drive/YouTube unlisted) with view access enabled
- [ ] Prompt documentation attached (`docs/PROMPT_DOCUMENTATION.md`)
- [ ] GitHub repository pushed (public) with README and run instructions
- [ ] Email to `jobs@gokulinfocare.com` with: video link, prompt doc link, repo link, one-paragraph summary
- [ ] Test all links in a private browser window before sending
