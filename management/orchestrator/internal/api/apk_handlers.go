package api

import (
	"io"
	"net/http"
	"path/filepath"
)

// POST /v1/instances/{id}/install-apk
// multipart form field "apk" (file) OR JSON is not used — file upload only.
func (s *Server) handleInstallAPK(w http.ResponseWriter, r *http.Request) {
	if err := r.ParseMultipartForm(512 << 20); err != nil { // 512MB
		http.Error(w, "multipart form required (field: apk)", http.StatusBadRequest)
		return
	}
	f, hdr, err := r.FormFile("apk")
	if err != nil {
		http.Error(w, "form field 'apk' required", http.StatusBadRequest)
		return
	}
	defer f.Close()
	data, err := io.ReadAll(f)
	if err != nil {
		http.Error(w, err.Error(), http.StatusBadRequest)
		return
	}
	name := filepath.Base(hdr.Filename)
	out, err := s.orch.InstallAPKBytes(r.Context(), r.PathValue("id"), data, name)
	if err != nil {
		http.Error(w, err.Error()+": "+out, http.StatusBadRequest)
		return
	}
	writeJSON(w, http.StatusOK, map[string]string{"status": "installed", "detail": out, "file": name})
}
