# Exploración — backend-foundation

**Fecha:** 2026-05-03
**Agente:** sdd-explore

---

## Estado Actual del Proyecto

El repositorio existe y contiene:
- `PRD.md` — documento de requisitos completo y actualizado.
- `docs/` — 5 documentos de arquitectura, MT, flujo de datos, estructura y roadmap.
- `stockfish/stockfish-windows-x86-64-avx2.exe` — binario disponible (~114 MB, AVX2).
- `agents/` — sub-agentes SDD bootstrapeados.
- `openspec/` — artifact store inicializado.

**No existe aún ningún código de producción.** El proyecto parte desde cero en términos de implementación.

---

## Análisis de la Fase a Implementar

### Alcance de `backend-foundation`

Según `docs/05-roadmap.md` Fase 1:
- Inicializar módulo Go en `backend/`.
- Parser PGN **custom** (no usa librerías externas — decisión ya tomada).
- Cliente Stockfish UCI vía subproceso (`os/exec`).
- Lexer CP Loss → `{M, E, H}` con filtro por Elo.
- Persistencia JSON (`data/history.json`).
- Servidor HTTP con `POST /api/analyze` y `GET /api/history`.

### Tecnologías Confirmadas
- **Go** — `net/http` estándar (sin frameworks externos).
- **Solo stdlib** — `os/exec`, `bufio`, `encoding/json`, `net/http`, `strings`, `strconv`.
- El MT, LLM y Frontend son **fuera del alcance** de esta fase.

---

## Dependencias Técnicas Identificadas

### 1. Parser PGN Custom
El formato PGN tiene partes bien definidas:
- **Header tags**: `[Key "Value"]` — para extraer Elo, jugadores, evento.
- **Movetext**: secuencia de jugadas en SAN + comentarios + variaciones.

**Complejidad:**
- El PGN tiene variaciones `(...)` y comentarios `{...}` que deben ignorarse.
- Los números de jugada (`1.`, `1...`) son solo separadores.
- El resultado (`1-0`, `0-1`, `1/2-1/2`, `*`) marca el fin.
- La conversión **SAN → FEN** es lo más complejo: requiere conocer el estado completo del tablero (turno, enroque, en-passant).

> ⚠️ **Riesgo Principal**: La conversión SAN → posiciones FEN requiere implementar la lógica de movimiento de piezas de ajedrez en Go desde cero, lo cual es no trivial. Hay que evaluar si se hace completo o se busca una alternativa.

### 2. Cliente Stockfish UCI
El binario ya existe en `stockfish/`. El protocolo UCI es texto plano sobre stdin/stdout:

```
→ uci          ← Inicializa
← uciok        ← Listo
→ isready
← readyok
→ position fen <FEN>
→ go depth 18
← info depth 18 ... score cp <N> ...
← bestmove <move>
```

**Punto crítico:** Stockfish evalúa desde el punto de vista del jugador que mueve. Para calcular CP Loss correctamente hay que:
1. Evaluar la posición **antes** de la jugada (con signo correcto según turno).
2. Evaluar la posición **después** de la jugada.
3. CP Loss = `eval_before - eval_after` (ajustado por turno).

### 3. Lexer con Elo
El filtro de apertura no es fijo en 10 jugadas. Depende del Elo:
- Elo < 1500: omitir primeras 15 jugadas.
- 1500–2200: omitir primeras 10 jugadas.
- 2200+: omitir primeras 6 jugadas (ya tienen más teoría propia).

### 4. Persistencia JSON
Simple: leer `history.json`, append del nuevo análisis, escribir de vuelta.
El JSON debe ser un array de objetos `AnalysisRecord`.

---

## Riesgos y Decisiones Abiertas

| # | Riesgo / Decisión | Severidad | Recomendación |
|---|-------------------|-----------|---------------|
| R1 | SAN → FEN es complejo; implementar desde cero podría llevar semanas | ALTA | Evaluar si podemos pasar las jugadas acumuladas a Stockfish usando `position startpos moves e2e4 e7e5 ...` en lugar de FEN |
| R2 | Stockfish devuelve eval en centi-peones desde perspectiva del jugador activo | MEDIA | Normalizar siempre desde perspectiva de Blancas o del jugador analizado |
| R3 | PGN puede no tener Elo; hay que manejar el fallback al input del usuario | BAJA | El campo `elo_white`/`elo_black` en el request del API cubre esto |
| R4 | Concurrencia: análisis de partidas largas puede tomar >30s | BAJA | Acceptable para MVP local; sin timeout agresivo |

---

## Preguntas Abiertas para sdd-propose

1. **R1 (crítico):** ¿Usamos `position startpos moves ...` (lista de moves en UCI notation) en vez de FEN para evitar implementar la lógica del tablero completa?
2. ¿El endpoint `GET /api/history` está en el alcance de esta fase o solo `POST /api/analyze`?
3. ¿La Máquina de Turing y el LLM son stubs o están completamente fuera del scope de esta fase?
