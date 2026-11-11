// Package pipeline implements the job-runner steps (design v1.1 §3, §10).
package pipeline

import (
	"context"
	"encoding/json"
	"errors"
	"fmt"
	"time"

	"codebaseintel/internal/config"
	"codebaseintel/internal/gemini"
	"codebaseintel/internal/github"
	"codebaseintel/internal/parserclient"
	"codebaseintel/internal/store"
)

type Pipeline struct {
	Store  *store.Store
	GitHub *github.Client
	Parser *parserclient.Client
	Gemini *gemini.Client
	Cfg    config.Config
}

func New(st *store.Store, gh *github.Client, pc *parserclient.Client, gm *gemini.Client, cfg config.Config) *Pipeline {
	return &Pipeline{Store: st, GitHub: gh, Parser: pc, Gemini: gm, Cfg: cfg}
}

type stepPayload struct {
	ArchivePath string `json:"archive_path,omitempty"`
	FileCount   int    `json:"file_count,omitempty"`
	Skipped     map[string]int `json:"skipped,omitempty"`
}

type permanentError struct{ msg string }

func (e permanentError) Error() string { return e.msg }

func permanent(format string, args ...any) error {
	return permanentError{msg: fmt.Sprintf(format, args...)}
}

func IsPermanent(err error) bool {
	var p permanentError
	return errors.As(err, &p)
}

// Run executes a single claimed job.
func (p *Pipeline) Run(ctx context.Context, job *store.Job) error {
	var payload stepPayload
	if len(job.Payload) > 0 {
		_ = json.Unmarshal(job.Payload, &payload)
	}
	switch job.Step {
	case "fetch_metadata":
		return p.fetchMetadata(ctx, job)
	case "download_archive":
		return p.downloadArchive(ctx, job, payload)
	case "extract":
		return p.extractStep(ctx, job, payload)
	case "parse":
		return p.parseStep(ctx, job, payload)
	case "embed":
		return p.embedStep(ctx, job)
	case "graph_and_summary":
		return p.graphStep(ctx, job)
	case "finalize":
		return p.finalizeStep(ctx, job)
	default:
		return permanent("unknown step %q", job.Step)
	}
}

// Work is the worker loop: claim → execute → complete/retry with backoff.
func (p *Pipeline) Work(ctx context.Context, workerID string) {
	for {
		select {
		case <-ctx.Done():
			return
		default:
		}
		job, err := p.Store.ClaimJob(ctx, workerID)
		if err != nil {
			time.Sleep(2 * time.Second)
			continue
		}
		// Defense in depth: never exceed the retry budget, even if a lease expired.
		if job.Attempts > job.MaxAttempts {
			p.Store.RecordEvent(ctx, job.AnalysisID, job.Step, "failed", map[string]any{"error": "retry budget exhausted"})
			_ = p.Store.Fail(ctx, job.AnalysisID, job.Step, "retry budget exhausted")
			_ = p.Store.MarkJobFailed(ctx, job.ID, "retry budget exhausted")
			continue
		}
		p.Store.RecordEvent(ctx, job.AnalysisID, job.Step, "started", nil)
		if err := p.Run(ctx, job); err != nil {
			if IsPermanent(err) {
				p.Store.RecordEvent(ctx, job.AnalysisID, job.Step, "failed", map[string]any{"error": err.Error()})
				_ = p.Store.Fail(ctx, job.AnalysisID, job.Step, err.Error())
				_ = p.Store.MarkJobFailed(ctx, job.ID, err.Error())
				continue
			}
			failed, ferr := p.Store.FailJob(ctx, job.ID, err.Error(), time.Duration(job.Attempts*5)*time.Second)
			p.Store.RecordEvent(ctx, job.AnalysisID, job.Step, "retried", map[string]any{
				"error": err.Error(), "attempt": job.Attempts, "permanent": failed && ferr == nil,
			})
			if failed {
				_ = p.Store.Fail(ctx, job.AnalysisID, job.Step, err.Error())
			}
			continue
		}
		_ = p.Store.CompleteJob(ctx, job.ID)
		p.Store.RecordEvent(ctx, job.AnalysisID, job.Step, "succeeded", nil)
	}
}
