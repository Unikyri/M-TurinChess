# ♟️ M-TurinChess

**M-TurinChess** is an advanced chess analysis platform that combines the tactical precision of **Stockfish** with the psychological profiling of **LLMs (Large Language Models)**. 

The system doesn't just look for the best moves; it evaluates the "humanity" of the player by analyzing suspicious moves through a hybrid pipeline involving a **Turing Machine simulator** and **NVIDIA NIM (Llama 3.3 70B)**.

## 🚀 Key Features

- **Hybrid Analysis Pipeline**: Integrates Stockfish evaluations with LLM-powered psychological peritaje.
- **AI-Powered Peritaje**: Uses Llama 3.3 (via NVIDIA NIM) to decide if a "perfect" move is naturally human, opening theory, or suspicious modulo assistance.
- **Batch Processing**: Optimized analysis that processes suspicious moves in parallel batches to minimize latency and API costs.
- **Turing Tape Logic**: A symbolic representation of the game's integrity, processed by a custom Turing Machine model.
- **Modern UI**: A sleek, dark-mode React interface for uploading PGNs and visualizing analysis results in real-time.

## 🛠️ Tech Stack

- **Backend**: Go (Golang)
- **Frontend**: React (TypeScript) + Tailwind CSS
- **Chess Engine**: Stockfish 16.1
- **LLM Integration**: NVIDIA NIM (Llama 3.3 70B Instruct)
- **Architecture**: Domain-Driven Design (DDD) with a custom Lexer and Turing Machine simulation.

## ⚙️ Setup & Installation

### Backend
1. Navigate to `backend/`
2. Create a `.env` file with:
   ```env
   LLM_PROVIDER=nvidia
   LLM_API_KEY=your_nvidia_api_key
   LLM_MODEL=meta/llama-3.3-70b-instruct
   ```
3. Ensure Stockfish is available in the `../stockfish/` directory.
4. Run: `go run ./cmd/server/`

### Frontend
1. Navigate to `frontend/`
2. Run: `npm install`
3. Run: `npm run dev`

## 🧠 How it Works

1. **Stockfish Analysis**: The game is analyzed at Depth 18 to identify "perfect" moves (CP Loss 0-10).
2. **Lexer Classification**: Moves are classified into Symbols (M: Module, H: Human, E: Standard).
3. **LLM Verification**: The LLM reviews all "M" moves with the full context of the game and can downgrade them to "H" or "E" if it finds a human explanation.
4. **Turing Machine**: The final tape is processed to determine the definitive verdict of the player's integrity.

---
*Developed for the Final Project - Languages 2026-1*
