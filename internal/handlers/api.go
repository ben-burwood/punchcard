package handlers

import (
	"database/sql"
	"encoding/json"
	"errors"
	"fmt"
	"net/http"
	"strings"
	"time"

	"github.com/ben-burwood/punchcard/internal/store"
)

type punchRequest struct {
	RunID *string `json:"run_id"`
	Name  *string `json:"name"`
}

func Punch(runs *store.JobRuns) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		var body punchRequest
		dec := json.NewDecoder(r.Body)
		if err := dec.Decode(&body); err != nil {
			writeError(w, http.StatusBadRequest, "Invalid JSON body")
			return
		}

		runID := ""
		if body.RunID != nil {
			runID = strings.TrimSpace(*body.RunID)
		}

		ctx := r.Context()

		if runID != "" {
			existing, err := runs.FindByID(ctx, runID)
			if err != nil {
				writeError(w, http.StatusInternalServerError, "database error")
				return
			}
			if existing != nil {
				if existing.StoppedAt != nil {
					writeError(w, http.StatusConflict, fmt.Sprintf("run_id '%s' is already stopped", runID))
					return
				}
				now := time.Now().UTC()
				if err := runs.Stop(ctx, runID, now); err != nil {
					if errors.Is(err, sql.ErrNoRows) {
						writeError(w, http.StatusConflict, fmt.Sprintf("run_id '%s' is already stopped", runID))
						return
					}
					writeError(w, http.StatusInternalServerError, "database error")
					return
				}
				writeJSON(w, http.StatusOK, map[string]any{
					"run_id":           existing.ID,
					"name":             existing.Name,
					"started_at":       formatWithZ(existing.StartedAt),
					"stopped_at":       formatWithZ(now),
					"duration_seconds": int(now.Sub(existing.StartedAt).Seconds()),
				})
				return
			}
		}

		name := ""
		if body.Name != nil {
			name = strings.TrimSpace(*body.Name)
		}
		if name == "" {
			writeError(w, http.StatusBadRequest, "'name' is required")
			return
		}

		id := runID
		if id == "" {
			id = newUUID()
		}
		started := time.Now().UTC()
		if err := runs.Create(ctx, store.JobRun{ID: id, Name: name, StartedAt: started}); err != nil {
			writeError(w, http.StatusInternalServerError, "database error")
			return
		}
		writeJSON(w, http.StatusCreated, map[string]any{
			"run_id":     id,
			"name":       name,
			"started_at": formatWithZ(started),
		})
	})
}
