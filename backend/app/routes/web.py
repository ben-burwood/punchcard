from robyn import Request, Robyn
from sqlalchemy import asc

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
