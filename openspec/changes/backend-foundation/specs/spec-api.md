# Spec — API HTTP

**Change:** backend-foundation
**Componente:** `internal/api/`

---

## Endpoints Requeridos

### `POST /api/analyze`

**Propósito:** Recibir un PGN, analizarlo y retornar el resultado completo.

**Request:**
```
Content-Type: multipart/form-data

Campos:
  pgn_file     file     (requerido) Archivo .pgn
  player_color string   (requerido) "white" | "black"
  elo_white    integer  (opcional)  Elo de blancas; 0 si no se provee
  elo_black    integer  (opcional)  Elo de negras; 0 si no se provee
  threshold    integer  (opcional)  Umbral de detección; default: 6
  depth        integer  (opcional)  Profundidad Stockfish; default: 18
```

**Response 200 OK:**
```json
{
  "id": "uuid-generado",
  "analyzed_at": "2026-05-03T04:00:00Z",
  "verdict": "MODULE_DETECTED" | "HUMAN_PLAYER",
  "suspicion_count": 8,
  "threshold": 6,
  "total_moves_analyzed": 30,
  "player_color": "white",
  "elo": 1800,
  "tape_input": ["M", "E", "M", "H"],
  "tape_output": ["1", "1", "B", "B"],
  "move_details": [
    {
      "move_number": 11,
      "san": "Nf3",
      "uci": "g1f3",
      "cp_loss": 3,
      "classification": "M",
      "best_move": "g1f3",
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

**Response 400 Bad Request:**
```json
{ "error": "descripción del problema" }
```

**Errores que deben retornar 400:**
- `pgn_file` no incluido en el request.
- `player_color` no es `"white"` ni `"black"`.
- El archivo .pgn no puede parsearse (formato inválido).

**Errores que deben retornar 500:**
- Stockfish no disponible.
- Error interno de análisis.

---

### `GET /api/history`

**Propósito:** Retornar el historial de análisis guardados.

**Response 200 OK:**
```json
[
  {
    "id": "uuid",
    "analyzed_at": "...",
    "verdict": "HUMAN_PLAYER",
    "suspicion_count": 2,
    "threshold": 6,
    "total_moves_analyzed": 28,
    "player_color": "white",
    "elo": 1500
  }
]
```

> El historial **no** incluye `move_details`, `tape_input`, `tape_output`, ni `mt_trace` para mantener el JSON compacto.

---

## CORS

- **Debe** incluir los headers CORS para permitir requests desde `http://localhost:5173`.
- **Debe** manejar preflight `OPTIONS` correctamente.
- Headers requeridos: `Access-Control-Allow-Origin`, `Access-Control-Allow-Methods`, `Access-Control-Allow-Headers`.

---

## Criterios de Aceptación

- [ ] `POST /api/analyze` sin `pgn_file` retorna 400.
- [ ] `POST /api/analyze` con `player_color = "rook"` retorna 400.
- [ ] `POST /api/analyze` con PGN válido retorna 200 con estructura completa.
- [ ] `GET /api/history` retorna array (vacío si no hay análisis).
- [ ] Preflight `OPTIONS /api/analyze` retorna 200 con headers CORS correctos.
