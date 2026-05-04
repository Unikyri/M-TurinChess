// Types from backend
export interface TraceStep {
	step: number;
	state: string;
	read_c1: string;
	read_c2: string;
	action: string;
	suspicion: number;
}

export interface MoveDetail {
	move_number: number;
	san: string;
	uci: string;
	cp_loss: number;
	classification: string;
	best_move: string;
	llm_flag: string | null;
}

export interface AnalysisResult {
	id: string;
	analyzed_at: string;
	verdict: string;
	suspicion_count: number;
	threshold: number;
	total_moves_analyzed: number;
	player_color: string;
	elo: number;
	tape_input: string[];
	tape_output: string[];
	move_details: MoveDetail[];
	mt_trace: TraceStep[];
}

export interface HistoryRecord {
	id: string;
	analyzed_at: string;
	verdict: string;
	suspicion_count: number;
	threshold: number;
	total_moves_analyzed: number;
	player_color: string;
	elo: number;
}

const API_BASE = 'http://localhost:8080/api';

export async function analyzePGN(formData: FormData): Promise<AnalysisResult> {
	const response = await fetch(`${API_BASE}/analyze`, {
		method: 'POST',
		body: formData, // browser automatically sets Content-Type for FormData
	});

	if (!response.ok) {
		let errorText = 'API Error';
		try {
			const errorData = await response.json();
			errorText = errorData.error || errorText;
		} catch (e) {
            errorText = await response.text();
        }
		throw new Error(errorText);
	}

	return response.json();
}

export async function getHistory(): Promise<HistoryRecord[]> {
	const response = await fetch(`${API_BASE}/history`);
	if (!response.ok) {
		throw new Error('Failed to fetch history');
	}
	return response.json();
}
