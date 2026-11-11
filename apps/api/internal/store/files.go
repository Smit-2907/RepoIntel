package store

import "context"

type FileRef struct {
	ID       string
	Path     string
	Language string
}

func (s *Store) ListFiles(ctx context.Context, analysisID string) ([]FileRef, error) {
	rows, err := s.Pool.Query(ctx, `select id, path, language from analysis_files where analysis_id=$1 and is_indexed`, analysisID)
	if err != nil {
		return nil, err
	}
	defer rows.Close()
	var out []FileRef
	for rows.Next() {
		var f FileRef
		if err := rows.Scan(&f.ID, &f.Path, &f.Language); err != nil {
			return nil, err
		}
		out = append(out, f)
	}
	return out, rows.Err()
}

type ImportEdge struct {
	ID       string
	FromPath string
	Raw      string
}

func (s *Store) ListImportEdges(ctx context.Context, analysisID string) ([]ImportEdge, error) {
	rows, err := s.Pool.Query(ctx, `
		select fi.id, f.path, fi.raw_specifier
		from file_imports fi join analysis_files f on f.id = fi.from_file_id
		where fi.analysis_id=$1`, analysisID)
	if err != nil {
		return nil, err
	}
	defer rows.Close()
	var out []ImportEdge
	for rows.Next() {
		var e ImportEdge
		if err := rows.Scan(&e.ID, &e.FromPath, &e.Raw); err != nil {
			return nil, err
		}
		out = append(out, e)
	}
	return out, rows.Err()
}

func (s *Store) SetImportResolution(ctx context.Context, id string, toFileID *string, kind string) error {
	_, err := s.Pool.Exec(ctx,
		`update file_imports set to_file_id=$2, kind=$3 where id=$1`, id, toFileID, kind)
	return err
}

type SymbolRef struct {
	Name  string
	Kind  string
	Path  string
	Start int
	End   int
}

// ChunksNearLine returns chunks for a file near a given line (stack-trace context).
func (s *Store) ChunksNearLine(ctx context.Context, analysisID, filePath string, line, radius, limit int) ([]RetrievedChunk, error) {
	rows, err := s.Pool.Query(ctx, `
		select c.id, f.path, coalesce(s.name,''), c.start_line, c.end_line, c.content, 0
		from code_chunks c
		join analysis_files f on f.id = c.file_id
		left join code_symbols s on s.id = c.symbol_id
		where c.analysis_id=$1 and f.path=$2
		  and c.start_line <= $3 + $4 and c.end_line >= $3 - $4
		order by abs(c.start_line - $3)
		limit $5`, analysisID, filePath, line, radius, limit)
	if err != nil {
		return nil, err
	}
	defer rows.Close()
	var out []RetrievedChunk
	for rows.Next() {
		var c RetrievedChunk
		if err := rows.Scan(&c.ChunkID, &c.FilePath, &c.SymbolName, &c.StartLine, &c.EndLine, &c.Content, &c.Score); err != nil {
			return nil, err
		}
		out = append(out, c)
	}
	return out, rows.Err()
}

type ImportCounts struct {
	Internal int
	External int
	Unknown  int
}

func (s *Store) ImportCounts(ctx context.Context, analysisID string) (ImportCounts, error) {
	var c ImportCounts
	err := s.Pool.QueryRow(ctx, `
		select
		  count(*) filter (where kind='internal'),
		  count(*) filter (where kind='external'),
		  count(*) filter (where kind='unknown')
		from file_imports where analysis_id=$1`, analysisID).
		Scan(&c.Internal, &c.External, &c.Unknown)
	return c, err
}

func (s *Store) GetArtifact(ctx context.Context, analysisID, kind string) (string, error) {
	var content string
	err := s.Pool.QueryRow(ctx, `
		select coalesce(content,'') from analysis_artifacts
		where analysis_id=$1 and kind=$2 order by created_at desc limit 1`, analysisID, kind).Scan(&content)
	return content, err
}

func (s *Store) InsertExport(ctx context.Context, analysisID, workspaceID, storagePath string, byteSize int) error {
	_, err := s.Pool.Exec(ctx, `
		insert into exports(analysis_id, workspace_id, storage_path, byte_size, expires_at)
		values($1,$2,$3,$4, now() + interval '7 days')`, analysisID, workspaceID, storagePath, byteSize)
	return err
}

func (s *Store) ListTopSymbols(ctx context.Context, analysisID string, limit int) ([]SymbolRef, error) {
	rows, err := s.Pool.Query(ctx, `
		select s.name, s.kind, f.path, s.start_line, s.end_line
		from code_symbols s join analysis_files f on f.id = s.file_id
		where s.analysis_id=$1
		order by (s.end_line - s.start_line) desc
		limit $2`, analysisID, limit)
	if err != nil {
		return nil, err
	}
	defer rows.Close()
	var out []SymbolRef
	for rows.Next() {
		var r SymbolRef
		if err := rows.Scan(&r.Name, &r.Kind, &r.Path, &r.Start, &r.End); err != nil {
			return nil, err
		}
		out = append(out, r)
	}
	return out, rows.Err()
}
