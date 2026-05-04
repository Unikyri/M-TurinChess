# Propuesta — turing-machine

**Fecha:** 2026-05-03
**Agente:** sdd-propose
**Change name:** `turing-machine`

---

## Problema

El stub actual en `internal/turing/machine.go` siempre retorna `SuspicionCount: 0`, haciendo que el veredicto sea siempre `HUMAN_PLAYER`. La MT real es el corazón del proyecto.

## Solución

Reemplazar el stub por la implementación completa de la MT de 2 cintas con la tabla de transiciones δ definida en `docs/02-turing-machine.md` y `exploration.md`.

## Alcance

### ✅ Incluido

| Componente | Archivo |
|------------|---------|
| Tipos de estado, símbolo, dirección | `internal/turing/machine.go` |
| Tabla de transiciones δ | `internal/turing/transitions.go` |
| Simulador `Run(tape []string) MTResult` | `internal/turing/machine.go` |
| Tests unitarios con trazas verificadas | `internal/turing/machine_test.go` |

### ❌ Excluido

- Modificación de ningún otro paquete (el contrato de `Run` no cambia).
- Visualización del trace en el frontend (Fase 3).

## Criterios de Éxito

1. `Run([]string{"M","M","H"})` → `SuspicionCount = 1`.
2. `Run([]string{"M","M","M","M","M","M"})` → `SuspicionCount = 6`.
3. `Run([]string{"H","H","H"})` → `SuspicionCount = 0` (no va negativo).
4. `Run([]string{})` → `SuspicionCount = 0`.
5. El `Trace` refleja cada paso con `State`, `ReadC1`, `ReadC2`, `Action`, `Suspicion`.
6. El package `turing` NO importa ningún otro paquete del proyecto.
