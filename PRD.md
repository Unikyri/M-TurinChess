# Documento de Requisitos del Producto (PRD)

Detección de Módulos en Ajedrez mediante Máquina de Turing (Aproximación de Kolmogorov)

**Fecha:** Mayo 2026
**Tecnologías Principales:** Go (Backend), Stockfish (Motor UCI), React (Frontend), Simulador MT personalizado. Base de datos JSON.

## 1. Visión General del Proyecto

Este proyecto es una herramienta de análisis de partidas de ajedrez (archivos PGN) que utiliza la Teoría de la Computación para detectar si un jugador humano recibió asistencia computacional (uso de módulo). En lugar de usar heurísticas tradicionales codificadas mediante condicionales básicos, el núcleo lógico recae enteramente en una Máquina de Turing Multicinta que modela una aproximación a la Complejidad de Kolmogorov, midiendo la entropía y compresibilidad algorítmica de la partida.

## 2. Objetivos

- **Académico:** Demostrar el uso práctico y la superioridad matemática de una Máquina de Turing (que posee memoria bidireccional) frente a Autómatas Finitos para llevar un seguimiento de estados fluctuantes (suma y resta de sospecha).
- **Técnico:** Desarrollar un backend robusto en Go que actúe como "traductor léxico", parseando PGNs, consultando a Stockfish vía UCI y formateando los datos para la Máquina de Turing.
- **Usuario:** Proveer una interfaz local en React donde un usuario pueda subir un PGN (solicitando el Elo de los jugadores si no se encuentra en el archivo) y recibir un veredicto visual y justificado. También se requiere persistencia de los análisis en una base de datos basada en archivos JSON.

## 3. Arquitectura del Sistema

El sistema se divide en tres capas principales (Analogía del Tribunal):

### Capa de Extracción (Los Testigos)
- **Stockfish:** Evalúa cada jugada humana comparándola con la jugada óptima computacional para extraer la "Pérdida de Centipeones" (CP Loss).
- **LLM (Gemini):** Actúa como perito psicológico. Toma en cuenta el Elo del jugador (ej. 1200 vs 2800) para juzgar si una jugada óptima (0 CP Loss) es humanamente incomprensible o justificada.

### Capa Léxica y Controlador (El Secretario - Go)
- Escrito en Go. Toma los valores de CP Loss y los traduce estrictamente a un alfabeto formal $\Sigma = \{M, E, H\}$.
- Inicializa y alimenta a la Máquina de Turing.
- Gestiona la persistencia de los análisis en un archivo JSON.

### Capa Lógica / Motor de Decisión (El Juez - Máquina de Turing)
- Un simulador codificado en Go que ejecuta estrictamente una 7-tupla matemática.
- No tiene conexión a internet ni a Stockfish. Solo procesa el alfabeto, suma y borra "unos" en su cinta de trabajo para emitir el veredicto final.

## 4. Definición Formal de la Máquina de Turing (El Eje Principal)

El sistema utiliza una MT determinista de 2 cintas.
- **Cinta 1 (Entrada):** Cadena generada por Go. Solo lectura. Movimiento siempre a la derecha (R) o pausa (S).
- **Cinta 2 (Memoria/Trabajo):** Cinta en blanco inicial. Utiliza codificación unaria (1) para medir la sospecha. Movimiento bidireccional (L, R, S) para agregar o borrar sospecha según el contexto.

### 4.1. La 7-tupla $M = (Q, \Sigma, \Gamma, \delta, q_0, B, F)$
- $Q$ (Estados): $\{q_0, q_{borra}, q_f\}$
- $\Sigma$ (Alfabeto de Entrada): $\{M, E, H\}$
  - $M$: Jugada de Módulo (0-10 CP Loss).
  - $E$: Jugada Estándar / Neutra (11-50 CP Loss o Jugada de Libro).
  - $H$: Jugada Humana / Mediocre (>50 CP Loss).
- $\Gamma$ (Alfabeto de Cinta): $\{M, E, H, 1, B\}$
- $q_0$: Estado Inicial.
- $B$: Símbolo Blanco.
- $F$: $\{q_f\}$ (Estado Final).

### 4.2. Tabla de Transiciones ($\delta$)

Formato: $\delta(Estado, Lee_{C1}, Lee_{C2}) = (NuevoEstado, Escribe_{C1}, Mov_{C1}, Escribe_{C2}, Mov_{C2})$

| Estado Actual | Lee Cinta 1 | Lee Cinta 2 | Nuevo Estado | Escribe C1 | Mov C1 | Escribe C2 | Mov C2 | Justificación Lógica |
|--------------|-------------|-------------|--------------|------------|--------|------------|--------|----------------------|
| $q_0$ | $M$ | $B$ | $q_0$ | $M$ | $R$ | $1$ | $R$ | Módulo detectado. Aumenta sospecha. |
| $q_0$ | $E$ | $B$ | $q_0$ | $E$ | $R$ | $B$ | $S$ | Jugada normal. Sin cambios. |
| $q_0$ | $H$ | $B$ | $q_{borra}$ | $H$ | $S$ | $B$ | $L$ | Error humano. Pausa C1, retrocede C2. |
| $q_{borra}$ | $H$ | $1$ | $q_0$ | $H$ | $R$ | $B$ | $S$ | Borra un '1' de sospecha y reanuda. |
| $q_{borra}$ | $H$ | $B$ | $q_0$ | $H$ | $R$ | $B$ | $S$ | Nada que borrar (sospecha en 0). Reanuda. |
| $q_0$ | $B$ | $B$ | $q_f$ | $B$ | $S$ | $B$ | $S$ | Fin de la partida. Acepta cadena. |

## 5. Criterios de Traducción en Go (El Lexer)

El backend en Go debe mapear la Pérdida de Centipeones (CP Loss) calculada por Stockfish al alfabeto $\Sigma$:
- **Filtro Inicial:** Omitir las primeras 10 jugadas (20 plies) para evitar falsos positivos de la teoría de aperturas. Teniendo en cuenta el Elo.
- **Clasificación:**
  - CP Loss <= 10: Traducir a $M$ (Módulo).
  - 11 <= CP Loss <= 50: Traducir a $E$ (Estándar).
  - CP Loss > 50: Traducir a $H$ (Humano/Mediocre).
- **Generación de Cadena:** Si el jugador hizo 5 jugadas, Go genera un string: `[M, E, M, M, H]`. Esta es la entrada de la Cinta 1.

## 6. Lógica del Veredicto (Post-Máquina)

Una vez que la Máquina de Turing llega a $q_f$:
- Go cuenta el número de 1s restantes en la Cinta 2.
- Umbral de Trampa: Se define empíricamente (ej. Umbral = 6) y es influenciado por el Elo.
- Decisión:
  - Si Conteo(1) >= Umbral: Veredicto = ASISTENCIA DE MÓDULO DETECTADA.
  - Si Conteo(1) < Umbral: Veredicto = JUGADOR HUMANO.

## 7. Requisitos Técnicos y Cronograma de Desarrollo

**Semana 1: Backend Base y Análisis de Datos (Go + Stockfish)**
- [ ] Levantar un servidor en Go.
- [ ] Integrar un parser PGN o desarrollarlo.
- [ ] Crear un subproceso en Go para ejecutar el binario de Stockfish mediante el protocolo UCI.
- [ ] Crear la función que evalúa una posición, obtiene el CP Loss y devuelve la letra ($M, E$ o $H$).

**Semana 2: El Eje Principal (Simulador de Máquina de Turing en Go)**
- [ ] Crucial: Crear una estructura (struct en Go) independiente que simule una Máquina de Turing. Debe aceptar estados, cintas (arrays de strings) y una tabla de transiciones.
- [ ] Implementar la función de ciclo Run() que lea los cabezales, aplique la regla y actualice las cintas iterativamente hasta llegar a $q_f$.
- [ ] Conectar la salida del parser de ajedrez (Semana 1) con la entrada de la Cinta 1 de la MT.

**Semana 3: Frontend (React) y Conexión Final**
- [ ] Crear la estructura de directorios separando `backend` y `frontend` (React).
- [ ] Crear una interfaz en React con un tablero interactivo (librería react-chessboard o similar).
- [ ] Crear un área de "Carga de PGN" que detecte el Elo o solicite al usuario que lo ingrese manualmente.
- [ ] Integrar la persistencia de datos (guardar los análisis en un archivo JSON en el backend).
- [ ] Conectar el Frontend con la API de Go.
- [ ] Mostrar los resultados: La interfaz debe mostrar cuántos '1' quedaron en la Cinta 2 y el veredicto final.