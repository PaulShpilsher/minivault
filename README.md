# MiniVault API

A minimal local REST API that simulates a core ModelVault feature: receiving a prompt and returning a generated response using a local LLM (Ollama with gemma:2b).

## Features
- POST `/generate` endpoint
- Uses locally running Ollama (for this POC, the gemma:2b model is chosen)
- AI Generation logs are stored in `logs/log.jsonl`, everything else is logged to console

## Requirements
- Go 1.24+
- [Ollama](https://ollama.com/) installed and running locally
- gemma:2b model pulled (`ollama pull gemma:2b`)  
  _Note: For the purpose of this proof of concept, the gemma:2b model is used. You can substitute another model if desired, but the `ollamaURL` and `ollamaModel` constants in `infrastructure/ollama.go` will need to be updated to reflect the change._
  
  TODO: make this configurable, e.g. via environment variables

## Setup

1. **Clone repo**

```bash
git clone https://github.com/PaulShpilsher/minivault.git
cd minivault
```

2. **Download and install Ollama**
- **Linux:**
```bash
curl -fsSL https://ollama.com/install.sh | sh
```
This script downloads the latest ollama binary and sets it up in /usr/local/bin.

3. **Start Ollama**
```bash
ollama serve &
ollama pull gemma:2b
```

4. **Run the server**

```bash
go run main.go
```

The API will be available at `http://localhost:8080`.

## Usage

POST to `/generate`:

```
curl -X POST http://localhost:8080/generate \
  -H 'Content-Type: application/json' \
  -d '{"prompt": "What is ModelVault?"}'
```

Response:
```
{"response": "..."}
```

## Logs

AI Generation interactions are logged to `logs/log.jsonl` in JSONL format.
Everything else is logged to console.

## Notes
- The API is minimal and synchronous for simplicity. For streaming/token-by-token output, see Ollama's streaming API and Go's `http.Flusher`.
- Only runs locally, no cloud LLMs are used.

## Improvements
- Add streaming responses for token-by-token output
- Add CLI or Postman collection for easier testing
- Add tests and error handling improvements

---

Made for the ModelVault take-home project.
