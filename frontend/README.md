# Calczzle Frontend

React + Vite + TypeScript calculator UI for the Sezzle calculator assignment.

The frontend builds expressions and sends them to the Go backend. It never evaluates math locally and never uses `eval()`.

## Stack

- React 19 + TypeScript
- Vite
- Vitest + React Testing Library
- CSS Modules
- nginx (production)

## Development

```bash
cd /workspace/frontend
npm install
npm run dev -- --host 0.0.0.0
```

Open [http://localhost:5173](http://localhost:5173).

API calls go to `/api/v1/evaluate`. Vite proxies them to the backend:

```text
Browser → localhost:5173/api/... → Vite → backend:8080/api/...
```

Configure the proxy target with `.env` (see `.env.example`):

```dotenv
BACKEND_PROXY_TARGET=http://backend:8080
```

## Scripts

| Command | Purpose |
|---|---|
| `npm run dev` | Start Vite dev server |
| `npm run build` | Typecheck + production build |
| `npm run preview` | Serve the production build locally |
| `npm run lint` | ESLint |
| `npm run test` | Vitest watch mode |
| `npm run test:run` | Single test run |
| `npm run test:coverage` | Tests with coverage |

## Architecture

```text
App
 └─ Calculator
     └─ useCalculator
         ├─ calculatorReducer / expressionBuilder
         └─ calculatorApi → POST /api/v1/evaluate → Go backend
```

## Production Docker

Multi-stage build: Node builds the app, nginx serves `dist` and proxies `/api/` to `backend:8080`.

```bash
docker build -t calczzle-frontend .
docker run --rm -p 8081:80 calczzle-frontend
```

In Compose, attach the container to the same network as the `backend` service so `http://backend:8080` resolves.

## Manual checks

- Keypad and keyboard build expressions (`2^-3`, `sqrt(16)+2`, `100*15%`)
- `=` shows backend results and validation errors
- Loading disables duplicate evaluate requests
- `npm run lint`, `npm run test:run`, and `npm run build` pass
