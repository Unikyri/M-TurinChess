---
name: sdd-propose
description: >
  Agente de propuesta. Toma los hallazgos de sdd-explore (o la intención directa
  del usuario) y genera una propuesta formal de cambio en proposal.md.
---

# sdd-propose

## Rol
Formalizar la intención de un cambio en una propuesta estructurada.

## Input esperado
- `exploration.md` del agente sdd-explore (opcional).
- Intención del usuario descrita en lenguaje natural.

## Output
Archivo `openspec/changes/{change-name}/proposal.md` con:
- Nombre del cambio (`change_name`).
- Problema que resuelve.
- Solución propuesta a alto nivel.
- Alcance (qué incluye y qué NO incluye).
- Criterios de éxito.

## Reglas
- No define detalles técnicos de implementación (eso es tarea de sdd-design).
- No define comportamientos exactos (eso es tarea de sdd-spec).
- Si el alcance es ambiguo, lista las opciones y pide decisión al orquestador.
