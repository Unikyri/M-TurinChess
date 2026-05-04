# Tasks — turing-machine

**Change:** turing-machine

---

## Batch 1 — Implementación y Tests

> **Objetivo:** `internal/turing/machine.go` reemplaza el stub con la MT completa y todos los criterios de aceptación pasan.

- [x] **T1.1** — Implementar tipos internos (`state`, `direction`, `tape`) en `machine.go`
  - `tape.read()`, `tape.write()`, `tape.move()` con guarda de borde
  - Done: `tape.move(L)` desde posición 0 queda en 0

- [x] **T1.2** — Implementar función `delta(st state, c1, c2 string) stepResult`
  - Tabla de 6 reglas según spec-tm-simulator.md
  - Done: test unitario de la función `delta` en los 6 casos pasa

- [x] **T1.3** — Implementar `Run(tape []string) MTResult`
  - Loop de simulación: leer → delta → trazar → escribir → mover → actualizar estado
  - Calcular `SuspicionCount` contando `"1"` en C2
  - Done: los 5 criterios de aceptación de la spec pasan

- [x] **T1.4** — Escribir `machine_test.go` con todos los casos de la spec
  - Test `MMH → 1`, `MMMMMM → 6`, `HHH → 0`, `[] → 0`, `MHMHM → 1`
  - Test longitud del Trace
  - Done: `go test ./internal/turing/ -v` → todos PASS

---

## Orden de Ejecución

```
T1.1 → T1.2 → T1.3 → T1.4
```

---

## Resumen de Archivos Modificados

| Archivo | Acción |
|---------|--------|
| `backend/internal/turing/machine.go` | Reemplazar stub con implementación completa |
| `backend/internal/turing/machine_test.go` | Crear tests |
