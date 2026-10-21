package api

import (
	"encoding/json"
	"fmt"
	"io"
	"net/http"
	"strings"

	"github.com/relentlessworks/barcodekit/internal/model"
)

// Handler holds the HTTP handlers for the barcodekit API.
type Handler struct {
	noAuth bool
}

// NewHandler creates a new API handler.
func NewHandler(noAuth bool) *Handler {
	return &Handler{noAuth: noAuth}
}

// Register routes on the given mux.
func (h *Handler) Register(mux *http.ServeMux) {
	mux.HandleFunc("/help", h.handleHelp)
	mux.HandleFunc("/.well-known/agent.md", h.handleHelp)
	mux.HandleFunc("/types", h.handleTypes)
	mux.HandleFunc("/formats", h.handleFormats)
	mux.HandleFunc("/generate", h.handleGenerate)
	mux.HandleFunc("/mcp", h.handleMCP)
	// Individual type endpoints
	for _, t := range model.SupportedTypes() {
		mux.HandleFunc("/"+string(t), h.handleGenerateType(t))
	}
}

// wantsJSON checks if the client wants JSON output.
func wantsJSON(r *http.Request) bool {
	if r.URL.Query().Get("format") == "json" {
		return true
	}
	return r.Header.Get("Accept") == "application/json"
}

// writeError writes an instructive error response.
func writeError(w http.ResponseWriter, r *http.Request, status int, msg, hint string) {
	if wantsJSON(r) {
		w.Header().Set("Content-Type", "application/json")
		w.WriteHeader(status)
		json.NewEncoder(w).Encode(map[string]string{
			"error": msg,
			"hint":  hint,
		})
		return
	}
	w.Header().Set("Content-Type", "text/plain; charset=utf-8")
	w.WriteHeader(status)
	fmt.Fprintf(w, "error: %s | hint: %s\n", msg, hint)
}

// writeRecord writes a successful response in plain text or JSON.
func writeRecord(w http.ResponseWriter, r *http.Request, key, value string) {
	if wantsJSON(r) {
		w.Header().Set("Content-Type", "application/json")
		json.NewEncoder(w).Encode(map[string]string{key: value})
		return
	}
	w.Header().Set("Content-Type", "text/plain; charset=utf-8")
	fmt.Fprintf(w, "%s=%s\n", key, value)
}

// handleHelp returns the self-documenting help text.
func (h *Handler) handleHelp(w http.ResponseWriter, r *http.Request) {
	helpText := `barcodekit — agentic-first barcode generation service

Generate barcodes in 10 formats: Code128, EAN-13, EAN-8, UPC-A, Code39, Code93, ITF, Codabar, MSI, Pharmacode.
Output as SVG (vector) or PNG (base64 data URI).

AUTH:
  No auth required in default mode. Use -no-auth flag to explicitly disable.

ENDPOINTS:

  POST /generate
    Body (JSON): {"type":"code128","data":"Hello","format":"svg","width":200,"height":60,"padding":10}
    Body (form): type=code128&data=Hello&format=svg&width=200&height=60&padding=10
    Returns: the barcode image (SVG text or PNG base64 data URI)

  GET /generate?type=code128&data=Hello&format=svg&width=200&height=60&padding=10
    Same as POST but via query params

  POST /{type}
    Shortcut: POST /code128 with data=Hello&format=svg
    Same params as /generate but type is in the URL

  GET /{type}?data=Hello&format=svg
    Same as POST but via query params

  GET /types
    List all supported barcode types

  GET /formats
    List all supported output formats

  POST /mcp
    MCP (Model Context Protocol) JSON-RPC 2.0 endpoint

PARAMETERS:
  type    - barcode type (code128, ean13, ean8, upca, code39, code93, itf, codabar, msi, pharmacode)
  data    - the data to encode (required)
  format  - output format: svg (default) or png
  width   - image width in pixels (default: 200)
  height  - image height in pixels (default: 60)
  padding - padding in pixels (default: 10)

BARCODE TYPES:
  code128     - any ASCII text (1-80 chars)
  ean13       - 12 or 13 digits (checksum auto-calculated if 12)
  ean8        - 7 or 8 digits (checksum auto-calculated if 7)
  upca        - 11 or 12 digits (checksum auto-calculated if 11)
  code39      - A-Z, 0-9, - . $ / + % SPACE (start/stop * added automatically)
  code93      - A-Z, 0-9, - . $ / + % SPACE
  itf         - even number of digits (0-9)
  codabar     - digits 0-9 and - $ : / . + (start/stop A-D added automatically)
  msi         - digits 0-9
  pharmacode  - integer 3-131070

OUTPUT FORMATS:
  svg - SVG vector image (returned as text)
  png - PNG raster image (returned as base64 data URI)

RESPONSE FORMAT:
  Plain text by default (one labeled line per response).
  JSON available via Accept: application/json header or ?format=json query param.

ERRORS:
  All 4xx errors include a hint telling you what to do next.
  Format: error: message | hint: what to do next

EXAMPLES:
  curl "http://localhost:8080/generate?type=code128&data=Hello&format=svg"
  curl -X POST http://localhost:8080/generate -d 'type=ean13&data=590123412345&format=svg'
  curl "http://localhost:8080/code39?data=HELLO&format=png"
  curl -H "Accept: application/json" "http://localhost:8080/generate?type=code128&data=Hello"
`
	w.Header().Set("Content-Type", "text/plain; charset=utf-8")
	fmt.Fprint(w, helpText)
}

// handleTypes returns all supported barcode types.
func (h *Handler) handleTypes(w http.ResponseWriter, r *http.Request) {
	if wantsJSON(r) {
		w.Header().Set("Content-Type", "application/json")
		types := make([]map[string]string, 0)
		for _, t := range model.SupportedTypes() {
			types = append(types, map[string]string{
				"type":        string(t),
				"description": model.TypeDescription(t),
				"data":        model.TypeDataPattern(t),
			})
		}
		json.NewEncoder(w).Encode(types)
		return
	}
	w.Header().Set("Content-Type", "text/plain; charset=utf-8")
	for _, t := range model.SupportedTypes() {
		fmt.Fprintf(w, "type=%s desc=%s data=%s\n", t, model.TypeDescription(t), model.TypeDataPattern(t))
	}
}

// handleFormats returns all supported output formats.
func (h *Handler) handleFormats(w http.ResponseWriter, r *http.Request) {
	if wantsJSON(r) {
		w.Header().Set("Content-Type", "application/json")
		formats := make([]string, 0)
		for _, f := range model.SupportedFormats() {
			formats = append(formats, string(f))
		}
		json.NewEncoder(w).Encode(formats)
		return
	}
	w.Header().Set("Content-Type", "text/plain; charset=utf-8")
	for _, f := range model.SupportedFormats() {
		fmt.Fprintf(w, "format=%s\n", f)
	}
}

// parseRequest extracts barcode parameters from request.
func parseRequest(r *http.Request, defaultType model.BarcodeType) (*model.BarcodeRequest, error) {
	req := &model.BarcodeRequest{
		Type:    defaultType,
		Format:  model.FormatSVG,
		Width:   200,
		Height:  60,
		Padding: 10,
	}

	// Try JSON body first
	if r.Header.Get("Content-Type") == "application/json" && r.Body != nil {
		body, err := io.ReadAll(r.Body)
		if err != nil {
			return nil, fmt.Errorf("failed to read request body")
		}
		if len(body) > 0 {
			var jsonReq model.BarcodeRequest
			if err := json.Unmarshal(body, &jsonReq); err != nil {
				return nil, fmt.Errorf("invalid JSON body")
			}
			if jsonReq.Type != "" {
				req.Type = jsonReq.Type
			}
			if jsonReq.Data != "" {
				req.Data = jsonReq.Data
			}
			if jsonReq.Format != "" {
				req.Format = jsonReq.Format
			}
			if jsonReq.Width > 0 {
				req.Width = jsonReq.Width
			}
			if jsonReq.Height > 0 {
				req.Height = jsonReq.Height
			}
			if jsonReq.Padding >= 0 {
				req.Padding = jsonReq.Padding
			}
		}
	}

	// Form data / query params override
	if err := r.ParseForm(); err == nil {
		if v := r.Form.Get("type"); v != "" {
			req.Type = model.BarcodeType(v)
		}
		if v := r.Form.Get("data"); v != "" {
			req.Data = v
		}
		if v := r.Form.Get("format"); v != "" && v != "json" {
			req.Format = model.OutputFormat(v)
		}
		if v := r.Form.Get("width"); v != "" {
			fmt.Sscanf(v, "%d", &req.Width)
		}
		if v := r.Form.Get("height"); v != "" {
			fmt.Sscanf(v, "%d", &req.Height)
		}
		if v := r.Form.Get("padding"); v != "" {
			fmt.Sscanf(v, "%d", &req.Padding)
		}
	}

	return req, nil
}

// handleGenerate handles POST/GET /generate
func (h *Handler) handleGenerate(w http.ResponseWriter, r *http.Request) {
	if r.Method != http.MethodGet && r.Method != http.MethodPost {
		writeError(w, r, http.StatusMethodNotAllowed, "method not allowed", "use GET or POST")
		return
	}

	req, err := parseRequest(r, "")
	if err != nil {
		writeError(w, r, http.StatusBadRequest, err.Error(), "check your request format")
		return
	}

	h.generateBarcode(w, r, req)
}

// handleGenerateType returns a handler for a specific barcode type.
func (h *Handler) handleGenerateType(t model.BarcodeType) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		if r.Method != http.MethodGet && r.Method != http.MethodPost {
			writeError(w, r, http.StatusMethodNotAllowed, "method not allowed", "use GET or POST")
			return
		}

		req, err := parseRequest(r, t)
		if err != nil {
			writeError(w, r, http.StatusBadRequest, err.Error(), "check your request format")
			return
		}
		req.Type = t

		h.generateBarcode(w, r, req)
	}
}

// generateBarcode generates and returns the barcode.
func (h *Handler) generateBarcode(w http.ResponseWriter, r *http.Request, req *model.BarcodeRequest) {
	// Validate type
	if !model.IsValidType(req.Type) {
		writeError(w, r, http.StatusBadRequest,
			fmt.Sprintf("unsupported barcode type: %s", req.Type),
			"use one of: "+strings.Join(typesToStrings(model.SupportedTypes()), ", "))
		return
	}

	// Validate data
	if req.Data == "" {
		writeError(w, r, http.StatusBadRequest,
			"missing data parameter",
			"provide the data to encode, e.g. data=Hello for type "+string(req.Type))
		return
	}

	// Validate format
	if !model.IsValidFormat(req.Format) {
		writeError(w, r, http.StatusBadRequest,
			fmt.Sprintf("unsupported output format: %s", req.Format),
			"use svg or png")
		return
	}

	// Encode barcode
	pattern, encoded, err := model.Encode(*req)
	if err != nil {
		writeError(w, r, http.StatusBadRequest, err.Error(), "check the data format for "+string(req.Type))
		return
	}

	// Render
	var output string
	switch req.Format {
	case model.FormatSVG:
		output = model.RenderSVG(pattern, req.Width, req.Height, req.Padding)
	case model.FormatPNG:
		output, err = model.RenderPNG(pattern, req.Width, req.Height, req.Padding)
		if err != nil {
			writeError(w, r, http.StatusInternalServerError, "failed to render PNG", "try SVG format instead")
			return
		}
		output = "data:image/png;base64," + output
	}

	// Return response
	if wantsJSON(r) {
		w.Header().Set("Content-Type", "application/json")
		json.NewEncoder(w).Encode(&model.BarcodeResult{
			Type:   req.Type,
			Data:   req.Data,
			Format: req.Format,
			Output: output,
			Width:  req.Width,
			Height: req.Height,
		})
		return
	}

	w.Header().Set("Content-Type", "text/plain; charset=utf-8")
	fmt.Fprintf(w, "type=%s data=%s format=%s encoded=%s width=%d height=%d\n", req.Type, req.Data, req.Format, encoded, req.Width, req.Height)
	fmt.Fprint(w, output)
	fmt.Fprintln(w)
}

func typesToStrings(types []model.BarcodeType) []string {
	result := make([]string, len(types))
	for i, t := range types {
		result[i] = string(t)
	}
	return result
}
