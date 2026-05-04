# Flujo de Datos — End to End

> Desde el archivo PGN hasta el veredicto final.

## Flujo Completo

```mermaid
sequenceDiagram
    participant U as Usuario
    participant FE as Frontend (React)
    participant API as Go API
    participant PGN as Parser PGN
    participant SF as Stockfish (UCI)
    participant LLM as Gemini LLM
    participant LEX as Lexer
    participant MT as Máquina de Turing
    participant V as Motor de Veredicto
    participant DB as JSON DB

    U->>FE: Sube archivo .pgn (+ Elo si falta)
    FE->>API: POST /api/analyze (multipart)
    
    API->>PGN: Parsear PGN
    PGN-->>API: Lista de jugadas + posiciones FEN
    
    loop Por cada jugada (después de filtro inicial por Elo)
        API->>SF: position fen ... / go depth N
        SF-->>API: bestmove + eval (centipawns)
        API->>API: Calcular CP Loss
    end
    
    opt Jugadas con CP Loss ≤ 10
        API->>LLM: ¿Es humanamente comprensible para un Elo X?
        LLM-->>API: Clasificación (humana/inhumana)
    end
    
    API->>LEX: CP Loss[] + Elo → Alfabeto {M, E, H}
    LEX-->>API: Cadena formal (ej: "MMEHMMH")
    
    API->>MT: Ejecutar MT con Cinta 1 = cadena
    MT-->>API: Estado Cinta 2 (conteo de 1s)
    
    API->>V: Conteo vs Umbral (dinámico)
    V-->>API: Veredicto + metadata
    
    API->>DB: Guardar análisis en JSON
    DB-->>API: OK
    
    API-->>FE: JSON Response
    FE-->>U: Visualización del resultado
```

---

## Estructura del Request

### `POST /api/analyze`

```
Content-Type: multipart/form-data

Fields:
  - pgn_file: archivo .pgn
  - player_color: "white" | "black"
  - elo_white: integer (detectado o ingresado)
  - elo_black: integer (detectado o ingresado)
  - threshold: integer (umbral configurable)
  - depth: integer (profundidad de análisis Stockfish, default: 18)
```

---

## Estructura del Response

```json
{
  "verdict": "MODULE_DETECTED" | "HUMAN_PLAYER",
  "suspicion_count": 8,
  "threshold": 6,
  "total_moves_analyzed": 30,
  "tape_input": ["M", "E", "M", "M", "H", "E", "M", "M", "..."],
  "tape_output": ["1", "1", "1", "1", "1", "1", "1", "1", "B", "B"],
  "move_details": [
    {
      "move_number": 11,
      "san": "Nf3",
      "cp_loss": 3,
      "classification": "M",
      "best_move": "Nf3",
      "llm_flag": null
    }
  ],
  "mt_trace": [
    {
      "step": 0,
      "state": "q0",
      "read_c1": "M",
      "read_c2": "B",
      "action": "write 1, move R",
      "suspicion": 1
    }
  ]
}
```

---

## Pipeline Interno (Go)

```mermaid
graph TD
    A["PGN File + Elo"] --> B["Parser PGN\n(custom)"]
    B --> C["Lista de Moves\n+ Posiciones FEN"]
    C --> D{"Ply > Threshold(Elo)?"}
    D -->|No| Skip["Omitir\n(apertura)"]
    D -->|Sí| E["Stockfish\nposition fen ...\ngo depth N"]
    E --> F["Eval antes\nvs Eval después"]
    F --> G["CP Loss = |eval_before - eval_after|"]
    G --> H{"CP Loss ≤ 10?"}
    H -->|Sí| I["Gemini LLM\n(Analiza con Elo)"]
    H -->|No| J["Clasificar directamente"]
    I --> K["Clasificación final"]
    J --> K
    K --> L["Lexer\nCP Loss → {M, E, H}"]
    L --> M_C["Cadena para Cinta 1"]
    M_C --> N["Simulador MT\nRun()"]
    N --> O["Conteo de 1s\nen Cinta 2"]
    O --> P{"Conteo ≥ Umbral?"}
    P -->|Sí| Q["🚨 MODULE_DETECTED"]
    P -->|No| R["✅ HUMAN_PLAYER"]
    Q --> S["Guardar en DB JSON"]
    R --> S
```
