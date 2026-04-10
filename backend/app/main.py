import os

from robyn import Response, Robyn

from app.auth import CookieGetter, SessionAuthHandler
from app.config import PORT
from app.database import init_db
from app.routes.api import register_api_routes
from app.routes.auth import register_auth_routes
from app.routes.web import register_web_routes

app = Robyn(__file__)
app.configure_authentication(SessionAuthHandler(token_getter=CookieGetter()))

STATIC_DIR = os.path.join(os.path.dirname(__file__), "..", "static")

if os.path.isdir(STATIC_DIR):
    app.serve_directory(
        route="/",
        directory_path=os.path.abspath(STATIC_DIR),
        index_file="index.html",
    )


@app.startup_handler
async def startup():
    init_db()


register_api_routes(app)
register_auth_routes(app)
register_web_routes(app)


@app.get("/health", const=True)
async def health():
    return Response(status_code=200, headers={"content-type": "text/plain"}, description="ok")


if __name__ == "__main__":
    app.start(port=PORT, host="0.0.0.0")
