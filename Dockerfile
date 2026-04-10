# Stage 1: Build Frontend
FROM node:latest AS build-frontend

WORKDIR /app/frontend

COPY frontend/package.json frontend/package-lock.json* ./

RUN npm ci

COPY frontend/ ./

RUN npm run build

# Stage 2: Build Backend
FROM ghcr.io/astral-sh/uv:python3.14-bookworm-slim AS build-backend

WORKDIR /app

COPY backend/pyproject.toml backend/uv.lock* ./

RUN uv sync --no-dev --no-install-project

# Stage 3: Runtime
FROM python:3.14-slim-bookworm AS runtime

WORKDIR /app

COPY --from=build-backend /app/.venv /app/.venv
COPY backend/app/ ./app/
COPY --from=build-frontend /app/frontend/dist ./static/
COPY --from=ghcr.io/tarampampam/microcheck:1 /bin/httpcheck /bin/httpcheck

ENV PATH="/app/.venv/bin:$PATH"
ENV DB_PATH=/data/punchcard.db

EXPOSE 8080
VOLUME ["/data"]

HEALTHCHECK --interval=30s --timeout=5s --start-period=10s --retries=3 \
    CMD ["httpcheck", "http://localhost:8080/health"]

CMD ["python", "-m", "app.main"]
