package api

import (
	"encoding/json"
	"net/http"

	"github.com/go-chi/chi/v5"
	"github.com/go-chi/chi/v5/middleware"
	"github.com/tcooper/pg-playground/app/db"
	"github.com/tcooper/pg-playground/app/replication"
	"github.com/tcooper/pg-playground/app/simulator"
)

func NewRouter(conns *db.Connections, m *simulator.Metrics) http.Handler {
	r := chi.NewRouter()
	r.Use(middleware.Logger)
	r.Use(middleware.Recoverer)

	r.Get("/health", handleHealth)
	r.Get("/stats", handleAllStats(conns, m))
	r.Get("/stats/lag", handleLag(conns))
	r.Get("/stats/tables", handleTables(conns))
	r.Get("/stats/slots", handleSlots(conns))
	r.Get("/stats/workers", handleWorkers(conns))
	r.Get("/stats/simulator", handleSimulator(m))

	return r
}

func handleHealth(w http.ResponseWriter, r *http.Request) {
	writeJSON(w, http.StatusOK, map[string]string{"status": "ok"})
}

func handleLag(conns *db.Connections) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		data, err := replication.GetLag(r.Context(), conns.Primary)
		if err != nil {
			writeError(w, err)
			return
		}
		writeJSON(w, http.StatusOK, data)
	}
}

func handleTables(conns *db.Connections) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		data, err := replication.GetSubscriptionState(r.Context(), conns.Replica)
		if err != nil {
			writeError(w, err)
			return
		}
		writeJSON(w, http.StatusOK, data)
	}
}

func handleSlots(conns *db.Connections) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		data, err := replication.GetSlotHealth(r.Context(), conns.Primary)
		if err != nil {
			writeError(w, err)
			return
		}
		writeJSON(w, http.StatusOK, data)
	}
}

func handleWorkers(conns *db.Connections) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		data, err := replication.GetWorkerInfo(r.Context(), conns.Primary, conns.Replica)
		if err != nil {
			writeError(w, err)
			return
		}
		writeJSON(w, http.StatusOK, data)
	}
}

func handleSimulator(m *simulator.Metrics) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		writeJSON(w, http.StatusOK, m.Snapshot())
	}
}

func handleAllStats(conns *db.Connections, m *simulator.Metrics) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		ctx := r.Context()

		lag, err := replication.GetLag(ctx, conns.Primary)
		if err != nil {
			writeError(w, err)
			return
		}
		tables, err := replication.GetSubscriptionState(ctx, conns.Replica)
		if err != nil {
			writeError(w, err)
			return
		}
		slots, err := replication.GetSlotHealth(ctx, conns.Primary)
		if err != nil {
			writeError(w, err)
			return
		}
		workers, err := replication.GetWorkerInfo(ctx, conns.Primary, conns.Replica)
		if err != nil {
			writeError(w, err)
			return
		}

		writeJSON(w, http.StatusOK, map[string]any{
			"lag":       lag,
			"tables":    tables,
			"slots":     slots,
			"workers":   workers,
			"simulator": m.Snapshot(),
		})
	}
}

func writeJSON(w http.ResponseWriter, status int, v any) {
	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(status)
	json.NewEncoder(w).Encode(v)
}

func writeError(w http.ResponseWriter, err error) {
	writeJSON(w, http.StatusInternalServerError, map[string]string{"error": err.Error()})
}
