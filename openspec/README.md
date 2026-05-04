# OpenSpec — Artifact Store

Este directorio contiene todos los artefactos del flujo **Spec-Driven Development (SDD)** para el proyecto M-TurinChess.

## Estructura

```
openspec/
├── changes/          # Cambios activos
│   ├── archive/      # Cambios completados y archivados
│   └── <change-name>/
│       ├── proposal.md
│       ├── specs/
│       ├── design.md
│       ├── tasks.md
│       └── verify-report.md
├── state.yaml        # Estado global del SDD
└── README.md         # Este archivo
```

## Flujo

```
explore → propose → spec + design → tasks → apply → verify → archive
```

Cada sub-agente lee de su capa de dependencia y escribe su output en `changes/{change-name}/`.
