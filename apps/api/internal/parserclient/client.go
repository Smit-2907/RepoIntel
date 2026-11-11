// Package parserclient calls the Python Tree-sitter parsing service (design v1.1 §3 step 5).
package parserclient

import (
	"bytes"
	"context"
	"encoding/json"
	"fmt"
	"io"
	"net/http"
	"time"
)

type Client struct {
	BaseURL string
	HTTP    *http.Client
}

func New(baseURL string) *Client {
	return &Client{BaseURL: baseURL, HTTP: &http.Client{Timeout: 2 * time.Minute}}
}

type File struct {
	Path     string `json:"path"`
	Language string `json:"language"`
	Content  string `json:"content"`
}

type FileInfo struct {
	Path            string `json:"path"`
	Language        string `json:"language"`
	SizeBytes       int    `json:"size_bytes"`
	LineCount       int    `json:"line_count"`
	ParseErrorCount int    `json:"parse_error_count"`
	IsIndexed       bool   `json:"is_indexed"`
}

type Symbol struct {
	Path         string `json:"path"`
	Name         string `json:"name"`
	QualifiedName string `json:"qualified_name"`
	Kind         string `json:"kind"`
	StartLine    int    `json:"start_line"`
	EndLine      int    `json:"end_line"`
	ParentName   string `json:"parent_name"`
	IsExported   bool   `json:"is_exported"`
}

type Import struct {
	FromPath      string `json:"from_path"`
	RawSpecifier  string `json:"raw_specifier"`
	Kind          string `json:"kind"`
}

type Chunk struct {
	Path          string `json:"path"`
	SymbolName    string `json:"symbol_name"`
	StartLine     int    `json:"start_line"`
	EndLine       int    `json:"end_line"`
	Content       string `json:"content"`
	ContentHash   string `json:"content_hash"`
	TokenEstimate int    `json:"token_estimate"`
}

type Meta struct {
	ParserVersion string `json:"parser_version"`
	QueryVersion  string `json:"query_version"`
	DurationMS    int    `json:"duration_ms"`
}

type Result struct {
	Files   []FileInfo `json:"files"`
	Symbols []Symbol   `json:"symbols"`
	Imports []Import   `json:"imports"`
	Chunks  []Chunk    `json:"chunks"`
	Meta    Meta       `json:"meta"`
}

func (c *Client) Parse(ctx context.Context, analysisID string, files []File, maxFiles, maxChunks int) (*Result, error) {
	body, err := json.Marshal(map[string]any{
		"analysis_id": analysisID,
		"files":       files,
		"max_files":   maxFiles,
		"max_chunks":  maxChunks,
	})
	if err != nil {
		return nil, err
	}
	req, err := http.NewRequestWithContext(ctx, http.MethodPost, c.BaseURL+"/parse", bytes.NewReader(body))
	if err != nil {
		return nil, err
	}
	req.Header.Set("Content-Type", "application/json")
	resp, err := c.HTTP.Do(req)
	if err != nil {
		return nil, err
	}
	defer resp.Body.Close()
	if resp.StatusCode != http.StatusOK {
		msg, _ := io.ReadAll(io.LimitReader(resp.Body, 2048))
		return nil, fmt.Errorf("parser service: %s: %s", resp.Status, string(msg))
	}
	var out Result
	if err := json.NewDecoder(resp.Body).Decode(&out); err != nil {
		return nil, err
	}
	return &out, nil
}
