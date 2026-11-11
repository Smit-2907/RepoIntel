package pipeline

import (
	"context"
	"fmt"
	"os"
	"path/filepath"
	"strings"

	"codebaseintel/internal/extract"
	"codebaseintel/internal/store"
)

func (p *Pipeline) fetchMetadata(ctx context.Context, job *store.Job) error {
	a, err := p.Store.GetAnalysis(ctx, job.AnalysisID)
	if err != nil {
		return err
	}
	parts := strings.SplitN(a.RepoFullName, "/", 2)
	if len(parts) != 2 {
		return permanent("invalid repository name %q", a.RepoFullName)
	}
	_ = p.Store.SetStatus(ctx, a.ID, "fetching", "fetching repository metadata")
	meta, err := p.GitHub.Metadata(ctx, parts[0], parts[1])
	if err != nil {
		return fmt.Errorf("fetch metadata: %w", err)
	}
	if meta.SizeKB > p.Cfg.MaxRepoKB {
		return permanent("repository is %d KB, above the %d KB limit", meta.SizeKB, p.Cfg.MaxRepoKB)
	}
	if err := p.Store.SetMetadata(ctx, a.ID, meta.DefaultBranch, meta.CommitSHA, meta.SizeKB); err != nil {
		return err
	}
	_ = p.Store.SetStatus(ctx, a.ID, "downloading", "downloading archive")
	return p.Store.EnqueueJob(ctx, a.ID, "download_archive", nil, 3)
}

func (p *Pipeline) downloadArchive(ctx context.Context, job *store.Job, payload stepPayload) error {
	a, err := p.Store.GetAnalysis(ctx, job.AnalysisID)
	if err != nil {
		return err
	}
	parts := strings.SplitN(a.RepoFullName, "/", 2)
	data, err := p.GitHub.DownloadTarball(ctx, parts[0], parts[1], a.CommitSHA, int64(p.Cfg.MaxRepoKB)*1024)
	if err != nil {
		return fmt.Errorf("download: %w", err)
	}
	dir := filepath.Join(p.Cfg.DataDir, a.ID)
	if err := os.MkdirAll(dir, 0o755); err != nil {
		return err
	}
	archive := filepath.Join(dir, "source.tar.gz")
	if err := os.WriteFile(archive, data, 0o644); err != nil {
		return err
	}
	_ = p.Store.SetStatus(ctx, a.ID, "extracting", "extracting source files")
	return p.Store.EnqueueJob(ctx, a.ID, "extract", stepPayload{ArchivePath: archive}, 1)
}

func (p *Pipeline) extractStep(ctx context.Context, job *store.Job, payload stepPayload) error {
	a, err := p.Store.GetAnalysis(ctx, job.AnalysisID)
	if err != nil {
		return err
	}
	if payload.ArchivePath == "" {
		return permanent("extract: missing archive path")
	}
	data, err := os.ReadFile(payload.ArchivePath)
	if err != nil {
		return err
	}
	files, skipped, err := extract.Tarball(data, extract.Limits{
		MaxFiles:      p.Cfg.MaxFiles,
		MaxFileBytes:  p.Cfg.MaxFileBytes,
		MaxTotalBytes: 200 << 20,
	})
	if err != nil {
		return fmt.Errorf("extract: %w", err)
	}
	dir := filepath.Join(p.Cfg.DataDir, a.ID, "files")
	for _, f := range files {
		dst := filepath.Join(dir, filepath.FromSlash(f.Path))
		if !strings.HasPrefix(dst, dir+string(os.PathSeparator)) {
			continue
		}
		if err := os.MkdirAll(filepath.Dir(dst), 0o755); err != nil {
			return err
		}
		if err := os.WriteFile(dst, f.Content, 0o644); err != nil {
			return err
		}
	}
	_ = p.Store.SetStatus(ctx, a.ID, "parsing", fmt.Sprintf("parsing %d files", len(files)))
	return p.Store.EnqueueJob(ctx, a.ID, "parse", stepPayload{FileCount: len(files), Skipped: skipped}, 2)
}
