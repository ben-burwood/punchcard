# Punchcard

Job Tracker Dashboard via API.

## Stack

- **Backend**: Robyn + SQLAlchemy (SQLite)
- **Frontend**: Vue + Tailwind CSS + daisyUI

## Quick Start

```bash
docker compose up --build
```

Dashboard available at `http://localhost:8080`. Set your API key in `.env` (see `.env.example`).

## Local Development

```bash
# Backend
cd backend
PUNCHCARD_API_KEY=dev uv run python -m app.main

# Frontend (separate terminal)
cd frontend
npm install
npm run dev
```

The Vite dev server on `:5173` proxies `/api` and `/web` to the backend on `:8080`.

## Environment Variables

| Variable         | Required | Default      | Description          |
|------------------|----------|--------------|----------------------|
| `PUNCHCARD_API_KEY` | Yes      | —            | API key for POST     |
| `DB_PATH` | No       | `punchcard.db`  | SQLite database path |
| `PORT`    | No       | `8080`       | Server port          |

## API

All endpoints require an `X-API-Key` header.

---

### POST /api/punch

Punch endpoint — punches in (start) or out (stop) based on the run's current state.

- **`run_id` absent or not found** → punch in (start); `name` is required
- **`run_id` found and running** → punch out (stop)
- **`run_id` found and already stopped** → 409

**Request body:**

```json
{
  "name": "my-pipeline",
  "run_id": "optional-uuid"
}
```

| Field | Required | Description |
|-------|----------|-------------|
| `name` | On start | Human-readable job name, used to group multiple runs |
| `run_id` | No | UUID for this run; generated if not provided (always a start) |

**Response `201` (punched in):**

```json
{
  "run_id": "a1b2c3d4-...",
  "name": "my-pipeline",
  "started_at": "2026-04-10T12:00:00Z"
}
```

**Response `200` (punched out):**

```json
{
  "run_id": "a1b2c3d4-...",
  "name": "my-pipeline",
  "started_at": "2026-04-10T12:00:00Z",
  "stopped_at": "2026-04-10T12:05:30Z",
  "duration_seconds": 330
}
```

**curl — punch in:**

```bash
curl -X POST http://localhost:8080/api/punch \
  -H "Content-Type: application/json" \
  -H "X-API-Key: your-api-key" \
  -d '{"name": "my-pipeline"}'
```

**curl — punch out:**

```bash
curl -X POST http://localhost:8080/api/punch \
  -H "Content-Type: application/json" \
  -H "X-API-Key: your-api-key" \
  -d '{"run_id": "a1b2c3d4-..."}'
```

**Python (httpx):**

```python
import httpx

client = httpx.Client(
    base_url="http://localhost:8080",
    headers={"X-API-Key": "your-api-key"},
)

# Punch in
r = client.post("/api/punch", json={"name": "my-pipeline"})
run_id = r.json()["run_id"]

# Punch out
r = client.post("/api/punch", json={"run_id": run_id})
print(r.json()["duration_seconds"])
```

---

### Error responses

| Status | Meaning |
|--------|---------|
| 400 | Missing or invalid request body, or `name` absent on a start |
| 401 | Missing or invalid `X-API-Key` |
| 409 | `run_id` exists but is already stopped |
