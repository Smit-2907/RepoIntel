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
