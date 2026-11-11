package config

import (
	"os"
	"strconv"
)

type Config struct {
	DatabaseURL   string
	ParserURL     string
	APIPort       string
	GitHubPAT     string
	GeminiAPIKey  string
	GeminiModel   string
	EmbeddingModel string
	EmbeddingDim  int
	MaxRepoKB     int
	MaxFiles      int
	MaxFileBytes  int
	MaxChunks     int
	DataDir       string
}

func Load() Config {
	return Config{
		DatabaseURL:    env("DATABASE_URL", "postgres://postgres:postgres@localhost:54322/codebaseintel?sslmode=disable"),
		ParserURL:      env("PARSER_URL", "http://localhost:8090"),
		APIPort:        env("API_PORT", "8080"),
		GitHubPAT:      env("GITHUB_PAT", ""),
		GeminiAPIKey:   env("GEMINI_API_KEY", ""),
		GeminiModel:    env("GEMINI_MODEL", "gemini-3.8-flash"),
		EmbeddingModel: env("EMBEDDING_MODEL", "gemini-embedding-001"),
		EmbeddingDim:   envInt("EMBEDDING_DIM", 1536),
		MaxRepoKB:      envInt("MAX_REPO_KB", 51200),
		MaxFiles:       envInt("MAX_FILES", 3000),
		MaxFileBytes:   envInt("MAX_FILE_BYTES", 524288),
		MaxChunks:      envInt("MAX_CHUNKS", 15000),
		DataDir:        env("DATA_DIR", "/data"),
	}
}

func env(key, def string) string {
	if v := os.Getenv(key); v != "" {
		return v
	}
	return def
}

func envInt(key string, def int) int {
	if v := os.Getenv(key); v != "" {
		if n, err := strconv.Atoi(v); err == nil {
			return n
		}
	}
	return def
}
