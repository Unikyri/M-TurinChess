import { useState } from 'react';
import { Uploader } from './components/Uploader';
import { Spinner } from './components/ui/Spinner';
import { Result } from './components/Result';
import { History } from './components/History';
import { analyzePGN, type AnalysisResult } from './api/client';

type AppState = "IDLE" | "ANALYZING" | "RESULT";

function App() {
	const [appState, setAppState] = useState<AppState>("IDLE");
	const [result, setResult] = useState<AnalysisResult | null>(null);
	const [error, setError] = useState<string | null>(null);
	// key to force History refresh when a new analysis is complete
	const [historyKey, setHistoryKey] = useState(0);

	const handleAnalyze = async (formData: FormData) => {
		setAppState("ANALYZING");
		setError(null);
		
		try {
			const res = await analyzePGN(formData);
			setResult(res);
			setAppState("RESULT");
			setHistoryKey(prev => prev + 1); // refresh history
		} catch (err: any) {
			setError(err.message || "Ocurrió un error inesperado al analizar.");
			setAppState("IDLE");
		}
	};

	const handleReset = () => {
		setAppState("IDLE");
		setResult(null);
	};

	return (
		<div className="min-h-screen bg-background relative overflow-x-hidden pb-12">
			{/* Decorative background blobs */}
			<div className="absolute top-[-10%] left-[-10%] w-[40%] h-[40%] bg-primary/20 rounded-full blur-[120px] pointer-events-none"></div>
			<div className="absolute bottom-[-10%] right-[-10%] w-[40%] h-[40%] bg-module/10 rounded-full blur-[120px] pointer-events-none"></div>

			{/* Header */}
			<header className="pt-12 pb-8 px-4 text-center relative z-10">
				<h1 className="text-4xl md:text-5xl font-extrabold text-transparent bg-clip-text bg-gradient-to-r from-primary to-emerald-400 mb-4 tracking-tight">
					M-TurinChess
				</h1>
				<p className="text-slate-400 max-w-xl mx-auto text-lg">
					Sistema de detección algorítmica de asistencia computacional usando Máquinas de Turing y Modelos de Lenguaje.
				</p>
			</header>

			{/* Main Content Area */}
			<main className="container mx-auto px-4 relative z-10">
				{error && (
					<div className="max-w-2xl mx-auto mb-6 p-4 glass-panel border-module/50 bg-module/10 text-rose-200 text-center rounded-lg">
						{error}
					</div>
				)}

				{appState === "IDLE" && (
					<div className="animate-in fade-in slide-in-from-bottom-4 duration-700">
						<Uploader onAnalyze={handleAnalyze} />
						<History key={historyKey} />
					</div>
				)}

				{appState === "ANALYZING" && (
					<div className="animate-in fade-in zoom-in-95 duration-500 mt-12">
						<Spinner />
					</div>
				)}

				{appState === "RESULT" && result && (
					<Result result={result} onReset={handleReset} />
				)}
			</main>
		</div>
	);
}

export default App;
