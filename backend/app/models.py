from datetime import datetime

from sqlalchemy import DateTime, String
from sqlalchemy.orm import Mapped, mapped_column

from app.database import Base


class JobRun(Base):
    __tablename__ = "job_runs"

    id: Mapped[str] = mapped_column(String(36), primary_key=True)
    name: Mapped[str] = mapped_column(String(255), nullable=False, index=True)
    started_at: Mapped[datetime] = mapped_column(DateTime, nullable=False)
    stopped_at: Mapped[datetime | None] = mapped_column(DateTime, nullable=True, default=None)

    @property
    def is_running(self) -> bool:
        return self.stopped_at is None
