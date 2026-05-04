# Diseño Técnico — turing-machine

**Change:** turing-machine

---

## Estructura de Archivos

```
internal/turing/
├── doc.go          (ya existe)
├── machine.go      ← tipos + Run() — reemplaza el stub
└── machine_test.go ← tests unitarios con trazas verificadas
```

> **Sin `transitions.go` separado:** la tabla δ es pequeña (6 reglas). Se codifica inline en `machine.go` usando un `map` o `switch` dentro de `step()`.

---

## Tipos Go

```go
// state representa los estados internos de la MT.
type state int
const (
    q0     state = iota
    qBorra
    qF
)

// direction representa el movimiento del cabezal.
type direction int
const (
    L direction = -1
    S direction = 0
    R direction = +1
)

// tape es una cinta dinámica que crece a la derecha según necesidad.
// El valor "" (string vacío) representa el símbolo blanco B.
type tape struct {
    cells []string
    head  int
}
func (t *tape) read() string      // retorna cells[head] o "" si fuera de rango
func (t *tape) write(sym string)  // escribe en cells[head], extiende si necesario
func (t *tape) move(d direction)  // mueve head; guarda: head >= 0 siempre

// stepResult es el resultado de aplicar una transición.
type stepResult struct {
    writeC1 string
    moveC1  direction
    writeC2 string
    moveC2  direction
    next    state
}
```

---

## Función `Run`

```go
func Run(input []string) MTResult {
    // Inicializar Cinta 1 con los símbolos de input + blanco final.
    // Inicializar Cinta 2 vacía (todo blancos).
    // Inicializar estado = q0.
    //
    // Loop:
    //   1. Leer C1 y C2
    //   2. Aplicar δ(state, c1, c2) → stepResult
    //   3. Registrar TraceStep
    //   4. Aplicar escrituras y movimientos
    //   5. Actualizar estado
    //   6. Si estado == qF → terminar
    //
    // Calcular SuspicionCount = contar "1" en C2
    // Retornar MTResult
}
```

---

## Función de Transición δ

```go
func delta(st state, c1, c2 string) stepResult {
    switch st {
    case q0:
        switch c1 {
        case "M": return stepResult{"M", R, "1", R, q0}
        case "E": return stepResult{"E", R, c2, S, q0}
        case "H": return stepResult{"H", R, c2, L, qBorra}
        default:  return stepResult{c1, S, c2, S, qF}  // B o desconocido → qf
        }
    case qBorra:
        switch c2 {
        case "1": return stepResult{c1, S, "B", S, q0}  // erase
        default:  return stepResult{c1, S, c2, S, q0}   // nothing to erase
        }
    }
    return stepResult{c1, S, c2, S, qF} // qF es absorbente
}
```

---

## Descripción de Acción para el Trace

```go
func actionDesc(st state, c1, c2 string, res stepResult) string {
    switch {
    case res.next == qF:           return "fin de cinta → qf"
    case st == q0 && c1 == "M":   return "sospecha: escribe 1, mueve C2 →"
    case st == q0 && c1 == "E":   return "neutral: sin cambio en C2"
    case st == q0 && c1 == "H":   return "posible borrado: mueve C2 ←"
    case st == qBorra && c2 == "1": return "borra 1: suspicion--"
    case st == qBorra && c2 != "1": return "count=0: nada que borrar"
    }
    return ""
}
```

---

## Complejidad

- **Tiempo:** O(n) — un paso por símbolo + 1 paso en qf.
- **Espacio:** O(n) — la cinta 2 tiene como máximo n celdas (una por M).
- **Aislamiento:** `internal/turing` usa solo la stdlib de Go. Cero dependencias externas.
