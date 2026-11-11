package store

import (
	"context"
	"fmt"
	"strings"
)

type PendingChunk struct {
	ID      string
	Content string
}

func (s *Store) ChunksMissingEmbeddings(ctx context.Context, analysisID string, limit int) ([]PendingChunk, error) {
	rows, err := s.Pool.Query(ctx, `
		select id, content from code_chunks
		where analysis_id=$1 and embedding is null
		order by created_at limit $2`, analysisID, limit)
	if err != nil {
		return nil, err
	}
	defer rows.Close()
	var out []PendingChunk
	for rows.Next() {
		var c PendingChunk
		if err := rows.Scan(&c.ID, &c.Content); err != nil {
			return nil, err
		}
		out = append(out, c)
	}
	return out, rows.Err()
}

func vectorLiteral(vec []float32) string {
	var b strings.Builder
	b.WriteByte('[')
	for i, v := range vec {
		if i > 0 {
			b.WriteByte(',')
		}
		fmt.Fprintf(&b, "%g", v)
	}
	b.WriteByte(']')
	return b.String()
}

func (s *Store) SetEmbedding(ctx context.Context, chunkID string, vec []float32) error {
	_, err := s.Pool.Exec(ctx, `update code_chunks set embedding=$2::vector where id=$1`, chunkID, vectorLiteral(vec))
	return err
}

type RetrievedChunk struct {
	ChunkID    string
	FilePath   string
	SymbolName string
	StartLine  int
	EndLine    int
	Content    string
	Score      float64
}

// HybridSearch calls the single auditable retrieval function (design §8.3).
// Weights allow ablations: stub embeddings carry no meaning, so callers can
// set semantic weight to 0 and rank purely by full-text relevance.
func (s *Store) HybridSearch(ctx context.Context, analysisID, queryText string, vec []float32,
	k int, ftWeight, semWeight float64) ([]RetrievedChunk, error) {
	rows, err := s.Pool.Query(ctx, `
		select chunk_id, file_path, coalesce(symbol_name,''), start_line, end_line, content, rrf_score
		from hybrid_search($1,$2,$3::vector,$4,$5,$6)`,
		analysisID, queryText, vectorLiteral(vec), k, ftWeight, semWeight)
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
