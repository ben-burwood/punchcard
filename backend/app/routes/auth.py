import json

from robyn import Request, Response, Robyn

from app.auth import create_session, delete_session
from app.config import DASHBOARD_PASSWORD, DASHBOARD_USER, SESSION_COOKIE


def register_auth_routes(app: Robyn):
    @app.get("/web/auth/status", auth_required=True)
    async def status(request: Request) -> Response:
        return Response(status_code=204, headers={"content-type": "application/json"}, description="")

    @app.post("/web/auth/login")
    async def login(request: Request) -> Response:
        try:
            body = json.loads(request.body)
        except (json.JSONDecodeError, TypeError):
            return Response(
                status_code=400,
                headers={"content-type": "application/json"},
                description='{"error": "Invalid JSON body"}',
            )

        username = body.get("username", "")
        password = body.get("password", "")

        if username != DASHBOARD_USER or password != DASHBOARD_PASSWORD:
            return Response(
                status_code=401,
                headers={"content-type": "application/json"},
                description='{"error": "Invalid credentials"}',
            )

        token = create_session()
        return Response(
            status_code=200,
            headers={
                "content-type": "application/json",
                "set-cookie": f"{SESSION_COOKIE}={token}; Path=/; HttpOnly; SameSite=Strict; Max-Age=86400",
            },
            description='{"ok": true}',
        )

    @app.post("/web/auth/logout")
    async def logout(request: Request) -> Response:
        cookie_header = request.headers.get("cookie") or ""
        for part in cookie_header.split(";"):
            part = part.strip()
            if part.startswith(f"{SESSION_COOKIE}="):
                delete_session(part.split("=", 1)[1])
                break

        return Response(
            status_code=200,
            headers={
                "content-type": "application/json",
                "set-cookie": f"{SESSION_COOKIE}=; Path=/; HttpOnly; SameSite=Strict; Max-Age=0",
            },
            description='{"ok": true}',
        )
