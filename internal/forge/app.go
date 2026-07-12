// Package forge is a minimal HTTP service used to demo the publish+verify
// workflow for the Forge platform.
package forge

import (
	"encoding/json"
	"net/http"
	"runtime"
	"sync/atomic"
	"time"
)

// App is a tiny HTTP application: /version, /health, /echo.
type App struct {
	startedAt time.Time
	requests  atomic.Uint64
}

// New returns a new App.
func New() *App { return &App{startedAt: time.Now().UTC()} }

// Handler returns the HTTP mux for the app.
func (a *App) Handler() http.Handler {
	mux := http.NewServeMux()
	mux.HandleFunc("GET /version", a.version)
	mux.HandleFunc("GET /health", a.health)
	mux.HandleFunc("POST /echo", a.echo)
	return mux
}

type versionResponse struct {
	Name      string `json:"name"`
	Version   string `json:"version"`
	GoVersion string `json:"go_version"`
	StartedAt string `json:"started_at"`
	Uptime    string `json:"uptime"`
}

func (a *App) version(w http.ResponseWriter, r *http.Request) {
	a.requests.Add(1)
	writeJSON(w, http.StatusOK, versionResponse{
		Name:      "forge-test-app",
		Version:   "0.1.0",
		GoVersion: runtime.Version(),
		StartedAt: a.startedAt.Format(time.RFC3339),
		Uptime:    time.Since(a.startedAt).Round(time.Second).String(),
	})
}

type healthResponse struct {
	Status   string `json:"status"`
	Requests uint64 `json:"requests"`
}

func (a *App) health(w http.ResponseWriter, r *http.Request) {
	writeJSON(w, http.StatusOK, healthResponse{
		Status:   "ok",
		Requests: a.requests.Load(),
	})
}

func (a *App) echo(w http.ResponseWriter, r *http.Request) {
	a.requests.Add(1)
	var body map[string]any
	dec := json.NewDecoder(r.Body)
	if err := dec.Decode(&body); err != nil {
		http.Error(w, "invalid json", http.StatusBadRequest)
		return
	}
	writeJSON(w, http.StatusOK, map[string]any{
		"received": body,
		"at":       time.Now().UTC().Format(time.RFC3339),
	})
}

func writeJSON(w http.ResponseWriter, status int, v any) {
	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(status)
	_ = json.NewEncoder(w).Encode(v)
}
