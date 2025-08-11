# Artifact — Technical Implementation (MVP)

This document describes the MVP technical design of the Artifact CLI trader. It focuses on how the system is built and how it executes trades safely via SnapTrade and market‑data vendors, driven by a chat interface and an MCP tool layer.

Note on naming: the product name is Artifact; the current CLI binary is `skybit` in v0.1. Replace `skybit` with `artifact` if/when an alias is introduced.

## 1. System architecture

Components:
- CLI/TUI (Go binary): interactive terminal UI with panes (chat, quotes, book, orders, positions, PnL, logs).
 - Chat/Planner (LLM): interprets user intent and produces a Plan (ordered MCP tool calls) and a human Preview.
   - Provider: OpenRouter (`https://openrouter.ai/api/v1`)
   - Default model: `anthropic/claude-3-opus` (Opus)
   - Fallbacks: `openai/gpt-4.1`, `anthropic/claude-3.5-sonnet`
   - Settings: temperature 0.1–0.2, strict JSON tool calling, budget and latency guardrails
- MCP tool layer: deterministic, callable tools for data, broker, trading, policy, and audit.
- Aggregator/Broker: SnapTrade for auth, accounts, positions, and live order submission. IBKR/Robinhood in read‑only at MVP.
- Market data vendors: Polygon, IEX Cloud, Nasdaq Basic, Databento; CLI renders live streams; model uses snapshots.
- Local state: short‑lived session token, config, policy file, append‑only audit log; no background daemons.

Data flow (high level):
1) CLI streams L1/L2 from vendor WS → renders panes.
2) User enters a plain‑English request → LLM plans MCP calls.
3) MCP fetches snapshots/positions → builds human Preview.
4) User approves → policy checks gate live submit → SnapTrade order → status streamed → audit logged.

## 2. Execution lifecycle

1) Login: `skybit login` starts device‑code flow; store short‑lived token.
2) Broker link: `skybit brokers connect` runs SnapTrade Connect (OAuth/MFA) in browser.
3) Data setup: `skybit data add <provider>` adds user API keys; CLI connects WS during sessions.
4) Session: `skybit` opens chat + market panes and maintains WS connections.
5) Plan & Preview: planner proposes tool calls and renders Preview (orders, fees, slippage, policy result).
6) Approval: inline confirm. If paper mode → simulate. If live → run policy checks then submit.
7) Submit & stream: place order(s) via SnapTrade; stream updates (accepted/filled/rejected/canceled).
8) Audit: append events (inputs, outputs, timing) to a tamper‑evident log; exportable.

Pseudocode (conceptual):
```text
while session_open:
  user_intent = read_chat()
  plan = llm.plan(user_intent)
  data = mcp.data.snapshot(plan.symbols)
  positions = mcp.broker.positions(active_account)
  preview = build_preview(plan, data, positions)
  show(preview)
  if approve():
    ensure(policy.check(preview))
    result = mcp.trade.place(preview.orders)
    audit.log({intent, plan, preview, result})
```

## 3. MCP tool specification (MVP)

Tool families and representative contracts (shape only; exact fields may vary):

- data.snapshot
  - Input: `{ symbols: ["AAPL", "MSFT", ...] }`
  - Output: `{ quotes: { "AAPL": { bid, ask, mid, ts }, ... } }`
- data.l2_snapshot
  - Input: `{ symbol: "AAPL", depth: 10 }`
  - Output: `{ bids: [...], asks: [...], ts }`
- broker.accounts
  - Output: `{ accounts: [{ id, name, type }...] }`
- broker.positions
  - Input: `{ account_id }`
  - Output: `{ positions: [{ symbol, qty, avg_price, mkt_value }...] }`
- trade.preview
  - Input: `{ account_id, orders: [{ symbol, side, qty, type, limit_price?, tif }] }`
  - Output: `{ orders: [...], est_avg_price, fees, slippage_bps, pnl_delta }`
- trade.place_order
  - Input: `{ account_id, order: { symbol, side, qty, type, limit_price?, tif, client_id } }`
  - Output: `{ broker_order_id, status, submitted_at }`
- trade.cancel / trade.cancel_all
  - Input: `{ account_id, broker_order_id }` / `{ account_id }`
  - Output: `{ status }`
- policy.check
  - Input: `{ preview, context: { time, account_limits, entitlements } }`
  - Output: `{ allowed: true|false, reason? }`
- audit.log
  - Input: `{ event, payload, ts }`
  - Output: `{ persisted: true }`

Example payloads:
```jsonc
// trade.preview
{
  "account_id": "acct_123",
  "orders": [
    { "symbol": "AAPL", "side": "BUY", "qty": 100, "type": "LIMIT", "limit_price": 225.10, "tif": "DAY" }
  ]
}
```
```jsonc
// trade.place_order response
{ "broker_order_id": "SNAP-abc123", "status": "ACCEPTED", "submitted_at": "2025-08-11T17:32:00Z" }
```

## 4. Connectivity & data

- Aggregator: SnapTrade is the live execution path (OAuth/MFA via Connect; tokens stored securely; scopes least‑privilege).
- Read‑only: IBKR, Robinhood for positions/quotes at MVP.
- Data vendors: user supplies API keys; entitlements are displayed; we avoid redistributing proprietary data.
- Streaming model (CLI): vendor WS → L1 NBBO (default) and optional L2 where licensed; heartbeats + backoff + gap detection and auto‑resubscribe.
- Snapshot model (LLM): planner uses `data.snapshot`/`data.l2_snapshot` for deterministic reasoning.

### LLM orchestration (OpenRouter)

- Base URL: `https://openrouter.ai/api/v1`
- Auth: `Authorization: Bearer ${OPENROUTER_API_KEY}`
- Headers: `HTTP-Referer`, `X-Title: Skybit CLI`
- Messages: OpenAI-style `messages` with `tools` (JSON Schemas matching MCP tools); expect `tool_calls` with arguments.
- Retries: exponential backoff on 429/5xx; per-request timeout; per-session budget cap.
- Observability: log tool graph and latencies (no secrets or raw tokens).

## 5. Policy engine

- Deterministic gate before live submit; runs on Preview.
- Minimal YAML policy (example):
```yaml
# ~/.skybit/policy.yaml
limits:
  max_notional_per_order: 25000
  max_position_qty: 1000
  trading_window:
    start: "09:30"
    end:   "16:00"
  deny_symbols: ["GME", "AMC"]
```
- Checks: notional, TIF/time window, symbol denylist, position/sector caps, paper‑mode default for new users.
- CLI helpers: `skybit policy open`, `skybit policy check`.

## 6. Diagnostics

- `skybit doctor` runs:
  - SnapTrade reachability + auth scopes
  - Vendor entitlements and WS reachability/latency
  - System clock drift and TLS
  - Config and environment sanity

## 7. Configuration

Default path: `~/.skybit/config.yaml`
```yaml
# Example
default_broker: snaptrade
paper_mode: true
risk_limits:
  max_position_size: 10000
  max_daily_trades: 50
data_providers:
  - polygon
  - iex
```

CLI environment variables:

- `SKYBIT_WS_URL` (default `wss://api.skybit.ai/artifact/v1`)
- `SKYBIT_TOKEN` (required) — session token for WS Auth
- `SKYBIT_VERSION` (optional)

## 8. Command surface (developer‑facing)

- Session & auth: `skybit`, `skybit login`, `skybit logout`, `skybit doctor`
- Brokers: `skybit brokers connect|list|unlink`
- Data: `skybit data add|list|remove <provider>`
- Trading (in chat): `quote <SYM>`, `book <SYM>`, `buy|sell …`, `cancel <ID>|all`, `positions`, `orders`, `pnl`
- Policy & config: `skybit policy open|check`, `skybit config show|set`
- Audit: `skybit audit tail|export`

Decision UX mirrors Claude Code:
- After `Preview`, the CLI accepts: `approve`, `cancel`, or `edit quantity=... limit_price=...` typed in chat.

## 9. Notes & scope boundaries (MVP)

In scope: SnapTrade equities live; IBKR read‑only; NBBO real‑time (user key) or delayed without key; optional L2; deterministic policy checks; paper engine; local cache.

Out of scope: options/futures execution; background automations/bots; team features; hosted cloud.
