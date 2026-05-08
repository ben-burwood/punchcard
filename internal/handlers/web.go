package handlers

import (
	"net/http"
	"strconv"

	"github.com/ben-burwood/punchcard/internal/store"
)

type runningItem struct {
	ID        string  `json:"id"`
	Name      string  `json:"name"`
	StartedAt string  `json:"started_at"`
	StoppedAt *string `json:"stopped_at"`
}

type historyItem struct {
	ID              string `json:"id"`
	Name            string `json:"name"`
	StartedAt       string `json:"started_at"`
	StoppedAt       string `json:"stopped_at"`
	DurationSeconds int    `json:"duration_seconds"`
}

type historyResponse struct {
	Total int           `json:"total"`
	Items []historyItem `json:"items"`
}

func Running(runs *store.JobRuns) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		list, err := runs.ListRunning(r.Context())
		if err != nil {
			writeError(w, http.StatusInternalServerError, "database error")
			return
		}
		out := make([]runningItem, 0, len(list))
		for _, run := range list {
			item := runningItem{
				ID:        run.ID,
				Name:      run.Name,
				StartedAt: formatNoZ(run.StartedAt),
			}
			if run.StoppedAt != nil {
				s := formatNoZ(*run.StoppedAt)
				item.StoppedAt = &s
			}
			out = append(out, item)
		}
		writeJSON(w, http.StatusOK, out)
	})
}

func History(runs *store.JobRuns) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		q := r.URL.Query()
		search := q.Get("search")
		sort := q.Get("sort")
		if sort == "" {
			sort = "stopped_at"
		}
		order := q.Get("order")
		if order == "" {
			order = "desc"
		}

		limit := 25
		if v := q.Get("limit"); v != "" {
			if n, err := strconv.Atoi(v); err == nil {
				limit = n
			}
		}
		if limit > 100 {
			limit = 100
		}
		if limit < 0 {
			limit = 0
		}

		offset := 0
		if v := q.Get("offset"); v != "" {
			if n, err := strconv.Atoi(v); err == nil {
				offset = n
			}
		}
		if offset < 0 {
			offset = 0
		}

		total, list, err := runs.ListHistory(r.Context(), search, sort, order, limit, offset)
		if err != nil {
			writeError(w, http.StatusInternalServerError, "database error")
			return
		}

		items := make([]historyItem, 0, len(list))
		for _, run := range list {
			if run.StoppedAt == nil {
				continue
			}
			items = append(items, historyItem{
				ID:              run.ID,
				Name:            run.Name,
				StartedAt:       formatNoZ(run.StartedAt),
				StoppedAt:       formatNoZ(*run.StoppedAt),
				DurationSeconds: int(run.StoppedAt.Sub(run.StartedAt).Seconds()),
			})
		}
		writeJSON(w, http.StatusOK, historyResponse{Total: total, Items: items})
	})
}
