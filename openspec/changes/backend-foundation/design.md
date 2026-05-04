# Diseño Técnico — backend-foundation

**Change:** backend-foundation
**Basado en:** proposal.md + specs/

---

## 1. Estructura de Módulo Go

```
backend/
├── cmd/server/main.go
├── internal/
│   ├── chess/
│   │   ├── pgn.go          // Parser PGN
│   │   ├── board.go        // Board mínimo (SAN → UCI)
│   │   ├── stockfish.go    // Cliente UCI
│   │   └── lexer.go        // CP Loss → {M, E, H}
│   ├── turing/
│   │   └── machine.go      // STUB para esta fase
│   ├── llm/
│   │   └── gemini.go       // STUB para esta fase
│   ├── analysis/
│   │   └── analyzer.go     // Orquestador del pipeline
│   ├── db/
│   │   └── storage.go      // JSON persistence
│   └── api/
│       ├── handler.go
│       ├── router.go
│       └── middleware.go
└── data/
    └── history.json        // Creado en primera escritura
```

---

## 2. Tipos de Dominio

Todos los tipos de dominio se definen en sus paquetes respectivos, sin capa de modelos separada (suficiente para este tamaño).

### `internal/chess`

```go
// pgn.go
type Game struct {
    White    string
    Black    string
    WhiteElo int
    BlackElo int
    Event    string
    Date     string
    Result   string
    Moves    []string // SAN notation
}

// stockfish.go
type EvalResult struct {
    Centipawns int
    IsMate     bool
    MateIn     int
}

type MoveAnalysis struct {
    MoveNumber int
    SAN        string
    UCIMove    string
    CPLoss     int
    BestMove   string // UCI notation
    EvalBefore EvalResult
    EvalAfter  EvalResult
}

// lexer.go
type Symbol = string // "M" | "E" | "H"

type LexerResult struct {
    Tape          []Symbol
    AnalyzedMoves []MoveAnalysis
}
```

### `internal/turing` (Stub)

```go
// machine.go
type MTResult struct {
    Tape2State    []string  // estado de la cinta 2
    SuspicionCount int
    Trace         []TraceStep
}

type TraceStep struct {
    Step      int
    State     string
    ReadC1    string
    ReadC2    string
    Action    string
    Suspicion int
}

// Run - STUB: retorna conteo=0, sin traza real
func Run(tape1 []Symbol) MTResult { ... }
```

### `internal/analysis`

```go
type AnalysisRequest struct {
    PGNBytes    []byte
    PlayerColor string // "white" | "black"
    EloWhite    int
    EloBlack    int
    Threshold   int
    Depth       int
}

type AnalysisResult struct {
    ID                 string
    AnalyzedAt         time.Time
    Verdict            string // "MODULE_DETECTED" | "HUMAN_PLAYER"
    SuspicionCount     int
    Threshold          int
    TotalMovesAnalyzed int
    PlayerColor        string
    Elo                int
    TapeInput          []string
    TapeOutput         []string
    MoveDetails        []MoveDetail
    MTTrace            []turing.TraceStep
}

type MoveDetail struct {
    MoveNumber     int
    SAN            string
    UCI            string
    CPLoss         int
    Classification string // "M" | "E" | "H"
    BestMove       string
    LLMFlag        *string // nil en esta fase
}
```

### `internal/db`

```go
type HistoryRecord struct {
    ID                 string    `json:"id"`
    AnalyzedAt         time.Time `json:"analyzed_at"`
    Verdict            string    `json:"verdict"`
    SuspicionCount     int       `json:"suspicion_count"`
    Threshold          int       `json:"threshold"`
    TotalMovesAnalyzed int       `json:"total_moves_analyzed"`
    PlayerColor        string    `json:"player_color"`
    Elo                int       `json:"elo"`
}

type Storage interface {
    Append(record HistoryRecord) error
    ReadAll() ([]HistoryRecord, error)
}
```

---

## 3. Flujo del Analizador (`analyzer.go`)

```
AnalysisRequest
    │
    ▼
chess.ParsePGN(req.PGNBytes)
    → Game{}
    │
    ▼
Determinar Elo del jugador analizado
  (req.EloWhite si color=white, req.EloBlack si color=black)
    │
    ▼
chess.NewStockfishClient(binaryPath)
    │
    ├── Para cada jugada del jugador (post-filtro apertura por Elo):
    │       stockfish.EvaluateBefore(moveList[:N-1])
    │       stockfish.EvaluateAfter(moveList[:N])
    │       → MoveAnalysis{CPLoss, BestMove, ...}
    │
    ▼
chess.ApplyLexer(moves, playerElo)
    → LexerResult{Tape, AnalyzedMoves}
    │
    ▼
turing.Run(tape)    // STUB en esta fase
    → MTResult{SuspicionCount=0, Trace=[]}
    │
    ▼
Calcular veredicto: suspicionCount >= threshold
    │
    ▼
db.Append(HistoryRecord{...})
    │
    ▼
AnalysisResult{}
```

---

## 4. Board Mínimo (`board.go`)

Estrategia para SAN → UCI:

```
Board mantiene:
  - pieces: map[Square]Piece   // qué pieza está en cada casilla
  - turn: Color                // blancas o negras
  - castlingRights: [4]bool    // K, Q, k, q
  - enPassantSquare: Square    // casilla de en-passant (si aplica)

Para resolver SAN:
  1. Identificar el tipo de pieza desde el SAN (primera letra mayúscula; ninguna = Peón).
  2. Identificar la casilla de destino (últimas 2 chars antes de `+`, `#`, `=`).
  3. Buscar en `pieces` la pieza del color correcto que pueda llegar a esa casilla.
  4. Manejar desambiguación: si hay 2 piezas que pueden ir, usar la fila/columna del SAN.
  5. Construir el string UCI: `<from><to>` + promoción si aplica.
```

**Piezas a implementar:** Rey, Reina, Torre, Alfil, Caballo, Peón (con movimiento, captura, en-passant y coronación).

---

## 5. Cliente Stockfish (`stockfish.go`)

### Patrón de Lectura UCI

```go
// Leer líneas hasta encontrar "bestmove"
// Guardar la última línea "info" con "score cp" o "score mate"
func readUntilBestMove(reader *bufio.Reader) (EvalResult, string, error) {
    var lastScore EvalResult
    var bestMove string
    for {
        line, _ := reader.ReadString('\n')
        line = strings.TrimSpace(line)
        if strings.HasPrefix(line, "info") && strings.Contains(line, "score") {
            lastScore = parseScore(line)
        }
        if strings.HasPrefix(line, "bestmove") {
            parts := strings.Fields(line)
            bestMove = parts[1]
            break
        }
    }
    return lastScore, bestMove, nil
}
```

---

## 6. Servidor HTTP (`cmd/server/main.go`)

```go
func main() {
    cfg := loadConfig()  // puerto, ruta stockfish
    storage := db.NewJSONStorage(cfg.HistoryPath)
    analyzer := analysis.NewAnalyzer(cfg.StockfishPath, storage)
    
    router := api.NewRouter(analyzer)
    
    log.Printf("Server listening on :%s", cfg.Port)
    http.ListenAndServe(":"+cfg.Port, router)
}
```

---

## 7. Configuración

Variables de entorno (con defaults):

| Variable | Default | Descripción |
|----------|---------|-------------|
| `PORT` | `8080` | Puerto del servidor |
| `STOCKFISH_PATH` | `./stockfish/stockfish-windows-x86-64-avx2.exe` | Ruta al binario |
| `HISTORY_PATH` | `./data/history.json` | Ruta al archivo JSON |
| `GEMINI_API_KEY` | `""` | API key de Gemini (stub en esta fase) |

---

## 8. Patrones Aplicados

| Patrón | Dónde | Justificación |
|--------|-------|---------------|
| **Dependency Injection** | `analyzer.go`, `main.go` | Facilita testing con stubs |
| **Interface para Storage** | `internal/db` | Permite cambiar la implementación sin tocar el analizador |
| **Struct config** | `main.go` | Centraliza configuración, evita magia por variables globales |
| **Early return** | Handlers HTTP | Limpieza del flujo de error sin anidamiento |
