// Package gemini wraps the Gemini API with a deterministic stub fallback for local demos.
package gemini

import (
	"bytes"
	"context"
	"crypto/sha256"
	"encoding/binary"
	"encoding/json"
	"fmt"
	"io"
	"math"
	"net/http"
	"sync/atomic"
	"time"
)

type QuotaError struct{ Detail string }

func (e QuotaError) Error() string { return "gemini quota exceeded: " + e.Detail }

type Client struct {
	APIKey         string
	Model          string
	EmbeddingModel string
	Dim            int
	HTTP           *http.Client
	quotaExceeded  atomic.Bool
}
func New(apiKey, model, embeddingModel string, dim int) *Client {
	return &Client{
		APIKey:         apiKey,
		Model:          model,
		EmbeddingModel: embeddingModel,
		Dim:            dim,
		HTTP:           &http.Client{Timeout: 20 * time.Second},
	}
}

// IsStub reports whether responses are generated without a working provider.
func (c *Client) IsStub() bool {
	return c.APIKey == "" || c.quotaExceeded.Load()
}

// Embed returns embeddings for texts in batches of 100. Without an API key it
// returns deterministic stub vectors (dev only; lexical search still works).
func (c *Client) Embed(ctx context.Context, texts []string, taskType string) ([][]float32, error) {
	if c.APIKey == "" || c.quotaExceeded.Load() {
		return StubEmbeddings(texts, c.Dim), nil
	}
	out := make([][]float32, 0, len(texts))
	for start := 0; start < len(texts); start += 100 {
		end := start + 100
		if end > len(texts) {
			end = len(texts)
		}
		batch, err := c.embedBatch(ctx, texts[start:end], taskType)
		if err != nil {
			return nil, err
		}
		out = append(out, batch...)
	}
	return out, nil
}

func (c *Client) embedBatch(ctx context.Context, texts []string, taskType string) ([][]float32, error) {
	reqs := make([]map[string]any, 0, len(texts))
	for _, t := range texts {
		reqs = append(reqs, map[string]any{
			"model": "models/" + c.EmbeddingModel,
			"content": map[string]any{"parts": []map[string]string{{"text": t}}},
			"taskType": taskType,
			"outputDimensionality": c.Dim,
		})
	}
	body, err := json.Marshal(map[string]any{"requests": reqs})
	if err != nil {
		return nil, err
	}
	url := fmt.Sprintf("https://generativelanguage.googleapis.com/v1beta/models/%s:batchEmbedContents?key=%s",
		c.EmbeddingModel, c.APIKey)
	req, err := http.NewRequestWithContext(ctx, http.MethodPost, url, bytes.NewReader(body))
	if err != nil {
		return nil, err
	}
	req.Header.Set("Content-Type", "application/json")
	resp, err := c.HTTP.Do(req)
	if err != nil {
		return nil, err
	}
	defer resp.Body.Close()
	if resp.StatusCode == http.StatusTooManyRequests {
		msg, _ := io.ReadAll(io.LimitReader(resp.Body, 4096))
		c.quotaExceeded.Store(true)
		return nil, QuotaError{Detail: string(msg)}
	}
	if resp.StatusCode >= 500 {
		msg, _ := io.ReadAll(io.LimitReader(resp.Body, 2048))
		c.quotaExceeded.Store(true)
		return nil, fmt.Errorf("gemini embeddings: %s: %s", resp.Status, string(msg))
	}
	if resp.StatusCode != http.StatusOK {
		msg, _ := io.ReadAll(io.LimitReader(resp.Body, 2048))
		return nil, fmt.Errorf("gemini embeddings: %s: %s", resp.Status, string(msg))
	}
	var parsed struct {
		Embeddings []struct {
			Values []float32 `json:"values"`
		} `json:"embeddings"`
	}
	if err := json.NewDecoder(resp.Body).Decode(&parsed); err != nil {
		return nil, err
	}
	if len(parsed.Embeddings) != len(texts) {
		return nil, fmt.Errorf("gemini embeddings: expected %d vectors, got %d", len(texts), len(parsed.Embeddings))
	}
	out := make([][]float32, len(parsed.Embeddings))
	for i, e := range parsed.Embeddings {
		out[i] = e.Values
	}
	return out, nil
}

// StubEmbeddings is deterministic per input text (same text -> same vector).
func StubEmbeddings(texts []string, dim int) [][]float32 {
	out := make([][]float32, len(texts))
	for i, t := range texts {
		vec := make([]float32, dim)
		h := sha256.Sum256([]byte(t))
		seed := binary.LittleEndian.Uint64(h[:8])
		for j := 0; j < dim; j++ {
			seed = seed*6364136223846793005 + 1442695040888963407
			vec[j] = float32(int64(seed>>33)%2000-1000) / 1000.0
		}
		var norm float64
		for _, v := range vec {
			norm += float64(v) * float64(v)
		}
		norm = math.Sqrt(norm)
		if norm > 0 {
			for j := range vec {
				vec[j] = float32(float64(vec[j]) / norm)
			}
		}
		out[i] = vec
	}
	return out
}
