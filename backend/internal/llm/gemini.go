package llm

import (
	"bytes"
	"encoding/json"
	"fmt"
	"io"
	"log"
	"net/http"
	"strings"
	"time"
)

// GeminiClient communicates with the Gemini API.
type GeminiClient struct {
	apiKey     string
	model      string
	httpClient *http.Client
}

func NewGeminiClient(apiKey, model string) *GeminiClient {
	if model == "" {
		model = "gemini-3.1-flash-lite-preview"
	}
	return &GeminiClient{
		apiKey: apiKey,
		model:  model,
		httpClient: &http.Client{
			Timeout: 60 * time.Second, // 15s was too short for Gemini under load
		},
	}
}

type geminiRequest struct {
	Contents []geminiContent `json:"contents"`
}

type geminiContent struct {
	Parts []geminiPart `json:"parts"`
}

type geminiPart struct {
	Text string `json:"text"`
}

type geminiResponse struct {
	Candidates []struct {
		Content struct {
			Parts []struct {
				Text string `json:"text"`
			} `json:"parts"`
		} `json:"content"`
	} `json:"candidates"`
}

// AnalyzeBatch implements the Client interface for Gemini.
func (gc *GeminiClient) AnalyzeBatch(sans []string, elo int, context string) map[string]string {
	results := make(map[string]string)
	if gc == nil || len(sans) == 0 {
		return results
	}

	prompt := buildBatchPrompt(sans, elo, context)
	resultText, err := gc.callAPIWithRetries(prompt)
	if err != nil {
		log.Printf("[gemini] API error for batch: %v", err)
		return results
	}

	return parseBatchResponse(*resultText)
}

func (gc *GeminiClient) callAPIWithRetries(prompt string) (*string, error) {
	// Strict pacing: Gemini free tier allows 14 RPM (1 request per ~4.3s).
	// Sleeping 5 seconds BEFORE every request guarantees we never exceed the rate limit,
	// regardless of success or failure of previous requests.
	time.Sleep(5 * time.Second)

	url := fmt.Sprintf("https://generativelanguage.googleapis.com/v1beta/models/%s:generateContent?key=%s", gc.model, gc.apiKey)
	reqBody := geminiRequest{Contents: []geminiContent{{Parts: []geminiPart{{Text: prompt}}}} }

	jsonBody, err := json.Marshal(reqBody)
	if err != nil {
		return nil, fmt.Errorf("marshal request: %w", err)
	}

	var lastErr error
	maxRetries := 5 // Increased retries to handle transient spikes
	backoff := 2 * time.Second

	for i := 0; i < maxRetries; i++ {
		resp, err := gc.httpClient.Post(url, "application/json", bytes.NewReader(jsonBody))
		if err != nil {
			lastErr = fmt.Errorf("http post: %w", err)
			time.Sleep(backoff)
			backoff *= 2
			continue
		}

		body, err := io.ReadAll(resp.Body)
		resp.Body.Close()
		
		if err != nil {
			lastErr = fmt.Errorf("read response: %w", err)
			time.Sleep(backoff)
			backoff *= 2
			continue
		}

		if resp.StatusCode == http.StatusTooManyRequests {
			lastErr = fmt.Errorf("status 429: %s", string(body[:min(len(body), 200)]))
			log.Printf("[gemini] Rate limited (429): %s. Retrying in %v...", string(body[:min(len(body), 200)]), backoff)
			time.Sleep(backoff)
			backoff *= 2
			continue
		}

		if resp.StatusCode != http.StatusOK {
			return nil, fmt.Errorf("status %d: %s", resp.StatusCode, string(body[:min(len(body), 200)]))
		}

		return parseGeminiResponse(body)
	}

	return nil, fmt.Errorf("max retries exceeded, last error: %v", lastErr)
}

func parseGeminiResponse(body []byte) (*string, error) {
	var resp geminiResponse
	if err := json.Unmarshal(body, &resp); err != nil {
		log.Printf("[gemini] raw body on unmarshal error: %s", string(body[:min(len(body), 300)]))
		return nil, fmt.Errorf("unmarshal response: %w", err)
	}

	if len(resp.Candidates) == 0 || len(resp.Candidates[0].Content.Parts) == 0 {
		log.Printf("[gemini] empty candidates — raw body: %s", string(body[:min(len(body), 300)]))
		return nil, fmt.Errorf("empty response")
	}

	text := strings.TrimSpace(resp.Candidates[0].Content.Parts[0].Text)
	text = strings.TrimPrefix(text, "```json")
	text = strings.TrimPrefix(text, "```")
	text = strings.TrimSuffix(text, "```")
	text = strings.TrimSpace(text)

	log.Printf("[gemini] raw LLM text: %s", text[:min(len(text), 200)])
	return &text, nil
}

func parseBatchResponse(text string) map[string]string {
	results := make(map[string]string)
	var arr []peritajeResultBatch
	if err := json.Unmarshal([]byte(text), &arr); err != nil {
		log.Printf("[gemini] Could not parse batch JSON, returning empty map. Err: %v", err)
		return results
	}

	for _, item := range arr {
		if item.Flag == "" {
			continue
		}
		flag := item.Flag
		if item.Reason != "" {
			flag = item.Flag + ": " + item.Reason
		}
		results[item.SAN] = flag
	}
	return results
}

type peritajeResultBatch struct {
	SAN    string `json:"san"`
	Flag   string `json:"flag"`
	Reason string `json:"reason"`
}

func buildBatchPrompt(sans []string, elo int, context string) string {
	sanList := strings.Join(sans, ", ")
	return fmt.Sprintf(`Eres un árbitro experto de ajedrez antitrampas.
Analiza las siguientes jugadas sospechosas: [%s].
Contexto (primeras jugadas para entender la apertura): %s
Elo del jugador: %d.

Clasifica CADA jugada estrictamente como:
- "M" (Jugada de Módulo): Si es tácticamente perfecta, inhumana, o antinatural para el Elo.
- "H" (Jugada Humana): Si es natural, tiene errores o imprecisiones típicas.
- "E" (Jugada Estándar): Si es teoría de aperturas evidente o un movimiento obvio/forzado.

Responde ÚNICAMENTE con un arreglo JSON válido, sin markdown ni comillas invertidas, con esta estructura:
[
  {
    "san": "e4",
    "flag": "M",
    "reason": "breve justificación de 10 palabras"
  }
]`, sanList, context, elo)
}
