package api

import (
	"context"
	"encoding/json"
	"net/http"
	"time"

	"github.com/vikramiiitm/unborn/management/orchestrator/internal/config"
	"github.com/vikramiiitm/unborn/management/orchestrator/internal/orchestrator"
	"github.com/vikramiiitm/unborn/management/orchestrator/internal/persona"
	"github.com/vikramiiitm/unborn/management/orchestrator/internal/vitality"
)

type Server struct {
	orch *orchestrator.Orchestrator
	cfg  *config.Config
	mux  *http.ServeMux
	http *http.Server
}

func NewServer(orch *orchestrator.Orchestrator, cfg *config.Config) *Server {
	s := &Server{orch: orch, cfg: cfg, mux: http.NewServeMux()}
	s.routes()
	return s
}

func (s *Server) routes() {
	s.mux.HandleFunc("GET /health", s.handleHealth)
	s.mux.HandleFunc("GET /{$}", s.handleDashboard)
	s.mux.HandleFunc("GET /dashboard", s.handleDashboard)
	s.mux.HandleFunc("GET /v1/instances", s.handleListInstances)
	s.mux.HandleFunc("POST /v1/instances", s.handleCreateInstance)
	s.mux.HandleFunc("GET /v1/instances/{id}", s.handleGetInstance)
	s.mux.HandleFunc("POST /v1/instances/{id}/stop", s.handleStopInstance)
	s.mux.HandleFunc("GET /v1/personas", s.handleListPersonas)
	s.mux.HandleFunc("POST /v1/personas", s.handleCreatePersona)
	s.mux.HandleFunc("GET /v1/personas/{id}", s.handleGetPersona)
	s.mux.HandleFunc("GET /v1/device-profiles", s.handleListDeviceProfiles)
	s.mux.HandleFunc("GET /v1/personas/{id}/next-action", s.handleNextAction)
	s.mux.HandleFunc("GET /v1/personas/{id}/vitality", s.handleGetVitality)
	s.mux.HandleFunc("GET /v1/vitality", s.handleListVitality)
	s.mux.HandleFunc("GET /v1/playbooks", s.handleListPlaybooks)
	s.mux.HandleFunc("POST /v1/playbooks", s.handleCreatePlaybook)
	s.mux.HandleFunc("POST /v1/playbooks/{id}/assign", s.handleAssignPlaybook)
	s.mux.HandleFunc("GET /v1/playbook-assignments", s.handleListAssignments)
	s.mux.HandleFunc("GET /v1/proxies", s.handleListProxies)
	s.mux.HandleFunc("PUT /v1/personas/{id}/proxy", s.handleSetProxy)
	s.mux.HandleFunc("GET /v1/personas/{id}/proxy", s.handleGetProxy)
	s.mux.HandleFunc("DELETE /v1/personas/{id}/proxy", s.handleDeleteProxy)
}

func (s *Server) ListenAndServe(addr string) error {
	s.http = &http.Server{
		Addr: addr, Handler: s.mux,
		ReadTimeout: 15 * time.Second, WriteTimeout: 15 * time.Second, IdleTimeout: 60 * time.Second,
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
	writeJSON(w, http.StatusOK, map[string]string{"status": "ok", "service": "unborn-orchestrator"})
}

func (s *Server) handleListInstances(w http.ResponseWriter, r *http.Request) {
	writeJSON(w, http.StatusOK, s.orch.ListInstances())
}

func (s *Server) handleCreateInstance(w http.ResponseWriter, r *http.Request) {
	var req struct {
		PersonaID string `json:"persona_id"`
		Simulated bool   `json:"simulated"`
	}
	if err := json.NewDecoder(r.Body).Decode(&req); err != nil || req.PersonaID == "" {
		http.Error(w, "persona_id is required", http.StatusBadRequest)
		return
	}
	simulated := true
	if r.URL.Query().Get("real") == "true" {
		simulated = false
	}
	inst, err := s.orch.CreateInstance(r.Context(), req.PersonaID, simulated || req.Simulated)
	if err != nil {
		if err == orchestrator.ErrPersonaNotFound {
			http.Error(w, err.Error(), http.StatusBadRequest)
			return
		}
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}
	writeJSON(w, http.StatusCreated, inst)
}

func (s *Server) handleGetInstance(w http.ResponseWriter, r *http.Request) {
	inst, ok := s.orch.GetInstance(r.PathValue("id"))
	if !ok {
		http.Error(w, "instance not found", http.StatusNotFound)
		return
	}
	writeJSON(w, http.StatusOK, inst)
}

func (s *Server) handleStopInstance(w http.ResponseWriter, r *http.Request) {
	if err := s.orch.StopInstance(r.Context(), r.PathValue("id")); err != nil {
		http.Error(w, err.Error(), http.StatusNotFound)
		return
	}
	writeJSON(w, http.StatusOK, map[string]string{"status": "stopped"})
}

func (s *Server) handleListPersonas(w http.ResponseWriter, r *http.Request) {
	list, err := s.orch.Personas().List(r.Context())
	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}
	writeJSON(w, http.StatusOK, list)
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
	created, err := s.orch.Personas().Create(r.Context(), p)
	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}
	s.orch.Vitality().Ensure(created.ID)
	writeJSON(w, http.StatusCreated, created)
}

func (s *Server) handleGetPersona(w http.ResponseWriter, r *http.Request) {
	p, err := s.orch.Personas().Get(r.Context(), r.PathValue("id"))
	if err != nil {
		if err == persona.ErrNotFound {
			http.Error(w, "persona not found", http.StatusNotFound)
			return
		}
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}
	writeJSON(w, http.StatusOK, p)
}

func (s *Server) handleListDeviceProfiles(w http.ResponseWriter, r *http.Request) {
	writeJSON(w, http.StatusOK, s.orch.ListDeviceProfiles())
}

func (s *Server) handleNextAction(w http.ResponseWriter, r *http.Request) {
	act, err := s.orch.NextAction(r.Context(), r.PathValue("id"))
	if err != nil {
		http.Error(w, err.Error(), http.StatusNotFound)
		return
	}
	writeJSON(w, http.StatusOK, act)
}

func (s *Server) handleGetVitality(w http.ResponseWriter, r *http.Request) {
	id := r.PathValue("id")
	score := s.orch.Vitality().Get(id)
	writeJSON(w, http.StatusOK, map[string]any{"score": score, "level": vitality.Level(score.Value)})
}

func (s *Server) handleListVitality(w http.ResponseWriter, r *http.Request) {
	writeJSON(w, http.StatusOK, s.orch.Vitality().List())
}

func (s *Server) handleListPlaybooks(w http.ResponseWriter, r *http.Request) {
	writeJSON(w, http.StatusOK, s.orch.Playbooks().List())
}

func (s *Server) handleCreatePlaybook(w http.ResponseWriter, r *http.Request) {
	var req struct {
		Name        string            `json:"name"`
		Description string            `json:"description"`
		Kind        string            `json:"kind"`
		Params      map[string]string `json:"params"`
	}
	if err := json.NewDecoder(r.Body).Decode(&req); err != nil || req.Name == "" {
		http.Error(w, "name required", http.StatusBadRequest)
		return
	}
	p := s.orch.Playbooks().Create(req.Name, req.Description, req.Kind, req.Params)
	writeJSON(w, http.StatusCreated, p)
}

func (s *Server) handleAssignPlaybook(w http.ResponseWriter, r *http.Request) {
	var req struct {
		PersonaID string `json:"persona_id"`
	}
	if err := json.NewDecoder(r.Body).Decode(&req); err != nil || req.PersonaID == "" {
		http.Error(w, "persona_id required", http.StatusBadRequest)
		return
	}
	a, err := s.orch.Playbooks().Assign(r.PathValue("id"), req.PersonaID)
	if err != nil {
		http.Error(w, err.Error(), http.StatusNotFound)
		return
	}
	writeJSON(w, http.StatusCreated, a)
}

func (s *Server) handleListAssignments(w http.ResponseWriter, r *http.Request) {
	personaID := r.URL.Query().Get("persona_id")
	writeJSON(w, http.StatusOK, s.orch.Playbooks().ListAssignments(personaID))
}

func (s *Server) handleListProxies(w http.ResponseWriter, r *http.Request) {
	writeJSON(w, http.StatusOK, s.orch.Proxies().List())
}

func (s *Server) handleSetProxy(w http.ResponseWriter, r *http.Request) {
	var req struct {
		Host     string `json:"host"`
		Port     int    `json:"port"`
		Type     string `json:"type"`
		Username string `json:"username"`
		Password string `json:"password"`
	}
	if err := json.NewDecoder(r.Body).Decode(&req); err != nil || req.Host == "" || req.Port == 0 {
		http.Error(w, "host and port required", http.StatusBadRequest)
		return
	}
	a := s.orch.Proxies().Set(r.PathValue("id"), req.Host, req.Port, req.Type, req.Username, req.Password)
	writeJSON(w, http.StatusOK, a)
}

func (s *Server) handleGetProxy(w http.ResponseWriter, r *http.Request) {
	a, ok := s.orch.Proxies().Get(r.PathValue("id"))
	if !ok {
		http.Error(w, "proxy not set", http.StatusNotFound)
		return
	}
	writeJSON(w, http.StatusOK, a)
}

func (s *Server) handleDeleteProxy(w http.ResponseWriter, r *http.Request) {
	if !s.orch.Proxies().Delete(r.PathValue("id")) {
		http.Error(w, "proxy not set", http.StatusNotFound)
		return
	}
	w.WriteHeader(http.StatusNoContent)
}

func writeJSON(w http.ResponseWriter, status int, v any) {
	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(status)
	_ = json.NewEncoder(w).Encode(v)
}
