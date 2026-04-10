import json
import uuid
from datetime import datetime, timezone

from robyn import Request, Robyn

from app.auth import require_api_key
from app.database import get_session
from app.models import JobRun


def register_api_routes(app: Robyn):
    @app.post("/api/punch")
    def punch(request: Request):
        auth_error = require_api_key(request)
        if auth_error:
            return auth_error

        try:
            body = json.loads(request.body)
        except (ValueError, TypeError):
            return ({"error": "Invalid JSON body"}, {}, 400)

        run_id = (body.get("run_id") or "").strip()

        with get_session() as session:
            run = session.get(JobRun, run_id) if run_id else None

            if run is not None:
                # Punch out
                if run.stopped_at is not None:
                    return ({"error": f"run_id '{run_id}' is already stopped"}, {}, 409)
                run.stopped_at = datetime.now(timezone.utc).replace(tzinfo=None)
                session.commit()
                return {
                    "run_id": run.id,
                    "name": run.name,
                    "started_at": run.started_at.isoformat() + "Z",
                    "stopped_at": run.stopped_at.isoformat() + "Z",
                    "duration_seconds": int((run.stopped_at - run.started_at).total_seconds()),
                }

            # Punch in
            name = body.get("name", "").strip()
            if not name:
                return ({"error": "'name' is required"}, {}, 400)
            run = JobRun(
                id=run_id or str(uuid.uuid4()),
                name=name,
                started_at=datetime.now(timezone.utc).replace(tzinfo=None),
            )
            session.add(run)
            session.commit()
            return ({"run_id": run.id, "name": run.name, "started_at": run.started_at.isoformat() + "Z"}, {}, 201)
