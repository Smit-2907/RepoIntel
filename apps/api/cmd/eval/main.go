// Command eval runs the 48-question bank against indexed analyses and writes a report.
package main

import (
	"bytes"
	"encoding/json"
	"flag"
	"fmt"
	"net/http"
	"os"
	"sort"
	"strings"
	"time"
)

type Benchmark struct {
	Repo       string `json:"repo"`
	AnalysisID string `json:"analysis_id"`
}

type Question struct {
	ID       string   `json:"id"`
	Repo     string   `json:"repo"`
	Category string   `json:"category"`
	Question string   `json:"question"`
	Expected []string `json:"expected"`
}

type Bank struct {
	Version    string      `json:"version"`
	Benchmarks []Benchmark `json:"benchmarks"`
	Questions  []Question  `json:"questions"`
}

type Result struct {
	ID        string
	Repo      string
	Category  string
	Question  string
	Status    int
	LatencyMS int
	Citations int
	Stub      bool
	Answer    string
	Error     string
}

func main() {
	bankPath := flag.String("bank", "docs/Evaluation_Bank.json", "path to the question bank")
	api := flag.String("api", "http://localhost:8081", "API base URL")
	out := flag.String("out", "docs/EVALUATION_REPORT.md", "report output path")
	flag.Parse()

	raw, err := os.ReadFile(*bankPath)
	if err != nil {
		fmt.Fprintf(os.Stderr, "read bank: %v\n", err)
		os.Exit(1)
	}
	var bank Bank
	if err := json.Unmarshal(raw, &bank); err != nil {
		fmt.Fprintf(os.Stderr, "parse bank: %v\n", err)
		os.Exit(1)
	}

	byRepo := map[string]string{}
	for _, b := range bank.Benchmarks {
		if b.AnalysisID != "" {
			byRepo[b.Repo] = b.AnalysisID
		}
	}

	results := make([]Result, 0, len(bank.Questions))
	for _, q := range bank.Questions {
		id := byRepo[q.Repo]
		if id == "" {
			results = append(results, Result{ID: q.ID, Repo: q.Repo, Category: q.Category, Question: q.Question, Error: "no analysis_id for " + q.Repo})
			fmt.Printf("%-6s %-18s SKIP (no analysis_id for %s)\n", q.ID, q.Category, q.Repo)
			continue
		}
		r := ask(*api, id, q)
		results = append(results, r)
		fmt.Printf("%-6s %-18s %6dms cites=%d stub=%v\n", q.ID, q.Category, r.LatencyMS, r.Citations, r.Stub)
	}

	if err := writeReport(*out, results); err != nil {
		fmt.Fprintf(os.Stderr, "write report: %v\n", err)
		os.Exit(1)
	}
	fmt.Printf("\nreport written to %s (%d questions)\n", *out, len(results))
}

func ask(api, analysisID string, q Question) Result {
	body, _ := json.Marshal(map[string]string{"question": q.Question})
	start := time.Now()
	resp, err := http.Post(fmt.Sprintf("%s/api/analyses/%s/questions", api, analysisID), "application/json", bytes.NewReader(body))
	if err != nil {
		return Result{ID: q.ID, Repo: q.Repo, Category: q.Category, Question: q.Question, Error: err.Error()}
	}
	defer resp.Body.Close()
	var parsed struct {
		AnswerMarkdown string `json:"answer_markdown"`
		Citations      []any  `json:"citations"`
		StubMode       bool   `json:"stub_mode"`
	}
	_ = json.NewDecoder(resp.Body).Decode(&parsed)
	return Result{
		ID: q.ID, Repo: q.Repo, Category: q.Category, Question: q.Question,
		Status: resp.StatusCode, LatencyMS: int(time.Since(start).Milliseconds()),
		Citations: len(parsed.Citations), Stub: parsed.StubMode, Answer: parsed.AnswerMarkdown,
	}
}

func writeReport(path string, results []Result) error {
	var b strings.Builder
	b.WriteString("# Evaluation Report\n\n")
	fmt.Fprintf(&b, "Generated: %s\n\n", time.Now().Format("2006-01-02 15:04"))

	type agg struct{ n, ok, cites, lat int }
	byCat := map[string]*agg{}
	for _, r := range results {
		a := byCat[r.Category]
		if a == nil {
			a = &agg{}
			byCat[r.Category] = a
		}
		a.n++
		if r.Error == "" && r.Status == 200 {
			a.ok++
			a.cites += r.Citations
			a.lat += r.LatencyMS
		}
	}
	cats := make([]string, 0, len(byCat))
	for c := range byCat {
		cats = append(cats, c)
	}
	sort.Strings(cats)

	b.WriteString("## Summary by category\n\n")
	b.WriteString("| Category | Questions | Succeeded | Avg citations | Avg latency (ms) |\n|---|---|---|---|---|\n")
	for _, c := range cats {
		a := byCat[c]
		avgCites, avgLat := 0, 0
		if a.ok > 0 {
			avgCites = a.cites / a.ok
			avgLat = a.lat / a.ok
		}
		fmt.Fprintf(&b, "| %s | %d | %d | %d | %d |\n", c, a.n, a.ok, avgCites, avgLat)
	}

	b.WriteString("\n## Results\n\n")
	b.WriteString("| ID | Category | Status | Latency (ms) | Citations | Stub | Answer (excerpt) |\n|---|---|---|---|---|---|---|\n")
	for _, r := range results {
		answer := strings.ReplaceAll(strings.TrimSpace(r.Answer), "\n", " ")
		if len(answer) > 160 {
			answer = answer[:160] + "…"
		}
		status := fmt.Sprintf("%d", r.Status)
		if r.Error != "" {
			status = "ERR: " + r.Error
		}
		fmt.Fprintf(&b, "| %s | %s | %s | %d | %d | %v | %s |\n",
			r.ID, r.Category, status, r.LatencyMS, r.Citations, r.Stub, answer)
	}

	b.WriteString("\n## Grading notes\n\n")
	b.WriteString("Scores are assigned by a human using `docs/Evaluation_Bank_v1.0.md` (0–3 rubric).\n")
	b.WriteString("Trap questions must be graded on abstention: any fabricated file reference is a safety failure.\n")
	return os.WriteFile(path, []byte(b.String()), 0o644)
}
