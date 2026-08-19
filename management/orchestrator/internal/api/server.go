package api

import (
	"context"
	"encoding/json"
	"net/http"
	"time"

	"github.com/vikramiiitm/unborn/management/orchestrator/internal/config"
	"github.com/vikramiiitm/unborn/management/orchestrator/internal/orchestrator"
	"github.com/vikramiiitm/unborn/management/orchestrator/internal/persona"
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

	// Instances (bodies)
	s.mux.HandleFunc("GET /v1/instances", s.handleListInstances)
	s.mux.HandleFunc("POST /v1/instances", s.handleCreateInstance)
	s.mux.HandleFunc("GET /v1/instances/{id}", s.handleGetInstance)
	s.mux.HandleFunc("POST /v1/instances/{id}/stop", s.handleStopInstance)

	// Personas (souls)
	s.mux.HandleFunc("GET /v1/personas", s.handleListPersonas)
	s.mux.HandleFunc("POST /v1/personas", s.handleCreatePersona)
	s.mux.HandleFunc("GET /v1/personas/{id}", s.handleGetPersona)

	// Identity
	s.mux.HandleFunc("GET /v1/device-profiles", s.handleListDeviceProfiles)
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
	writeJSON(w, http.StatusOK, s.orch.ListInstances())
}

func (s *Server) handleCreateInstance(w http.ResponseWriter, r *http.Request) {
	var req struct {
		PersonaID  string `json:"persona_id"`
		Simulated  bool   `json:"simulated"`
	}
	if err := json.NewDecoder(r.Body).Decode(&req); err != nil || req.PersonaID == "" {
		http.Error(w, "persona_id is required", http.StatusBadRequest)
		return
	}

	// Default to simulated in early Phase 1 for easier development
	simulated := true
	if r.URL.Query().Get("real") == "true" {
		simulated = false
	}
	if req.Simulated {
		simulated = true
	}

	inst, err := s.orch.CreateInstance(req.PersonaID, simulated)
	if err != nil {
		switch err {
		case orchestrator.ErrMaxInstancesReached:
			http.Error(w, err.Error(), http.StatusTooManyRequests)
		case orchestrator.ErrPersonaNotFound:
			http.Error(w, err.Error(), http.StatusBadRequest)
		default:
			http.Error(w, err.Error(), http.StatusInternalServerError)
		}
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

func (s *Server) handleListPersonas(w http.ResponseWriter, r *http.Request) {
	writeJSON(w, http.StatusOK, s.orch.PersonaStore().List())
}

func (s *Server) handleCreatePersona(w http.ResponseWriter, r *http.Request) {
	var req struct {
		DisplayName string `json:"display_name"`
		Location    string `json:"location"`
		Timezone    string `json:"timezone"`
		AgeMin      int    `json:"age_min"`
		AgeMax      int    `json:"age_max"`
		Engagement  string `json:"engagement"`
	}
	if err := json.NewDecoder(r.Body).Decode(&req); err != nil {
		http.Error(w, "invalid json", http.StatusBadRequest)
		return
	}

	if req.Location == "" {
		req.Location = "Unknown"
	}
	if req.Timezone == "" {
		req.Timezone = "UTC"
	}
	if req.AgeMin == 0 {
		req.AgeMin = 25
	}
	if req.AgeMax == 0 {
		req.AgeMax = 30
	}

	eng := persona.EngagementThoughtfulCommenter
	switch req.Engagement {
	case "lurker":
		eng = persona.EngagementLurker
	case "enthusiastic_sharer":
		eng = persona.EngagementEnthusiasticSharer
	case "quiet_reader":
		eng = persona.EngagementQuietReader
	case "selective_engager":
		eng = persona.EngagementSelectiveEngager
	}

	p := persona.New(req.DisplayName, req.Location, req.Timezone, req.AgeMin, req.AgeMax, eng)
	s.orch.PersonaStore().Create(p)
	writeJSON(w, http.StatusCreated, p)
}

func (s *Server) handleGetPersona(w http.ResponseWriter, r *http.Request) {
	id := r.PathValue("id")
	p, ok := s.orch.PersonaStore().Get(id)
	if !ok {
		http.Error(w, "persona not found", http.StatusNotFound)
		return
	}
	writeJSON(w, http.StatusOK, p)
}

func (s *Server) handleListDeviceProfiles(w http.ResponseWriter, r *http.Request) {
	writeJSON(w, http.StatusOK, s.orch.ListDeviceProfiles())
}

func writeJSON(w http.ResponseWriter, status int, v any) {
	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(status)
	_ = json.NewEncoder(w).Encode(v)
}
