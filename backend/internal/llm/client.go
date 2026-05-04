package llm

import "log"

// Client defines the interface for AI move analysis.
type Client interface {
	AnalyzeBatch(sans []string, elo int, context string) map[string]string
}

// Config holds the configuration for the LLM Provider.
type Config struct {
	Provider string // e.g., "gemini", "nvidia"
	APIKey   string
	Model    string
}

// NewClient is a factory that returns the appropriate LLM Client based on config.
func NewClient(cfg Config) Client {
	if cfg.APIKey == "" {
		log.Println("[llm] No API key provided — LLM analysis disabled")
		return nil
	}

	switch cfg.Provider {
	case "gemini":
		return NewGeminiClient(cfg.APIKey, cfg.Model)
	case "nvidia":
		return NewNvidiaClient(cfg.APIKey, cfg.Model)
	default:
		log.Printf("[llm] Unknown provider '%s', LLM analysis disabled", cfg.Provider)
		return nil
	}
}
