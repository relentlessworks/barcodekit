package api

import (
	"encoding/json"
	"fmt"
	"io"
	"net/http"

	"github.com/relentlessworks/barcodekit/internal/model"
)

// MCPRequest represents a JSON-RPC 2.0 request.
type MCPRequest struct {
	JSONRPC string          `json:"jsonrpc"`
	ID      json.RawMessage `json:"id"`
	Method  string          `json:"method"`
	Params  json.RawMessage `json:"params,omitempty"`
}

// MCPResponse represents a JSON-RPC 2.0 response.
type MCPResponse struct {
	JSONRPC string          `json:"jsonrpc"`
	ID      json.RawMessage `json:"id"`
	Result  interface{}     `json:"result,omitempty"`
	Error   *MCPError       `json:"error,omitempty"`
}

// MCPError represents a JSON-RPC error.
type MCPError struct {
	Code    int    `json:"code"`
	Message string `json:"message"`
}

// handleMCP handles MCP JSON-RPC 2.0 requests.
func (h *Handler) handleMCP(w http.ResponseWriter, r *http.Request) {
	if r.Method != http.MethodPost {
		writeError(w, r, http.StatusMethodNotAllowed, "method not allowed", "use POST for MCP")
		return
	}

	body, err := io.ReadAll(r.Body)
	if err != nil {
		writeMCPError(w, nil, -32700, "parse error")
		return
	}

	var req MCPRequest
	if err := json.Unmarshal(body, &req); err != nil {
		writeMCPError(w, nil, -32700, "parse error")
		return
	}

	w.Header().Set("Content-Type", "application/json")

	switch req.Method {
	case "initialize":
		writeMCPResult(w, req.ID, map[string]interface{}{
			"protocolVersion": "2024-11-05",
			"capabilities": map[string]interface{}{
				"tools": map[string]interface{}{},
			},
			"serverInfo": map[string]interface{}{
				"name":    "barcodekit",
				"version": "0.1.0",
			},
		})

	case "tools/list":
		tools := []map[string]interface{}{}
		for _, t := range model.SupportedTypes() {
			tools = append(tools, map[string]interface{}{
				"name":        "generate_" + string(t),
				"description": fmt.Sprintf("Generate a %s barcode", model.TypeDescription(t)),
				"inputSchema": map[string]interface{}{
					"type": "object",
					"properties": map[string]interface{}{
						"data": map[string]interface{}{
							"type":        "string",
							"description": model.TypeDataPattern(t),
						},
						"format": map[string]interface{}{
							"type":        "string",
							"enum":        []string{"svg", "png"},
							"default":     "svg",
							"description": "Output format",
						},
						"width": map[string]interface{}{
							"type":        "integer",
							"default":     200,
							"description": "Image width in pixels",
						},
						"height": map[string]interface{}{
							"type":        "integer",
							"default":     60,
							"description": "Image height in pixels",
						},
						"padding": map[string]interface{}{
							"type":        "integer",
							"default":     10,
							"description": "Padding in pixels",
						},
					},
					"required": []string{"data"},
				},
			})
		}
		writeMCPResult(w, req.ID, map[string]interface{}{"tools": tools})

	case "tools/call":
		var params struct {
			Name      string          `json:"name"`
			Arguments json.RawMessage `json:"arguments"`
		}
		if err := json.Unmarshal(req.Params, &params); err != nil {
			writeMCPError(w, req.ID, -32602, "invalid params")
			return
		}

		// Parse tool name to get barcode type
		var bcType model.BarcodeType
		for _, t := range model.SupportedTypes() {
			if params.Name == "generate_"+string(t) {
				bcType = t
				break
			}
		}
		if bcType == "" {
			writeMCPError(w, req.ID, -32601, fmt.Sprintf("unknown tool: %s", params.Name))
			return
		}

		var args struct {
			Data    string `json:"data"`
			Format  string `json:"format"`
			Width   int    `json:"width"`
			Height  int    `json:"height"`
			Padding int    `json:"padding"`
		}
		if err := json.Unmarshal(params.Arguments, &args); err != nil {
			writeMCPError(w, req.ID, -32602, "invalid arguments")
			return
		}

		if args.Data == "" {
			writeMCPError(w, req.ID, -32602, "data is required")
			return
		}

		bcReq := model.BarcodeRequest{
			Type:    bcType,
			Data:    args.Data,
			Format:  model.FormatSVG,
			Width:   200,
			Height:  60,
			Padding: 10,
		}
		if args.Format == "png" {
			bcReq.Format = model.FormatPNG
		}
		if args.Width > 0 {
			bcReq.Width = args.Width
		}
		if args.Height > 0 {
			bcReq.Height = args.Height
		}
		if args.Padding >= 0 {
			bcReq.Padding = args.Padding
		}

		pattern, encoded, err := model.Encode(bcReq)
		if err != nil {
			writeMCPError(w, req.ID, -32603, err.Error())
			return
		}

		var output string
		switch bcReq.Format {
		case model.FormatSVG:
			output = model.RenderSVG(pattern, bcReq.Width, bcReq.Height, bcReq.Padding)
		case model.FormatPNG:
			output, err = model.RenderPNG(pattern, bcReq.Width, bcReq.Height, bcReq.Padding)
			if err != nil {
				writeMCPError(w, req.ID, -32603, "failed to render PNG")
				return
			}
			output = "data:image/png;base64," + output
		}

		writeMCPResult(w, req.ID, map[string]interface{}{
			"content": []map[string]interface{}{
				{
					"type": "text",
					"text": fmt.Sprintf("type=%s data=%s format=%s encoded=%s\n%s", bcType, args.Data, bcReq.Format, encoded, output),
				},
			},
		})

	default:
		writeMCPError(w, req.ID, -32601, fmt.Sprintf("method not found: %s", req.Method))
	}
}

func writeMCPResult(w http.ResponseWriter, id json.RawMessage, result interface{}) {
	json.NewEncoder(w).Encode(MCPResponse{
		JSONRPC: "2.0",
		ID:      id,
		Result:  result,
	})
}

func writeMCPError(w http.ResponseWriter, id json.RawMessage, code int, msg string) {
	json.NewEncoder(w).Encode(MCPResponse{
		JSONRPC: "2.0",
		ID:      id,
		Error: &MCPError{
			Code:    code,
			Message: msg,
		},
	})
}
