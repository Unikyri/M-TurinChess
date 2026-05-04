import { type MoveDetail } from "../api/client";
import { MessageSquare, Bot, ArrowRight, User } from "lucide-react";

interface MoveListProps {
	moves: MoveDetail[];
}

export function MoveList({ moves }: MoveListProps) {
	const getBadgeStyles = (classification: string) => {
		switch (classification) {
			case "M": return "bg-module/20 text-module border-module/30";
			case "E": return "bg-slate-700 text-slate-300 border-slate-600";
			case "H": return "bg-human/20 text-human border-human/30";
			default: return "bg-slate-800 text-slate-400 border-slate-700";
		}
	};

	const getLlmIcon = (flag: string) => {
		if (flag.includes("natural")) return <User className="w-4 h-4 text-human" />;
		if (flag.includes("suspicious") || flag.includes("inhuman")) return <Bot className="w-4 h-4 text-module" />;
		return <MessageSquare className="w-4 h-4 text-primary" />;
	};

	return (
		<div className="glass-panel overflow-hidden">
			<div className="px-6 py-4 border-b border-slate-700/50 bg-slate-800/50">
				<h3 className="font-semibold text-slate-200">Detalle de Jugadas y Peritaje LLM</h3>
			</div>
			
			<div className="overflow-x-auto">
				<table className="w-full text-left border-collapse">
					<thead>
						<tr className="bg-slate-800/30 text-xs uppercase tracking-wider text-slate-400 border-b border-slate-700/50">
							<th className="px-6 py-4 font-medium">Nº</th>
							<th className="px-6 py-4 font-medium">Jugada</th>
							<th className="px-6 py-4 font-medium">CP Loss</th>
							<th className="px-6 py-4 font-medium">Clase</th>
							<th className="px-6 py-4 font-medium min-w-[300px]">Peritaje LLM (Solo 'M')</th>
						</tr>
					</thead>
					<tbody className="divide-y divide-slate-700/30">
						{moves.map((m) => (
							<tr key={m.move_number} className="hover:bg-slate-800/40 transition-colors">
								<td className="px-6 py-4 text-slate-500 font-mono text-sm">{m.move_number}</td>
								<td className="px-6 py-4">
									<div className="flex items-center space-x-2">
										<span className="font-bold text-slate-200">{m.san}</span>
										{m.cp_loss > 0 && m.san !== m.best_move && (
											<>
												<ArrowRight className="w-3 h-3 text-slate-600" />
												<span className="text-xs text-slate-500 font-mono" title="Best move">({m.best_move})</span>
											</>
										)}
									</div>
								</td>
								<td className="px-6 py-4">
									<span className={`font-mono text-sm ${m.cp_loss <= 10 ? 'text-module' : 'text-slate-400'}`}>
										{m.cp_loss}
									</span>
								</td>
								<td className="px-6 py-4">
									<span className={`inline-flex items-center justify-center px-2.5 py-0.5 rounded border text-xs font-bold ${getBadgeStyles(m.classification)}`}>
										{m.classification}
									</span>
								</td>
								<td className="px-6 py-4">
									{m.llm_flag ? (
										<div className="flex items-start space-x-2 text-sm">
											<div className="mt-0.5">{getLlmIcon(m.llm_flag)}</div>
											<span className={m.llm_flag.includes("natural") ? "text-slate-300" : "text-rose-300 font-medium"}>
												{m.llm_flag}
											</span>
										</div>
									) : (
										<span className="text-xs text-slate-600 italic">-</span>
									)}
								</td>
							</tr>
						))}
					</tbody>
				</table>
			</div>
		</div>
	);
}
