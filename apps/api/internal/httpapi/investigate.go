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

// investigate handles POST /api/analyses/{id}/investigations.
func (s *Server) investigate(w http.ResponseWriter, r *http.Request) {
	started := time.Now()
	id := r.PathValue("id")

	var body struct {
		StackTrace string `json:"stack_trace"`
	}
	if err := json.NewDecoder(http.MaxBytesReader(w, r.Body, 30<<10)).Decode(&body); err != nil {
		writeErr(w, http.StatusBadRequest, "INVALID_INPUT", "body must be {\"stack_trace\": \"...\"}")
		return
	}
	trace := strings.TrimSpace(body.StackTrace)
	if len(trace) < 10 || len(trace) > 20000 {
		writeErr(w, http.StatusUnprocessableEntity, "INVALID_INPUT", "stack_trace must be 10-20000 characters")
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

	frames := parseStackFrames(trace)
	files, err := s.Store.ListFiles(r.Context(), a.ID)
	if err != nil {
		writeErr(w, http.StatusInternalServerError, "INTERNAL", "could not load index")
		return
	}

	var chunks []store.RetrievedChunk
	matched := 0
	for _, fr := range frames {
		if p, ok := matchFile(files, fr.File); ok {
			matched++
			near, err := s.Store.ChunksNearLine(r.Context(), a.ID, p, fr.Line, 40, 2)
			if err == nil {
				chunks = append(chunks, near...)
			}
			if len(chunks) >= 12 {
				break
			}
		}
	}

	// Augment with lexical retrieval over the trace text when frames are sparse.
	if len(chunks) < 4 {
		vecs, err := s.Gemini.Embed(r.Context(), []string{trace}, "RETRIEVAL_QUERY")
		if err == nil {
			ft, sem := 1.0, 1.0
			if s.Gemini.IsStub() {
				sem = 0.0
			}
			if more, err := s.Store.HybridSearch(r.Context(), a.ID, trace, vecs[0], 8, ft, sem); err == nil {
				chunks = append(chunks, more...)
			}
		}
	}
	chunks = dedupeChunks(chunks, 12)

	ctxChunks := make([]gemini.ContextChunk, len(chunks))
	for i, c := range chunks {
		ctxChunks[i] = gemini.ContextChunk{
			FilePath: c.FilePath, SymbolName: c.SymbolName,
			StartLine: c.StartLine, EndLine: c.EndLine, Content: c.Content,
		}
	}

	question := "Investigate this stack trace. Identify the most likely cause, the exact code involved, and concrete next steps.\n\n" + trace
	ans, err := s.Gemini.Answer(r.Context(), question, ctxChunks)
	if err != nil {
		log.Printf("investigation fallback: %v", err)
		ans = gemini.StubAnswer(question, ctxChunks)
		ans.Unknowns = append(ans.Unknowns,
			"LLM provider unavailable ("+err.Error()+"); extractive fallback used (no LLM reasoning).")
	}

	valid, dropped := validateCitations(ans.Citations, chunks)
	if len(valid) == 0 && !ans.Stub {
		ans.AnswerMarkdown = "Insufficient evidence in the indexed snapshot to investigate this trace."
		ans.Unknowns = append(ans.Unknowns, "No citation could be grounded in retrieved code.")
	}
	if dropped > 0 {
		ans.Unknowns = append(ans.Unknowns, "Some citations were dropped because they did not match the index.")
	}

	convID, err := s.Store.CreateConversation(r.Context(), a.ID, a.WorkspaceID, "investigation")
	if err != nil {
		writeErr(w, http.StatusInternalServerError, "INTERNAL", "could not persist conversation")
		return
	}
	_, _ = s.Store.InsertMessage(r.Context(), convID, "user", trace, nil, nil, nil, nil, "", 0, 0, 0)
	latency := int(time.Since(started).Milliseconds())
	msgID, err := s.Store.InsertMessage(r.Context(), convID, "assistant", ans.AnswerMarkdown,
		valid, ans.Facts, ans.Hypotheses, ans.Unknowns, ans.Model, ans.InputTokens, ans.OutputTokens, latency)
	if err != nil {
		log.Printf("investigation message: %v", err)
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
		"parsed_frames":  len(frames),
		"matched_frames": matched,
		"stub_mode":      ans.Stub,
	})
}
