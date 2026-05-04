# Tasks — gemini-integration

**Change:** gemini-integration

---

## Batch 1 — Cliente Gemini + Integración en Analyzer

> **Objetivo:** Reemplazar el stub LLM con un cliente Gemini real y conectarlo al pipeline de análisis.

- [x] **T1.1** — Implementar `internal/llm/gemini.go`: `GeminiClient`
  - Struct con `apiKey`, `httpClient` (timeout 10s), `model` (gemini-2.0-flash)
  - `NewGeminiClient(apiKey string) *GeminiClient` — retorna nil si apiKey vacía
  - `Analyze(san string, elo int, context string) *string` — prompt + HTTP POST + parse respuesta
  - Graceful degradation: si error HTTP o timeout → retorna nil (no rompe el pipeline)
  - Done: compila sin error

- [x] **T1.2** — Diseñar el prompt de peritaje
  - Prompt: "Eres un perito de ajedrez. Evalúa si la jugada {SAN} en el contexto {context} es natural para un jugador de Elo {elo}. Responde SOLO con un JSON: {\"flag\": \"natural\" | \"suspicious\" | \"inhuman_for_elo\", \"reason\": \"breve explicación\"}"
  - Parsear respuesta JSON del campo `text` de Gemini
  - Done: test manual confirma respuesta parseable

- [x] **T1.3** — Integrar Gemini en `internal/analysis/analyzer.go`
  - Agregar `geminiClient *llm.GeminiClient` al struct `Analyzer`
  - En Step 7 (build move details): si `Classification == "M"` y `geminiClient != nil`, llamar `Analyze()`
  - El `llm_flag` se asigna al `MoveDetail.LLMFlag`
  - Done: `go build ./...` OK

- [x] **T1.4** — Actualizar `cmd/server/main.go`
  - Pasar `cfg.GeminiAPIKey` a `NewAnalyzer()`
  - Done: servidor arranca sin error con y sin API key

---

## Archivos Modificados

| Archivo | Acción |
|---------|--------|
| `internal/llm/gemini.go` | Reescribir (reemplaza stub) |
| `internal/analysis/analyzer.go` | Modificar (agregar GeminiClient + llamadas) |
| `cmd/server/main.go` | Modificar (pasar API key) |
