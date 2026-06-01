import React, { useState } from 'react';

interface TapeState {
  cells: string[];
  head: number;
  minIdx: number;
}

interface TuringConsoleProps {
  tape1: TapeState | null;
  tape2: TapeState | null;
  tape3: TapeState | null;
  activeState: string;
  stepNumber: number;
  totalSteps: number;
  onStepChange?: (step: number) => void;
}

const CELL_WIDTH = 48; // 48px
const GAP = 4; // gap-1 (4px)
const STEP_SIZE = CELL_WIDTH + GAP;

const TapeRow: React.FC<{
  title: string;
  tape: TapeState | null;
  isCounterTape?: boolean;
  onDragStep?: (delta: number) => void;
}> = ({ title, tape, isCounterTape, onDragStep }) => {
  const [isDragging, setIsDragging] = useState(false);
  const [startX, setStartX] = useState(0);

  if (!tape || !tape.cells || tape.cells.length === 0) {
    return (
      <div className="flex flex-col gap-1.5 my-4">
        <div className="flex justify-between text-xs font-semibold text-neutral-500">
          <span>{title}</span>
          <span className="font-mono text-neutral-450">Vacía</span>
        </div>
        <div className="relative w-full h-16 bg-neutral-100 border border-neutral-200 rounded-xl flex items-center justify-center text-neutral-450 text-sm">
          Sin actividad en la cinta
        </div>
      </div>
    );
  }

  const { cells, head } = tape;
  const offset = head * STEP_SIZE + CELL_WIDTH / 2;

  const handleMouseDown = (e: React.MouseEvent) => {
    e.preventDefault();
    setIsDragging(true);
    setStartX(e.clientX);
  };

  const handleMouseMove = (e: React.MouseEvent) => {
    if (!isDragging) return;
    const diffX = e.clientX - startX;
    const threshold = 35; // 35 píxeles de arrastre equivalen a 1 paso
    if (Math.abs(diffX) >= threshold) {
      const delta = Math.round(diffX / threshold);
      if (onDragStep) {
        onDragStep(-delta); // Arrastrar a la izquierda avanza pasos
      }
      setStartX(e.clientX);
    }
  };

  const handleMouseUpOrLeave = () => {
    setIsDragging(false);
  };

  return (
    <div className="flex flex-col gap-2 my-4">
      <div className="flex justify-between items-center text-xs font-medium text-neutral-500 px-1">
        <span className="font-bold tracking-wider text-neutral-700 text-[13px]">{title}</span>
        <span className="font-mono bg-neutral-200 px-2 py-0.5 rounded text-[10px] text-neutral-700 font-bold">
          Cabezal: {head + tape.minIdx} (índice {head})
        </span>
      </div>

      <div
        onMouseDown={handleMouseDown}
        onMouseMove={handleMouseMove}
        onMouseUp={handleMouseUpOrLeave}
        onMouseLeave={handleMouseUpOrLeave}
        className={`relative w-full h-20 bg-sky-50/20 border rounded-xl overflow-hidden flex items-center select-none transition-shadow ${
          isDragging
            ? 'cursor-grabbing border-sky-400 shadow-inner'
            : 'cursor-grab border-neutral-200 hover:border-sky-300 hover:bg-sky-50/30'
        }`}
        title="Haz clic sostenido y arrastra horizontalmente para avanzar o retroceder la simulación"
      >
        {/* Center alignment guide */}
        <div className="absolute left-1/2 top-0 bottom-0 w-[52px] -ml-[26px] border-l border-r border-dashed border-sky-400/40 bg-sky-500/5 pointer-events-none z-10" />
        <div className="absolute left-1/2 top-0 -translate-x-1/2 z-20 pointer-events-none">
          <div className="w-0 h-0 border-l-[6px] border-r-[6px] border-t-[6px] border-l-transparent border-r-transparent border-t-sky-500" />
        </div>
        <div className="absolute left-1/2 bottom-0 -translate-x-1/2 z-20 pointer-events-none">
          <div className="w-0 h-0 border-l-[6px] border-r-[6px] border-b-[6px] border-l-transparent border-r-transparent border-b-sky-500" />
        </div>

        {/* Sliding tape content */}
        <div
          className="absolute left-1/2 flex gap-1 h-12 items-center tape-slider pointer-events-none"
          style={{
            transform: `translateX(-${offset}px)`,
          }}
        >
          {cells.map((cell, idx) => {
            const isActive = idx === head;
            const isDelimiter = cell === '$';
            const isBlank = cell === '_';
            const isOne = cell === '1';

            let cellClass = "w-12 h-12 flex items-center justify-center text-base font-mono border rounded-lg transition-all duration-200 shrink-0 select-none ";
            if (isActive) {
              cellClass += "border-sky-500 bg-sky-100 text-sky-800 font-black shadow-md shadow-sky-500/10 scale-105 z-10";
            } else if (isDelimiter) {
              cellClass += "border-neutral-300 bg-neutral-200 text-neutral-800 font-bold";
            } else if (isBlank) {
              cellClass += "border-neutral-200 bg-neutral-50/50 text-neutral-400";
            } else if (isOne && isCounterTape) {
              cellClass += "border-sky-300 bg-sky-50 text-sky-700 font-bold";
            } else {
              cellClass += "border-neutral-200 bg-white text-neutral-800";
            }

            return (
              <div key={idx} className={cellClass} title={`Índice: ${idx + tape.minIdx}`}>
                {isBlank ? '•' : cell}
              </div>
            );
          })}
        </div>
      </div>
    </div>
  );
};

export const TuringConsole: React.FC<TuringConsoleProps> = ({
  tape1,
  tape2,
  tape3,
  activeState,
  stepNumber,
  totalSteps,
  onStepChange,
}) => {
  const handleDragStep = (delta: number) => {
    if (onStepChange) {
      const nextStep = Math.max(0, Math.min(totalSteps - 1, stepNumber + delta));
      onStepChange(nextStep);
    }
  };

  return (
    <div className="bg-white border border-neutral-200 rounded-2xl p-6 shadow-sm">
      <div className="flex justify-between items-center border-b border-neutral-100 pb-3 mb-2">
        <div className="flex items-center gap-2">
          <span className="flex h-2.5 w-2.5 rounded-full bg-sky-550 animate-pulse" />
          <h3 className="text-base font-bold tracking-wide uppercase text-neutral-800">
            Consola de la Máquina de Turing
          </h3>
        </div>
        <div className="flex gap-4 text-xs font-bold text-neutral-600">
          <div>
            Estado: <span className="font-mono font-bold text-sky-700 bg-sky-50 px-2.5 py-1 rounded-md border border-sky-100">{activeState}</span>
          </div>
          <div className="flex items-center">
            Paso: <span className="font-mono text-neutral-800 ml-1">{stepNumber}</span> <span className="mx-1">/</span> <span className="font-mono text-neutral-450">{totalSteps > 0 ? totalSteps - 1 : 0}</span>
          </div>
        </div>
      </div>

      <div className="space-y-1">
        <TapeRow title="Cinta 1: Historial de FENs" tape={tape1} onDragStep={handleDragStep} />
        <TapeRow title="Cinta 2: FEN de Jugada Actual" tape={tape2} onDragStep={handleDragStep} />
        <TapeRow title="Cinta 3: Conteo Unario de Coincidencias" tape={tape3} isCounterTape={true} onDragStep={handleDragStep} />
      </div>
    </div>
  );
};
