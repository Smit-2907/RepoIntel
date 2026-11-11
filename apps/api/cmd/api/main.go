package main

import (
	"context"
	"log"
	"net/http"
	"time"

	"codebaseintel/internal/config"
	"codebaseintel/internal/gemini"
	"codebaseintel/internal/httpapi"
	"codebaseintel/internal/store"
)

func main() {
	ctx := context.Background()
	cfg := config.Load()

	st, err := store.New(ctx, cfg.DatabaseURL)
	if err != nil {
		log.Fatalf("database: %v", err)
	}
	defer st.Close()

	gm := gemini.New(cfg.GeminiAPIKey, cfg.GeminiModel, cfg.EmbeddingModel, cfg.EmbeddingDim)
	srv := httpapi.New(st, gm, cfg)

	httpSrv := &http.Server{
		Addr:              ":" + cfg.APIPort,
		Handler:           srv.Routes(),
		ReadHeaderTimeout: 10 * time.Second,
	}
	mode := "stub"
	if cfg.GeminiAPIKey != "" {
		mode = "gemini"
	}
	log.Printf("api listening on :%s (llm mode: %s)", cfg.APIPort, mode)
	log.Fatal(httpSrv.ListenAndServe())
}
