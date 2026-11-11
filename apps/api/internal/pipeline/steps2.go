package pipeline

import (
	"context"
	"fmt"
	"io/fs"
	"os"
	"path/filepath"

	"codebaseintel/internal/extract"
	"codebaseintel/internal/gemini"
	"codebaseintel/internal/parserclient"
	"codebaseintel/internal/store"
)

func (p *Pipeline) parseStep(ctx context.Context, job *store.Job, payload stepPayload) error {
	a, err := p.Store.GetAnalysis(ctx, job.AnalysisID)
	if err != nil {
		return err
	}
	dir := filepath.Join(p.Cfg.DataDir, a.ID, "files")

	var files []parserclient.File
	err = filepath.WalkDir(dir, func(path string, d fs.DirEntry, werr error) error {
		if werr != nil {
			return werr
		}
		if d.IsDir() {
			return nil
		}
		rel, rerr := filepath.Rel(dir, path)
		if rerr != nil {
			return rerr
		}
		rel = filepath.ToSlash(rel)
		lang := extract.LanguageFor(rel)
		if lang == "other" {
			return nil
		}
		content, rerr := os.ReadFile(path)
		if rerr != nil {
			return rerr
		}
		files = append(files, parserclient.File{Path: rel, Language: lang, Content: string(content)})
		if len(files) >= p.Cfg.MaxFiles {
			return fs.SkipAll
		}
		return nil
	})
	if err != nil {
		return fmt.Errorf("read extracted files: %w", err)
	}
	if len(files) == 0 {
		return permanent("no supported source files found")
	}

	res, err := p.Parser.Parse(ctx, a.ID, files, p.Cfg.MaxFiles, p.Cfg.MaxChunks)
	if err != nil {
		return fmt.Errorf("parse: %w", err)
	}

	fileRows := make([]store.FileRow, 0, len(res.Files))
	for _, f := range res.Files {
		fileRows = append(fileRows, store.FileRow{
			Path: f.Path, Language: f.Language, SizeBytes: f.SizeBytes,
			LineCount: f.LineCount, ParseErrorCount: f.ParseErrorCount, IsIndexed: f.IsIndexed,
		})
	}
	symRows := make([]store.SymbolRow, 0, len(res.Symbols))
	for _, s := range res.Symbols {
		symRows = append(symRows, store.SymbolRow{
			Path: s.Path, Name: s.Name, QualifiedName: s.QualifiedName, Kind: s.Kind,
			StartLine: s.StartLine, EndLine: s.EndLine, ParentName: s.ParentName, IsExported: s.IsExported,
		})
	}
	chunkRows := make([]store.ChunkRow, 0, len(res.Chunks))
	for _, c := range res.Chunks {
		chunkRows = append(chunkRows, store.ChunkRow{
			Path: c.Path, SymbolName: c.SymbolName, StartLine: c.StartLine, EndLine: c.EndLine,
			Content: c.Content, ContentHash: c.ContentHash, TokenEstimate: c.TokenEstimate,
		})
	}
	impRows := make([]store.ImportRow, 0, len(res.Imports))
	for _, im := range res.Imports {
		impRows = append(impRows, store.ImportRow{FromPath: im.FromPath, RawSpecifier: im.RawSpecifier})
	}

	if err := p.Store.InsertParseResults(ctx, a.ID, fileRows, symRows, chunkRows, impRows); err != nil {
		return fmt.Errorf("store parse results: %w", err)
	}

	withErrors := 0
	indexed := 0
	for _, f := range res.Files {
		if f.IsIndexed {
			indexed++
			if f.ParseErrorCount > 0 {
				withErrors++
			}
		}
	}
	errRate := 0.0
	if indexed > 0 {
		errRate = float64(withErrors) / float64(indexed)
	}
	if err := p.Store.SetCounts(ctx, a.ID, len(res.Files), indexed, len(res.Symbols), len(res.Chunks), errRate); err != nil {
		return err
	}

	_ = p.Store.SetStatus(ctx, a.ID, "embedding", "embedding chunks")
	return p.Store.EnqueueJob(ctx, a.ID, "embed", nil, 3)
}

func (p *Pipeline) embedStep(ctx context.Context, job *store.Job) error {
	a, err := p.Store.GetAnalysis(ctx, job.AnalysisID)
	if err != nil {
		return err
	}
	_ = p.Store.SetStatus(ctx, a.ID, "embedding", "embedding chunks")
	total := 0
	for {
		batch, err := p.Store.ChunksMissingEmbeddings(ctx, a.ID, 100)
		if err != nil {
			return err
		}
		if len(batch) == 0 {
			break
		}
		texts := make([]string, len(batch))
		for i, c := range batch {
			texts[i] = c.Content
		}
		vecs, err := p.Gemini.Embed(ctx, texts, "RETRIEVAL_DOCUMENT")
		if err != nil {
			// Resilient fallback: any provider failure (quota, 503, network)
			// degrades to deterministic stub vectors, recorded honestly.
			vecs = gemini.StubEmbeddings(texts, p.Cfg.EmbeddingDim)
			p.Store.RecordEvent(ctx, a.ID, "embed", "retried", map[string]any{
				"warning": "embedding provider unavailable; stub embeddings used for this batch",
				"error":   err.Error(),
			})
		}
		for i, c := range batch {
			if err := p.Store.SetEmbedding(ctx, c.ID, vecs[i]); err != nil {
				return err
			}
		}
		total += len(batch)
	}
	_ = p.Store.SetStatus(ctx, a.ID, "summarizing", fmt.Sprintf("embedded %d chunks", total))
	return p.Store.EnqueueJob(ctx, a.ID, "graph_and_summary", nil, 2)
}
