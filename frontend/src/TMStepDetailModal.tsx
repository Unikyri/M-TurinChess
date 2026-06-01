import React, { useEffect, useRef } from 'react';
import { X, ArrowRight } from 'lucide-react';

interface TapeState {
  cells: string[];
  head: number;
  minIdx: number;
}

interface StepTrace {
  step: number;
  state: string;
  tape1: TapeState;
  tape2: TapeState;
  tape3: TapeState;
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

interface TMStepDetailModalProps {
  isOpen: boolean;
  onClose: () => void;
  steps: StepTrace[];
  activeStep: number;
  onSelectStep: (stepIndex: number) => void;
}

export const TMStepDetailModal: React.FC<TMStepDetailModalProps> = ({
  isOpen,
  onClose,
  steps,
  activeStep,
  onSelectStep,
}) => {
  const activeRowRef = useRef<HTMLTableRowElement>(null);
  const tableContainerRef = useRef<HTMLDivElement>(null);

  // Auto-scroll the active step into view inside the list
  useEffect(() => {
    if (isOpen && activeRowRef.current && tableContainerRef.current) {
      const container = tableContainerRef.current;
      const row = activeRowRef.current;
      
      const rowTop = row.offsetTop;
      const rowHeight = row.offsetHeight;
      const containerHeight = container.offsetHeight;
      
      container.scrollTop = rowTop - containerHeight / 2 + rowHeight / 2;
    }
  }, [isOpen, activeStep]);

  if (!isOpen) return null;

  const currentStep = steps[activeStep];

  return (
    <div className="fixed inset-0 bg-neutral-900/60 backdrop-blur-sm flex items-center justify-center z-50 p-4 animate-in fade-in duration-200">
      <div className="bg-white w-full max-w-5xl h-[85vh] rounded-2xl shadow-2xl border border-neutral-200 flex flex-col overflow-hidden animate-in zoom-in-95 duration-200">
        
        {/* Header */}
        <div className="flex justify-between items-center px-6 py-4 border-b border-neutral-150 shrink-0">
          <div>
            <h2 className="text-lg font-bold text-neutral-800 flex items-center gap-2">
              Depurador de Transiciones de la Máquina de Turing
            </h2>
            <p className="text-xs text-neutral-500">
              Inspecciona paso a paso los cambios de estado y los movimientos de los cabezales
            </p>
          </div>
          <button
            onClick={onClose}
            className="p-1.5 hover:bg-neutral-100 text-neutral-500 hover:text-neutral-800 rounded-lg transition-colors"
          >
            <X className="w-5 h-5" />
          </button>
        </div>

        {/* Modal Content Split Grid */}
        <div className="flex-1 overflow-hidden grid grid-cols-1 lg:grid-cols-3">
          
          {/* Left panel: Active Step Details (1/3 width) */}
          <div className="lg:border-r border-neutral-150 p-6 flex flex-col gap-6 overflow-y-auto bg-neutral-50/50">
            <div>
              <span className="text-xs font-bold text-sky-700 uppercase tracking-wider block mb-1">
                Análisis del Paso Activo
              </span>
              <h3 className="text-2xl font-black text-neutral-800">
                Paso #{activeStep}
              </h3>
            </div>

            {/* Transition Formula Card */}
            {currentStep && (
              <div className="bg-white border border-neutral-200 p-4 rounded-xl shadow-sm space-y-4">
                <div className="text-xs font-bold text-neutral-400 uppercase tracking-wide">
                  Transición Formal
                </div>
                
                {/* Transition Delta Formulation */}
                <div className="font-mono text-xs overflow-x-auto bg-neutral-50 p-3 rounded-lg border border-neutral-150 leading-relaxed text-neutral-800">
                  <div>
                    δ(<span className="text-sky-700 font-bold">{currentStep.state}</span>, [
                    <span className="text-blue-500 font-bold">'{currentStep.read1 === '_' ? '•' : currentStep.read1}'</span>,{' '}
                    <span className="text-blue-500 font-bold">'{currentStep.read2 === '_' ? '•' : currentStep.read2}'</span>,{' '}
                    <span className="text-blue-500 font-bold">'{currentStep.read3 === '_' ? '•' : currentStep.read3}'</span>
                    ]) =
                  </div>
                  <div className="mt-2 pl-4 flex items-center gap-1.5">
                    <ArrowRight className="w-3.5 h-3.5 text-neutral-400" />
                    <span>
                      (<span className="text-sky-700 font-bold">{currentStep.nextState}</span>, [
                      <span className="text-emerald-600 font-semibold">'{currentStep.write1 === '_' ? '•' : currentStep.write1}'</span>,{' '}
                      <span className="text-emerald-600 font-semibold">'{currentStep.write2 === '_' ? '•' : currentStep.write2}'</span>,{' '}
                      <span className="text-emerald-600 font-semibold">'{currentStep.write3 === '_' ? '•' : currentStep.write3}'</span>
                      ], [
                      <span className="text-orange-500 font-bold">{currentStep.dir1}</span>,{' '}
                      <span className="text-orange-500 font-bold">{currentStep.dir2}</span>,{' '}
                      <span className="text-orange-500 font-bold">{currentStep.dir3}</span>
                      ])
                    </span>
                  </div>
                </div>

                {/* Readable description of current state action */}
                <div className="text-xs text-neutral-600 space-y-1.5 pt-2">
                  <div className="font-bold text-neutral-700">Acción del Estado:</div>
                  <p className="italic">
                    {currentStep.state === 'q_init' && 'Inicializa la Cinta 3 con el símbolo delimitador ($).'}
                    {currentStep.state === 'q_init_write1' && 'Escribe la primera marca unaria (1) en la Cinta 3.'}
                    {currentStep.state === 'q_cmp' && 'Comparando caracteres del historial (Cinta 1) con el FEN actual (Cinta 2).'}
                    {currentStep.state === 'q_rewindC2_skip' && 'Omite el delimitador de FEN y retrocede el cabezal de la Cinta 2.'}
                    {currentStep.state === 'q_rebobinarC2' && 'Rebobina la Cinta 2 hasta su delimitador inicial ($) para buscar en el siguiente historial.'}
                    {currentStep.state === 'q_saltarC1' && 'Avanza el cabezal del historial (Cinta 1) para prepararse a comparar el siguiente FEN.'}
                    {currentStep.state === 'q_rewindC3' && 'Búsqueda finalizada. Rebobinando la Cinta 3 para contar las coincidencias halladas.'}
                    {currentStep.state === 'q_countC3_1' && 'Cinta 3 leyó la coincidencia número 1. Avanzando a la siguiente.'}
                    {currentStep.state === 'q_countC3_2' && 'Cinta 3 leyó la coincidencia número 2. Avanzando a la siguiente.'}
                    {currentStep.state === 'q_countC3_3' && 'Cinta 3 leyó la coincidencia número 3. ¡Triple repetición detectada! Transición a ACEPTAR.'}
                    {currentStep.state === 'q_accept' && 'La máquina de Turing aceptó la entrada. Repetición confirmada.'}
                    {currentStep.state === 'q_reject' && 'La máquina de Turing se detuvo en estado de rechazo. No se halló triple repetición.'}
                  </p>
                </div>
              </div>
            )}

            {/* Quick stats / overview */}
            <div className="space-y-3 mt-auto">
              <div className="text-xs font-bold text-neutral-400 uppercase tracking-wide">
                Estadísticas de Simulación
              </div>
              <div className="grid grid-cols-2 gap-2 text-xs">
                <div className="bg-neutral-100 p-2.5 rounded-lg">
                  <div className="text-neutral-500 font-semibold">Pasos Totales</div>
                  <div className="text-base font-bold text-neutral-800 font-mono">
                    {steps.length}
                  </div>
                </div>
                <div className="bg-neutral-100 p-2.5 rounded-lg">
                  <div className="text-neutral-500 font-semibold">Veredicto</div>
                  <div className={`text-base font-bold font-mono ${
                    steps[steps.length - 1]?.nextState === 'q_accept'
                      ? 'text-emerald-600'
                      : 'text-red-500'
                  }`}>
                    {steps[steps.length - 1]?.nextState === 'q_accept' ? 'ACEPTAR' : 'RECHAZAR'}
                  </div>
                </div>
              </div>
            </div>
          </div>

          {/* Right panel: Full Steps Table (2/3 width) */}
          <div ref={tableContainerRef} className="lg:col-span-2 overflow-y-auto relative border-t lg:border-t-0 border-neutral-150">
            <table className="w-full text-left border-collapse">
              <thead className="sticky top-0 bg-neutral-50 text-neutral-500 text-xs font-bold border-b border-neutral-200 z-30">
                <tr>
                  <th className="py-3 px-4 font-mono w-16">Paso</th>
                  <th className="py-3 px-4 w-32">Estado</th>
                  <th className="py-3 px-4 w-28">Leer [1,2,3]</th>
                  <th className="py-3 px-4 w-28">Escribir [1,2,3]</th>
                  <th className="py-3 px-4 w-24">Dir [1,2,3]</th>
                  <th className="py-3 px-4">Sig. Estado</th>
                </tr>
              </thead>
              <tbody className="divide-y divide-neutral-100 text-sm font-mono">
                {steps.map((step, idx) => {
                  const isCurrent = idx === activeStep;
                  const isAccept = step.nextState === 'q_accept';
                  const isReject = step.nextState === 'q_reject';

                  let rowClass = "hover:bg-neutral-50 cursor-pointer transition-colors ";
                  if (isCurrent) {
                    rowClass += "bg-sky-50 border-l-4 border-l-sky-500";
                  } else if (isAccept) {
                    rowClass += "bg-emerald-50";
                  } else if (isReject) {
                    rowClass += "bg-red-50";
                  }

                  return (
                    <tr
                      key={idx}
                      ref={isCurrent ? activeRowRef : null}
                      onClick={() => onSelectStep(idx)}
                      className={rowClass}
                    >
                      <td className="py-2.5 px-4 font-bold text-neutral-400 text-xs">
                        {idx}
                      </td>
                      <td className="py-2.5 px-4 font-semibold text-neutral-700">
                        {step.state}
                      </td>
                      <td className="py-2.5 px-4 text-blue-600">
                        [{step.read1 === '_' ? '•' : step.read1},{step.read2 === '_' ? '•' : step.read2},{step.read3 === '_' ? '•' : step.read3}]
                      </td>
                      <td className="py-2.5 px-4 text-emerald-600">
                        [{step.write1 === '_' ? '•' : step.write1},{step.write2 === '_' ? '•' : step.write2},{step.write3 === '_' ? '•' : step.write3}]
                      </td>
                      <td className="py-2.5 px-4 text-orange-500 font-bold">
                        [{step.dir1},{step.dir2},{step.dir3}]
                      </td>
                      <td className={`py-2.5 px-4 ${
                        isAccept ? 'text-emerald-600 font-bold' : 
                        isReject ? 'text-red-555 font-bold' : 
                        'text-neutral-500'
                      }`}>
                        {step.nextState}
                      </td>
                    </tr>
                  );
                })}
              </tbody>
            </table>
          </div>

        </div>
      </div>
    </div>
  );
};
