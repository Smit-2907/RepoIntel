package main

import (
	"context"
	"log"
	"os"
	"os/signal"
	"syscall"

	"codebaseintel/internal/config"
	"codebaseintel/internal/gemini"
	"codebaseintel/internal/github"
	"codebaseintel/internal/parserclient"
	"codebaseintel/internal/pipeline"
	"codebaseintel/internal/store"
)

func main() {
	cfg := config.Load()
	ctx, cancel := signal.NotifyContext(context.Background(), os.Interrupt, syscall.SIGTERM)
	defer cancel()

	st, err := store.New(ctx, cfg.DatabaseURL)
	if err != nil {
		log.Fatalf("database: %v", err)
	}
	defer st.Close()

	pl := pipeline.New(
		st,
		github.New(cfg.GitHubPAT),
		parserclient.New(cfg.ParserURL),
		gemini.New(cfg.GeminiAPIKey, cfg.GeminiModel, cfg.EmbeddingModel, cfg.EmbeddingDim),
		cfg,
	)
	mode := "stub"
	if cfg.GeminiAPIKey != "" {
		mode = "gemini"
	}
	log.Printf("worker started (parser: %s, llm mode: %s)", cfg.ParserURL, mode)
	pl.Work(ctx, "worker-1")
}
