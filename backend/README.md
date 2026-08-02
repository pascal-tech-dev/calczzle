# Calczzle Backend

Go microservice that evaluates arithmetic expressions over a small REST API.

## What it does

- Accepts expressions such as `3 + 4 - 7 * 2`, `2 ^ 3 ^ 2`, `sqrt(81)`, `100 * 15%`
- Returns a numeric result or a structured error
- Stays framework-light: Fiber for HTTP, standard library everywhere else

## Quick start

```bash
# local
make run
# or
go run ./cmd/api

# tests
make test
make cover
make check   # fmt + vet + test
```

Default listen address: `http://localhost:8080` (`PORT` env, default `8080`).

### Docker

Prefer the root Compose stack with the frontend:

```bash
# from repo root
make up
```

Standalone API image:

```bash
docker build -t calczzle-backend .
docker run --rm -p 8080:8080 calczzle-backend
```

Coverage without a local Go install (from repo root):

```bash
make cover-backend
```

## API

### Health

```http
GET /health
```

```json
{ "status": "ok" }
```

### Evaluate

```http
POST /api/v1/evaluate
Content-Type: application/json
```

```json
{ "expression": "3 + 4 - 7 * 2" }
```

Success:

```json
{ "result": -7 }
```

Error:

```json
{
  "error": {
    "code": "DIVISION_BY_ZERO",
    "message": "Division by zero is not allowed."
  }
}
```

### Error codes

| Code | HTTP | When |
|------|------|------|
| `INVALID_REQUEST` | 400 | Malformed JSON body |
| `EMPTY_EXPRESSION` | 400 | Missing / blank expression |
| `INVALID_EXPRESSION` | 400 | Syntax / parse / evaluate shape errors |
| `UNSUPPORTED_FUNCTION` | 400 | Unknown function name |
| `REQUEST_TOO_LARGE` | 413 | Body over 4 KiB limit |
| `DIVISION_BY_ZERO` | 422 | `/ 0` |
| `INVALID_SQUARE_ROOT` | 422 | `sqrt` of a negative number |
| `INTERNAL_SERVER_ERROR` | 500 | Unexpected failures |

## Architecture

Feature-first packages with a strict dependency direction:

```text
cmd/api
   ↓
internal/app          composition root (Fiber, middleware, wiring)
   ↓
controller            HTTP bind / transport validation / JSON
   ↓
service               use-case orchestration
   ↓
expression            Tokenizer → Parser → Evaluator
```

Shared helpers:

- `internal/platform/httpx` — central domain→HTTP error mapping
- `pkg` — generic `Stack` / `Queue` used by the expression engine

Calculator evaluation pipeline:

```text
"3 + 4 - 7 * 2"
        ↓ Tokenize
[3, +, 4, -, 7, *, 2]
        ↓ ToPostfix (shunting-yard)
[3, 4, +, 7, 2, *, -]
        ↓ EvaluatePostfix (stack)
-7
```

## Design decisions

### Why this layering?

Controllers must not know about precedence or arithmetic. The expression package must not know about Fiber or HTTP status codes. That keeps the engine unit-testable and lets HTTP concerns stay in one place (`httpx` + controllers).

### Why Fiber v3?

The guide chose Fiber for a small, fast REST service. Middleware is limited to what we need: recover, logger, a custom error handler, and a body limit. No DI framework, ORM, or auth stack — none of those are required for this API.

### Why shunting-yard + postfix?

Infix → postfix with the shunting-yard algorithm gives a clear split:

1. **Parser** owns precedence, associativity, parentheses, functions, unary minus, `%`
2. **Evaluator** only walks postfix tokens with a stack

That separation makes operator rules easy to test without mixing in float arithmetic.

### Why `pkg.Stack` / `pkg.Queue`?

Parser and evaluator both need LIFO (and the parser uses a FIFO output queue). Putting small generic collections in `pkg` avoids re-implementing slice push/pop logic in multiple places and keeps the algorithm code readable.

### Why sentinel domain errors?

Math failures (`ErrDivisionByZero`, `ErrInvalidSquareRoot`, …) are typed so `httpx` can map them with `errors.Is` to stable public codes. Clients never see stack traces or parser internals.

### Why clean floats in the controller?

IEEE-754 binary floats produce noisy JSON like `0.10500000000000001`. We round-trip through `strconv.FormatFloat(..., 'g', 15)` before responding so the API returns calculator-friendly values without pulling in a decimal library.

### Why a 4 KiB body limit?

Calculator expressions are tiny. Cap request bodies early (`BodyLimit`) so oversized payloads are rejected before tokenization. Default Fiber limit (4 MiB) is unnecessarily large for this service.

### Why no database / Redis / auth?

The product is a pure evaluation API. Persistence and identity are out of scope and would only add operational surface area.

### Supported expression surface

| Feature | Notes |
|---------|--------|
| `+ - * / ^` | `^` is right-associative and higher precedence than unary minus |
| Unary `+` / `-` | Unary plus is a no-op; unary minus becomes internal `u-` |
| `( )` | Grouping |
| `sqrt(...)` | Negative input → `INVALID_SQUARE_ROOT` |
| postfix `%` | `n%` → `n / 100` (so `100 * 15%` = `15`) |

## Project layout

```text
cmd/api/                 process entrypoint
internal/app/            Fiber app construction
internal/config/         env config (PORT)
internal/health/         GET /health
internal/calculator/     HTTP + service + expression engine
internal/platform/httpx/ API error envelope + mapping
pkg/                     reusable Stack / Queue
docs/backend-guide.md    implementation guide used during development
```

## Configuration

| Variable | Default | Meaning |
|----------|---------|---------|
| `PORT` | `8080` | HTTP listen port |

## Development notes

- Module path: `pascal-tech-dev/calczzle-backend`
- Go version: see `go.mod` (currently 1.26.x)
- Implementation detail and phase order live in [`docs/backend-guide.md`](docs/backend-guide.md)

## What’s intentionally out of scope (for now)

- CORS (add when a browser frontend is not behind a same-origin proxy)
- Authentication / rate limiting beyond body size
- Persisting calculation history
- Multi-argument functions beyond `sqrt`
