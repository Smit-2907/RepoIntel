// Package store is the Postgres data layer (design v1.1 §7–§8).
package store

import (
	"context"
	"time"

	"github.com/jackc/pgx/v5/pgxpool"
)

type Store struct{ Pool *pgxpool.Pool }

func New(ctx context.Context, url string) (*Store, error) {
	pool, err := pgxpool.New(ctx, url)
	if err != nil {
		return nil, err
	}
	if err := pool.Ping(ctx); err != nil {
		return nil, err
	}
	return &Store{Pool: pool}, nil
}

func (s *Store) Close() { s.Pool.Close() }

type Analysis struct {
	ID               string
	WorkspaceID      string
	RepoFullName     string
	DefaultBranch    string
	CommitSHA        string
	Status           string
	StatusDetail     string
	SizeKB           int
	FileCount        int
	IndexedFileCount int
	SymbolCount      int
	ChunkCount       int
	ParseErrorRate   float64
	CreatedAt        time.Time
}

// Slice only: one anonymous workspace. Replaced by per-session tokens later.
const devTokenHash = "dev-workspace"

func (s *Store) CreateAnalysis(ctx context.Context, owner, repo, htmlURL string) (Analysis, error) {
	tx, err := s.Pool.Begin(ctx)
	if err != nil {
		return Analysis{}, err
	}
	defer tx.Rollback(ctx)

	var wsID string
	if err := tx.QueryRow(ctx,
		`insert into workspaces(token_hash) values($1)
		 on conflict (token_hash) do update set last_seen_at = now() returning id`, devTokenHash).Scan(&wsID); err != nil {
		return Analysis{}, err
	}

	full := owner + "/" + repo
	var repoID string
	if err := tx.QueryRow(ctx,
		`insert into repositories(owner, name, full_name, html_url) values($1,$2,$3,$4)
		 on conflict (full_name) do update set html_url = excluded.html_url returning id`,
		owner, repo, full, htmlURL).Scan(&repoID); err != nil {
		return Analysis{}, err
	}

	var a Analysis
	a.RepoFullName = full
	a.WorkspaceID = wsID
	if err := tx.QueryRow(ctx,
		`insert into analyses(repository_id, workspace_id, default_branch, commit_sha, status)
		 values($1,$2,'',repeat('0',40),'queued') returning id, status, created_at`,
		repoID, wsID).Scan(&a.ID, &a.Status, &a.CreatedAt); err != nil {
		return Analysis{}, err
	}

	if _, err := tx.Exec(ctx,
		`insert into jobs(analysis_id, step, max_attempts) values($1,'fetch_metadata',3)`, a.ID); err != nil {
		return Analysis{}, err
	}
	return a, tx.Commit(ctx)
}

func (s *Store) GetAnalysis(ctx context.Context, id string) (Analysis, error) {
	var a Analysis
	err := s.Pool.QueryRow(ctx, `
		select a.id, a.workspace_id, r.full_name, a.default_branch, a.commit_sha, a.status,
		       coalesce(a.status_detail,''), coalesce(a.size_kb,0), a.file_count, a.indexed_file_count,
		       a.symbol_count, a.chunk_count, coalesce(a.parse_error_rate,0), a.created_at
		from analyses a join repositories r on r.id = a.repository_id where a.id = $1`, id).
		Scan(&a.ID, &a.WorkspaceID, &a.RepoFullName, &a.DefaultBranch, &a.CommitSHA, &a.Status,
			&a.StatusDetail, &a.SizeKB, &a.FileCount, &a.IndexedFileCount, &a.SymbolCount,
			&a.ChunkCount, &a.ParseErrorRate, &a.CreatedAt)
	return a, err
}

func (s *Store) SetStatus(ctx context.Context, id, status, detail string) error {
	_, err := s.Pool.Exec(ctx, `update analyses set status=$2, status_detail=$3 where id=$1`, id, status, detail)
	return err
}

func (s *Store) SetMetadata(ctx context.Context, id, branch, sha string, sizeKB int) error {
	_, err := s.Pool.Exec(ctx,
		`update analyses set default_branch=$2, commit_sha=$3, size_kb=$4, started_at=coalesce(started_at, now()) where id=$1`,
		id, branch, sha, sizeKB)
	return err
}

func (s *Store) SetCounts(ctx context.Context, id string, files, indexed, symbols, chunks int, errRate float64) error {
	_, err := s.Pool.Exec(ctx,
		`update analyses set file_count=$2, indexed_file_count=$3, symbol_count=$4, chunk_count=$5, parse_error_rate=$6 where id=$1`,
		id, files, indexed, symbols, chunks, errRate)
	return err
}

func (s *Store) Finish(ctx context.Context, id, status, detail string) error {
	_, err := s.Pool.Exec(ctx,
		`update analyses set status=$2, status_detail=$3, finished_at=now() where id=$1`, id, status, detail)
	return err
}

func (s *Store) Fail(ctx context.Context, id, step, msg string) error {
	_, err := s.Pool.Exec(ctx,
		`update analyses set status='failed', error_step=$2, error_message=$3, finished_at=now() where id=$1`,
		id, step, msg)
	return err
}
