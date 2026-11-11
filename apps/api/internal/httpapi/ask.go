package httpapi

import (
	"encoding/json"
	"log"
	"net/http"
	"strings"
	"time"

	"codebaseintel/internal/gemini"
	"codebaseintel/internal/store"
)

type citationOut struct {
	FilePath  string `json:"file_path"`
	StartLine int    `json:"start_line"`
	EndLine   int    `json:"end_line"`
	Symbol    string `json:"symbol,omitempty"`
	Kind      string `json:"kind"`
}

func (s *Server) askQuestion(w http.ResponseWriter, r *http.Request) {
	started := time.Now()
	id := r.PathValue("id")

	var body struct {
		Question string `json:"question"`
	}
	if err := json.NewDecoder(http.MaxBytesReader(w, r.Body, 20<<10)).Decode(&body); err != nil {
		writeErr(w, http.StatusBadRequest, "INVALID_INPUT", "body must be {\"question\": \"...\"}")
		return
	}
	question := strings.TrimSpace(body.Question)
	if len(question) == 0 || len(question) > 2000 {
		writeErr(w, http.StatusUnprocessableEntity, "INVALID_INPUT", "question must be 1-2000 characters")
		return
	}

	a, err := s.Store.GetAnalysis(r.Context(), id)
	if err != nil {
		writeErr(w, http.StatusNotFound, "NOT_FOUND", "analysis not found")
		return
	}
	if a.Status != "ready" && a.Status != "partially_ready" {
		writeErr(w, http.StatusConflict, "ANALYSIS_NOT_READY", "analysis is not ready yet (status: "+a.Status+")")
		return
	}

	vecs, err := s.Gemini.Embed(r.Context(), []string{question}, "RETRIEVAL_QUERY")
	if err != nil {
		log.Printf("embed question: %v", err)
		writeErr(w, http.StatusBadGateway, "UPSTREAM_ERROR", "embedding provider failed")
		return
	}
	ftWeight, semWeight := 1.0, 1.0
	if s.Gemini.IsStub() {
		// Stub vectors are deterministic noise; rely on lexical relevance.
		semWeight = 0.0
	}
	chunks, err := s.Store.HybridSearch(r.Context(), a.ID, question, vecs[0], 12, ftWeight, semWeight)
	if err != nil {
		log.Printf("hybrid search: %v", err)
		writeErr(w, http.StatusInternalServerError, "INTERNAL", "retrieval failed")
		return
	}

	ctxChunks := make([]gemini.ContextChunk, len(chunks))
	for i, c := range chunks {
		ctxChunks[i] = gemini.ContextChunk{
			FilePath: c.FilePath, SymbolName: c.SymbolName,
			StartLine: c.StartLine, EndLine: c.EndLine, Content: c.Content,
		}
	}
	ans, err := s.Gemini.Answer(r.Context(), question, ctxChunks)
	if err != nil {
		// Any provider failure (quota, 503, network) degrades to the extractive
		// fallback, clearly labelled in the answer.
		log.Printf("answer fallback: %v", err)
		ans = gemini.StubAnswer(question, ctxChunks)
		ans.Unknowns = append(ans.Unknowns,
			"LLM provider unavailable ("+err.Error()+"); extractive fallback used (no LLM reasoning).")
	}

	valid, dropped := validateCitations(ans.Citations, chunks)
	if len(valid) == 0 && !ans.Stub {
		ans.AnswerMarkdown = "Insufficient evidence in the indexed snapshot to answer this question."
		ans.Facts = nil
		ans.Hypotheses = nil
		ans.Unknowns = append(ans.Unknowns, "The generated answer could not be grounded in retrieved code; nothing is asserted.")
	}
	if dropped > 0 {
		ans.Unknowns = append(ans.Unknowns, "Some citations were dropped because they did not match the retrieved index.")
	}

	convID, err := s.Store.CreateConversation(r.Context(), a.ID, a.WorkspaceID, "qa")
	if err != nil {
		log.Printf("conversation: %v", err)
		writeErr(w, http.StatusInternalServerError, "INTERNAL", "could not persist conversation")
		return
	}
	latency := int(time.Since(started).Milliseconds())
	if _, err := s.Store.InsertMessage(r.Context(), convID, "user", question, nil, nil, nil, nil, "", 0, 0, 0); err != nil {
		log.Printf("user message: %v", err)
	}
	msgID, err := s.Store.InsertMessage(r.Context(), convID, "assistant", ans.AnswerMarkdown,
		valid, ans.Facts, ans.Hypotheses, ans.Unknowns, ans.Model, ans.InputTokens, ans.OutputTokens, latency)
	if err != nil {
		log.Printf("assistant message: %v", err)
	}

	writeJSON(w, http.StatusOK, map[string]any{
		"message_id":      msgID,
		"answer_markdown": ans.AnswerMarkdown,
		"citations":       valid,
		"facts":           ans.Facts,
		"hypotheses":      ans.Hypotheses,
		"unknowns":        ans.Unknowns,
		"usage": map[string]int{
			"input_tokens":  ans.InputTokens,
			"output_tokens": ans.OutputTokens,
			"latency_ms":    latency,
		},
		"stub_mode": ans.Stub,
	})
}

// validateCitations keeps only citations that match a retrieved chunk's file and line range.
func validateCitations(citations []gemini.Citation, chunks []store.RetrievedChunk) ([]citationOut, int) {
	out := make([]citationOut, 0, len(citations))
	dropped := 0
	for _, c := range citations {
		ok := false
		for _, ch := range chunks {
			if ch.FilePath == c.FilePath && c.StartLine >= ch.StartLine && c.EndLine <= ch.EndLine {
				ok = true
				break
			}
		}
		if !ok {
			dropped++
			continue
		}
		kind := c.Kind
		if kind != "fact" && kind != "hypothesis" {
			kind = "fact"
		}
		out = append(out, citationOut{
			FilePath: c.FilePath, StartLine: c.StartLine, EndLine: c.EndLine, Symbol: c.Symbol, Kind: kind,
		})
	}
	return out, dropped
}
