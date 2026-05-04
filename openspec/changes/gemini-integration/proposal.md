# Propuesta — gemini-integration

**Fecha:** 2026-05-03
**Change name:** `gemini-integration`

---

## Problema

El stub `llm.Analyze()` siempre retorna `nil`. Las jugadas clasificadas como M (perfectas) no tienen contexto cualitativo sobre si son naturales para el Elo del jugador.

## Solución

Reemplazar el stub por un cliente HTTP que consulta la API de Gemini 2.0 Flash con un prompt especializado en ajedrez.

## Alcance

### ✅ Incluido

| Componente | Archivo |
|------------|---------|
| Cliente Gemini con prompt y parsing | `internal/llm/gemini.go` |
| Integración en el Analyzer | `internal/analysis/analyzer.go` (inyectar API key + llamar LLM en jugadas M) |
| Actualización de `main.go` | Pasar `GEMINI_API_KEY` al Analyzer |

### ❌ Excluido

- Cambios en el Lexer o la MT (el LLM flag es metadata, no altera la cinta).
- Cache de respuestas (optimización futura).

## Criterios de Éxito

1. `GeminiClient.Analyze("Nf3", 1200)` retorna un `*string` no-nil con la clasificación.
2. Si `GEMINI_API_KEY` está vacía, el LLM se desactiva silenciosamente (no crashea).
3. Si Gemini falla (timeout, error HTTP), la jugada queda con `llm_flag: null` y el pipeline continúa.
4. El Analyzer solo invoca Gemini para jugadas con `Classification == "M"`.
5. `go build ./...` compila sin error.
