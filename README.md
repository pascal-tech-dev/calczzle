# Calczzle

Calculator assignment: React frontend + Go evaluation API.

## Prerequisites

- [Docker](https://docs.docker.com/get-docker/)
- [Docker Compose](https://docs.docker.com/compose/install/) (Compose V2: `docker compose`)

Local Go or Node installs are optional. Running the app, tests, and coverage via the root `Makefile` only needs Docker.

## Quick start (Docker Compose)

From the repo root:

```bash
make up
# or: docker compose up --build -d
```

Then open [http://localhost:9091](http://localhost:9091).

| Environment | Frontend | Backend |
|-------------|----------|---------|
| Dev Containers | http://localhost:5173 | http://localhost:8080 |
| Production Compose | http://localhost:9091 | http://localhost:9090 |

Stop the stack:

```bash
make down
```

Nginx proxies `/api/` to `backend:8080` on the Compose network, so the UI at `:9091` can call the API without CORS.

## Tests & coverage (Docker)

Runs in ephemeral containers (`docker run --rm`). No local Go or Node install required.

```bash
make test               # backend + frontend → */test-results/
make test-backend       # → backend/test-results/report.txt + report.json
make test-frontend      # → frontend/test-results/report.txt + report.json

make cover              # backend + frontend with coverage reports
make cover-backend      # → backend/coverage.out, backend/coverage.html
make cover-frontend     # → frontend/coverage/
```

Test and coverage artifacts are gitignored.


## Development

Day-to-day work uses Dev Containers under [`.devcontainer/`](.devcontainer/).

- Backend details: [`backend/README.md`](backend/README.md)
- Frontend details: [`frontend/README.md`](frontend/README.md)
