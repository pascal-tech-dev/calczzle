# Calczzle

Calculator assignment: React frontend + Go evaluation API.

## Quick start (Docker Compose)

Requires Docker. From the repo root:

```bash
make up
# or: docker compose up --build -d
```

Then open [http://localhost](http://localhost).

| Service  | URL |
|----------|-----|
| App UI   | http://localhost |
| API      | http://localhost:8080 |
| Health   | http://localhost:8080/health |

Stop the stack:

```bash
make down
```

The frontend nginx container proxies `/api/` to the `backend` service on the Compose network, so the browser only talks to port 80.

## Test coverage (Docker)

Runs tests in ephemeral containers (`docker run --rm`), writes reports on the host, then exits:

```bash
make cover              # backend + frontend
make cover-backend      # → backend/coverage.out, backend/coverage.html
make cover-frontend     # → frontend/coverage/
```

No local Go or Node install required for these targets.

## Development

Day-to-day work uses Dev Containers under [`.devcontainer/`](.devcontainer/).

- Backend details: [`backend/README.md`](backend/README.md)
- Frontend details: [`frontend/README.md`](frontend/README.md)
