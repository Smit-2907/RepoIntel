package store

import "context"

func (s *Store) CreateConversation(ctx context.Context, analysisID, workspaceID, kind string) (string, error) {
	var id string
	err := s.Pool.QueryRow(ctx,
		`insert into conversations(analysis_id, workspace_id, kind) values($1,$2,$3) returning id`,
		analysisID, workspaceID, kind).Scan(&id)
	return id, err
}

func (s *Store) InsertMessage(ctx context.Context, conversationID, role, content string,
	citations, facts, hypotheses, unknowns any, model string, inTok, outTok, latencyMS int) (string, error) {
	var id string
	err := s.Pool.QueryRow(ctx, `
		insert into messages(conversation_id, role, content, citations, facts, hypotheses, unknowns, model, input_tokens, output_tokens, latency_ms)
		values($1,$2,$3,$4,$5,$6,$7,$8,$9,$10,$11) returning id`,
		conversationID, role, content, citations, facts, hypotheses, unknowns, model, inTok, outTok, latencyMS).Scan(&id)
	return id, err
}

func (s *Store) InsertArtifact(ctx context.Context, analysisID, kind, content, model string) error {
	_, err := s.Pool.Exec(ctx,
		`insert into analysis_artifacts(analysis_id, kind, content, model) values($1,$2,$3,$4)`,
		analysisID, kind, content, model)
	return err
}
