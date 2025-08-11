package tools

import (
    "context"
    "time"
)

// data.snapshot: returns simple stub quotes for now
func DataSnapshot() Handler {
    return func(ctx context.Context, args map[string]any) (any, error) {
        symsAny, ok := args["symbols"].([]any)
        if !ok {
            return map[string]any{"quotes": map[string]any{}}, nil
        }
        quotes := map[string]any{}
        now := time.Now().UTC().Format(time.RFC3339)
        for _, s := range symsAny {
            sym, _ := s.(string)
            if sym == "" { continue }
            quotes[sym] = map[string]any{
                "bid": 225.05,
                "ask": 225.15,
                "mid": 225.10,
                "ts":  now,
            }
        }
        return map[string]any{"quotes": quotes}, nil
    }
}


