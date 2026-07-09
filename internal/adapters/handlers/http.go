package handlers

import (
	"bytes"
	"github.com/cc-andres-portillo/cv-parser/internal/core/ports"
	"encoding/json"
	"io"
	"net/http"
	"path/filepath"
	"strings"
)

type HTTPCVHandler struct {
	cvService ports.CVServicePort
}

func NewHTTPCVHandler(service ports.CVServicePort) *HTTPCVHandler {
	return &HTTPCVHandler{cvService: service}
}

func (h *HTTPCVHandler) ParseCVHandler(w http.ResponseWriter, r *http.Request) {
	w.Header().Set("Access-Control-Allow-Origin", "*")
	w.Header().Set("Access-Control-Allow-Methods", "POST, OPTIONS")
	w.Header().Set("Access-Control-Allow-Headers", "Content-Type")

	if r.Method == "OPTIONS" {
		w.WriteHeader(http.StatusOK)
		return
	}

	if r.Method != "POST" {
		http.Error(w, "Método no permitido", http.StatusMethodNotAllowed)
		return
	}

	r.ParseMultipartForm(5 << 20) // Máximo 5MB
	file, header, err := r.FormFile("cv")
	if err != nil {
		http.Error(w, "Llave 'cv' requerida en el formulario", http.StatusBadRequest)
		return
	}
	defer file.Close()

	ext := strings.ToLower(filepath.Ext(header.Filename))
	
	// Leer el archivo a un buffer en memoria para pasarlo como io.ReaderAt
	var buf bytes.Buffer
	size, err := io.Copy(&buf, file)
	if err != nil {
		http.Error(w, "Error leyendo el buffer de archivo", http.StatusInternalServerError)
		return
	}
	readerAt := bytes.NewReader(buf.Bytes())

	// Invocar el caso de uso del núcleo hexagonal
	resultado, err := h.cvService.ParseCV(r.Context(), readerAt, size, ext)
	if err != nil {
		http.Error(w, err.Error(), http.StatusUnprocessableEntity)
		return
	}

	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(http.StatusOK)
	json.NewEncoder(w).Encode(resultado)
}
