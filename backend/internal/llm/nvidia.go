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

// NvidiaClient communicates with the NVIDIA NIM API for chess move peritaje.
// A nil NvidiaClient is safe to use — all methods return nil gracefully.
type NvidiaClient struct {
	apiKey     string
	model      string
	httpClient *http.Client
}

// NewNvidiaClient creates a client for the NVIDIA NIM API.
func NewNvidiaClient(apiKey, model string) *NvidiaClient {
	if model == "" {
		model = "deepseek-ai/deepseek-v4-pro"
	}
	return &NvidiaClient{
		apiKey: apiKey,
		model:  model,
		httpClient: &http.Client{
			Timeout: 20 * time.Second,
		},
	}
}

// --- OpenAI API request/response types (used by NVIDIA) ---

type openAIRequest struct {
	Model       string          `json:"model"`
	MaxTokens   int             `json:"max_tokens,omitempty"`
	Temperature float32         `json:"temperature,omitempty"`
	Messages    []openAIMessage `json:"messages"`
}

type openAIMessage struct {
	Role    string `json:"role"`
	Content string `json:"content"`
}

type openAIResponse struct {
	Choices []struct {
		Message struct {
			Content string `json:"content"`
		} `json:"message"`
	} `json:"choices"`
}

// peritajeResult is the expected JSON structure from the LLM.
type peritajeResult struct {
	Flag   string `json:"flag"`   // "natural" | "suspicious" | "inhuman_for_elo"
	Reason string `json:"reason"` // brief explanation
}

func (nc *NvidiaClient) AnalyzeBatch(sans []string, elo int, context string) map[string]string {
	results := make(map[string]string)
	if nc == nil || len(sans) == 0 {
		return results
	}

	prompt := buildNvidiaBatchPrompt(sans, elo, context)

	resultText, err := nc.callAPI(prompt)
	if err != nil {
		log.Printf("[nvidia] API error for batch: %v", err)
		return results
	}

	return parseNvidiaBatchResponse(*resultText)
}

func buildNvidiaBatchPrompt(sans []string, elo int, context string) string {
	sanList := strings.Join(sans, ", ")
	return fmt.Sprintf(`Eres un Gran Maestro y experto en psicología del ajedrez y detección de motores (anti-cheating).
Se te entrega una lista de jugadas que Stockfish evaluó como "perfectas" (CP Loss de 0-10).

CONTEXTO DE LA PARTIDA (Primeras jugadas):
%s

JUGADAS A EVALUAR (Separadas por comas): [%s]
ELO DEL JUGADOR: %d

Tu tarea es EVALUAR Y DECIDIR de forma concluyente si un humano de este Elo (%d) encontraría la jugada de forma natural. 
ESTÁ PROHIBIDO responder con frases ambiguas como "requiere mayor análisis", "necesita un estudio profundo" o "podría ser agresiva". DEBES dar un veredicto definitivo.

Para cada jugada, piensa:
1. Si el Elo es bajo (< 1500) y la jugada requiere un cálculo táctico brillante de 5 movimientos en el medio juego: Es sospechoso ("M").
2. Si es una captura obvia, jaque mate en 1, o desarrollo natural, incluso si es "perfecta", es fácil para un humano: Es humano ("H").
3. Si pertenece a teoría de aperturas estándar: Es estándar ("E").

Clasifica CADA jugada estrictamente como:
- "M" (Módulo/Sospechosa): Tácticamente demasiado compleja para su Elo.
- "H" (Humana): Natural, obvia, forzada, o fácil de encontrar para su Elo.
- "E" (Estándar): Teoría de aperturas.

Responde ÚNICAMENTE con un arreglo JSON válido, sin bloques de código markdown:
[
  {
    "san": "e4",
    "flag": "M",
    "reason": "Veredicto directo: Demasiado compleja para su Elo. (Máximo 15 palabras)"
  }
]`, context, sanList, elo, elo)
}

func (nc *NvidiaClient) callAPI(prompt string) (*string, error) {
	url := "https://integrate.api.nvidia.com/v1/chat/completions"

	reqBody := openAIRequest{
		Model:       nc.model,
		MaxTokens:   2048,
		Temperature: 0.2,
		Messages: []openAIMessage{
			{Role: "user", Content: prompt},
		},
	}

	jsonBody, err := json.Marshal(reqBody)
	if err != nil {
		return nil, fmt.Errorf("marshal request: %w", err)
	}

	var lastErr error
	maxRetries := 3
	backoff := 2 * time.Second

	for i := 0; i < maxRetries; i++ {
		req, err := http.NewRequest("POST", url, bytes.NewReader(jsonBody))
		if err != nil {
			return nil, fmt.Errorf("create request: %w", err)
		}
		req.Header.Set("Authorization", "Bearer "+nc.apiKey)
		req.Header.Set("Content-Type", "application/json")

		resp, err := nc.httpClient.Do(req)
		if err != nil {
			lastErr = fmt.Errorf("http do: %w", err)
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
			log.Printf("[nvidia] Rate limited (429), retrying in %v...", backoff)
			time.Sleep(backoff)
			backoff *= 2
			continue
		}

		if resp.StatusCode != http.StatusOK {
			return nil, fmt.Errorf("status %d: %s", resp.StatusCode, string(body[:min(len(body), 200)]))
		}

		// Small delay to prevent hitting rate limits on the next immediate call
		time.Sleep(500 * time.Millisecond)

		return parseNvidiaResponse(body)
	}

	return nil, fmt.Errorf("max retries exceeded, last error: %v", lastErr)
}

func parseNvidiaResponse(body []byte) (*string, error) {
	var resp openAIResponse
	if err := json.Unmarshal(body, &resp); err != nil {
		return nil, fmt.Errorf("unmarshal response: %w", err)
	}

	if len(resp.Choices) == 0 {
		return nil, fmt.Errorf("empty choices from NVIDIA API")
	}

	text := strings.TrimSpace(resp.Choices[0].Message.Content)

	// Strip markdown code fences if LLM wraps the JSON.
	text = strings.TrimPrefix(text, "```json")
	text = strings.TrimPrefix(text, "```")
	text = strings.TrimSuffix(text, "```")
	text = strings.TrimSpace(text)

	return &text, nil
}

func parseNvidiaBatchResponse(text string) map[string]string {
	results := make(map[string]string)
	
	// Try parsing it
	var arr []peritajeNvidiaResultBatch
	if err := json.Unmarshal([]byte(text), &arr); err != nil {
		log.Printf("[nvidia] Could not parse batch JSON. Err: %v\nRaw text: %s", err, text)
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

type peritajeNvidiaResultBatch struct {
	SAN    string `json:"san"`
	Flag   string `json:"flag"`
	Reason string `json:"reason"`
}

func min(a, b int) int {
	if a < b {
		return a
	}
	return b
}
