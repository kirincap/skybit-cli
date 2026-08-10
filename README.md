# Skybit CLI — Minimal Agent via Crush + MCP

Natural‑language trading with explicit approvals. We run Charm’s Crush as the agent shell and expose Skybit tools via a local MCP HTTP server. The model uses OpenRouter (Opus) and calls tools deterministically; you review a Preview and then type `approve` to execute.

## How it works

1. Login — `skybit login` starts a device‑code browser flow; the CLI holds a short‑lived session token.
2. Link a broker — `skybit brokers connect` opens SnapTrade Connect; you OAuth/MFA with your broker once.
3. Add data — `skybit data add polygon|iex|nasdaq|databento` (your API keys). The CLI streams from vendor websockets; the model uses snapshots for determinism.
4. Chat — `skybit` opens a session. You type plain English; the model proposes a Plan (tool calls) and a Preview (orders, impact, fees if enabled, policy checks).
5. Approve — you confirm inline; we submit via SnapTrade and stream status; everything is logged.
6. Exit — quitting the CLI closes the websocket; no watchers, no daemons.

> **Callout — MCP in this product**
> MCP is the tool layer the LLM calls (e.g., `trade.place_order`, `data.snapshot`). MCP is not a copilot; the chat model provides intent, MCP provides callable tools.

## Features

- 🤖 **Natural Language Trading** - Express trading intent in plain English
- 📊 **Real-time Market Data** - NBBO quotes and L2 depth from multiple vendors  
- 🔒 **Safety First** - Preview all orders before execution with risk checks
- 🏦 **Multi-Broker Support** - SnapTrade aggregator with IBKR and Robinhood
- 📝 **Full Audit Trail** - Tamper-evident logging of all trading activity
- 🚀 **Fast & Lightweight** - Native Go binary, no background processes

## Quick Start

1) Start the local MCP server (embedded in skybit)
```bash
go run cmd/skybit/main.go
```

2) Option A — direct to OpenRouter (dev/early prod)
```bash
export OPENROUTER_API_KEY=your_openrouter_key
crush --config docs/crush-openrouter.json
```

Option B — via Skybit proxy (recommended for production)
```bash
export OPENROUTER_API_KEY=your_team_key SKYBIT_PROXY_API_KEY=your_skybit_token
go run cmd/skybit-proxy/main.go &
export SKYBIT_TOKEN=your_skybit_token
crush --config docs/crush.json
```

3) Chat examples
```text
Buy 100 AAPL at limit 225.10 DAY
approve
cancel all
```

## Usage

### Basic Commands

```bash
# Session management
skybit login              # Authenticate via browser
skybit logout             # End session
skybit doctor            # Diagnose connectivity issues

# Broker management  
skybit brokers connect    # Link new broker account
skybit brokers list      # Show connected accounts
skybit brokers unlink    # Remove broker connection

# Data providers
skybit data add <provider>     # Add market data source
skybit data list              # Show configured providers
skybit data remove <provider>  # Remove data source

# Trading (in chat interface)
quote AAPL               # Get current quote
book AAPL               # Show order book depth
buy 100 AAPL at 225.50  # Place limit buy order
sell all TSLA           # Sell entire position
cancel all              # Cancel all open orders
positions               # View current holdings
orders                 # List open orders
pnl                    # Show profit/loss
```

### Trading Examples

```bash
# Simple market-like order (uses midpoint)
> Buy 200 shares of NVDA

# Limit order with time constraint  
> Buy 500 AAPL at 225.10, good for the day

# Percentage-based position adjustment
> Trim 20% of my TSLA position if it breaks 225

# Risk-managed entry
> Buy $10k worth of SPY, cap slippage at 5 basis points
```

## Safety Features

- **Preview Before Execute** - Review order details, fees, and impact before confirming
- **Policy Engine** - Customizable risk rules (position limits, time windows, etc.)
- **Paper Trading** - Test strategies without real money (default for new users)
- **Panic Cancel** - Instantly cancel all pending orders
- **No Automation** - System only acts on explicit user approval

## System Requirements

- **OS**: macOS 11+, Linux (glibc 2.31+), Windows 10+
- **Memory**: 256MB RAM minimum
- **Network**: Stable internet connection for real-time data
- **Terminal**: Any modern terminal emulator with UTF-8 support

## Configuration

- `OPENROUTER_API_KEY` (required)
- `OPENROUTER_MODEL` (default `anthropic/claude-3-opus`)
- `SNAPTRADE_CLIENT_ID` (required for live orders)
- `SNAPTRADE_CLIENT_SECRET` (required for live orders)
- `SNAPTRADE_ENV` (default `sandbox`)
- `SNAPTRADE_ACCOUNT_ID` (optional; default is first linked account)

### Environment variables

- `SKYBIT_WS_URL` (default `wss://api.skybit.ai/artifact/v1`)
- `SKYBIT_TOKEN` (required) — session token used in WS Auth
- `SKYBIT_VERSION` (optional) — client version sent in Auth

Backend LLM orchestration uses OpenRouter with Opus by default (see docs/technical-implementation.md). No configuration needed in the CLI.

## Development

### Building from Source

```bash
# Clone repository
git clone https://github.com/skybit/cli
cd cli

# Install dependencies
go mod download

# Build binary
go build -o skybit ./cmd/skybit

# Run tests
go test ./...
```

### Project Structure

```
skybit-cli/
├── cmd/skybit/          # Main entry point
├── internal/
│   ├── app/            # Application logic
│   ├── broker/         # Broker integrations
│   ├── chat/           # LLM interface
│   ├── market/         # Market data handling
│   ├── policy/         # Risk management
│   ├── tui/            # Terminal UI components
│   └── ws/             # WebSocket client
├── pkg/                # Public packages
├── proto/              # Protocol definitions
└── docs/               # Documentation
```

## Documentation

- [Planner Orchestration (OpenRouter)](docs/planner-openrouter.md)
- [MCP Tool Schemas](docs/mcp-tools-schemas.json)
- [Crush sample config](docs/crush.json)

### Using with Crush

1. Start the local MCP server:
   ```bash
   go run cmd/skybit-mcp/main.go
   ```
2. Export your OpenRouter key and run Crush in this repo with our config:
   ```bash
   export OPENROUTER_API_KEY=your_key
   crush --config docs/crush.json
   ```
3. Chat normally. The model will call our MCP tools for preview/approve flows.

## Contributing

We welcome contributions! Please see [CONTRIBUTING.md](CONTRIBUTING.md) for guidelines.

### Development Setup

1. Fork the repository
2. Create a feature branch (`git checkout -b feature/amazing-feature`)
3. Commit your changes (`git commit -m 'Add amazing feature'`)
4. Push to the branch (`git push origin feature/amazing-feature`)
5. Open a Pull Request

## Security

- **Zero Custody** - We never hold your funds or keys
- **Encrypted Storage** - Credentials stored in system keychain
- **Least Privilege** - Minimal API permissions requested
- **Open Source** - Full code transparency

Report security issues to security@skybit.ai

## Support

- **Documentation**: [docs.skybit.ai](https://docs.skybit.ai)
- **Discord**: [discord.gg/skybit](https://discord.gg/skybit)
- **Email**: support@skybit.ai
- **Issues**: [GitHub Issues](https://github.com/skybit/cli/issues)

## License

MIT License - see [LICENSE](LICENSE) file for details

## Disclaimer

This software is provided for informational purposes only. Trading involves substantial risk of loss. Past performance does not guarantee future results. Always do your own research and consult with qualified financial advisors.