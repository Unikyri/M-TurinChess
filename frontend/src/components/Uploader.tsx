import { useState } from "react";
import { FileText, Settings2, User } from "lucide-react";

interface UploaderProps {
	onAnalyze: (data: FormData) => void;
}

export function Uploader({ onAnalyze }: UploaderProps) {
	const [pgnText, setPgnText] = useState<string>("");
	const [color, setColor] = useState<"white" | "black">("white");
	const [elo, setElo] = useState<string>("");
	const [threshold, setThreshold] = useState<string>("6");

	// Auto-calculate a sensible threshold based on ELO.
	// Lower ELO → lower threshold: even a few perfect moves from a weak player is suspicious.
	// Higher ELO → higher threshold: strong players naturally play many near-perfect moves.
	const autoThreshold = (eloVal: number): number => {
		if (eloVal < 800)  return 3;
		if (eloVal < 1000) return 4;
		if (eloVal < 1200) return 5;
		if (eloVal < 1500) return 6;
		if (eloVal < 1800) return 8;
		if (eloVal < 2000) return 10;
		if (eloVal < 2200) return 13;
		return 16;
	};

	const handleEloChange = (val: string) => {
		setElo(val);
		const n = parseInt(val, 10);
		if (!isNaN(n) && n > 0) {
			setThreshold(String(autoThreshold(n)));
		}
	};

	const handleSubmit = (e: React.FormEvent) => {
		e.preventDefault();
		if (!pgnText.trim()) return;

		const formData = new FormData();
		formData.append("pgn_text", pgnText);
		formData.append("player_color", color);
		
		if (elo) {
			formData.append(color === "white" ? "elo_white" : "elo_black", elo);
		}
		if (threshold) {
			formData.append("threshold", threshold);
		}

		onAnalyze(formData);
	};


	return (
		<div className="w-full max-w-2xl mx-auto glass-panel p-8">
			<div className="text-center mb-8">
				<h2 className="text-2xl font-bold text-slate-100 mb-2">Ingresa tu Partida</h2>
				<p className="text-slate-400">Pega el texto PGN de tu partida para detectar asistencia computacional.</p>
			</div>

			<form onSubmit={handleSubmit} className="space-y-6">
				{/* Textarea Area */}
				<div className="space-y-2">
					<label className="flex items-center space-x-2 text-sm font-medium text-slate-300">
						<FileText className="w-4 h-4" />
						<span>Texto PGN</span>
					</label>
					<textarea 
						placeholder="[Event &quot;FIDE World Cup 2023&quot;]&#10;[Site &quot;Baku AZE&quot;]&#10;...&#10;1. e4 e5 2. Nf3 Nc6..."
						className="w-full h-48 bg-slate-800 border border-slate-700 rounded-xl p-4 text-slate-200 focus:ring-2 focus:ring-primary focus:border-transparent outline-none transition-all font-mono text-sm resize-y custom-scrollbar"
						value={pgnText}
						onChange={(e) => setPgnText(e.target.value)}
						required
					/>
				</div>

				{/* Configuration Options */}
				<div className="grid grid-cols-1 md:grid-cols-3 gap-6">
					<div className="space-y-2">
						<label className="flex items-center space-x-2 text-sm font-medium text-slate-300">
							<User className="w-4 h-4" />
							<span>Color a Evaluar</span>
						</label>
						<select 
							className="w-full bg-slate-800 border border-slate-700 rounded-lg px-3 py-2 text-slate-200 focus:ring-2 focus:ring-primary focus:border-transparent outline-none transition-all"
							value={color}
							onChange={(e) => setColor(e.target.value as "white" | "black")}
						>
							<option value="white">Blancas</option>
							<option value="black">Negras</option>
						</select>
					</div>

					<div className="space-y-2">
						<label className="flex items-center space-x-2 text-sm font-medium text-slate-300">
							<Settings2 className="w-4 h-4" />
							<span>Elo (Opcional)</span>
						</label>
						<input 
							type="number" 
							placeholder="Ej: 1500"
							className="w-full bg-slate-800 border border-slate-700 rounded-lg px-3 py-2 text-slate-200 focus:ring-2 focus:ring-primary focus:border-transparent outline-none transition-all"
							value={elo}
							onChange={(e) => handleEloChange(e.target.value)}
						/>
					</div>

					<div className="space-y-2">
						<label className="flex items-center space-x-2 text-sm font-medium text-slate-300">
							<Settings2 className="w-4 h-4" />
							<span>Umbral MT</span>
						</label>
						<input 
							type="number" 
							placeholder="Def: 6"
							className="w-full bg-slate-800 border border-slate-700 rounded-lg px-3 py-2 text-slate-200 focus:ring-2 focus:ring-primary focus:border-transparent outline-none transition-all"
							value={threshold}
							onChange={(e) => setThreshold(e.target.value)}
						/>
					</div>
				</div>

				<button 
					type="submit" 
					disabled={!pgnText.trim()}
					className="w-full bg-primary hover:bg-indigo-600 text-white font-semibold py-3 rounded-lg transition-colors disabled:opacity-50 disabled:cursor-not-allowed mt-4"
				>
					Iniciar Análisis M-Turin
				</button>
			</form>
		</div>
	);
}
