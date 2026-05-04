# Spec — Cliente Stockfish UCI

**Change:** backend-foundation
**Componente:** `internal/chess/stockfish.go` + `internal/chess/board.go`

---

## Comportamiento Esperado

El cliente **debe** gestionar el ciclo de vida del subproceso Stockfish y calcular el **CP Loss** por cada jugada de un jugador dado.

---

## Ciclo de Vida del Subproceso

1. **Inicio:** `exec.Command(binaryPath)` con pipes stdin/stdout.
2. **Handshake:** Enviar `uci\n` → esperar línea `uciok`.
3. **Ready check:** Enviar `isready\n` → esperar `readyok`.
4. **Configuración:** Enviar opciones de rendimiento si se necesita.
5. **Análisis:** Bucle por cada jugada (ver abajo).
6. **Cierre:** Enviar `quit\n` al terminar.

---

## Protocolo de Evaluación por Jugada

Para calcular el CP Loss de la jugada N del jugador analizado:

### Paso 1 — Evaluar ANTES de la jugada N
```
→ position startpos moves <m1> <m2> ... <m(N-1)>
→ go depth <depth>
← info ... score cp <eval_before> ...
← bestmove <move>
```

### Paso 2 — Evaluar DESPUÉS de la jugada N
```
→ position startpos moves <m1> <m2> ... <m(N-1)> <mN>
→ go depth <depth>
← info ... score cp <eval_after> ...
← bestmove <move>
```

### Cálculo de CP Loss
```
// La evaluación siempre es desde el punto de vista del jugador que mueve.
// Después de la jugada N, es el turno del oponente, así que el score se invierte.
cp_loss = eval_before - (-eval_after)
cp_loss = eval_before + eval_after   // si ambos están en perspectiva del jugador que movió
cp_loss = max(0, cp_loss)            // nunca negativo
```

> **Nota sobre mate:** Si Stockfish reporta `score mate N`, se convierte a `cp = 10000 - N` (mate en N = casi perfecto). Si es `score mate -N` (el oponente lo mata), `cp = -(10000 - N)`.

---

## Interfaz Requerida

```go
type EvalResult struct {
    ScoreCentipawns int
    IsMate          bool
    MateIn          int  // positivo = gana, negativo = pierde
}

type MoveAnalysis struct {
    MoveNumber int
    SAN        string
    UCIMove    string
    CPLoss     int
    BestMove   string
}
```

---

## Reglas

- **Debe** leer la **última** línea `info ... score` antes del `bestmove` (es la de mayor depth).
- **Debe** manejar `score mate N` como cp equivalente.
- **Debe** retornar error si el proceso Stockfish no responde en `30s`.
- **No debe** reiniciar Stockfish entre jugadas; usa una sola instancia por análisis.
- **Debe** cerrar el subproceso al terminar el análisis (incluso si hay error).

---

## SAN → UCI (Board mínimo)

El Board minimalista **debe**:
- Inicializarse desde la posición inicial estándar.
- Aplicar movimientos UCI para actualizar el estado.
- Resolver SAN → UCI con la posición actual.

El Board **no debe**:
- Validar movimientos ilegales.
- Implementar lógica de jaque, jaque mate, tablas por repetición, etc.

---

## Criterios de Aceptación

- [ ] Dado el PGN de una partida conocida, el CP Loss de la jugada 1 de Stockfish coincide con el calculado manualmente.
- [ ] Si Stockfish no está en la ruta indicada, retorna error claro: `"stockfish binary not found at <path>"`.
- [ ] El subproceso se cierra correctamente después del análisis (sin procesos zombie).
- [ ] Un score `mate 2` se convierte a `cp = 9998`.
