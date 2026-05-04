# Tasks — frontend-ui

**Change:** frontend-ui

---

## Batch 1 — Scaffolding e Infraestructura

- [x] **T1.1** — Inicializar proyecto React + Vite en la carpeta `/frontend` (`npx -y create-vite@latest frontend --template react-ts`).
- [x] **T1.2** — Instalar dependencias base (`tailwindcss`, `postcss`, `autoprefixer`, `framer-motion`, `lucide-react`, `clsx`, `tailwind-merge`).
- [x] **T1.3** — Configurar TailwindCSS (`tailwind.config.js` e `index.css`) con la paleta de colores Dark Theme.
- [x] **T1.4** — Crear `src/api/client.ts` con las interfaces de TypeScript (AnalysisResult, MoveDetail) y las funciones `analyzePGN` y `getHistory`.

## Batch 2 — UI: Formulario de Subida y Estados

- [x] **T2.1** — Crear componente `Uploader.tsx`: Input de archivo (PGN), input para Elo, selector de color, y botón de Enviar.
- [x] **T2.2** — Crear componente `Spinner.tsx` (estado de carga "Procesando con Stockfish & Gemini...").
- [x] **T2.3** — Implementar lógica en `App.tsx` para manejar el estado (`IDLE` -> `ANALYZING` -> `RESULT`).

## Batch 3 — UI: Visualización de Resultados y Máquina de Turing

- [x] **T3.1** — Crear componente `Verdict.tsx`: Tarjeta grande que muestra si es Humano o Módulo, con colores respectivos.
- [x] **T3.2** — Crear componente `TuringTape.tsx`: Visualización de las cintas 1 (M,E,H) y 2 (1s).
- [x] **T3.3** — Crear componente `MoveList.tsx`: Tabla que itere sobre `move_details`, mostrando SAN, clasificación y justificación LLM.
- [x] **T3.4** — Ensamblar en `Result.tsx`.

## Batch 4 — UI: Historial y Pulido

- [x] **T4.1** — Crear componente `History.tsx` que liste los registros pasados.
- [x] **T4.2** — Integrar `History.tsx` en el layout principal de `App.tsx`.
- [x] **T4.3** — Ajustes finales de estilo, animaciones de entrada con Framer Motion y revisión visual general.
