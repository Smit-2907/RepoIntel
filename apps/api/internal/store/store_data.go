package store

import (
	"context"
	"fmt"
)

type FileRow struct {
	Path            string
	Language        string
	SizeBytes       int
	LineCount       int
	ParseErrorCount int
	IsIndexed       bool
}

type SymbolRow struct {
	Path          string
	Name          string
	QualifiedName string
	Kind          string
	StartLine     int
	EndLine       int
	ParentName    string
	IsExported    bool
}

type ChunkRow struct {
	Path          string
	SymbolName    string
	StartLine     int
	EndLine       int
	Content       string
	ContentHash   string
	TokenEstimate int
}

type ImportRow struct {
	FromPath     string
	RawSpecifier string
}

// InsertParseResults writes files, symbols, chunks and import edges in one transaction.
func (s *Store) InsertParseResults(ctx context.Context, analysisID string, files []FileRow,
	symbols []SymbolRow, chunks []ChunkRow, imports []ImportRow) error {

	tx, err := s.Pool.Begin(ctx)
	if err != nil {
		return err
	}
	defer tx.Rollback(ctx)

	fileIDs := make(map[string]string, len(files))
	for _, f := range files {
		var id string
		if err := tx.QueryRow(ctx, `
			insert into analysis_files(analysis_id, path, language, size_bytes, line_count, parse_error_count, is_indexed)
			values($1,$2,$3,$4,$5,$6,$7) returning id`,
			analysisID, f.Path, f.Language, f.SizeBytes, f.LineCount, f.ParseErrorCount, f.IsIndexed).Scan(&id); err != nil {
			return err
		}
		fileIDs[f.Path] = id
	}

	symbolIDs := make(map[string]string, len(symbols))
	for _, sym := range symbols {
		fid, ok := fileIDs[sym.Path]
		if !ok {
			continue
		}
		var id string
		if err := tx.QueryRow(ctx, `
			insert into code_symbols(analysis_id, file_id, name, qualified_name, kind, start_line, end_line, is_exported)
			values($1,$2,$3,$4,$5,$6,$7,$8) returning id`,
			analysisID, fid, sym.Name, sym.QualifiedName, sym.Kind, sym.StartLine, sym.EndLine, sym.IsExported).Scan(&id); err != nil {
			return err
		}
		symbolIDs[fmt.Sprintf("%s|%d|%d", sym.Path, sym.StartLine, sym.EndLine)] = id
	}

	for _, c := range chunks {
		fid, ok := fileIDs[c.Path]
		if !ok {
			continue
		}
		var symID *string
		if id, ok := symbolIDs[fmt.Sprintf("%s|%d|%d", c.Path, c.StartLine, c.EndLine)]; ok {
			symID = &id
		}
		if _, err := tx.Exec(ctx, `
			insert into code_chunks(analysis_id, file_id, symbol_id, start_line, end_line, content, content_hash, token_estimate)
			values($1,$2,$3,$4,$5,$6,$7,$8) on conflict do nothing`,
			analysisID, fid, symID, c.StartLine, c.EndLine, c.Content, c.ContentHash, c.TokenEstimate); err != nil {
			return err
		}
	}

	for _, imp := range imports {
		fid, ok := fileIDs[imp.FromPath]
		if !ok {
			continue
		}
		if _, err := tx.Exec(ctx, `
			insert into file_imports(analysis_id, from_file_id, raw_specifier, kind)
			values($1,$2,$3,'unknown') on conflict do nothing`,
			analysisID, fid, imp.RawSpecifier); err != nil {
			return err
		}
	}
	return tx.Commit(ctx)
}
