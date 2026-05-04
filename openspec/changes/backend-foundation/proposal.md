# Propuesta — backend-foundation

**Fecha:** 2026-05-03
**Agente:** sdd-propose
**Change name:** `backend-foundation`
**Basado en:** `exploration.md`

---

## Problema que Resuelve

El proyecto M-TurinChess no tiene aún ningún código de producción. Sin un backend funcional, no es posible analizar partidas, comunicarse con Stockfish, ni generar la cadena de entrada para la Máquina de Turing. Esta fase sienta las bases de todo el sistema.

---

## Solución Propuesta

Implementar el **núcleo del backend en Go** que procesa un archivo PGN de extremo a extremo y retorna un resultado de análisis. La MT real y el LLM se incluyen como **stubs** para que el pipeline completo sea funcional y testeable desde el día 1.

### Decisión de Arquitectura Clave (R1 resuelto)

En lugar de implementar la lógica completa del tablero de ajedrez (SAN → FEN), usaremos el protocolo UCI de forma más inteligente:

- **`position startpos moves <lista-uci>`**: Stockfish mantiene el estado del tablero internamente. Solo necesitamos convertir SAN → notación UCI (ej: `Nf3` → `g1f3`), que es una transformación de texto mucho más simple si tenemos el tablero en memoria de forma básica.
- Alternativa aun más simple: parsear las jugadas en **notación larga algebraica** del PGN cuando esté disponible, o mantener un mapa de piezas mínimo solo para resolver ambigüedades SAN.

> **Decisión final:** Implementar un `Board` minimalista en Go que solo rastrée la posición de las piezas para resolver SAN → UCI. No implementamos reglas completas de ajedrez; delegamos la validez al propio Stockfish.

---

## Alcance

### ✅ Incluido

| Componente | Archivo | Descripción |
|------------|---------|-------------|
| Módulo Go | `backend/go.mod` | Inicialización del módulo |
| Entrada del servidor | `backend/cmd/server/main.go` | HTTP server, config |
| Parser PGN | `backend/internal/chess/pgn.go` | Lee .pgn, extrae metadata y jugadas |
| Board minimalista | `backend/internal/chess/board.go` | Rastreo de posición para SAN→UCI |
| Cliente Stockfish | `backend/internal/chess/stockfish.go` | UCI via subproceso, CP Loss |
| Lexer | `backend/internal/chess/lexer.go` | CP Loss + Elo → `{M, E, H}` |
| Analizador | `backend/internal/analysis/analyzer.go` | Orquesta el pipeline completo |
| MT Stub | `backend/internal/turing/machine.go` | Stub que retorna conteo=0 |
| LLM Stub | `backend/internal/llm/gemini.go` | Stub que retorna `null` |
| Persistencia JSON | `backend/internal/db/json_storage.go` | CRUD sobre `data/history.json` |
| API HTTP | `backend/internal/api/` | `POST /api/analyze`, `GET /api/history` |
| CORS Middleware | `backend/internal/api/middleware.go` | Para que React pueda conectarse |

### ❌ Excluido de esta fase

- Implementación real de la Máquina de Turing (→ Fase 2).
- Integración real con Gemini LLM (→ Fase 2).
- Frontend React (→ Fase 3).
- Modo "jugar localmente" (→ Fase 4).

---

## Criterios de Éxito

1. `POST /api/analyze` acepta un `.pgn` y retorna el JSON de respuesta completo (con `tape_input`, `move_details`, veredicto stub).
2. El CP Loss se calcula correctamente usando Stockfish en al menos una partida de prueba.
3. El resultado se persiste en `data/history.json`.
4. `GET /api/history` retorna los análisis guardados.
5. El servidor arranca con `go run ./cmd/server` desde `backend/`.
6. CORS habilitado para `localhost:5173` (puerto por defecto de Vite).
7. Tests unitarios para el Parser PGN y el Lexer pasan.

---

## Riesgos Residuales

| Riesgo | Mitigación |
|--------|-----------|
| Resolución SAN ambigua (ej: dos caballos pueden ir a f3) | El Board minimalista rastrea posiciones; los casos ambiguos se resuelven por archivo/fila |
| Stockfish puede tardar por jugada en depth alto | Default `depth 18`; configurable por el usuario desde el request |
| `history.json` puede corromperse si dos requests llegan simultáneamente | Para MVP local, acceso secuencial es suficiente. Mutex si se necesita |
