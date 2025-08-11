package tools

import (
    "context"
    "encoding/json"
    "os"
    "path/filepath"
    "time"
)

func PolicyCheck() Handler {
    return func(ctx context.Context, args map[string]any) (any, error) {
        return map[string]any{"allowed": true}, nil
    }
}

func AuditLog() Handler {
    return func(ctx context.Context, args map[string]any) (any, error) {
        home, _ := os.UserHomeDir()
        path := filepath.Join(home, ".skybit", "audit.jsonl")
        _ = os.MkdirAll(filepath.Dir(path), 0o755)
        f, err := os.OpenFile(path, os.O_CREATE|os.O_APPEND|os.O_WRONLY, 0o644)
        if err != nil { return map[string]any{"persisted": false}, nil }
        defer f.Close()
        payload := map[string]any{"ts": time.Now().UTC().Format(time.RFC3339), "event": args["event"], "payload": args["payload"]}
        b, _ := json.Marshal(payload)
        _, _ = f.Write(append(b, '\n'))
        return map[string]any{"persisted": true}, nil
    }
}


