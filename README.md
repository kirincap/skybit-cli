# Skybit CLI Trader

> **Claude‑Code‑style trading in your terminal:** type your intent → get a plan/analysis → preview the orders → approve to execute.

Trade equities from any terminal through a natural‑language chat interface. The LLM interprets intent and uses MCP tools to draft, simulate, and execute brokered orders—only after your approval.

## Features

- 🤖 **Natural Language Trading** - Express trading intent in plain English
- 📊 **Real-time Market Data** - NBBO quotes and L2 depth from multiple vendors  
- 🔒 **Safety First** - Preview all orders before execution with risk checks
- 🏦 **Multi-Broker Support** - SnapTrade aggregator with IBKR and Robinhood
- 📝 **Full Audit Trail** - Tamper-evident logging of all trading activity
- 🚀 **Fast & Lightweight** - Native Go binary, no background processes

## Quick Start

### Installation

#### macOS/Linux via Homebrew
```bash
brew tap skybit/tap
brew install skybit-cli
```

#### Direct Download
Download the latest binary from [GitHub Releases](https://github.com/skybit/cli/releases)

### Setup

1. **Login with device-code authentication**
   ```bash
   skybit login
   ```

2. **Connect your broker**
   ```bash
   skybit brokers connect
   ```

3. **Add market data provider**
   ```bash
   skybit data add polygon
   # Enter your API key when prompted
   ```

4. **Start trading**
   ```bash
   skybit
   # Type: Buy 100 shares of AAPL at limit 225.50
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

Configuration file location: `~/.skybit/config.yaml`

```yaml
# Example configuration
theme: dark
default_broker: ibkr
data_providers:
  - polygon
  - iex
paper_mode: false
risk_limits:
  max_position_size: 10000
  max_daily_trades: 50
```

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