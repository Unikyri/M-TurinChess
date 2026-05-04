package main

import (
	"log"
	"net/http"
	"os"

	"github.com/Unikyri/M-TurinChess/backend/internal/analysis"
	"github.com/Unikyri/M-TurinChess/backend/internal/api"
	"github.com/joho/godotenv"
	"github.com/Unikyri/M-TurinChess/backend/internal/db"
	"github.com/Unikyri/M-TurinChess/backend/internal/llm"
)

// Config holds all runtime configuration loaded from environment variables.
type Config struct {
	Port          string
	StockfishPath string
	HistoryPath   string
	LLMProvider   string
	LLMAPIKey     string
	LLMModel      string
}

// loadConfig reads configuration from environment variables with sensible defaults.
func loadConfig() Config {
	_ = godotenv.Load()
	return Config{
		Port:          getEnv("PORT", "8080"),
		StockfishPath: getEnv("STOCKFISH_PATH", "../stockfish/stockfish-windows-x86-64-avx2.exe"),
		HistoryPath:   getEnv("HISTORY_PATH", "./data/history.json"),
		LLMProvider:   getEnv("LLM_PROVIDER", "gemini"),
		LLMAPIKey:     getEnv("LLM_API_KEY", ""),
		LLMModel:      getEnv("LLM_MODEL", "gemini-3.1-flash-lite"),
	}
}

func getEnv(key, defaultValue string) string {
	if v, ok := os.LookupEnv(key); ok && v != "" {
		return v
	}
	return defaultValue
}

func main() {
	cfg := loadConfig()

	// Wire up the dependency graph.
	storage := db.NewJSONStorage(cfg.HistoryPath)
	
	llmCfg := llm.Config{
		Provider: cfg.LLMProvider,
		APIKey:   cfg.LLMAPIKey,
		Model:    cfg.LLMModel,
	}
	llmClient := llm.NewClient(llmCfg)
	
	analyzer := analysis.NewAnalyzer(cfg.StockfishPath, storage, llmClient)
	handler := api.NewHandler(analyzer, storage)
	router := api.NewRouter(handler)

	log.Printf("[M-TurinChess] Backend listening on :%s", cfg.Port)
	log.Printf("[M-TurinChess] Stockfish   : %s", cfg.StockfishPath)
	log.Printf("[M-TurinChess] History DB  : %s", cfg.HistoryPath)

	if err := http.ListenAndServe(":"+cfg.Port, router); err != nil {
		log.Fatalf("[M-TurinChess] Server failed: %v", err)
	}
}
