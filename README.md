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

### POST /api/feedback

Submit feedback. Requires `X-API-Key` header.

**curl:**

```bash
curl -X POST http://localhost:8080/api/report \
  -H "Content-Type: application/json" \
  -H "X-API-Key: your-api-key" \
  -d '{ }'
```

**Python (httpx):**

```python
import httpx

response = httpx.post(
    "http://localhost:8080/api/report",
    headers={"X-API-Key": "your-api-key"},
    json={
    },
)
print(response.json())
```
