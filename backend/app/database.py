from sqlalchemy import create_engine, event
from sqlalchemy.orm import DeclarativeBase, Session

from app.config import DB_PATH

engine = create_engine(f"sqlite:///{DB_PATH}", echo=False)


@event.listens_for(engine, "connect")
def set_sqlite_pragma(dbapi_connection, connection_record):
    cursor = dbapi_connection.cursor()
    cursor.execute("PRAGMA journal_mode=WAL")
    cursor.execute("PRAGMA foreign_keys=ON")
    cursor.close()


class Base(DeclarativeBase):
    pass


def init_db():
    import app.models  # noqa: F401 — ensure models are registered with Base.metadata

    Base.metadata.create_all(engine)


def get_session() -> Session:
    return Session(engine)
