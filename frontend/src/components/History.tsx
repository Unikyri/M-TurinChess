import { useEffect, useState } from "react";
import { getHistory, type HistoryRecord } from "../api/client";
import { Clock, CheckCircle, AlertTriangle } from "lucide-react";

export function History() {
	const [history, setHistory] = useState<HistoryRecord[]>([]);
	const [loading, setLoading] = useState(true);

	useEffect(() => {
		loadHistory();
	}, []);

	const loadHistory = async () => {
		try {
			const data = await getHistory();
			// Sort newest first
			const sorted = data.sort((a, b) => new Date(b.analyzed_at).getTime() - new Date(a.analyzed_at).getTime());
			setHistory(sorted);
		} catch (e) {
			console.error("Error loading history:", e);
		} finally {
			setLoading(false);
		}
	};

	if (loading) {
		return <div className="text-center text-slate-500 py-8">Cargando historial...</div>;
	}

	if (history.length === 0) {
		return null; // Don't show anything if history is empty
	}

	return (
		<div className="w-full max-w-4xl mx-auto mt-16 mb-8">
			<div className="flex items-center space-x-2 mb-6">
				<Clock className="w-5 h-5 text-slate-400" />
				<h2 className="text-xl font-bold text-slate-200">Análisis Recientes</h2>
			</div>

			<div className="grid grid-cols-1 md:grid-cols-2 lg:grid-cols-3 gap-4">
				{history.map((record) => {
					const isModule = record.verdict === "MODULE_DETECTED";
					const date = new Date(record.analyzed_at).toLocaleString();

					return (
						<div key={record.id} className="glass-panel p-5 hover:bg-surface/90 transition-colors flex flex-col justify-between">
							<div className="flex justify-between items-start mb-4">
								<div className="text-xs font-mono text-slate-500">{record.id.split('-')[0]}</div>
								<div className="text-xs text-slate-400">{date}</div>
							</div>
							
							<div className="flex items-center space-x-3 mb-4">
								<div className={`p-2 rounded-full ${isModule ? "bg-module/20 text-module" : "bg-human/20 text-human"}`}>
									{isModule ? <AlertTriangle className="w-5 h-5" /> : <CheckCircle className="w-5 h-5" />}
								</div>
								<div>
									<h3 className={`font-bold ${isModule ? "text-module" : "text-human"}`}>
										{isModule ? "Módulo" : "Humano"}
									</h3>
									<p className="text-sm text-slate-400 capitalize">{record.player_color} ({record.elo})</p>
								</div>
							</div>

							<div className="mt-auto border-t border-slate-700/50 pt-3 flex justify-between text-sm">
								<span className="text-slate-400">Sospecha MT:</span>
								<span className="font-mono text-slate-200">{record.suspicion_count} / {record.threshold}</span>
							</div>
						</div>
					);
				})}
			</div>
		</div>
	);
}
