# Diseño Técnico — frontend-ui

**Change:** frontend-ui

---

## Stack
- **Herramienta de Build:** Vite
- **Framework:** React 18 (TypeScript)
- **Estilos:** TailwindCSS
- **Animaciones:** Framer Motion (para transiciones suaves de componentes y alertas).
- **Iconos:** Lucide React (minimalistas y limpios).

## Estructura de Directorios
El frontend vivirá en la carpeta `/frontend` a nivel raíz del proyecto (hermano de `/backend`).

```
frontend/
├── index.html
├── package.json
├── tailwind.config.js
├── postcss.config.js
├── src/
│   ├── main.tsx
│   ├── App.tsx          # Layout y orquestador de estado
│   ├── index.css        # Configuración de Tailwind y vars CSS
│   ├── api/
│   │   └── client.ts    # Fetch API wrappers (analyze, history)
│   ├── components/
│   │   ├── Uploader.tsx # Drag & drop y form
│   │   ├── Result.tsx   # Dashboard de resultados (Verdict, Tapes)
│   │   ├── History.tsx  # Sidebar o sección de historial
│   │   └── ui/          # Componentes genéricos (Spinner, Card, Badge)
│   └── types/
│       └── index.ts     # Tipos TypeScript para AnalysisResult
```

## Paleta de Colores (Tema Dark "Premium")
- **Background:** Slate oscuro (`#0f172a` o similar).
- **Cards/Panels:** Slate con opacidad y borde (Glassmorphism sutil).
- **Acentos:**
  - Éxito (Human): Verde esmeralda o teal (`#10b981`).
  - Peligro (Módulo): Rojo rosado o crimson (`#f43f5e`).
  - Info: Azul índigo (`#6366f1`).
- **Texto:** Gris muy claro (`#f8fafc`).

## Flujo de Estado (App.tsx)
La aplicación será una SPA muy simple, manejando su estado central en `App.tsx`:

```typescript
type AppState = "IDLE" | "ANALYZING" | "RESULT";

// State hooks
const [appState, setAppState] = useState<AppState>("IDLE");
const [result, setResult] = useState<AnalysisResult | null>(null);
const [history, setHistory] = useState<HistoryRecord[]>([]);

// Handlers
const handleAnalyze = async (formData: FormData) => {
    setAppState("ANALYZING");
    try {
        const res = await apiAnalyze(formData);
        setResult(res);
        setAppState("RESULT");
        loadHistory(); // Refrescar historial
    } catch (e) {
        // Manejar error
        setAppState("IDLE");
    }
}
```

## Visualización de la Máquina de Turing
- Las cintas se renderizarán como secuencias horizontales de pequeñas "celdas" (cajas cuadradas).
- La Cinta 1 (Símbolos M, E, H) tendrá colores diferentes para cada símbolo (ej. M=Rojo, E=Gris, H=Verde).
- La Cinta 2 (Conteo de 1s) representará la sospecha acumulada.
