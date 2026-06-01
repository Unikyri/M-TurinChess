import React from 'react';
import {
  Play,
  Pause,
  ChevronLeft,
  ChevronRight,
  ChevronsLeft,
  ChevronsRight,
  Terminal,
} from 'lucide-react';

interface PlaybackControlsProps {
  isPlaying: boolean;
  onPlayPause: () => void;
  onStepForward: () => void;
  onStepBackward: () => void;
  onFirstStep: () => void;
  onLastStep: () => void;
  speed: number; // in milliseconds
  onSpeedChange: (speed: number) => void;
  onToggleDebugger: () => void;
  activeStep: number;
  totalSteps: number;
}

export const PlaybackControls: React.FC<PlaybackControlsProps> = ({
  isPlaying,
  onPlayPause,
  onStepForward,
  onStepBackward,
  onFirstStep,
  onLastStep,
  speed,
  onSpeedChange,
  onToggleDebugger,
  activeStep,
  totalSteps,
}) => {
  return (
    <div className="bg-white border border-neutral-205 rounded-2xl p-5 shadow-sm flex flex-col md:flex-row md:items-center justify-between gap-4">
      {/* Simulation Speed & Slider */}
      <div className="flex items-center gap-3 min-w-[220px]">
        <span className="text-xs font-bold text-neutral-500 whitespace-nowrap">
          Velocidad: {speed}ms
        </span>
        <input
          type="range"
          min="50"
          max="1000"
          step="50"
          value={speed}
          onChange={(e) => onSpeedChange(Number(e.target.value))}
          className="w-full h-1.5 bg-neutral-200 rounded-lg appearance-none cursor-pointer accent-sky-600"
        />
      </div>

      {/* Control Buttons */}
      <div className="flex items-center justify-center gap-1 bg-neutral-50 p-1.5 rounded-xl border border-neutral-100">
        <button
          onClick={onFirstStep}
          disabled={activeStep === 0}
          className="p-2 text-neutral-600 hover:text-neutral-950 hover:bg-neutral-150 rounded-lg disabled:opacity-30 disabled:pointer-events-none transition-colors"
          title="Regresar al inicio"
        >
          <ChevronsLeft className="w-4 h-4" />
        </button>

        <button
          onClick={onStepBackward}
          disabled={activeStep === 0}
          className="p-2 text-neutral-600 hover:text-neutral-950 hover:bg-neutral-150 rounded-lg disabled:opacity-30 disabled:pointer-events-none transition-colors"
          title="Paso anterior"
        >
          <ChevronLeft className="w-4 h-4" />
        </button>

        <button
          onClick={onPlayPause}
          disabled={totalSteps <= 1}
          className="mx-1 p-3 bg-sky-600 hover:bg-sky-700 active:bg-sky-800 text-white rounded-full shadow-md shadow-sky-500/20 disabled:opacity-50 disabled:pointer-events-none transition-colors duration-150"
          title={isPlaying ? 'Pausar' : 'Reproducir automáticamente'}
        >
          {isPlaying ? <Pause className="w-5 h-5 fill-current" /> : <Play className="w-5 h-5 fill-current ml-0.5" />}
        </button>

        <button
          onClick={onStepForward}
          disabled={activeStep >= totalSteps - 1}
          className="p-2 text-neutral-600 hover:text-neutral-950 hover:bg-neutral-150 rounded-lg disabled:opacity-30 disabled:pointer-events-none transition-colors"
          title="Siguiente paso"
        >
          <ChevronRight className="w-4 h-4" />
        </button>

        <button
          onClick={onLastStep}
          disabled={activeStep >= totalSteps - 1}
          className="p-2 text-neutral-600 hover:text-neutral-950 hover:bg-neutral-150 rounded-lg disabled:opacity-30 disabled:pointer-events-none transition-colors"
          title="Saltar al final"
        >
          <ChevronsRight className="w-4 h-4" />
        </button>
      </div>

      {/* Debugger Toggle */}
      <div className="flex items-center justify-end">
        <button
          onClick={onToggleDebugger}
          className="flex items-center gap-2 px-3 py-2 bg-neutral-100 hover:bg-neutral-200 active:bg-neutral-250 text-neutral-700 rounded-xl text-xs font-semibold tracking-wide border border-neutral-200/50 transition-colors"
        >
          <Terminal className="w-3.5 h-3.5" />
          <span>Depurador de Transiciones</span>
        </button>
      </div>
    </div>
  );
};
