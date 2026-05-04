# Spec — Simulador de Máquina de Turing de 2 Cintas

**Change:** turing-machine
**Componente:** `internal/turing/machine.go` + `internal/turing/transitions.go`

---

## Contrato de Entrada/Salida

```go
// Input
tape []string  // símbolos de Cinta 1: "M" | "E" | "H"

// Output
MTResult{
    SuspicionCount int         // número de 1s en Cinta 2 al finalizar
    Tape2State     []string    // estado final de C2 (slice de "1" y "B")
    Trace          []TraceStep // registro de cada paso de ejecución
}
```

---

## Comportamiento del Simulador

### Estados
| Estado | Semántica |
|--------|-----------|
| `q0` | Estado inicial y de lectura normal |
| `q_borra` | Borrado de un `1` en Cinta 2 |
| `qf` | Estado final (acepta) |

### Tabla de Transiciones δ

| Estado | C1 | C2 | → C1 escribe | C1 mueve | C2 escribe | C2 mueve | → Estado |
|--------|----|----|-------------|---------|-----------|---------|---------|
| q0 | M | B | M | R | 1 | R | q0 |
| q0 | E | B | E | R | B | S | q0 |
| q0 | H | B | H | R | B | L* | q_borra |
| q0 | B | * | B | S | * | S | qf |
| q_borra | * | 1 | * | S | B | S | q0 |
| q_borra | * | B | * | S | B | S | q0 |

> **L\*:** Si el cabezal C2 está en la posición 0, el movimiento L se convierte en S (guarda de borde).

---

## Invariante de Cinta 2

El cabezal de C2 siempre apunta al **primer blanco a la derecha del último `1`**. La cinta actúa como una pila unaria.

---

## Casos de Borde

| Caso | Comportamiento |
|------|---------------|
| `tape = []` | `SuspicionCount = 0`, `Tape2State = []` |
| Todos `H`, count ya en 0 | Count permanece 0 (guarda de borde) |
| Todos `E` | `SuspicionCount = 0` |
| Todos `M` | `SuspicionCount = len(tape)` |

---

## TraceStep

```go
type TraceStep struct {
    Step      int    // índice 0-based
    State     string // "q0" | "q_borra" | "qf"
    ReadC1    string // símbolo leído en C1
    ReadC2    string // símbolo leído en C2
    Action    string // descripción en lenguaje natural
    Suspicion int    // conteo de 1s en C2 DESPUÉS de este paso
}
```

---

## Criterios de Aceptación

- [ ] `Run(["M","M","H"])` → `SuspicionCount = 1`
- [ ] `Run(["M","M","M","M","M","M"])` → `SuspicionCount = 6`
- [ ] `Run(["H","H","H"])` → `SuspicionCount = 0`
- [ ] `Run([])` → `SuspicionCount = 0`
- [ ] `Run(["M","H","M","H","M"])` → `SuspicionCount = 1`
- [ ] El Trace tiene exactamente `len(tape) + 1` pasos (los `len(tape)` de análisis + el paso final en `qf`)
- [ ] `internal/turing` NO importa ningún paquete externo al propio paquete
