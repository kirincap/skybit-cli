# Planner Orchestration via OpenRouter

This document specifies how the backend planner/orchestrator uses OpenRouter as the primary LLM to turn `Chat` into `Plan` and `Preview`, and then executes `Approve/Cancel` via MCP tools.

## Provider
- Base URL: `https://openrouter.ai/api/v1`
- Auth: `Authorization: Bearer ${OPENROUTER_API_KEY}`
- Default model: `anthropic/claude-3-opus`
- Fallbacks: `openai/gpt-4.1`, `anthropic/claude-3.5-sonnet`
- Headers: `HTTP-Referer` (repo URL), `X-Title: Skybit CLI`

## Messages and tools
- API: OpenAI-compatible `/chat/completions`
- Submit `messages` and `tools` (strict JSON Schemas; see `docs/mcp-tools-schemas.json`). Tools include live SnapTrade-backed endpoints: `broker.accounts`, `broker.positions`, `trade.place_order.live`, `trade.cancel.live`, `trade.cancel_all.live`.
- `tool_choice`: `auto`
- Settings: `temperature: 0.1`, `top_p: 0.9`, `max_tokens`: tuned per stage

## Tool loop (server-side)
1. Receive `Chat` (WS) with user text and lightweight context flags.
2. Call OpenRouter with system prompt, user message, and tool schemas.
3. If the model emits `tool_calls`, execute matching MCP tool(s) via adapters:
   - Validate args against schema; on error, return structured tool error.
   - Record tool result, redact secrets, append as `role: tool` message.
4. Resume completion; if a plan can be summarized, emit `Plan` (WS) early.
5. After relevant tools complete, build a `Preview` (WS) with human-readable summary and machine fields (orders, impact, policy, quote).
6. Wait for `Approve`/`Cancel` from client; on `Approve`, execute `trade.place_order` with `client_order_id` and `idem_key` and stream `Status` updates.

## Request example (OpenRouter)
```json
{
  "model": "anthropic/claude-3-opus",
  "temperature": 0.1,
  "messages": [
    {"role": "system", "content": "You are a trading planner. Use tools deterministically and return safe previews before any execution."},
    {"role": "user", "content": "Buy 200 shares of AAPL at limit 225.10, good for the day"}
  ],
  "tools": [ /* see docs/mcp-tools-schemas.json */ ],
  "tool_choice": "auto",
  "stream": true
}
```

## Fallback and retries
- Retries: exponential backoff on 429/5xx (jittered, cap at 30s)
- Fallback policy: downgrade to `openai/gpt-4.1` then `anthropic/claude-3.5-sonnet` when latency SLA or budget exceeded
- Timeouts: per-request timeout; cancel outstanding calls on disconnect

## Safety and budgets
- Determinism: low `temperature`, validated tool arguments, schema-constrained outputs
- Budgets: per-request token caps; per-session cost ceiling; disable streaming rationale when near cap
- Redaction: remove PII, secrets, account IDs before echoing to the model or logs

## Observability
- Metrics: model latency, tool latency, preview TTFB, token usage, cost/session, error rates
- Logs: structured tool graph (names, arg hashes, timings), envelope IDs, session IDs (no secrets)
- Tracing: spans for Chat → Planner → Tool Calls → Preview → Approve → Status

## Configuration
- `OPENROUTER_API_KEY` (required)
- `OPENROUTER_BASE_URL` (default `https://openrouter.ai/api/v1`)
- `OPENROUTER_PRIMARY_MODEL` (default `anthropic/claude-3-opus`)
- `OPENROUTER_FALLBACK_MODELS` (comma-separated)
- Budget caps: `LLM_MAX_TOKENS`, `LLM_SESSION_BUDGET_USD`, latency SLOs

## Error handling
- Tool schema violation → return tool error to the model; regenerate with corrected args
- Planner error → emit `Error` (WS) with recoverable=true and suggestions
- OpenRouter error → retry; on persistent failure, fallback then emit `Error`

## Streaming behavior
- Stream textual plan rationale if helpful, but ensure `Plan` WS emits a compact machine-readable plan
- Pause model stream while tools execute; resume with `role: tool` messages
- Emit `Preview` WS as soon as human reviewable content is ready

## Security
- Secrets stored in secret manager; never logged
- TLS everywhere; strict timeouts and circuit breakers
- Validate all user content before passing to the model; escape/normalize symbols and quantities
