package api

import (
	"encoding/json"
	"net/http"
)

func (s *Server) handleInputTap(w http.ResponseWriter, r *http.Request) {
	var req struct {
		X int `json:"x"`
		Y int `json:"y"`
	}
	if err := json.NewDecoder(r.Body).Decode(&req); err != nil {
		http.Error(w, "invalid json", http.StatusBadRequest)
		return
	}
	if err := s.orch.InputTap(r.Context(), r.PathValue("id"), req.X, req.Y); err != nil {
		http.Error(w, err.Error(), http.StatusBadRequest)
		return
	}
	writeJSON(w, http.StatusOK, map[string]any{"status": "ok", "action": "tap", "x": req.X, "y": req.Y})
}

func (s *Server) handleInputSwipe(w http.ResponseWriter, r *http.Request) {
	var req struct {
		X1 int `json:"x1"`
		Y1 int `json:"y1"`
		X2 int `json:"x2"`
		Y2 int `json:"y2"`
		Ms int `json:"ms"`
	}
	if err := json.NewDecoder(r.Body).Decode(&req); err != nil {
		http.Error(w, "invalid json", http.StatusBadRequest)
		return
	}
	if err := s.orch.InputSwipe(r.Context(), r.PathValue("id"), req.X1, req.Y1, req.X2, req.Y2, req.Ms); err != nil {
		http.Error(w, err.Error(), http.StatusBadRequest)
		return
	}
	writeJSON(w, http.StatusOK, map[string]string{"status": "ok", "action": "swipe"})
}

func (s *Server) handleInputKey(w http.ResponseWriter, r *http.Request) {
	var req struct {
		Keycode int `json:"keycode"`
	}
	if err := json.NewDecoder(r.Body).Decode(&req); err != nil || req.Keycode == 0 {
		http.Error(w, "keycode required", http.StatusBadRequest)
		return
	}
	if err := s.orch.InputKey(r.Context(), r.PathValue("id"), req.Keycode); err != nil {
		http.Error(w, err.Error(), http.StatusBadRequest)
		return
	}
	writeJSON(w, http.StatusOK, map[string]any{"status": "ok", "keycode": req.Keycode})
}

func (s *Server) handleInputText(w http.ResponseWriter, r *http.Request) {
	var req struct {
		Text string `json:"text"`
	}
	if err := json.NewDecoder(r.Body).Decode(&req); err != nil || req.Text == "" {
		http.Error(w, "text required", http.StatusBadRequest)
		return
	}
	if err := s.orch.InputText(r.Context(), r.PathValue("id"), req.Text); err != nil {
		http.Error(w, err.Error(), http.StatusBadRequest)
		return
	}
	writeJSON(w, http.StatusOK, map[string]string{"status": "ok", "action": "text"})
}
