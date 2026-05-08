package store

import (
	"context"
	"database/sql"
	"errors"
	"fmt"
	"time"
)

const tsLayout = "2006-01-02T15:04:05.000000Z"

type JobRun struct {
	ID        string
	Name      string
	StartedAt time.Time
	StoppedAt *time.Time
}

type JobRuns struct {
	db *sql.DB
}

func NewJobRuns(db *sql.DB) *JobRuns {
	return &JobRuns{db: db}
}

func formatTS(t time.Time) string {
	return t.UTC().Format(tsLayout)
}

func parseTS(s string) (time.Time, error) {
	if t, err := time.Parse(tsLayout, s); err == nil {
		return t, nil
	}
	return time.Parse(time.RFC3339Nano, s)
}

func (j *JobRuns) Create(ctx context.Context, run JobRun) error {
	_, err := j.db.ExecContext(ctx,
		`INSERT INTO job_runs (id, name, started_at, stopped_at) VALUES (?, ?, ?, NULL)`,
		run.ID, run.Name, formatTS(run.StartedAt),
	)
	return err
}

func (j *JobRuns) FindByID(ctx context.Context, id string) (*JobRun, error) {
	row := j.db.QueryRowContext(ctx,
		`SELECT id, name, started_at, stopped_at FROM job_runs WHERE id = ?`, id)
	return scanRow(row.Scan)
}

func (j *JobRuns) Stop(ctx context.Context, id string, at time.Time) error {
	res, err := j.db.ExecContext(ctx,
		`UPDATE job_runs SET stopped_at = ? WHERE id = ? AND stopped_at IS NULL`,
		formatTS(at), id,
	)
	if err != nil {
		return err
	}
	n, err := res.RowsAffected()
	if err != nil {
		return err
	}
	if n == 0 {
		return sql.ErrNoRows
	}
	return nil
}

func (j *JobRuns) ListRunning(ctx context.Context) ([]JobRun, error) {
	rows, err := j.db.QueryContext(ctx,
		`SELECT id, name, started_at, stopped_at FROM job_runs WHERE stopped_at IS NULL ORDER BY started_at ASC`)
	if err != nil {
		return nil, err
	}
	defer rows.Close()
	return collect(rows)
}

var validSort = map[string]string{
	"name":       "name",
	"started_at": "started_at",
	"stopped_at": "stopped_at",
	"duration":   "(julianday(stopped_at) - julianday(started_at))",
}

func (j *JobRuns) ListHistory(ctx context.Context, search, sort, order string, limit, offset int) (int, []JobRun, error) {
	sortCol, ok := validSort[sort]
	if !ok {
		sortCol = validSort["stopped_at"]
	}
	if order != "asc" && order != "desc" {
		order = "desc"
	}

	where := "stopped_at IS NOT NULL"
	args := []any{}
	if search != "" {
		where += " AND LOWER(name) LIKE LOWER(?)"
		args = append(args, "%"+search+"%")
	}

	var total int
	if err := j.db.QueryRowContext(ctx,
		"SELECT COUNT(*) FROM job_runs WHERE "+where, args...,
	).Scan(&total); err != nil {
		return 0, nil, err
	}

	q := fmt.Sprintf(
		"SELECT id, name, started_at, stopped_at FROM job_runs WHERE %s ORDER BY %s %s LIMIT ? OFFSET ?",
		where, sortCol, order,
	)
	args = append(args, limit, offset)
	rows, err := j.db.QueryContext(ctx, q, args...)
	if err != nil {
		return 0, nil, err
	}
	defer rows.Close()

	items, err := collect(rows)
	if err != nil {
		return 0, nil, err
	}
	return total, items, nil
}

func scanRow(scan func(...any) error) (*JobRun, error) {
	var run JobRun
	var startedStr string
	var stoppedStr sql.NullString
	if err := scan(&run.ID, &run.Name, &startedStr, &stoppedStr); err != nil {
		if errors.Is(err, sql.ErrNoRows) {
			return nil, nil
		}
		return nil, err
	}
	t, err := parseTS(startedStr)
	if err != nil {
		return nil, err
	}
	run.StartedAt = t
	if stoppedStr.Valid {
		st, err := parseTS(stoppedStr.String)
		if err != nil {
			return nil, err
		}
		run.StoppedAt = &st
	}
	return &run, nil
}

func collect(rows *sql.Rows) ([]JobRun, error) {
	out := []JobRun{}
	for rows.Next() {
		run, err := scanRow(rows.Scan)
		if err != nil {
			return nil, err
		}
		if run != nil {
			out = append(out, *run)
		}
	}
	return out, rows.Err()
}
