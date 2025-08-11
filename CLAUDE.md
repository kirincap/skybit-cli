# CLAUDE.md - Skybit CLI Trading Assistant Instructions

## Project Context

This is the **Skybit CLI Trader** - a terminal-based natural language trading interface that enables users to trade equities through conversational commands. The system uses LLM interpretation with MCP tools to draft, simulate, and execute trades through broker integrations, requiring explicit user approval before any execution.

## Core Architecture

### Technology Stack
- **Language**: Go 1.21+
- **TUI Framework**: Charmbracelet (Bubble Tea, Bubbles, Lip Gloss, Glamour)
- **Communication**: WebSocket client to backend API
- **Data Protocol**: JSON over WebSocket with defined message envelopes
- **Broker Integration**: SnapTrade API (aggregator)
- **Market Data**: Multiple vendor support (Polygon, IEX, NASDAQ, Databento)

### Project Structure
```
skybit-cli/
├── cmd/skybit/          # Main entry point
├── internal/
│   ├── app/            # Core application logic
│   ├── broker/         # Broker integrations (SnapTrade)
│   ├── chat/           # LLM interface & message handling
│   ├── market/         # Market data vendor adapters
│   ├── policy/         # Risk management & policy engine
│   ├── tui/            # Terminal UI components
│   │   ├── components/ # Reusable UI widgets
│   │   ├── views/      # Main application views
│   │   └── theme/      # Styling & colors
│   └── ws/             # WebSocket client implementation
├── pkg/                # Public packages
│   ├── types/          # Shared type definitions
│   └── utils/          # Common utilities
├── proto/              # Protocol definitions (JSON schemas)
└── docs/               # Documentation
```

## Development Guidelines

### Code Standards

1. **Go Conventions**
   - Follow standard Go formatting (use `gofmt`)
   - Use meaningful variable names (no single letters except loop indices)
   - Error handling: always check errors, wrap with context
   - Keep functions under 50 lines when possible

2. **TUI Development**
   - Use Bubble Tea's Model-Update-View pattern consistently
   - Component state should be immutable (create new state in Update)
   - Keep view rendering pure (no side effects)
   - Use Lip Gloss for all styling (no hardcoded ANSI codes)

3. **WebSocket Communication**
   - All messages follow the envelope structure
   - Implement automatic reconnection with exponential backoff
   - Handle heartbeats to detect connection health
   - Queue messages during reconnection

### Key Implementation Details

#### WebSocket Protocol
```go
type Envelope struct {
    Type    string      `json:"type"`    // Message type
    ID      string      `json:"id"`      // UUID
    TS      time.Time   `json:"ts"`      // RFC3339 timestamp
    Session string      `json:"session"` // Session UUID
    Data    interface{} `json:"data"`    // Payload
}
```

Message types: `Auth`, `Chat`, `Plan`, `Preview`, `Approve`, `Cancel`, `Status`, `Heartbeat`, `Error`

#### TUI Layout
- **Two-pane layout**: Chat/Preview (left), Context (right)
- **Key bindings**: `A` (Approve), `E` (Edit), `C` (Cancel), `Tab` (switch panes)
- **Status indicators**: Connection state, market data freshness, order status

#### Trading Flow
1. User enters natural language command
2. Backend generates Plan (tool calls)
3. Backend creates Preview with order details
4. User reviews Preview (orders, fees, impact)
5. User approves/edits/cancels
6. System executes and streams status updates

### Safety & Security Considerations

1. **Never log sensitive data** (API keys, session tokens, order details in production)
2. **Validate all user input** before sending to backend
3. **Implement rate limiting** for API calls
4. **Show clear risk warnings** in Preview
5. **Default to paper trading** for new users
6. **Require explicit confirmation** for all trades

### Testing Approach

1. **Unit tests** for all business logic
2. **Integration tests** for WebSocket communication
3. **Golden trace tests** for deterministic trading scenarios
4. **Manual testing** of TUI interactions
5. **Load testing** for concurrent WebSocket connections

### Common Tasks

#### Adding a New Command
1. Define command in `internal/app/commands.go`
2. Add parsing logic in `internal/chat/parser.go`
3. Implement handler in appropriate package
4. Update help text in `internal/tui/components/help.go`
5. Add tests

#### Adding a Market Data Vendor
1. Create adapter in `internal/market/vendors/`
2. Implement `MarketDataProvider` interface
3. Add configuration in `internal/app/config.go`
4. Update vendor selection logic
5. Test with real and mock data

#### Modifying TUI Components
1. Update model in `internal/tui/models/`
2. Modify update logic for state changes
3. Adjust view rendering
4. Test keyboard navigation
5. Verify theme consistency

### Performance Guidelines

1. **Minimize allocations** in hot paths (message handling, rendering)
2. **Use buffered channels** for async communication
3. **Cache market data** with appropriate TTLs
4. **Batch UI updates** to avoid excessive redraws
5. **Profile regularly** with pprof

### Error Handling

1. **User-facing errors**: Clear, actionable messages
2. **Network errors**: Automatic retry with backoff
3. **Trading errors**: Show in Preview, prevent execution
4. **Critical errors**: Log, notify user, graceful shutdown

### Release Process

1. Run full test suite
2. Build binaries for all platforms
3. Sign macOS binary (notarization)
4. Create release notes
5. Update Homebrew formula
6. Tag release in Git

## MCP Tool Integration

The backend provides MCP tools that the CLI interacts with:

- `trade.*` - Order placement and management
- `data.*` - Market data queries
- `risk.*` - Risk checks and validation
- `portfolio.*` - Position and balance queries
- `notify.*` - User notifications

## Development Commands

```bash
# Run locally
go run cmd/skybit/main.go

# Run with mock backend
SKYBIT_ENV=mock go run cmd/skybit/main.go

# Run tests
go test ./...

# Build binary
go build -o skybit cmd/skybit/main.go

# Generate mocks
go generate ./...

# Lint code
golangci-lint run

# Format code
gofmt -w .
```

## Debugging Tips

1. **Enable debug logging**: `SKYBIT_DEBUG=true skybit`
2. **Inspect WebSocket traffic**: Use browser dev tools with test backend
3. **TUI debugging**: Log to file, not stdout (breaks terminal)
4. **Performance**: Use pprof endpoints in debug mode
5. **State issues**: Add state snapshots to debug log

## Important Patterns

### State Management
- Central state store in `internal/app/state.go`
- State updates trigger UI refreshes
- Immutable state updates (no direct mutations)

### Event Handling
- User input → Command → Action → State Update → UI Refresh
- WebSocket message → Handler → State Update → UI Refresh
- Background tasks use channels for communication

### Error Recovery
- Network: Exponential backoff reconnection
- Trading: Rollback on partial failure
- UI: Graceful degradation, show error state

## CI/CD Configuration

GitHub Actions workflow handles:
1. Testing on push/PR
2. Building release binaries
3. Creating GitHub releases
4. Updating Homebrew tap
5. Notifying Discord on release

## Contact & Resources

- Backend API docs: `wss://api.skybit.ai/artifact/v1`
- Design specs: See `/docs/architecture.md`
- Team chat: Internal Slack #skybit-cli
- Issue tracking: GitHub Issues

## Quick Reference

### File Locations
- Config: `~/.skybit/config.yaml`
- Logs: `~/.skybit/logs/`
- Cache: `~/.skybit/cache/`

### Environment Variables
- `SKYBIT_ENV` - Environment (prod/staging/mock)
- `SKYBIT_DEBUG` - Enable debug logging
- `SKYBIT_CONFIG` - Custom config path
- `SKYBIT_WS_URL` - Override WebSocket URL

### Common Issues
1. **Connection refused**: Check firewall, verify API status
2. **Auth failure**: Refresh token with `skybit login`
3. **Slow updates**: Check network latency with `skybit doctor`
4. **UI glitches**: Update terminal, check UTF-8 support