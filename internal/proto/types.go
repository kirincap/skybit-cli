package proto

import "encoding/json"

// Minimal stubs to represent messages exchanged over WS/tool-plane.

type PlanRequest struct {
	Prompt string `json:"prompt"`
}

type PreviewResponse struct {
	PlanID string `json:"plan_id"`
	Text   string `json:"text"`
}

type ApproveRequest struct {
	PlanID string `json:"plan_id"`
}

type StatusMessage struct {
	Kind string          `json:"kind"`
	Data json.RawMessage `json:"data"`
}
