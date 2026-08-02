# Calczzle

Calculator assignment: a React frontend that builds expressions and a Go backend that evaluates them over REST. The browser never evaluates math locally and never uses `eval()`.

## Prerequisites

- [Docker](https://docs.docker.com/get-docker/)
- [Docker Compose](https://docs.docker.com/compose/install/) (Compose V2: `docker compose`)

Local Go or Node installs are optional. Running the app, tests, and coverage via the root `Makefile` only needs Docker.

## Setup

1. Clone the repository.
2. Confirm Docker is running (`docker version`, `docker compose version`).
3. From the repo root, start the stack (builds images on first run):

```bash
make up
# or: docker compose up --build -d
```

4. Open the UI at [http://localhost:9091](http://localhost:9091).

Stop everything with `make down`.

| Environment | Frontend | Backend |
|-------------|----------|---------|
| Dev Containers | http://localhost:5173 | http://localhost:8080 |
| Production Compose | http://localhost:9091 | http://localhost:9090 |

Production host ports use `90xx` so they do not clash with Dev Containers (`5173` / `8080`).

## How to run the frontend and backend

### Recommended: Compose (both services)

```bash
make up
```

- UI: http://localhost:9091 (nginx serves the built React app and proxies `/api/` to the backend)
- API: http://localhost:9090

### Dev Containers (day-to-day development)

Open the backend and frontend Dev Containers under [`.devcontainer/`](.devcontainer/).

- Backend: `cd /workspace/backend && go run ./cmd/api` → http://localhost:8080
- Frontend: `cd /workspace/frontend && npm install && npm run dev -- --host 0.0.0.0` → http://localhost:5173  
  Vite proxies `/api` to `http://backend:8080` (see `frontend/.env.example`).

### Standalone containers

```bash
# backend
docker build -t calczzle-backend ./backend
docker run --rm -p 9090:8080 calczzle-backend

# frontend (needs the same Compose/Docker network as a service named backend)
docker build -t calczzle-frontend ./frontend
docker run --rm -p 9091:80 calczzle-frontend
```

More detail: [`backend/README.md`](backend/README.md), [`frontend/README.md`](frontend/README.md).

## API examples

### Health

```bash
curl -s http://localhost:9090/health
```

```json
{ "status": "ok" }
```

### Evaluate (success)

```bash
curl -s -X POST http://localhost:9090/api/v1/evaluate \
  -H 'Content-Type: application/json' \
  -d '{"expression":"3 + 4 - 7 * 2"}'
```

```json
{ "result": -7 }
```

Through the frontend nginx proxy (same path the UI uses):

```bash
curl -s -X POST http://localhost:9091/api/v1/evaluate \
  -H 'Content-Type: application/json' \
  -d '{"expression":"sqrt(16)+2"}'
```

```json
{ "result": 6 }
```

Other example expressions: `2 ^ 3 ^ 2`, `100 * 15%`, `2^-3`.

### Evaluate (error)

```bash
curl -s -X POST http://localhost:9090/api/v1/evaluate \
  -H 'Content-Type: application/json' \
  -d '{"expression":"1/0"}'
```

```json
{
  "error": {
    "code": "DIVISION_BY_ZERO",
    "message": "Division by zero is not allowed."
  }
}
```

Common error codes: `EMPTY_EXPRESSION`, `INVALID_EXPRESSION`, `UNSUPPORTED_FUNCTION`, `DIVISION_BY_ZERO`, `INVALID_SQUARE_ROOT`, `REQUEST_TOO_LARGE`. Full table in [`backend/README.md`](backend/README.md).

## Tests & coverage

```bash
make test               # → */test-results/
make test-backend
make test-frontend

make cover              # coverage reports
make cover-backend      # → backend/coverage.out, coverage.html
make cover-frontend     # → frontend/coverage/
```

Artifacts land under `*/test-results/` and coverage paths above. No local Go/Node required for these targets.

## Design decisions and assumptions

- **Backend is the source of truth for math.** The frontend only builds the expression string and displays results or errors. No client-side evaluation, no `eval` / `Function`.
- **Custom expression engine (not a library / not `eval`).** Evaluation runs entirely in `internal/calculator/expression` as an explicit pipeline:

  ```text
  expression string
        → Tokenizer
        → ToPostfix (Dijkstra shunting-yard: infix → RPN/postfix)
        → EvaluatePostfix (stack-based postfix evaluation)
        → numeric result
  ```

  Shunting-yard owns precedence, associativity, parentheses, unary minus, `%`, and `sqrt`. The postfix evaluator only walks RPN tokens with a value stack (`pkg.Stack` / `pkg.Queue`). That split keeps parsing rules testable without mixing them into HTTP or UI code.
- **Single evaluate endpoint.** `POST /api/v1/evaluate` with `{ "expression": "..." }` keeps the contract small and assignment-focused.
- **Same-origin `/api` in the browser.** Dev uses the Vite proxy; production uses nginx → `backend:8080`. Relative URLs avoid baking hostnames into the JS bundle and avoid CORS for the normal UI path.
- **HTTP stays outside the engine.** Controllers/`httpx` map typed domain errors to stable API codes. Clients never see postfix tokens, stack state, or parser internals.
- **Operator rules (assumptions).** `^` is right-associative and binds tighter than unary minus; postfix `%` means `/ 100` (so `100 * 15%` = `15`); only `sqrt` is supported as a function.
- **Separate networks for environments.** Dev Containers keep `5173`/`8080`; Compose publishes `9091`/`9090` so both stacks can run without port fights.
- **Docker-first DX for reviewers.** Clone + `make up` is enough to try the product; `make test` / `make cover` run in ephemeral containers.


## Development

Day-to-day work uses Dev Containers under [`.devcontainer/`](.devcontainer/).

- Backend: [`backend/README.md`](backend/README.md)
- Frontend: [`frontend/README.md`](frontend/README.md)
