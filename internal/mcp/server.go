package mcp

import (
    "encoding/json"
    "fmt"
    "log"
    "net"
    "net/http"
    "time"

    "github.com/kirincap/skybit-cli/internal/tools"
)

type request struct {
    Name string                 `json:"name"`
    Args map[string]any         `json:"args"`
}

type response struct {
    OK    bool        `json:"ok"`
    Error string      `json:"error,omitempty"`
    Data  interface{} `json:"data,omitempty"`
}

// NewMux returns an HTTP handler for the MCP tool endpoint.
func NewMux(reg *tools.Registry) http.Handler {
    mux := http.NewServeMux()
    mux.HandleFunc("/mcp", func(w http.ResponseWriter, r *http.Request) {
        if r.Method != http.MethodPost {
            w.WriteHeader(http.StatusMethodNotAllowed)
            return
        }
        var req request
        if err := json.NewDecoder(r.Body).Decode(&req); err != nil {
            writeJSON(w, http.StatusBadRequest, response{OK: false, Error: "invalid json"})
            return
        }
        ctx := r.Context()
        ctx, cancel := time.WithTimeout(ctx, 10*time.Second)
        defer cancel()
        data, err := reg.Call(ctx, req.Name, req.Args)
        if err != nil {
            writeJSON(w, http.StatusOK, response{OK: false, Error: err.Error()})
            return
        }
        writeJSON(w, http.StatusOK, response{OK: true, Data: data})
    })
    return mux
}

// Start starts an MCP HTTP server on the given addr (":port" or "127.0.0.1:0"), returning the actual addr in use.
func Start(addr string, reg *tools.Registry) (string, func() error, error) {
    mux := NewMux(reg)
    server := &http.Server{Handler: mux}
    ln, err := net.Listen("tcp", addr)
    if err != nil { return "", nil, err }
    actual := ln.Addr().String()
    go func() {
        log.Printf("skybit mcp server listening on %s", actual)
        _ = server.Serve(ln)
    }()
    stop := func() error { return server.Close() }
    return actual, stop, nil
}

func writeJSON(w http.ResponseWriter, code int, v any) {
    w.Header().Set("Content-Type", "application/json")
    w.WriteHeader(code)
    if err := json.NewEncoder(w).Encode(v); err != nil {
        fmt.Println("encode error:", err)
    }
}


