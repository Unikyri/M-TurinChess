# Spec — Lexer (CP Loss → Alfabeto)

**Change:** backend-foundation
**Componente:** `internal/chess/lexer.go`

---

## Comportamiento Esperado

El Lexer **debe** convertir una lista de análisis de jugadas en la cadena de entrada formal `{M, E, H}` para la Máquina de Turing, aplicando el filtro de apertura según Elo.

---

## Contrato de Entrada

```go
type LexerInput struct {
    Moves       []MoveAnalysis // resultado del cliente Stockfish
    PlayerColor string         // "white" | "black"
    PlayerElo   int            // Elo del jugador analizado
}
```

---

## Contrato de Salida

```go
type Symbol string // "M" | "E" | "H"

type LexerOutput struct {
    Tape          []Symbol       // cadena para Cinta 1 de la MT
    FilteredMoves []MoveAnalysis // solo las jugadas analizadas (post-filtro)
}
```

---

## Reglas de Filtrado (Apertura)

El Lexer **debe** omitir las primeras N jugadas del jugador analizado según su Elo:

| Rango de Elo | Jugadas a omitir (del jugador analizado) |
|-------------|------------------------------------------|
| 0 – 1499    | 15 jugadas                               |
| 1500 – 2199 | 10 jugadas                               |
| 2200+       | 6 jugadas                                |

> "Jugadas del jugador analizado" = plies del color analizado, no plies totales.

---

## Reglas de Clasificación

Sobre las jugadas post-filtro:

| CP Loss | Símbolo | Significado |
|---------|---------|-------------|
| 0 – 10  | `M` | Módulo — jugada perfecta |
| 11 – 50 | `E` | Estándar — jugada humana razonable |
| > 50    | `H` | Humano — error evidente |

---

## Casos de Borde

| Caso | Comportamiento |
|------|----------------|
| `PlayerElo == 0` (no disponible) | Usar el filtro más conservador: 15 jugadas |
| Partida con menos jugadas que el filtro | `Tape` vacío, sin error |
| CP Loss exactamente en el límite (10, 50) | Usar los rangos inclusivos del límite inferior |

---

## Criterios de Aceptación

- [ ] Un CP Loss de `10` clasifica como `M`.
- [ ] Un CP Loss de `11` clasifica como `E`.
- [ ] Un CP Loss de `50` clasifica como `E`.
- [ ] Un CP Loss de `51` clasifica como `H`.
- [ ] Con Elo `1200`, se omiten las primeras 15 jugadas del jugador.
- [ ] Con Elo `1800`, se omiten las primeras 10 jugadas del jugador.
- [ ] Con Elo `2500`, se omiten las primeras 6 jugadas del jugador.
- [ ] Con Elo `0`, se omiten las primeras 15 jugadas del jugador.
