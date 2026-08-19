package api

import (
	"context"
	"encoding/json"
	"net/http"
	"time"

	"github.com/vikramiiitm/unborn/management/orchestrator/internal/config"
	"github.com/vikramiiitm/unborn/management/orchestrator/internal/orchestrator"
)

type Server struct {
	orch *orchestrator.Orchestrator
	cfg  *config.Config
	mux  *http.ServeMux
	http *http.Server
}

func NewServer(orch *orchestrator.Orchestrator, cfg *config.Config) *Server {
	s := &Server{
		orch: orch,
		cfg:  cfg,
		mux:  http.NewServeMux(),
	}

	s.routes()
	return s
}

func (s *Server) routes() {
	s.mux.HandleFunc("GET /health", s.handleHealth)
	s.mux.HandleFunc("GET /v1/instances", s.handleListInstances)
	s.mux.HandleFunc("POST /v1/instances", s.handleCreateInstance)
	s.mux.HandleFunc("GET /v1/instances/{id}", s.handleGetInstance)
	s.mux.HandleFunc("POST /v1/instances/{id}/stop", s.handleStopInstance)
}

func (s *Server) ListenAndServe(addr string) error {
	s.http = &http.Server{
		Addr:         addr,
		Handler:      s.mux,
		ReadTimeout:  10 * time.Second,
		WriteTimeout: 10 * time.Second,
		IdleTimeout:  60 * time.Second,
	}
	return s.http.ListenAndServe()
}

func (s *Server) Shutdown(ctx context.Context) error {
	if s.http != nil {
		return s.http.Shutdown(ctx)
	}
	return nil
}

func (s *Server) handleHealth(w http.ResponseWriter, r *http.Request) {
	writeJSON(w, http.StatusOK, map[string]string{
		"status":  "ok",
		"service": "unborn-orchestrator",
	})
}

func (s *Server) handleListInstances(w http.ResponseWriter, r *http.Request) {
	instances := s.orch.ListInstances()
	writeJSON(w, http.StatusOK, instances)
}

func (s *Server) handleCreateInstance(w http.ResponseWriter, r *http.Request) {
	var req struct {
		PersonaID string `json:"persona_id"`
	}
	if err := json.NewDecoder(r.Body).Decode(&req); err != nil || req.PersonaID == "" {
		http.Error(w, "persona_id is required", http.StatusBadRequest)
		return
	}

	inst, err := s.orch.CreateInstance(req.PersonaID)
	if err != nil {
		if err == orchestrator.ErrMaxInstancesReached {
			http.Error(w, err.Error(), http.StatusTooManyRequests)
			return
		}
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}

	writeJSON(w, http.StatusCreated, inst)
}

func (s *Server) handleGetInstance(w http.ResponseWriter, r *http.Request) {
	id := r.PathValue("id")
	inst, ok := s.orch.GetInstance(id)
	if !ok {
		http.Error(w, "instance not found", http.StatusNotFound)
		return
	}
	writeJSON(w, http.StatusOK, inst)
}

func (s *Server) handleStopInstance(w http.ResponseWriter, r *http.Request) {
	id := r.PathValue("id")
	if err := s.orch.StopInstance(id); err != nil {
		if err == orchestrator.ErrInstanceNotFound {
			http.Error(w, err.Error(), http.StatusNotFound)
			return
		}
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}
	writeJSON(w, http.StatusOK, map[string]string{"status": "stopped"})
}

func writeJSON(w http.ResponseWriter, status int, v any) {
	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(status)
	_ = json.NewEncoder(w).Encode(v)
}
