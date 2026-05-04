import { Loader2 } from "lucide-react";

interface SpinnerProps {
	message?: string;
}

export function Spinner({ message = "Procesando con Stockfish & LLM..." }: SpinnerProps) {
	return (
		<div className="flex flex-col items-center justify-center p-12 space-y-6">
			<div className="relative">
				{/* Outer glowing ring */}
				<div className="absolute inset-0 rounded-full blur-xl bg-primary/30 animate-pulse"></div>
				{/* Spinning icon */}
				<Loader2 className="w-16 h-16 text-primary animate-spin relative z-10" />
			</div>
			
			<div className="text-center space-y-2">
				<h3 className="text-xl font-semibold text-slate-100">Analizando Partida</h3>
				<p className="text-sm text-slate-400 max-w-sm">
					{message}
				</p>
				<p className="text-xs text-slate-500 mt-4 animate-pulse">
					Este proceso puede tomar varios segundos dependiendo de la profundidad.
				</p>
			</div>
		</div>
	);
}
