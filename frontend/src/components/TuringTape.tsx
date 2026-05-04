import { Database, Binary } from "lucide-react";

interface TuringTapeProps {
	inputTape: string[];
	outputTape: string[];
}

export function TuringTape({ inputTape, outputTape }: TuringTapeProps) {
	const getCellColor = (sym: string) => {
		switch (sym) {
			case "M": return "bg-module text-white border-module";
			case "E": return "bg-slate-600 text-slate-200 border-slate-500";
			case "H": return "bg-human text-white border-human";
			case "1": return "bg-primary text-white border-primary shadow-[0_0_10px_rgba(99,102,241,0.5)]";
			default: return "bg-slate-800 text-slate-500 border-slate-700";
		}
	};

	return (
		<div className="space-y-6">
			{/* Cinta 1: Entrada */}
			<div className="glass-panel p-6">
				<div className="flex items-center space-x-3 mb-4">
					<Database className="w-5 h-5 text-slate-400" />
					<h3 className="font-semibold text-slate-200">Cinta 1: Entrada Lexicalizada</h3>
				</div>
				<div className="flex space-x-2 overflow-x-auto pb-4 custom-scrollbar">
					{inputTape.map((sym, idx) => (
						<div 
							key={idx} 
							className={`flex-shrink-0 w-12 h-12 flex items-center justify-center rounded-md border-2 font-bold text-lg transition-transform hover:scale-110 ${getCellColor(sym)}`}
						>
							{sym}
						</div>
					))}
					{/* Infinite tape illusion */}
					<div className="flex-shrink-0 w-12 h-12 flex items-center justify-center rounded-md border-2 border-slate-800/50 bg-slate-800/20 text-slate-700">...</div>
				</div>
				<div className="flex justify-between mt-2 text-xs text-slate-500 font-mono">
					<span>Index: 0</span>
					<span>Length: {inputTape.length}</span>
				</div>
			</div>

			{/* Cinta 2: Memoria / Pila */}
			<div className="glass-panel p-6">
				<div className="flex items-center space-x-3 mb-4">
					<Binary className="w-5 h-5 text-slate-400" />
					<h3 className="font-semibold text-slate-200">Cinta 2: Memoria de Sospecha (MT Stack)</h3>
				</div>
				<div className="flex space-x-2 overflow-x-auto pb-4 custom-scrollbar">
					{outputTape.length === 0 ? (
						<div className="text-slate-500 italic py-2">La cinta está vacía (B)</div>
					) : (
						outputTape.map((sym, idx) => (
							<div 
								key={idx} 
								className={`flex-shrink-0 w-12 h-12 flex items-center justify-center rounded-md border-2 font-bold text-lg animate-pulse ${getCellColor(sym)}`}
							>
								{sym}
							</div>
						))
					)}
					<div className="flex-shrink-0 w-12 h-12 flex items-center justify-center rounded-md border-2 border-slate-800/50 bg-slate-800/20 text-slate-700">...</div>
				</div>
			</div>
		</div>
	);
}
