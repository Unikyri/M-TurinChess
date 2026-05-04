# Spec — Frontend UI

**Change:** frontend-ui
**Directorio:** `frontend/`

---

## Requisitos Funcionales

### 1. Carga de Archivo (Upload)
- Zona de Drag & Drop para archivos `.pgn`.
- Input opcional para especificar el Elo del jugador (sobrescribe la detección automática del PGN).
- Selector del color del jugador a evaluar (Blancas / Negras).
- Botón "Analizar".
- Estado de carga (Loading) que desactive interacciones, ya que el análisis backend puede tardar varios segundos (Stockfish + LLM).

### 2. Vista de Resultados (Analysis Result)
Debe recibir y procesar el objeto `AnalysisResult` devuelto por el backend.
- **Veredicto:** Tarjeta destacada con el resultado (`MODULE_DETECTED` en rojo/alerta, `HUMAN_PLAYER` en verde/seguro).
- **Resumen:** Total de movimientos, umbral, cantidad de jugadas sospechosas (`suspicion_count`).
- **Máquina de Turing:**
  - Visualización de la "Cinta de Entrada" (Símbolos M, E, H).
  - Visualización de la "Cinta de Salida" (Conteo de 1s, B).
- **Detalle de Jugadas:** Una tabla o lista que muestra:
  - Número de movimiento.
  - Jugada SAN.
  - Mejor jugada sugerida (BestMove).
  - CP Loss.
  - Clasificación (M, E, H).
  - Justificación LLM (si aplica).

### 3. Historial (History)
- Consumir el endpoint `/api/history`.
- Listado de análisis previos mostrando ID, fecha, veredicto y color del jugador.

---

## Interacción con la API

- `POST http://localhost:8080/api/analyze`
  - Content-Type: `multipart/form-data`
  - Campos: `pgn_file` (File), `player_color` (String), `elo_white` (int, opcional), `elo_black` (int, opcional), `threshold` (int, opcional), `depth` (int, opcional).
  
- `GET http://localhost:8080/api/history`
  - Retorna un array de objetos JSON con el resumen histórico.

---

## Criterios de Aceptación
- [ ] La UI permite cargar correctamente un `.pgn` al backend sin errores de CORS.
- [ ] La aplicación responde correctamente a los estados de la petición (idle, loading, success, error).
- [ ] El veredicto final es el elemento más prominente de la interfaz tras el análisis.
- [ ] La cinta de la Máquina de Turing se visualiza correctamente.
- [ ] La experiencia visual es inmersiva y moderna (Dark mode, tipografías legibles).
