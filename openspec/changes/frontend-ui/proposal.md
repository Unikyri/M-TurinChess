# Propuesta — frontend-ui

**Fecha:** 2026-05-03
**Change name:** `frontend-ui`

---

## Problema
Actualmente la única forma de interactuar con el motor M-TurinChess es a través de peticiones cURL o Postman al backend. Falta la cara visible del proyecto ("El Tribunal").

## Solución
Construir una Single Page Application (SPA) en React que actúe como dashboard central. La aplicación permitirá la carga de archivos, mostrará estados de carga elegantes y visualizará los resultados complejos (Turing Trace, Evaluaciones, Veredicto) en una interfaz moderna y "premium".

## Alcance

### ✅ Incluido
- Scaffolding de Vite + React (TypeScript).
- Configuración de TailwindCSS y variables CSS de la paleta.
- Diseño e implementación de la vista principal (Drag & Drop + Formulario).
- Diseño e implementación de la vista de Resultados (Veredicto, Cinta, Detalles).
- Diseño e implementación del Historial de análisis.
- Integración con el backend (`fetch` a `localhost:8080/api/analyze` y `/api/history`).

### ❌ Excluido
- Un tablero de ajedrez interactivo pieza por pieza (por ahora, nos enfocamos en el análisis matemático y el texto SAN). Si el tiempo lo permite, se podría agregar `react-chessboard` después, pero el scope inicial es mostrar los datos de la MT y Lexer.
- Websockets (la API actual es síncrona).

## Criterios de Éxito
1. La aplicación inicia con `npm run dev` en el puerto `5173`.
2. El usuario puede subir un PGN y recibir el análisis sin errores de CORS.
3. El veredicto final se destaca visualmente con una animación atractiva.
4. Se puede consultar el historial de partidas analizadas.
5. El diseño se ve moderno (Dark Mode, Tailwind, buena tipografía).
