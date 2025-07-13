# MiniVault API

A minimal, production-grade REST API for prompt generation using a local LLM (Ollama with gemma:2b). Built with Domain-Driven Design (DDD) principles, clean architecture, and robust testing.

---

## ⚡ Quickstart

1. **Install Go 1.24+ and [Ollama](https://ollama.com/)**
2. **Pull the model:** `ollama pull gemma:2b`
3. **Run Ollama:** `ollama serve &`
4. **Start API:**
```bash
go run cmd/main.go
```
5. **Test:**
```bash
curl -X POST http://localhost:8080/generate \
  -H 'Content-Type: application/json' \
  -d '{"prompt": "What is ModelVault?"}'
```

---

## 🏛️ Architecture

**Layered, DDD-inspired structure:**

- **Client**  
  Sends HTTP requests to the API.

- **API Layer (`api/`)**  
  Parses incoming HTTP requests and validates input.  
  Delegates business logic to the Application Layer.  
  Formats and returns HTTP responses.

- **Server & Middleware (`server/`)**  
  Sets up the HTTP server, routes, and middleware (body limit, panic recovery, etc).

- **Application Layer (`usecases/`)**  
  Orchestrates business logic.  
  Implements domain interfaces (ports).  
  Calls into infrastructure for side effects.

- **Domain Layer (`domain/`)**  
  Defines core business entities, validation, and interfaces (ports) for all dependencies.

- **Infrastructure Layer (`infrastructure/`)**  
  Implements adapters for logging and LLM (Ollama) API.  
  Handles external communication and persistence.

- **Config (`config/`)**  
  Centralized configuration management. Loads environment variables and provides settings to the rest of the application.

- **Mocks (`mocks/`)**  
  Test doubles for all ports/interfaces, used in unit tests.

**Request Flow Example:**  
Client → API Handler → Usecase (Application) → Domain Validation → Infrastructure (Ollama, Logger) → Response

### DDD & Clean Architecture
- **Domain Layer**: Business entities, validation, and interfaces (`domain/`)
- **Use Cases**: Application logic, orchestrating domain and infrastructure (`usecases/`)
- **Infrastructure**: Adapters for logging and LLM API (`infrastructure/`)
- **API Layer**: HTTP handlers, request/response formatting (`api/`)
- **Server**: HTTP server and middleware (`server/`)
- **Config**: Centralized configuration management (`config/`)
- **Mocks**: Test doubles for all ports (`mocks/`)

---

## 📂 Code Structure

```
minivault/
├── api/                # HTTP handlers
├── cmd/                # Entry point (main.go)
├── config/             # Centralized configuration management
├── domain/             # Entities, validation, ports (interfaces)
├── infrastructure/     # Adapters: logging, LLM (Ollama)
├── mocks/              # Generated/test mocks
├── server/             # Server and middleware
├── usecases/           # Application/business logic
├── logs/               # Interaction logs (created at runtime)
├── go.mod, go.sum      # Go dependencies
└── README.md           # This file
```

---

## 🚀 API Documentation

### POST `/generate`
Generate a response for a given prompt using the local LLM.

- **URL:** `/generate`
- **Method:** `POST`
- **Content-Type:** `application/json`

#### Request Body
```json
{
  "prompt": "What is ModelVault?"
}
```

#### Response
```json
{
  "response": "..."
}
```

#### Error Responses
| Code | Description                | Example message         |
|------|----------------------------|------------------------|
| 400  | Invalid JSON / Validation  | "Invalid JSON" / "Validation error" |
| 405  | Method not allowed         | "Method not allowed"   |
| 500  | Internal error             | "Failed to generate response" |

- All responses include an `X-Request-ID` header for tracing.

#### Example Usage
```bash
curl -X POST http://localhost:8080/generate \
  -H 'Content-Type: application/json' \
  -d '{"prompt": "What is ModelVault?"}'
```

---

## ⚙️ Configuration

MiniVault uses environment variables (optionally loaded from a `.env` file) for configuration. The following settings are available:

| Variable         | Default                                 | Description                                                      |
|------------------|-----------------------------------------|------------------------------------------------------------------|
| MINIVAULT_PORT   | `:8080`                                 | The port/address the API server listens on                       |
| OLLAMA_URL       | `http://localhost:11434/api/chat`       | The URL for the Ollama chat API                                  |
| OLLAMA_MODEL     | `gemma:2b`                              | The Ollama model to use for generation (must be installed)       |

You can create a `.env` file in the project root to override these defaults:

---

## 🧪 Testing

- All business logic, HTTP handlers, and infrastructure are covered by unit tests.
- Mocks for all ports/interfaces (`mocks/` directory)
- Example: `go test ./...`
- Key tests:
  - Success and error cases for `/generate`
  - Input validation, HTTP codes, and edge cases
  - Logging and infrastructure behavior

---

## 📜 Logging

- **Generation interactions**: JSONL format, saved to `logs/log.jsonl`
- **Errors, warnings, info**: Console (with timestamps)
- Uses [zerolog](https://github.com/rs/zerolog) for structured logging

```
MINIVAULT_PORT=:8080
OLLAMA_URL=http://localhost:11434/api/chat
OLLAMA_MODEL=gemma:2b
```

> **Note:** The value of `OLLAMA_MODEL` must match a model that is installed in your local Ollama instance. For example, if you set `OLLAMA_MODEL=llama2:7b`, you must have run `ollama pull llama2:7b` beforehand.

---

## ℹ️ Notes
- Minimal, synchronous API for clarity. For streaming, see Ollama docs and Go's `http.Flusher`.
- All logic is local; no cloud LLMs are used.

---

## 🛠️ Improvements & TODOs
- [ ] Streaming responses (token-by-token)
- [ ] Make model/endpoint configurable via env vars
- [ ] Add CLI or Postman collection for easier testing
- [ ] Add more endpoints (health, status, etc.)
- [ ] Expand test coverage (integration, infra)
- [ ] Enhance error handling and observability


_Made for the ModelVault take-home project._
