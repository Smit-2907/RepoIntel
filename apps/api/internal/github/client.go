// Package github talks to the GitHub REST API and codeload (design v1.1 §3).
package github

import (
	"context"
	"encoding/json"
	"fmt"
	"io"
	"net/http"
	"regexp"
	"time"
)

var repoURLRe = regexp.MustCompile(`^https://github\.com/([A-Za-z0-9-]{1,39})/([A-Za-z0-9._-]{1,100})(?:\.git)?/?$`)

type Client struct {
	Token string
	HTTP  *http.Client
}

func New(token string) *Client {
	return &Client{Token: token, HTTP: &http.Client{Timeout: 30 * time.Second}}
}

// ParseRepoURL validates and normalizes a GitHub URL (SSRF guard: github.com only).
func ParseRepoURL(raw string) (owner, repo string, err error) {
	m := repoURLRe.FindStringSubmatch(raw)
	if m == nil {
		return "", "", fmt.Errorf("invalid repository URL: must be https://github.com/{owner}/{repo}")
	}
	return m[1], m[2], nil
}

type Metadata struct {
	DefaultBranch string `json:"default_branch"`
	SizeKB        int    `json:"size"`
	Language      string `json:"language"`
	Archived      bool   `json:"archived"`
	CommitSHA     string `json:"-"`
}

type commitInfo struct {
	SHA string `json:"sha"`
}

func (c *Client) get(ctx context.Context, url string, out any) error {
	req, err := http.NewRequestWithContext(ctx, http.MethodGet, url, nil)
	if err != nil {
		return err
	}
	req.Header.Set("Accept", "application/vnd.github+json")
	req.Header.Set("User-Agent", "codebase-intel")
	if c.Token != "" {
		req.Header.Set("Authorization", "Bearer "+c.Token)
	}
	resp, err := c.HTTP.Do(req)
	if err != nil {
		return err
	}
	defer resp.Body.Close()
	if resp.StatusCode != http.StatusOK {
		body, _ := io.ReadAll(io.LimitReader(resp.Body, 2048))
		return fmt.Errorf("github %s: %s: %s", url, resp.Status, string(body))
	}
	return json.NewDecoder(resp.Body).Decode(out)
}

// Metadata fetches repository metadata and resolves the default-branch head SHA.
func (c *Client) Metadata(ctx context.Context, owner, repo string) (Metadata, error) {
	var meta Metadata
	base := fmt.Sprintf("https://api.github.com/repos/%s/%s", owner, repo)
	if err := c.get(ctx, base, &meta); err != nil {
		return meta, err
	}
	var commit commitInfo
	url := fmt.Sprintf("%s/commits/%s", base, meta.DefaultBranch)
	if err := c.get(ctx, url, &commit); err != nil {
		return meta, err
	}
	meta.CommitSHA = commit.SHA
	return meta, nil
}

// DownloadTarball streams a codeload archive pinned to a commit SHA, with a hard byte cap.
func (c *Client) DownloadTarball(ctx context.Context, owner, repo, sha string, maxBytes int64) ([]byte, error) {
	url := fmt.Sprintf("https://codeload.github.com/%s/%s/tar.gz/%s", owner, repo, sha)
	req, err := http.NewRequestWithContext(ctx, http.MethodGet, url, nil)
	if err != nil {
		return nil, err
	}
	req.Header.Set("User-Agent", "codebase-intel")
	if c.Token != "" {
		req.Header.Set("Authorization", "Bearer "+c.Token)
	}
	resp, err := c.HTTP.Do(req)
	if err != nil {
		return nil, err
	}
	defer resp.Body.Close()
	if resp.StatusCode != http.StatusOK {
		return nil, fmt.Errorf("codeload: %s", resp.Status)
	}
	data, err := io.ReadAll(io.LimitReader(resp.Body, maxBytes+1))
	if err != nil {
		return nil, err
	}
	if int64(len(data)) > maxBytes {
		return nil, fmt.Errorf("archive exceeds limit of %d bytes", maxBytes)
	}
	return data, nil
}
