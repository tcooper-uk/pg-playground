package api

import (
	"encoding/json"
	"net/http"

	"github.com/go-chi/chi/v5"
	"github.com/go-chi/chi/v5/middleware"
	"github.com/tcooper/pg-playground/app/config"
	"github.com/tcooper/pg-playground/app/db"
	"github.com/tcooper/pg-playground/app/replication"
	"github.com/tcooper/pg-playground/app/simulator"
)

func NewRouter(conns *db.Connections, m *simulator.Metrics, ws *simulator.WeightStore, rs *simulator.RateStore) http.Handler {
	r := chi.NewRouter()
	r.Use(middleware.Logger)
	r.Use(middleware.Recoverer)

	r.Get("/health", handleHealth)
	r.Get("/stats", handleAllStats(conns, m, ws, rs))
	r.Get("/stats/lag", handleLag(conns))
	r.Get("/stats/tables", handleTables(conns))
	r.Get("/stats/slots", handleSlots(conns))
	r.Get("/stats/workers", handleWorkers(conns))
	r.Get("/stats/simulator", handleSimulator(m))
	r.Get("/config/weights", handleGetWeights(ws))
	r.Put("/config/weights", handleSetWeights(ws))
	r.Get("/config/rate", handleGetRate(rs))
	r.Put("/config/rate", handleSetRate(rs))
	r.Post("/replication/pause", handleReplicationPause(conns))
	r.Post("/replication/resume", handleReplicationResume(conns))

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

func handleGetRate(rs *simulator.RateStore) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		writeJSON(w, http.StatusOK, map[string]int{"rate_per_second": rs.Get()})
	}
}

func handleSetRate(rs *simulator.RateStore) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		var body struct {
			RatePerSecond int `json:"rate_per_second"`
		}
		if err := json.NewDecoder(r.Body).Decode(&body); err != nil {
			writeJSON(w, http.StatusBadRequest, map[string]string{"error": err.Error()})
			return
		}
		rs.Set(body.RatePerSecond)
		writeJSON(w, http.StatusOK, map[string]int{"rate_per_second": rs.Get()})
	}
}

func handleReplicationPause(conns *db.Connections) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		if err := replication.PauseSubscription(r.Context(), conns.Replica, "dvdrental_sub"); err != nil {
			writeError(w, err)
			return
		}
		writeJSON(w, http.StatusOK, map[string]string{"status": "paused"})
	}
}

func handleReplicationResume(conns *db.Connections) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		if err := replication.ResumeSubscription(r.Context(), conns.Replica, "dvdrental_sub"); err != nil {
			writeError(w, err)
			return
		}
		writeJSON(w, http.StatusOK, map[string]string{"status": "active"})
	}
}

func handleGetWeights(ws *simulator.WeightStore) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		writeJSON(w, http.StatusOK, ws.Get())
	}
}

func handleSetWeights(ws *simulator.WeightStore) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		var wc config.WeightConfig
		if err := json.NewDecoder(r.Body).Decode(&wc); err != nil {
			writeJSON(w, http.StatusBadRequest, map[string]string{"error": err.Error()})
			return
		}
		ws.Set(wc)
		writeJSON(w, http.StatusOK, ws.Get())
	}
}

func handleAllStats(conns *db.Connections, m *simulator.Metrics, ws *simulator.WeightStore, rs *simulator.RateStore) http.HandlerFunc {
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
		subscriptions, err := replication.GetSubscriptionStatus(ctx, conns.Replica)
		if err != nil {
			writeError(w, err)
			return
		}

		writeJSON(w, http.StatusOK, map[string]any{
			"lag":           lag,
			"tables":        tables,
			"slots":         slots,
			"workers":       workers,
			"simulator":     m.Snapshot(),
			"weights":       ws.Get(),
			"rate":          rs.Get(),
			"subscriptions": subscriptions,
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
