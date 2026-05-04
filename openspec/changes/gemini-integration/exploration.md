# Exploración — gemini-integration

**Fecha:** 2026-05-03

---

## Contexto

El flujo de datos (03-data-flow.md) define un paso opcional **entre** el cálculo de CP Loss y el Lexer:

> **Si CP Loss ≤ 10 → consultar Gemini: "¿Es esta jugada humanamente comprensible para un Elo X?"**

Si Gemini responde que **no** es comprensible para ese Elo (jugada demasiado perfecta para un humano de 1200), 
la clasificación podría ajustarse. El stub actual retorna `nil` para todas las jugadas.

## Rol de Gemini en el Pipeline

1. **Cuándo se invoca:** Solo para jugadas clasificadas como `M` (CP Loss 0-10).
2. **Qué pregunta:** "¿Un jugador de Elo X encontraría esta jugada de forma natural?"
3. **Qué retorna:** Un string descriptivo (`"inhuman_for_elo"`, `"natural"`) o `nil` si no aplica.
4. **Efecto en el veredicto:** El `llm_flag` se agrega a `MoveDetail` para visualización en frontend.
   No altera directamente el Lexer ni la MT — es metadata cualitativa complementaria.

## API de Gemini

- **Endpoint:** `https://generativelanguage.googleapis.com/v1beta/models/gemini-2.0-flash:generateContent`
- **Auth:** Query param `?key=API_KEY` (via env var `GEMINI_API_KEY`)
- **Formato:** JSON con `contents[].parts[].text`
- **Modelo:** `gemini-2.0-flash` (rápido, bajo costo)

## Decisiones de Diseño

1. **No bloquear el pipeline:** Si Gemini falla (rate limit, timeout), el análisis continúa sin LLM flag.
2. **Rate limiting:** Un call por jugada M, máximo ~15 calls por partida.
3. **Prompt engineering:** El prompt debe ser conciso y pedir respuesta estructurada (JSON).
4. **Timeout:** 10 segundos por call.
