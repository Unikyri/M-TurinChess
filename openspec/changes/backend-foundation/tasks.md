# Tasks — backend-foundation

**Change:** backend-foundation
**Basado en:** specs/ + design.md

---

## Batch 1 — Inicialización del Proyecto

> **Objetivo:** El módulo Go existe, compila y el servidor arranca vacío.

- [x] **T1.1** — Crear carpeta `backend/` y ejecutar `go mod init github.com/Unikyri/M-TurinChess/backend`
  - Archivos: `backend/go.mod`
  - Done: `go mod tidy` sale sin errores

- [x] **T1.2** — Crear la estructura de directorios vacíos
  - Crear: `backend/cmd/server/`, `backend/internal/chess/`, `backend/internal/turing/`, `backend/internal/llm/`, `backend/internal/analysis/`, `backend/internal/db/`, `backend/internal/api/`, `backend/data/`
  - Done: árbol de directorios existe

- [x] **T1.3** — Crear `backend/cmd/server/main.go` con servidor HTTP mínimo
  - Servidor que escuche en `PORT` (default `8080`) y responda `200 OK` en `/health`
  - Done: `go run ./cmd/server` arranca sin error y `GET /health` retorna `{"status":"ok"}`

---

## Batch 2 — Parser PGN

> **Objetivo:** Dado un archivo `.pgn`, producir un `Game` struct con metadata y jugadas SAN.

- [x] **T2.1** — Implementar `internal/chess/pgn.go`: función `ParsePGN(data []byte) (Game, error)`
  - Parsear tags `[Key "Value"]`
  - Extraer: White, Black, WhiteElo, BlackElo, Event, Date, Result
  - Done: test unitario con PGN real pasa

- [x] **T2.2** — Implementar parseo del movetext en `pgn.go`
  - Eliminar comentarios `{...}`, variaciones `(...)`, NAGs `$N`
  - Eliminar números de jugada y resultado
  - Retornar `[]string` de jugadas SAN limpias
  - Done: dado PGN con variaciones y comentarios, `Game.Moves` contiene solo jugadas de la línea principal

- [x] **T2.3** — Tests unitarios para `pgn.go`
  - Test: PGN completo con Elo → extrae campos correctamente
  - Test: PGN sin Elo → `WhiteElo == 0`
  - Test: PGN con comentarios y variaciones → `Moves` solo tiene línea principal
  - Test: archivo vacío → error no-nil
  - Done: `go test ./internal/chess/` pasa

---

## Batch 3 — Board Mínimo + Cliente Stockfish

> **Objetivo:** Dado un PGN parseado, calcular CP Loss por jugada comunicándose con Stockfish.

- [x] **T3.1** — Implementar `internal/chess/board.go`: Board mínimo para SAN → UCI
  - Struct `Board` con mapa de piezas y turno
  - Método `Init()`: posición inicial estándar
  - Método `ApplyUCI(move string)`: aplica un movimiento UCI al estado
  - Método `SANToUCI(san string) (string, error)`: convierte SAN a UCI
  - Done: test unitario `SANToUCI("e4") == "e2e4"` (en posición inicial) pasa

- [x] **T3.2** — Implementar `internal/chess/stockfish.go`: `StockfishClient`
  - Struct `StockfishClient` con proceso, stdin writer, stdout reader
  - Método `Start(binaryPath string) error`: inicia el proceso, hace handshake UCI
  - Método `Evaluate(uciMoves []string, depth int) (EvalResult, string, error)`: evalúa una posición
  - Método `Close()`: envía `quit\n` y espera que el proceso termine
  - Done: test de integración `Evaluate([]string{}, 10)` retorna eval de la posición inicial

- [x] **T3.3** — Implementar `AnalyzeGame` en `stockfish.go`
  - Función que recibe `Game`, `playerColor`, `depth` y retorna `[]MoveAnalysis`
  - Bucle: para cada jugada del jugador analizado, llamar `Evaluate` antes y después
  - Calcular `CPLoss = max(0, evalBefore + evalAfter)` (perspectiva normalizada)
  - Manejar `score mate N` → cp equivalente
  - Done: dado PGN de partida conocida, los primeros 3 CP Loss coinciden con análisis manual

---

## Batch 4 — Lexer

> **Objetivo:** Dado `[]MoveAnalysis` y Elo, producir `[]Symbol` para la Cinta 1.

- [x] **T4.1** — Implementar `internal/chess/lexer.go`: función `ApplyLexer(moves []MoveAnalysis, playerElo int) LexerResult`
  - Calcular N jugadas a omitir según Elo (ver spec-lexer.md)
  - Clasificar CP Loss: 0-10→M, 11-50→E, >50→H
  - Done: test unitario pasa los 8 criterios de aceptación del spec

---

## Batch 5 — Stubs, DB y Analizador

> **Objetivo:** Pipeline completo de extremo a extremo con stubs para MT y LLM.

- [x] **T5.1** — Implementar stub `internal/turing/machine.go`
  - Función `Run(tape []string) MTResult` que retorna `MTResult{SuspicionCount: 0, Tape2State: []string{}, Trace: []TraceStep{}}`
  - Done: compila sin error

- [x] **T5.2** — Implementar stub `internal/llm/gemini.go`
  - Función `Analyze(san, elo string) *string` que retorna `nil`
  - Done: compila sin error

- [x] **T5.3** — Implementar `internal/db/storage.go`
  - Interface `Storage` con métodos `Append` y `ReadAll`
  - Struct `JSONStorage` que implementa `Storage` leyendo/escribiendo `history.json`
  - `ReadAll`: si el archivo no existe, retorna slice vacío sin error
  - `Append`: usa mutex para escrituras seguras
  - Done: test unitario append + readAll retorna el registro guardado

- [x] **T5.4** — Implementar `internal/analysis/analyzer.go`
  - Struct `Analyzer` con dependencias: stockfishPath, storage
  - Método `Analyze(req AnalysisRequest) (AnalysisResult, error)` que orquesta el pipeline completo
  - Generar UUID para el ID del análisis (`"crypto/rand"` o `"math/rand"`)
  - Done: test de integración con PGN real retorna `AnalysisResult` con campos poblados

---

## Batch 6 — API HTTP

> **Objetivo:** Servidor HTTP expone los endpoints según spec-api.md con CORS habilitado.

- [x] **T6.1** — Implementar `internal/api/middleware.go`
  - Middleware CORS que agrega los headers para `localhost:5173`
  - Maneja preflight `OPTIONS`
  - Done: `OPTIONS /api/analyze` retorna 200 con headers correctos

- [x] **T6.2** — Implementar `internal/api/handler.go`
  - `HandleAnalyze`: parsea multipart, valida campos, llama `analyzer.Analyze()`, retorna JSON
  - `HandleHistory`: llama `storage.ReadAll()`, retorna JSON
  - Manejo correcto de errores 400 y 500
  - Done: test de integración con request completo retorna 200

- [x] **T6.3** — Implementar `internal/api/router.go`
  - Registrar rutas: `POST /api/analyze`, `GET /api/history`, `GET /health`
  - Aplicar middleware CORS
  - Done: rutas accesibles

- [x] **T6.4** — Actualizar `cmd/server/main.go`
  - Leer config de variables de entorno con defaults
  - Inyectar dependencias: `JSONStorage` → `Analyzer` → `Router`
  - Done: `go run ./cmd/server` arranca, `POST /api/analyze` con PGN real retorna veredicto

---

## Orden de Ejecución para `sdd-apply`

```
Batch 1 → Batch 2 → Batch 3 → Batch 4 → Batch 5 → Batch 6
```

Cada batch debe compilar y pasar sus tests antes de avanzar al siguiente.

---

## Resumen de Archivos a Crear

| Archivo | Batch |
|---------|-------|
| `backend/go.mod` | 1 |
| `backend/cmd/server/main.go` | 1, 6 |
| `backend/internal/chess/pgn.go` | 2 |
| `backend/internal/chess/pgn_test.go` | 2 |
| `backend/internal/chess/board.go` | 3 |
| `backend/internal/chess/stockfish.go` | 3 |
| `backend/internal/chess/lexer.go` | 4 |
| `backend/internal/chess/lexer_test.go` | 4 |
| `backend/internal/turing/machine.go` | 5 |
| `backend/internal/llm/gemini.go` | 5 |
| `backend/internal/db/storage.go` | 5 |
| `backend/internal/db/storage_test.go` | 5 |
| `backend/internal/analysis/analyzer.go` | 5 |
| `backend/internal/api/middleware.go` | 6 |
| `backend/internal/api/handler.go` | 6 |
| `backend/internal/api/router.go` | 6 |
| `backend/data/` (dir) | 1 |
