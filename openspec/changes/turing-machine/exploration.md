# Exploración — turing-machine

**Fecha:** 2026-05-03
**Agente:** sdd-explore

---

## Estado Actual

- `internal/turing/machine.go` existe como **stub** (retorna `SuspicionCount: 0`).
- El resto del pipeline (Lexer → MT → Verdict) ya está integrado en `analyzer.go`.
- La MT real es el último componente faltante para tener un veredicto significativo.

---

## Especificación Formal (desde docs/02-turing-machine.md)

**MT de 2 Cintas, Determinista.**

- **Q** = `{q0, q_borra, qf}`
- **Σ (Cinta 1)** = `{M, E, H, B}` — M=Módulo, E=Estándar, H=Humano, B=Blanco
- **Γ (Cinta 2)** = `{1, B}` — unario, B=Blanco
- **Estado inicial:** `q0`
- **Estado de aceptación:** `qf`

### Tabla de Transiciones δ

| Estado | C1 | C2 | C1 escribe | C1 mueve | C2 escribe | C2 mueve | Siguiente |
|--------|----|----|-----------|---------|-----------|---------|-----------|
| q0     | M  | B  | M         | R       | 1         | R       | q0        |
| q0     | E  | B  | E         | R       | B         | S       | q0        |
| q0     | H  | B  | H         | R       | B         | L*      | q_borra   |
| q0     | B  | B  | B         | S       | B         | S       | qf        |
| q_borra| *  | 1  | *         | S       | B         | S       | q0        |
| q_borra| *  | B  | *         | S       | B         | S       | q0        |

> **L\* con guarda de borde:** si el cabezal de C2 está en la posición 0, `L` equivale a `S` (no sale del tape).

### Semántica

- **M** (jugada perfecta): sospecha incrementa → escribe `1`, mueve C2 derecha.
- **E** (jugada estándar): neutral → C2 no cambia.
- **H** (jugada humana/error): sospecha decrementa → mueve C2 izquierda, `q_borra` borra el `1` (si existe).
- **qf** al leer `B` en C1 → fin de análisis.
- **Veredicto:** `count(1s en C2) >= umbral` → `MODULE_DETECTED`.

---

## Invariante de C2

El cabezal de C2 siempre apunta al **primer blanco después del último `1`**. C2 actúa como una pila (stack) en unario.

---

## Preguntas Resueltas

1. **Cómo reportar el conteo:** contar los `1`s en C2 al llegar a `qf`.
2. **Borde cuando count=0 y llega H:** el guarda impide ir a posición negativa; `q_borra` ve B y no hace nada.
3. **Trace:** cada paso se registra como `TraceStep` para visualización en el frontend.
