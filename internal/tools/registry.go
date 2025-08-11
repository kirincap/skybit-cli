package tools

import (
    "context"
    "fmt"
)

// Handler signature for a tool adapter.
type Handler func(ctx context.Context, args map[string]any) (any, error)

// Registry maps tool names to handlers.
type Registry struct {
    handlers map[string]Handler
}

func NewRegistry() *Registry {
    return &Registry{handlers: make(map[string]Handler)}
}

func (r *Registry) Register(name string, h Handler) {
    r.handlers[name] = h
}

func (r *Registry) Names() []string {
    names := make([]string, 0, len(r.handlers))
    for k := range r.handlers {
        names = append(names, k)
    }
    return names
}

func (r *Registry) Call(ctx context.Context, name string, args map[string]any) (any, error) {
    h, ok := r.handlers[name]
    if !ok {
        return nil, fmt.Errorf("unknown tool: %s", name)
    }
    return h(ctx, args)
}


