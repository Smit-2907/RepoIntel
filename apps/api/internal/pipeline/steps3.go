package pipeline

import (
	"context"
	"fmt"
	"os"
	"path/filepath"
	"sort"
	"strings"

	"codebaseintel/internal/store"
)

func (p *Pipeline) graphStep(ctx context.Context, job *store.Job) error {
	a, err := p.Store.GetAnalysis(ctx, job.AnalysisID)
	if err != nil {
		return err
	}
	_ = p.Store.SetStatus(ctx, a.ID, "summarizing", "building import graph and summary")
	if err := p.resolveImports(ctx, a.ID); err != nil {
		return fmt.Errorf("import graph: %w", err)
	}
	summary, err := p.buildSummary(ctx, a)
	if err != nil {
		return err
	}
	if err := p.Store.InsertArtifact(ctx, a.ID, "architecture_summary", summary, p.Gemini.Model); err != nil {
		return err
	}
	return p.Store.EnqueueJob(ctx, a.ID, "finalize", nil, 3)
}

func (p *Pipeline) buildSummary(ctx context.Context, a store.Analysis) (string, error) {
	files, err := p.Store.ListFiles(ctx, a.ID)
	if err != nil {
		return "", err
	}
	syms, err := p.Store.ListTopSymbols(ctx, a.ID, 25)
	if err != nil {
		return "", err
	}
	langs := map[string]int{}
	entry := make([]string, 0, 12)
	for _, f := range files {
		langs[f.Language]++
		switch strings.ToLower(filepath.Base(f.Path)) {
		case "main.go", "main.py", "app.py", "index.ts", "index.js", "server.ts", "__init__.py":
			entry = append(entry, f.Path)
		}
	}
	sort.Strings(entry)

	var b strings.Builder
	fmt.Fprintf(&b, "# Architecture summary — %s\n\n", a.RepoFullName)
	fmt.Fprintf(&b, "- Snapshot: `%s` @ `%s`\n", a.DefaultBranch, a.CommitSHA)
	fmt.Fprintf(&b, "- Files: %d (%d indexed) · Symbols: %d · Chunks: %d · Parse-error rate: %.1f%%\n\n",
		a.FileCount, a.IndexedFileCount, a.SymbolCount, a.ChunkCount, a.ParseErrorRate*100)
	b.WriteString("## Languages\n\n")
	for lang, n := range langs {
		fmt.Fprintf(&b, "- %s: %d files\n", lang, n)
	}
	b.WriteString("\n## Entry-point candidates\n\n")
	if len(entry) == 0 {
		b.WriteString("- (none detected)\n")
	}
	for i, e := range entry {
		if i >= 12 {
			break
		}
		fmt.Fprintf(&b, "- `%s`\n", e)
	}
	b.WriteString("\n## Key symbols (largest first)\n\n")
	for _, s := range syms {
		fmt.Fprintf(&b, "- `%s` (%s) — `%s` lines %d-%d\n", s.Name, s.Kind, s.Path, s.Start, s.End)
	}
	return b.String(), nil
}

func (p *Pipeline) finalizeStep(ctx context.Context, job *store.Job) error {
	a, err := p.Store.GetAnalysis(ctx, job.AnalysisID)
	if err != nil {
		return err
	}
	status := "ready"
	detail := fmt.Sprintf("indexed %d files, %d symbols, %d chunks", a.IndexedFileCount, a.SymbolCount, a.ChunkCount)
	if a.ParseErrorRate > 0.6 {
		status = "partially_ready"
		detail = fmt.Sprintf("partially ready: %.0f%% of indexed files had parse errors", a.ParseErrorRate*100)
	}
	if err := p.Store.Finish(ctx, a.ID, status, detail); err != nil {
		return err
	}
	_ = os.RemoveAll(filepath.Join(p.Cfg.DataDir, a.ID))
	return nil
}
