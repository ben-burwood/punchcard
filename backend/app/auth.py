import secrets
import time

from robyn import Request, Response
from robyn.authentication import AuthenticationHandler
from robyn.robyn import Identity

from app.config import API_KEY, SESSION_COOKIE

# In-memory session store: token -> expiry timestamp
_sessions: dict[str, float] = {}
SESSION_TTL = 86400  # 24 hours


def require_api_key(request: Request) -> Response | None:
    """Check X-API-Key Header. Returns an error Response if invalid, None if OK."""
    key = request.headers.get("x-api-key") or ""
    if key != API_KEY:
        return Response(
            status_code=401,
            headers={"content-type": "application/json"},
            description='{"error": "Invalid or missing API key"}',
        )
    return None


def create_session() -> str:
    """Create a new session token with 24h expiry"""
    now = time.time()
    expired = [t for t, exp in _sessions.items() if exp < now]
    for t in expired:
        del _sessions[t]

    token = secrets.token_hex(32)
    _sessions[token] = now + SESSION_TTL
    return token


def delete_session(token: str) -> None:
    """Remove a session token"""
    _sessions.pop(token, None)


class CookieGetter:
    """Extract Session Token from SESSION_COOKIE Cookie"""

    scheme = "Cookie"

    def get_token(self, request: Request) -> str | None:
        cookie_header = request.headers.get("cookie") or ""
        for part in cookie_header.split(";"):
            part = part.strip()
            if part.startswith(f"{SESSION_COOKIE}="):
                return part.split("=", 1)[1]
        return None


class SessionAuthHandler(AuthenticationHandler):
    """Validates Session Cookie Yokens for dashboard routes"""

    def authenticate(self, request: Request) -> Identity | None:
        token = self.token_getter.get_token(request)
        if not token or token not in _sessions:
            return None

        if _sessions[token] < time.time():
            del _sessions[token]
            return None

        return Identity(claims={"authenticated": "true"})
