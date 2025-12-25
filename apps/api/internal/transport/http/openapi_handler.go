package http

import (
	"encoding/json"
	"net/http"

	"github.com/go-chi/render"
)

func (s *Server) handleOpenAPISpec(w http.ResponseWriter, r *http.Request) {
	spec := GenerateOpenAPISpec()
	
	w.Header().Set("Content-Type", "application/json")
	json.NewEncoder(w).Encode(spec)
}

func (s *Server) handleOpenAPIYAML(w http.ResponseWriter, r *http.Request) {
	spec := GenerateOpenAPISpec()
	
	// Convert to YAML
	data, err := spec.MarshalJSON()
	if err != nil {
		render.Status(r, http.StatusInternalServerError)
		render.JSON(w, r, map[string]string{"error": "Failed to generate spec"})
		return
	}
	
	w.Header().Set("Content-Type", "application/json")
	w.Write(data)
}

