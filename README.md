# MiniVault API

A minimal local REST API that simulates a core ModelVault feature: receiving a prompt and returning a generated response using a local LLM (Ollama with gemma:2b).

## Features
- POST `/generate` endpoint
- Uses locally running Ollama (gemma:2b model)
- Logs all input/output to `logs/log.jsonl`
- Minimal dependencies (pure Go stdlib)

## Requirements
- Go 1.18+
- [Ollama](https://ollama.com/) installed and running locally
- gemma:2b model pulled (`ollama pull gemma:2b`)

## Setup

1. **Clone repo**

```bash
git clone <your_repo_url>
cd minivault
```

2. **Start Ollama**

```bash
ollama serve &
ollama pull gemma:2b
```

3. **Run the server**

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

All interactions are logged to `logs/log.jsonl` in JSONL format.

## Notes
- The API is minimal and synchronous for simplicity. For streaming/token-by-token output, see Ollama's streaming API and Go's `http.Flusher`.
- Only runs locally, no cloud LLMs are used.

## Improvements
- Add streaming responses for token-by-token output
- Add CLI or Postman collection for easier testing
- Add tests and error handling improvements

---

Made for the ModelVault take-home project.
