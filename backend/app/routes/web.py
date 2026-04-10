import json

from robyn import Request, Robyn
from sqlalchemy import asc, desc, func

from app.database import get_session
from app.models import JobRun


def register_web_routes(app: Robyn):
    @app.get("/web/jobs/running", auth_required=True)
    def running_jobs(request: Request):
        with get_session() as session:
            runs = session.query(JobRun).filter(JobRun.stopped_at.is_(None)).order_by(asc(JobRun.started_at)).all()
            return [
                {
                    "id": run.id,
                    "name": run.name,
                    "started_at": run.started_at.isoformat(),
                    "stopped_at": run.stopped_at.isoformat() if run.stopped_at else None,
                }
                for run in runs
            ]

    @app.get("/web/jobs/history", auth_required=True)
    def job_history(request: Request):
        search = request.query_params.get("search", "").strip()
        sort = request.query_params.get("sort", "stopped_at")
        order = request.query_params.get("order", "desc")
        limit = min(int(request.query_params.get("limit", "25")), 100)
        offset = max(int(request.query_params.get("offset", "0")), 0)

        if sort not in ("name", "started_at", "stopped_at", "duration"):
            sort = "stopped_at"
        if order not in ("asc", "desc"):
            order = "desc"

        duration_expr = func.julianday(JobRun.stopped_at) - func.julianday(JobRun.started_at)

        sort_col = {
            "name": JobRun.name,
            "started_at": JobRun.started_at,
            "stopped_at": JobRun.stopped_at,
            "duration": duration_expr,
        }[sort]

        order_fn = asc if order == "asc" else desc

        with get_session() as session:
            base = session.query(JobRun).filter(JobRun.stopped_at.isnot(None))
            if search:
                base = base.filter(JobRun.name.ilike(f"%{search}%"))

            total = base.count()
            runs = base.order_by(order_fn(sort_col)).limit(limit).offset(offset).all()

            items = [
                {
                    "id": run.id,
                    "name": run.name,
                    "started_at": run.started_at.isoformat(),
                    "stopped_at": run.stopped_at.isoformat(),
                    "duration_seconds": int((run.stopped_at - run.started_at).total_seconds()),
                }
                for run in runs
            ]

        return {"total": total, "items": items}
