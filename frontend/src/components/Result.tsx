import { type AnalysisResult } from "../api/client";
import { Verdict } from "./Verdict";
import { TuringTape } from "./TuringTape";
import { MoveList } from "./MoveList";
import { RotateCcw } from "lucide-react";

interface ResultProps {
	result: AnalysisResult;
	onReset: () => void;
}

export function Result({ result, onReset }: ResultProps) {
	return (
		<div className="animate-in fade-in slide-in-from-bottom-8 duration-700 max-w-5xl mx-auto space-y-8 w-full">
			{/* Header Action */}
			<div className="flex justify-between items-end">
				<div>
					<h2 className="text-2xl font-bold text-slate-100">Reporte de Análisis</h2>
					<p className="text-slate-400">ID: <span className="font-mono text-xs">{result.id.split('-')[0]}</span> • Jugador: <span className="capitalize text-slate-200">{result.player_color}</span> (Elo {result.elo})</p>
				</div>
				<button 
					onClick={onReset}
					className="flex items-center space-x-2 bg-slate-800 hover:bg-slate-700 text-slate-200 py-2 px-4 rounded-lg transition-colors border border-slate-700"
				>
					<RotateCcw className="w-4 h-4" />
					<span>Nuevo Análisis</span>
				</button>
			</div>

			{/* Verdict Card */}
			<Verdict 
				verdict={result.verdict} 
				suspicionCount={result.suspicion_count} 
				threshold={result.threshold}
				totalMoves={result.total_moves_analyzed}
			/>

			{/* Turing Machine Visualization */}
			<TuringTape 
				inputTape={result.tape_input} 
				outputTape={result.tape_output} 
			/>

			{/* Lexical and LLM Detail Table */}
			<MoveList moves={result.move_details} />
		</div>
	);
}
