# Roadmap de Desarrollo

> Cronograma adaptado del PRD con detalle de tareas.

---

## Fase 1 — Backend Base y Análisis de Datos (Semana 1)

> **Objetivo:** Go puede recibir un PGN, comunicarse con Stockfish, y producir CP Loss por jugada.

### Tareas

- [ ] **1.1** Inicializar módulo Go en `backend/` (`go mod init github.com/Unikyri/M-TurinChess/backend`)
- [ ] **1.2** Crear estructura del proyecto (`cmd/`, `internal/`, `data/`)
- [ ] **1.3** Implementar Parser PGN custom (`internal/chess/pgn.go`)
  - Leer archivo .pgn
  - Extraer metadata (jugadores, resultado, evento, **Elo**)
  - Extraer lista de jugadas en notación SAN
  - Generar posiciones FEN a partir de las jugadas
- [ ] **1.4** Implementar cliente Stockfish UCI (`internal/chess/stockfish.go`)
  - Iniciar subproceso con el binario local
  - Handshake UCI (`uci` → `uciok`)
  - Enviar posiciones y recibir evaluaciones (`position fen ...` / `go depth ...`)
  - Calcular CP Loss por jugada
- [ ] **1.5** Implementar Lexer (`internal/chess/lexer.go`)
  - Adaptar filtro inicial dependiendo del Elo de los jugadores
  - Clasificar CP Loss → `{M, E, H}`
  - Generar cadena formal para la Cinta 1
- [ ] **1.6** Implementar persistencia JSON (`internal/db/json_storage.go`)
  - Crear, leer y agregar registros a `data/history.json`
- [ ] **1.7** Levantar servidor HTTP básico (`internal/api/`)
  - Endpoint `POST /api/analyze`
- [ ] **1.8** Tests unitarios para Parser PGN y Lexer

---

## Fase 2 — Simulador de Máquina de Turing (Semana 2)

> **Objetivo:** La MT procesa cadenas del alfabeto y emite el conteo correcto de sospecha.

### Tareas

- [ ] **2.1** Implementar struct `Machine` (`internal/turing/machine.go`)
  - Cintas como slices de strings
  - Cabezales como índices
  - Estado actual
- [ ] **2.2** Implementar tabla de transiciones (`internal/turing/transitions.go`)
  - Las 6 reglas del PRD actualizadas (alfabeto `{M, E, H}`)
- [ ] **2.3** Implementar método `Run()`
  - Ciclo: leer cabezales → buscar regla → aplicar → repetir
  - Condición de parada: estado = $q_f$
- [ ] **2.4** Implementar traza de ejecución
  - Registrar cada paso para visualización en frontend
- [ ] **2.5** Conectar pipeline (`internal/analysis/analyzer.go`)
- [ ] **2.6** Implementar lógica de veredicto (conteo vs umbral dinámico basado en Elo/ritmo)
- [ ] **2.7** Integrar Gemini LLM (`internal/llm/gemini.go`)
  - LLM evalúa si una jugada de módulo (M) es comprensible para el Elo del jugador
- [ ] **2.8** Tests unitarios para la MT

---

## Fase 3 — Frontend en React y Conexión Final (Semana 3)

> **Objetivo:** Interfaz funcional SPA donde el usuario interactúa, ve resultados y persistencia.

### Tareas

- [ ] **3.1** Inicializar proyecto Vite/React en `frontend/`
- [ ] **3.2** Implementar área de carga de PGN
  - Drag & drop + botón de selección
  - Si el PGN no contiene Elo, solicitar input manual
- [ ] **3.3** Integrar tablero interactivo (react-chessboard)
  - Navegación por jugadas sincronizada con análisis
- [ ] **3.4** Panel de configuración
  - Selector de color del jugador analizado
  - Slider de umbral de detección
  - Selector de profundidad Stockfish
- [ ] **3.5** Visualización de resultados
  - Veredicto final con animación
  - Traza de la Máquina de Turing (paso a paso)
  - Estado visual de la Cinta 2 (como una pila)
- [ ] **3.6** Conectar React con API Go (`axios` o `fetch`)
- [ ] **3.7** Vista de Historial (Base de datos JSON)
  - Cargar los análisis previos

---

## Dependencias Externas

| Dependencia | Tipo | Propósito |
|-------------|------|-----------|
| Stockfish 17 | Binario local | Motor de evaluación UCI |
| Gemini API | API REST | Análisis de jugadas inhumanas por Elo |
| React | Frontend | SPA Framework |
| react-chessboard | JS Component | Tablero visual en frontend |
