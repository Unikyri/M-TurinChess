import { useState, useEffect, useRef, memo } from 'react';
import { Chessboard } from 'react-chessboard';
import {
  ChevronLeft,
  ChevronRight,
  AlertCircle,
  CheckCircle,
  XCircle,
  HelpCircle,
  Loader2,
  Play,
  Pause,
  ChevronsLeft,
  ChevronsRight,
} from 'lucide-react';
import { TuringConsole } from './TuringConsole';
import { PlaybackControls } from './PlaybackControls';
import { TMStepDetailModal } from './TMStepDetailModal';

// Types matching backend API (compacted cells)
interface TapeStateCompact {
  head: number;
  minIdx: number;
}

interface StepTraceCompact {
  step: number;
  state: string;
  tape1: TapeStateCompact;
  tape2: TapeStateCompact;
  tape3: TapeStateCompact;
  read1: string;
  read2: string;
  read3: string;
  write1: string;
  write2: string;
  write3: string;
  dir1: string;
  dir2: string;
  dir3: string;
  nextState: string;
}

interface SimulationResultCompact {
  steps?: StepTraceCompact[];
  finalState: string;
  accepted: boolean;
  tape1Initial: string;
  tape2Initial: string;
  tape3Initial: string;
}

interface AnalyzeResponse {
  moves: string[];
  simplifiedFens: string[];
  simulations: SimulationResultCompact[];
  repetitionDetected: boolean;
  repetitionMoveIndex: number;
}

// Expanded formats for frontend rendering
interface ExpandedTapeState {
  cells: string[];
  head: number;
  minIdx: number;
}

interface ExpandedStepTrace {
  step: number;
  state: string;
  tape1: ExpandedTapeState;
  tape2: ExpandedTapeState;
  tape3: ExpandedTapeState;
  read1: string;
  read2: string;
  read3: string;
  write1: string;
  write2: string;
  write3: string;
  dir1: string;
  dir2: string;
  dir3: string;
  nextState: string;
}

// Interactive Board Memoized to isolate from TM step re-renders (Performance Optimization)
const InteractiveBoard = memo(({ fen }: { fen: string }) => {
  return (
    <Chessboard
      options={{
        position: fen,
        allowDragging: false,
        boardStyle: {
          borderRadius: '16px',
          boxShadow: '0 8px 30px rgba(0, 0, 0, 0.08)',
          width: '100%',
        },
        darkSquareStyle: { backgroundColor: '#769656' },
        lightSquareStyle: { backgroundColor: '#eeeed2' },
      }}
    />
  );
});

InteractiveBoard.displayName = 'InteractiveBoard';

// Reconstruct cells dynamically from the compact representation
function expandSimulationResult(compactResult: SimulationResultCompact): ExpandedStepTrace[] {
  if (!compactResult.steps || compactResult.steps.length === 0) return [];

  const steps = compactResult.steps;
  const numSteps = steps.length;

  const tape1Cells = new Map<number, string>();
  const tape2Cells = new Map<number, string>();
  const tape3Cells = new Map<number, string>();

  // Fill tapes 1 and 2 with their initial content
  for (let i = 0; i < compactResult.tape1Initial.length; i++) {
    tape1Cells.set(i, compactResult.tape1Initial[i]);
  }
  for (let i = 0; i < compactResult.tape2Initial.length; i++) {
    tape2Cells.set(i, compactResult.tape2Initial[i]);
  }

  const expandedSteps: ExpandedStepTrace[] = [];

  // Track the minimum and maximum indices used across simulation steps
  let t1Min = 0, t1Max = Math.max(0, compactResult.tape1Initial.length - 1);
  let t2Min = 0, t2Max = Math.max(0, compactResult.tape2Initial.length - 1);
  let t3Min = 0, t3Max = 0;
  tape3Cells.set(0, "_");

  for (let k = 0; k < numSteps; k++) {
    const step = steps[k];

    t1Min = Math.min(t1Min, step.tape1.minIdx, step.tape1.head);
    t1Max = Math.max(t1Max, step.tape1.head);

    t2Min = Math.min(t2Min, step.tape2.minIdx, step.tape2.head);
    t2Max = Math.max(t2Max, step.tape2.head);

    t3Min = Math.min(t3Min, step.tape3.minIdx, step.tape3.head);
    t3Max = Math.max(t3Max, step.tape3.head);

    const getCellsArray = (cellsMap: Map<number, string>, min: number, max: number) => {
      const arr: string[] = [];
      for (let i = min; i <= max; i++) {
        arr.push(cellsMap.get(i) || "_");
      }
      return arr;
    };

    const t1Arr = getCellsArray(tape1Cells, t1Min, t1Max);
    const t2Arr = getCellsArray(tape2Cells, t2Min, t2Max);
    const t3Arr = getCellsArray(tape3Cells, t3Min, t3Max);

    expandedSteps.push({
      step: step.step,
      state: step.state,
      tape1: {
        cells: t1Arr,
        head: step.tape1.head - t1Min,
        minIdx: t1Min,
      },
      tape2: {
        cells: t2Arr,
        head: step.tape2.head - t2Min,
        minIdx: t2Min,
      },
      tape3: {
        cells: t3Arr,
        head: step.tape3.head - t3Min,
        minIdx: t3Min,
      },
      read1: step.read1,
      read2: step.read2,
      read3: step.read3,
      write1: step.write1,
      write2: step.write2,
      write3: step.write3,
      dir1: step.dir1,
      dir2: step.dir2,
      dir3: step.dir3,
      nextState: step.nextState,
    });

    // Write to tapes for the next step simulation
    tape1Cells.set(step.tape1.head, step.write1);
    tape2Cells.set(step.tape2.head, step.write2);
    tape3Cells.set(step.tape3.head, step.write3);
  }

  return expandedSteps;
}

// Example PGNs
const REPETITION_PGN = `[Event "Example Repetition"]
[Result "*"]

1. e4 e5 2. Nf3 Nc6 3. Bc4 Bc5 4. O-O O-O 5. d3 d6 6. c3 a6 7. Bb3 Ba7 8. h3 h6 9. Re1 Re8 10. Nbd2 Be6 11. Bc2 Bd7 12. Nf1 Be6 13. Ng3 Bd7 14. Nf1 Be6 15. Ng3 Bd7 16. Nf1 Be6 17. Ng3 Bd7 *`;

const SCHOLARS_MATE_PGN = `[Event "Scholar's Mate"]
[Result "1-0"]

1. e4 e5 2. Qh5 Nc6 3. Bc4 Nf6 4. Qxf7#`;

const SIMPLE_REPETITION_SHORT = `1. Nf3 Nf6 2. Ng1 Ng8 3. Nf3 Nf6 4. Ng1 Ng8`;

export default function App() {
  const [pgn, setPgn] = useState(REPETITION_PGN);
  const [loading, setLoading] = useState(false);
  const [error, setError] = useState<string | null>(null);
  const [data, setData] = useState<AnalyzeResponse | null>(null);

  // Reference for scrolling to the Interactive Board Section on analysis finish
  const boardSectionRef = useRef<HTMLDivElement>(null);

  // Chess game playback states (Move timeline)
  const [selectedMoveIndex, setSelectedMoveIndex] = useState<number>(-1); // -1 is starting position
  const [isPlayingGame, setIsPlayingGame] = useState<boolean>(false);
  const [gamePlaybackSpeed, setGamePlaybackSpeed] = useState<number>(1500); // ms per move

  // Turing Machine simulation states for selected move (on-demand loading)
  const [activeSimulation, setActiveSimulation] = useState<ExpandedStepTrace[] | null>(null);
  const [loadingTM, setLoadingTM] = useState<boolean>(false);
  const [tmCache, setTmCache] = useState<Record<number, SimulationResultCompact>>({});

  const [activeTMStep, setActiveTMStep] = useState<number>(0);
  const [isPlayingTM, setIsPlayingTM] = useState<boolean>(false);
  const [tmPlaybackSpeed, setTmPlaybackSpeed] = useState<number>(150); // delay in ms
  const [isDebuggerOpen, setIsDebuggerOpen] = useState<boolean>(false);

  // Autoplay Chess Game moves
  useEffect(() => {
    let intervalId: any;
    if (isPlayingGame && data && data.moves && data.moves.length > 0) {
      intervalId = setInterval(() => {
        setSelectedMoveIndex((prev) => {
          if (prev >= data.moves.length - 1) {
            setIsPlayingGame(false);
            return prev;
          }
          return prev + 1;
        });
      }, gamePlaybackSpeed);
    }
    return () => {
      if (intervalId) clearInterval(intervalId);
    };
  }, [isPlayingGame, gamePlaybackSpeed, data]);

  // Load detailed Turing Machine simulation on-demand when selected move changes
  useEffect(() => {
    if (!data || selectedMoveIndex < 0) {
      setActiveSimulation(null);
      setActiveTMStep(0);
      setIsPlayingTM(false);
      return;
    }

    const loadTMSimulation = async () => {
      // Check local cache
      if (tmCache[selectedMoveIndex]) {
        const expanded = expandSimulationResult(tmCache[selectedMoveIndex]);
        setActiveSimulation(expanded);
        setActiveTMStep(0);
        setIsPlayingTM(false);
        return;
      }

      setLoadingTM(true);
      setActiveSimulation(null);
      setActiveTMStep(0);
      setIsPlayingTM(false);

      try {
        const history = data.simplifiedFens.slice(0, selectedMoveIndex + 1);
        const current = data.simplifiedFens[selectedMoveIndex + 1];

        const response = await fetch('http://localhost:8080/api/simulate', {
          method: 'POST',
          headers: {
            'Content-Type': 'application/json',
          },
          body: JSON.stringify({ history, current }),
        });

        const resBody = await response.json();

        if (!response.ok) {
          throw new Error(resBody.error || `Error de simulación: ${response.status}`);
        }

        // Cache the simulation result
        setTmCache((prev) => ({ ...prev, [selectedMoveIndex]: resBody }));

        const expanded = expandSimulationResult(resBody);
        setActiveSimulation(expanded);
      } catch (err: any) {
        console.error('Error cargando simulación detallada:', err);
        setError(err.message || 'No se pudo conectar con el servidor de simulación.');
      } finally {
        setLoadingTM(false);
      }
    };

    loadTMSimulation();
  }, [selectedMoveIndex, data]);

  // Autoplay Turing Machine steps for the currently selected move
  useEffect(() => {
    let intervalId: any;
    if (isPlayingTM && activeSimulation && activeSimulation.length > 0) {
      intervalId = setInterval(() => {
        setActiveTMStep((prev) => {
          if (prev >= activeSimulation.length - 1) {
            setIsPlayingTM(false);
            return prev;
          }
          return prev + 1;
        });
      }, tmPlaybackSpeed);
    }
    return () => {
      if (intervalId) clearInterval(intervalId);
    };
  }, [isPlayingTM, tmPlaybackSpeed, activeSimulation]);

  // Reset active TM step when the selected move changes
  const handleSelectMove = (index: number) => {
    setSelectedMoveIndex(index);
  };

  // Run analysis via Go API
  const handleAnalyze = async () => {
    setLoading(true);
    setError(null);
    setData(null);
    setSelectedMoveIndex(-1);
    setActiveSimulation(null);
    setTmCache({});
    setActiveTMStep(0);
    setIsPlayingTM(false);
    setIsPlayingGame(false);

    try {
      const response = await fetch('http://localhost:8080/api/analyze', {
        method: 'POST',
        headers: {
          'Content-Type': 'application/json',
        },
        body: JSON.stringify({ pgn }),
      });

      const resBody = await response.json();

      if (!response.ok) {
        throw new Error(resBody.error || `Error del servidor: ${response.status}`);
      }

      setData(resBody);
      
      // Auto select the first move after analysis if available
      if (resBody.moves && resBody.moves.length > 0) {
        setSelectedMoveIndex(0);
      }

      // Smooth scroll to the interactive section after a tiny layout timeout
      setTimeout(() => {
        boardSectionRef.current?.scrollIntoView({ behavior: 'smooth' });
      }, 100);
    } catch (err: any) {
      console.error(err);
      setError(err.message || 'Error de conexión con el servidor. Verifica que se encuentre ejecutando en el puerto 8080.');
    } finally {
      setLoading(false);
    }
  };

  // Helper to extract the board part of a FEN
  const getBoardFen = (fullFen: string) => {
    if (!fullFen) return 'rnbqkbnr/pppppppp/8/8/8/8/PPPPPPPP/RNBQKBNR';
    return fullFen.split('|')[0];
  };

  // Format move pairs for rendering (e.g. 1. e4 e5)
  const renderMoveTimeline = () => {
    if (!data || !data.moves) return null;

    const moves = data.moves;
    const pairs: Array<{ number: number; white: string; whiteIdx: number; black?: string; blackIdx?: number }> = [];

    for (let i = 0; i < moves.length; i += 2) {
      pairs.push({
        number: Math.floor(i / 2) + 1,
        white: moves[i],
        whiteIdx: i,
        black: moves[i + 1] ? moves[i + 1] : undefined,
        blackIdx: moves[i + 1] ? i + 1 : undefined,
      });
    }

    return (
      <div className="space-y-1.5 pr-1 text-sm">
        {/* Start Position row */}
        <button
          onClick={() => handleSelectMove(-1)}
          className={`w-full text-left px-3.5 py-2.5 rounded-xl font-mono transition-colors font-bold ${
            selectedMoveIndex === -1
              ? 'bg-sky-600 text-white font-extrabold shadow-md'
              : 'hover:bg-neutral-100 text-neutral-700 bg-neutral-50'
          }`}
        >
          [POSICIÓN INICIAL]
        </button>

        {pairs.map((pair) => {
          const isWhiteSelected = selectedMoveIndex === pair.whiteIdx;
          const isBlackSelected = pair.blackIdx !== undefined && selectedMoveIndex === pair.blackIdx;

          const isWhiteRepetition = data.repetitionDetected && data.repetitionMoveIndex === pair.whiteIdx;
          const isBlackRepetition = data.repetitionDetected && pair.blackIdx !== undefined && data.repetitionMoveIndex === pair.blackIdx;

          return (
            <div key={pair.number} className="flex items-center gap-2 py-1 border-b border-neutral-150">
              <span className="text-xs font-mono text-neutral-450 w-7 text-right shrink-0">
                {pair.number}.
              </span>
              
              {/* White Move */}
              <button
                onClick={() => handleSelectMove(pair.whiteIdx)}
                className={`flex-1 text-left px-3 py-2 rounded-lg font-mono transition-all font-semibold ${
                  isWhiteSelected
                    ? 'bg-sky-600 text-white font-extrabold shadow-sm'
                    : isWhiteRepetition
                    ? 'bg-amber-100 text-amber-800 font-black border border-amber-300'
                    : 'hover:bg-neutral-100 text-neutral-800 bg-white border border-neutral-200'
                }`}
              >
                {pair.white}
                {isWhiteRepetition && ' ⚠️'}
              </button>

              {/* Black Move */}
              {pair.black !== undefined && pair.blackIdx !== undefined ? (
                <button
                  onClick={() => handleSelectMove(pair.blackIdx!)}
                  className={`flex-1 text-left px-3 py-2 rounded-lg font-mono transition-all font-semibold ${
                    isBlackSelected
                      ? 'bg-sky-600 text-white font-extrabold shadow-sm'
                      : isBlackRepetition
                      ? 'bg-amber-100 text-amber-800 font-black border border-amber-300'
                      : 'hover:bg-neutral-100 text-neutral-800 bg-white border border-neutral-200'
                  }`}
                >
                  {pair.black}
                  {isBlackRepetition && ' ⚠️'}
                </button>
              ) : (
                <div className="flex-1" />
              )}
            </div>
          );
        })}
      </div>
    );
  };

  // Get FEN to display on chessboard
  const currentFen = data
    ? selectedMoveIndex === -1
      ? data.simplifiedFens[0]
      : data.simplifiedFens[selectedMoveIndex + 1]
    : 'rnbqkbnr/pppppppp/8/8/8/8/PPPPPPPP/RNBQKBNR';

  // Get active TM simulation details
  const currentTMStepTrace = activeSimulation ? activeSimulation[activeTMStep] : null;
  const totalTMSteps = activeSimulation ? activeSimulation.length : 0;
  const currentSimulationCompact = data && selectedMoveIndex >= 0 ? data.simulations[selectedMoveIndex] : null;

  return (
    <div className="min-h-screen bg-neutral-50 text-neutral-850 flex flex-col font-sans">

      {/* TOP HEADER */}
      <header className="border-b border-neutral-200 bg-white/80 backdrop-blur-md sticky top-0 z-40 px-6 py-4 flex justify-between items-center shrink-0 shadow-sm">
        <div className="flex items-center gap-3">
          <div className="bg-sky-600 text-white p-2.5 rounded-xl shadow-md shadow-sky-500/25">
            <span className="font-mono font-black text-base tracking-tighter">MT</span>
          </div>
          <div>
            <h1 className="text-lg font-black tracking-wider uppercase text-neutral-900 leading-tight">
              Detector de Triple Repetición en Ajedrez
            </h1>
            <p className="text-[11px] text-neutral-500 uppercase tracking-widest font-bold">
              Analizador de Lenguaje Formal mediante Máquina de Turing de 3 Cintas
            </p>
          </div>
        </div>
      </header>

      {/* MAIN LAYOUT */}
      <main className="flex-1 flex flex-col gap-6 py-6 overflow-y-auto">
        
        {/* SECCIÓN 1: Entrada de Partida (PGN) */}
        <section className="max-w-7xl w-full mx-auto px-6 shrink-0">
          <div className="bg-white border border-neutral-200 rounded-2xl p-6 shadow-sm space-y-5">
            <div className="flex justify-between items-center">
              <h2 className="text-xs font-bold tracking-wider text-neutral-500 uppercase">
                Entrada de Notación PGN
              </h2>
              <span title="Pega una notación PGN estándar de ajedrez para iniciar el análisis.">
                <HelpCircle className="w-5 h-5 text-neutral-400 cursor-help hover:text-neutral-600" />
              </span>
            </div>

            <textarea
              value={pgn}
              onChange={(e) => setPgn(e.target.value)}
              placeholder="Pega la partida en PGN aquí..."
              rows={9}
              className="w-full text-sm font-mono p-4 bg-white border border-neutral-200 rounded-xl focus:outline-none focus:ring-2 focus:ring-sky-500 focus:border-transparent resize-none leading-relaxed text-neutral-900"
            />

            {/* Presets Row */}
            <div className="flex flex-wrap gap-2">
              <button
                onClick={() => setPgn(REPETITION_PGN)}
                className="px-3.5 py-2 bg-neutral-100 hover:bg-neutral-200 rounded-lg text-xs font-semibold tracking-wide transition-colors text-neutral-700"
              >
                Cargar Repetición
              </button>
              <button
                onClick={() => setPgn(SIMPLE_REPETITION_SHORT)}
                className="px-3.5 py-2 bg-neutral-100 hover:bg-neutral-200 rounded-lg text-xs font-semibold tracking-wide transition-colors text-neutral-700"
              >
                Cargar Rep. Corta
              </button>
              <button
                onClick={() => setPgn(SCHOLARS_MATE_PGN)}
                className="px-3.5 py-2 bg-neutral-100 hover:bg-neutral-200 rounded-lg text-xs font-semibold tracking-wide transition-colors text-neutral-700"
              >
                Cargar Mate del Pastor
              </button>
            </div>

            {/* Analyze Button */}
            <button
              onClick={handleAnalyze}
              disabled={loading || !pgn.trim()}
              className="w-full py-4 bg-sky-600 hover:bg-sky-700 active:bg-sky-800 text-white rounded-xl text-sm font-bold uppercase tracking-wider shadow-md shadow-sky-500/20 disabled:opacity-50 disabled:pointer-events-none transition-colors flex items-center justify-center gap-2"
            >
              {loading ? (
                <>
                  <Loader2 className="animate-spin w-5 h-5" />
                  <span>Analizando Partida...</span>
                </>
              ) : (
                <span>Analizar Partida</span>
              )}
            </button>
          </div>

          {/* Error Banner */}
          {error && (
            <div className="bg-red-50 border border-red-200 text-red-800 p-4 rounded-xl flex gap-3 text-xs leading-normal mt-4">
              <AlertCircle className="w-5 h-5 shrink-0 mt-0.5" />
              <div>
                <span className="font-bold">Análisis Fallido:</span> {error}
              </div>
            </div>
          )}
        </section>

        {/* SECCIÓN 2: Tablero Interactivo y Línea de Tiempo */}
        {data && (
          <section ref={boardSectionRef} className="max-w-7xl w-full mx-auto px-6 grid grid-cols-1 lg:grid-cols-12 gap-6 scroll-mt-20 shrink-0">
            
            {/* Tablero Interactivo (8/12 Ancho) */}
            <div className="lg:col-span-8 bg-white border border-neutral-200 rounded-2xl p-6 shadow-sm flex flex-col gap-5 justify-between">
              
              <div className="flex justify-between items-center text-xs border-b border-neutral-100 pb-3">
                <span className="font-bold text-neutral-500 uppercase tracking-wide">
                  Tablero Interactivo de Ajedrez
                </span>
                <span className="font-mono text-neutral-750 bg-neutral-100 px-3 py-1 rounded font-bold text-xs">
                  {selectedMoveIndex === -1 ? 'Posición inicial' : `Jugada #${selectedMoveIndex + 1}: ${data?.moves[selectedMoveIndex]}`}
                </span>
              </div>

              {/* Chessboard Container - Agrandado para mayor visibilidad */}
              <div className="w-full max-w-[460px] aspect-square mx-auto my-4 chessboard-container rounded-2xl overflow-hidden border border-neutral-200 shadow-md">
                <InteractiveBoard fen={getBoardFen(currentFen)} />
              </div>

              {/* Chess Game Playback Controls (Tablero con Reproductor Automático) */}
              <div className="flex flex-col gap-3 border-t border-neutral-150 pt-4">
                <div className="flex justify-between items-center text-xs font-bold text-neutral-500 px-1">
                  <span>Reproducción automática de la partida</span>
                  <span className="font-mono text-neutral-600 bg-neutral-100 px-2 py-0.5 rounded">
                    Jugada {selectedMoveIndex + 1} de {data.moves.length}
                  </span>
                </div>

                <div className="flex flex-col sm:flex-row items-center justify-between gap-3 bg-neutral-50 p-2.5 rounded-xl border border-neutral-150">
                  {/* Game Playback Speed Slider */}
                  <div className="flex items-center gap-2 min-w-[170px] w-full sm:w-auto">
                    <span className="text-[11px] font-bold text-neutral-500 whitespace-nowrap">
                      Intervalo: {gamePlaybackSpeed}ms
                    </span>
                    <input
                      type="range"
                      min="500"
                      max="3000"
                      step="250"
                      value={gamePlaybackSpeed}
                      onChange={(e) => setGamePlaybackSpeed(Number(e.target.value))}
                      className="w-full h-1 bg-neutral-200 rounded-lg appearance-none cursor-pointer accent-sky-600"
                    />
                  </div>

                  {/* Playback Buttons */}
                  <div className="flex items-center justify-center gap-1 bg-white p-1 rounded-lg border border-neutral-200">
                    <button
                      onClick={() => handleSelectMove(-1)}
                      disabled={selectedMoveIndex === -1}
                      className="p-1.5 text-neutral-600 hover:text-neutral-900 hover:bg-neutral-100 rounded disabled:opacity-30 transition-colors"
                      title="Posición inicial"
                    >
                      <ChevronsLeft className="w-4.5 h-4.5" />
                    </button>

                    <button
                      onClick={() => handleSelectMove(selectedMoveIndex - 1)}
                      disabled={selectedMoveIndex === -1}
                      className="p-1.5 text-neutral-600 hover:text-neutral-900 hover:bg-neutral-100 rounded disabled:opacity-30 transition-colors"
                      title="Jugada anterior"
                    >
                      <ChevronLeft className="w-4.5 h-4.5" />
                    </button>

                    <button
                      onClick={() => setIsPlayingGame(!isPlayingGame)}
                      disabled={data.moves.length === 0}
                      className="mx-1 p-2.5 bg-sky-600 hover:bg-sky-700 text-white rounded-full shadow transition-colors"
                      title={isPlayingGame ? 'Pausar partida' : 'Reproducir partida'}
                    >
                      {isPlayingGame ? <Pause className="w-4.5 h-4.5 fill-current" /> : <Play className="w-4.5 h-4.5 fill-current ml-0.5" />}
                    </button>

                    <button
                      onClick={() => handleSelectMove(selectedMoveIndex + 1)}
                      disabled={selectedMoveIndex >= data.moves.length - 1}
                      className="p-1.5 text-neutral-600 hover:text-neutral-900 hover:bg-neutral-100 rounded disabled:opacity-30 transition-colors"
                      title="Siguiente jugada"
                    >
                      <ChevronRight className="w-4.5 h-4.5" />
                    </button>

                    <button
                      onClick={() => handleSelectMove(data.moves.length - 1)}
                      disabled={selectedMoveIndex >= data.moves.length - 1}
                      className="p-1.5 text-neutral-600 hover:text-neutral-900 hover:bg-neutral-100 rounded disabled:opacity-30 transition-colors"
                      title="Saltar al final"
                    >
                      <ChevronsRight className="w-4.5 h-4.5" />
                    </button>
                  </div>
                </div>
              </div>
            </div>

            {/* Línea de Tiempo de Jugadas (4/12 Ancho) */}
            <div className="lg:col-span-4 bg-white border border-neutral-200 rounded-2xl p-6 shadow-sm flex flex-col gap-4">
              <h2 className="text-xs font-bold tracking-wider text-neutral-500 uppercase border-b border-neutral-100 pb-3">
                Línea de Tiempo de Jugadas
              </h2>
              <div className="flex-1 overflow-y-auto max-h-[380px] lg:max-h-[500px] custom-scrollbar">
                {renderMoveTimeline()}
              </div>

              {/* Verdict Status Info Card inside timeline column */}
              <div className={`border rounded-xl p-4 flex flex-col gap-2 ${
                data.repetitionDetected
                  ? 'bg-emerald-50/50 border-emerald-200'
                  : 'bg-neutral-50 border-neutral-200'
              }`}>
                <div className="flex items-center gap-1.5">
                  {data.repetitionDetected ? (
                    <CheckCircle className="w-4.5 h-4.5 text-emerald-600 shrink-0" />
                  ) : (
                    <XCircle className="w-4.5 h-4.5 text-neutral-400 shrink-0" />
                  )}
                  <span className="text-[10px] font-bold tracking-wider uppercase text-neutral-500">
                    Veredicto
                  </span>
                </div>
                <div className={`text-xl font-black ${
                  data.repetitionDetected ? 'text-emerald-600' : 'text-neutral-800'
                }`}>
                  {data.repetitionDetected ? 'ACEPTADO' : 'RECHAZADO'}
                </div>
              </div>
            </div>

          </section>
        )}

        {/* SECCIÓN 3: Consola de la Máquina de Turing */}
        {data && selectedMoveIndex >= 0 && (
          <section className="max-w-7xl w-full mx-auto px-6 flex flex-col gap-4 shrink-0">
            
            {/* Visual Tape Display */}
            {loadingTM ? (
              <div className="bg-white border border-neutral-200 rounded-2xl p-12 shadow-sm flex flex-col items-center justify-center gap-3">
                <Loader2 className="animate-spin text-sky-650 w-8 h-8" />
                <span className="text-sm font-semibold text-neutral-500">
                  Cargando simulación de la Máquina de Turing...
                </span>
              </div>
            ) : (
              <>
                <TuringConsole
                  tape1={currentTMStepTrace?.tape1 || null}
                  tape2={currentTMStepTrace?.tape2 || null}
                  tape3={currentTMStepTrace?.tape3 || null}
                  activeState={currentTMStepTrace?.state || 'Halt'}
                  stepNumber={activeTMStep}
                  totalSteps={totalTMSteps}
                  onStepChange={(stepIdx) => {
                    setActiveTMStep(stepIdx);
                    setIsPlayingTM(false); // Pause autocomplete during drag
                  }}
                />

                {/* Playback Controls & Transition details */}
                <div className="space-y-3">
                  {currentTMStepTrace && (
                    <div className="flex flex-wrap items-center justify-between text-xs font-mono bg-neutral-100 border border-neutral-200 px-4 py-3 rounded-xl text-neutral-800">
                      <div className="flex items-center gap-2">
                        <span className="font-bold text-sky-700">Transición Activa:</span>
                        <span>
                          δ({currentTMStepTrace.state}, [{currentTMStepTrace.read1 === '_' ? '•' : currentTMStepTrace.read1}, {currentTMStepTrace.read2 === '_' ? '•' : currentTMStepTrace.read2}, {currentTMStepTrace.read3 === '_' ? '•' : currentTMStepTrace.read3}])
                          {' '}→{' '}
                          ({currentTMStepTrace.nextState}, [{currentTMStepTrace.write1 === '_' ? '•' : currentTMStepTrace.write1}, {currentTMStepTrace.write2 === '_' ? '•' : currentTMStepTrace.write2}, {currentTMStepTrace.write3 === '_' ? '•' : currentTMStepTrace.write3}], [{currentTMStepTrace.dir1}, {currentTMStepTrace.dir2}, {currentTMStepTrace.dir3}])
                        </span>
                      </div>
                      <div className="font-semibold">
                        Veredicto de Jugada: <span className={`font-black ${
                          currentSimulationCompact?.accepted ? 'text-emerald-600' : 'text-red-500'
                        }`}>
                          {currentSimulationCompact?.accepted ? 'ACEPTADO' : 'RECHAZADO'}
                        </span>
                      </div>
                    </div>
                  )}

                  <PlaybackControls
                    isPlaying={isPlayingTM}
                    onPlayPause={() => setIsPlayingTM(!isPlayingTM)}
                    onStepForward={() => setActiveTMStep((prev) => Math.min(prev + 1, totalTMSteps - 1))}
                    onStepBackward={() => setActiveTMStep((prev) => Math.max(prev - 1, 0))}
                    onFirstStep={() => setActiveTMStep(0)}
                    onLastStep={() => setActiveTMStep(totalTMSteps - 1)}
                    speed={tmPlaybackSpeed}
                    onSpeedChange={setTmPlaybackSpeed}
                    onToggleDebugger={() => setIsDebuggerOpen(true)}
                    activeStep={activeTMStep}
                    totalSteps={totalTMSteps}
                  />
                </div>
              </>
            )}
          </section>
        )}

      </main>

      {/* STEP DEBUGGER MODAL */}
      {data && activeSimulation && (
        <TMStepDetailModal
          isOpen={isDebuggerOpen}
          onClose={() => setIsDebuggerOpen(false)}
          steps={activeSimulation as any}
          activeStep={activeTMStep}
          onSelectStep={(stepIdx) => setActiveTMStep(stepIdx)}
        />
      )}

    </div>
  );
}
