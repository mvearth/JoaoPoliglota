package handler

import (
	"encoding/json"
	"log"
	"net/http"

	"joao_poliglota/translation"
)

// maxBodyBytes caps the size of an incoming request body to guard against
// oversized payloads.
const maxBodyBytes = 1 << 20 // 1 MiB

// Handler serves the translation HTTP endpoints.
type Handler struct {
	translator translation.Translator
}

// New builds a Handler around a translator.
func New(t translation.Translator) *Handler {
	return &Handler{translator: t}
}

// Translate decodes a TranslationDictionary, translates it and writes the
// result as JSON.
func (h *Handler) Translate(w http.ResponseWriter, r *http.Request) {
	r.Body = http.MaxBytesReader(w, r.Body, maxBodyBytes)

	decoder := json.NewDecoder(r.Body)
	decoder.DisallowUnknownFields()

	var input translation.TranslationDictionary
	if err := decoder.Decode(&input); err != nil {
		log.Printf("handler: decoding request: %v", err)
		writeError(w, http.StatusBadRequest, "invalid request body")
		return
	}

	if input.OutputLang == "" {
		writeError(w, http.StatusBadRequest, "outputLang is required")
		return
	}
	if len(input.Dictionary) == 0 {
		writeError(w, http.StatusBadRequest, "dictionary must not be empty")
		return
	}

	output, err := h.translator.Translate(r.Context(), input)
	if err != nil {
		log.Printf("handler: translating: %v", err)
		writeError(w, http.StatusInternalServerError, "translation failed")
		return
	}

	writeJSON(w, http.StatusOK, output)
}

func writeJSON(w http.ResponseWriter, status int, payload interface{}) {
	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(status)
	if err := json.NewEncoder(w).Encode(payload); err != nil {
		log.Printf("handler: encoding response: %v", err)
	}
}

func writeError(w http.ResponseWriter, status int, message string) {
	writeJSON(w, status, map[string]string{"error": message})
}
