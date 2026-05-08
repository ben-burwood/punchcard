package store

import (
	"database/sql"

	_ "modernc.org/sqlite"
)

const schema = `
CREATE TABLE IF NOT EXISTS job_runs (
    id          TEXT    PRIMARY KEY,
    name        TEXT    NOT NULL,
    started_at  TEXT    NOT NULL,
    stopped_at  TEXT
);
CREATE INDEX IF NOT EXISTS ix_job_runs_name ON job_runs(name);
`

func Open(dbPath string) (*sql.DB, error) {
	dsn := "file:" + dbPath + "?_pragma=journal_mode(WAL)&_pragma=foreign_keys(ON)"
	db, err := sql.Open("sqlite", dsn)
	if err != nil {
		return nil, err
	}
	db.SetMaxOpenConns(1)

	if _, err := db.Exec(schema); err != nil {
		db.Close()
		return nil, err
	}
	return db, nil
}
