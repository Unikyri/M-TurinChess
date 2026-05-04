import { AlertTriangle, CheckCircle, ShieldAlert, ShieldCheck } from "lucide-react";

interface VerdictProps {
	verdict: string;
	suspicionCount: number;
	threshold: number;
	totalMoves: number;
}

export function Verdict({ verdict, suspicionCount, threshold, totalMoves }: VerdictProps) {
	const isModule = verdict === "MODULE_DETECTED";

	return (
		<div className={`p-8 rounded-2xl border-2 transition-all duration-500 relative overflow-hidden
			${isModule 
				? "border-module/50 bg-module/10 shadow-[0_0_30px_rgba(244,63,94,0.2)]" 
				: "border-human/50 bg-human/10 shadow-[0_0_30px_rgba(16,185,129,0.2)]"
			}`}
		>
			{/* Watermark icon */}
			<div className="absolute -right-8 -top-8 opacity-5">
				{isModule ? <ShieldAlert className="w-64 h-64 text-module" /> : <ShieldCheck className="w-64 h-64 text-human" />}
			</div>

			<div className="relative z-10 flex flex-col md:flex-row items-center justify-between gap-6">
				<div className="flex items-center space-x-6">
					<div className={`p-4 rounded-full ${isModule ? "bg-module/20 text-module" : "bg-human/20 text-human"}`}>
						{isModule ? <AlertTriangle className="w-10 h-10" /> : <CheckCircle className="w-10 h-10" />}
					</div>
					<div>
						<h2 className="text-sm font-bold tracking-widest text-slate-400 uppercase mb-1">Veredicto Oficial</h2>
						<h1 className={`text-4xl font-extrabold ${isModule ? "text-module" : "text-human"}`}>
							{isModule ? "Módulo Detectado" : "Jugador Humano"}
						</h1>
					</div>
				</div>

				<div className="flex space-x-4">
					<div className="glass-panel px-6 py-4 text-center">
						<p className="text-sm text-slate-400 mb-1">Sospecha MT</p>
						<p className="text-2xl font-bold text-slate-100">
							{suspicionCount} <span className="text-lg text-slate-500 font-normal">/ {threshold}</span>
						</p>
					</div>
					<div className="glass-panel px-6 py-4 text-center">
						<p className="text-sm text-slate-400 mb-1">Jugadas Evaluadas</p>
						<p className="text-2xl font-bold text-slate-100">{totalMoves}</p>
					</div>
				</div>
			</div>
		</div>
	);
}
