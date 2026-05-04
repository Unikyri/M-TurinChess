# Exploración — frontend-ui

**Fecha:** 2026-05-03

## Contexto
El backend de M-TurinChess está 100% operativo. Recibe archivos PGN, los analiza vía Stockfish, clasifica las jugadas (Lexer), evalúa el contexto con Gemini y simula una Máquina de Turing de 2 cintas para dar un veredicto final. Todo se expone mediante una API REST en `:8080` (con CORS habilitado para `http://localhost:5173`).

El siguiente gran paso es crear la interfaz gráfica para el usuario. Según el PRD, la interfaz no debe ser un MVP simple, sino una experiencia visual inmersiva ("premium", "estética atractiva", animaciones).

## Requerimientos Visuales y Funcionales (PRD & Guidelines)
1. **Subida de Archivos:** Drag & drop para archivos `.pgn`.
2. **Entrada de Elo:** Si el PGN no tiene Elo (o se quiere forzar), debe existir un modo de ingresarlo manualmente.
3. **Visualización de Resultados:**
   - Veredicto claro: `MODULE_DETECTED` (Alerta/Peligro) vs `HUMAN_PLAYER` (Éxito/Seguro).
   - "Cinta 1" (Secuencia Lógica) y "Cinta 2" (Conteo de Suspicion).
   - Detalle de jugadas: Mostrar SAN, CP Loss, clasificación (M/E/H) y la nota del LLM si existe.
   - Trace de la Máquina de Turing (Visualizador paso a paso).
4. **Historial:** Ver análisis previos.
5. **Estética "Premium":**
   - Paleta de colores curada (Dark Mode por defecto, glassmorphism).
   - Tipografía moderna (Inter, Roboto, o Outfit).
   - Micro-interacciones y animaciones sutiles.

## Stack Tecnológico Propuesto
- **Framework:** React + Vite (`create-vite`)
- **Estilos:** TailwindCSS (para velocidad, consistencia y sistema de diseño) + Framer Motion (para animaciones premium).
- **Iconos:** Lucide React.
- **Ruteo:** React Router (opcional, o manejo de vistas simple si es SPA pequeña).

## Desafíos
- El estado de análisis puede tomar tiempo (Stockfish profundo + Gemini). Se necesita feedback visual continuo (loaders o estados intermedios) o un buen loader de "Procesando...". En la API actual es sincrónico, así que la petición bloqueará. Un buen Spinner es vital.
- Presentar el Trace de la Máquina de Turing de forma que no abrume, pero demuestre el valor "académico" del proyecto.
