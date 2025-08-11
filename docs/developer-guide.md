# Skybit CLI Developer Guide

## Getting Started

### Prerequisites

- Go 1.21 or higher
- Git
- Make (optional but recommended)
- Protocol Buffers compiler (for schema generation)

### Initial Setup

```bash
# Clone the repository
git clone https://github.com/skybit/cli.git skybit-cli
cd skybit-cli

# Install dependencies
go mod download

# Install development tools
go install github.com/golangci/golangci-lint/cmd/golangci-lint@latest
go install github.com/vektra/mockery/v2@latest

# Run initial build
go build -o skybit cmd/skybit/main.go

# Verify installation
./skybit --version
```

## Project Structure

### Directory Layout

```
skybit-cli/
├── cmd/skybit/             # Application entry point
│   └── main.go            # Main function, initialization
├── internal/              # Private application code
│   ├── app/              # Core application logic
│   │   ├── app.go        # Main application struct
│   │   ├── commands.go   # Command definitions
│   │   ├── config.go     # Configuration management
│   │   └── state.go      # State management
│   ├── broker/           # Broker integrations
│   │   ├── interface.go  # Broker interface definition
│   │   ├── snaptrade.go  # SnapTrade implementation
│   │   └── mock.go       # Mock broker for testing
│   ├── chat/             # Chat/LLM interface
│   │   ├── parser.go     # Command parsing
│   │   ├── processor.go  # Message processing
│   │   └── history.go    # Chat history management
│   ├── market/           # Market data handling
│   │   ├── interface.go  # Market data interface
│   │   ├── vendors/      # Vendor implementations
│   │   └── cache.go      # Data caching layer
│   ├── policy/           # Risk and policy management
│   │   ├── engine.go     # Policy engine
│   │   ├── rules.go      # Rule definitions
│   │   └── validator.go  # Order validation
│   ├── tui/              # Terminal UI
│   │   ├── app.go        # Bubble Tea app
│   │   ├── components/   # Reusable UI components
│   │   ├── views/        # Application views
│   │   └── theme/        # Styling and colors
│   └── ws/               # WebSocket client
│       ├── client.go     # WebSocket client
│       ├── handler.go    # Message handlers
│       └── reconnect.go  # Reconnection logic
├── pkg/                   # Public packages
│   ├── types/            # Shared type definitions
│   └── utils/            # Utility functions
├── proto/                # Protocol definitions
│   ├── messages.json     # JSON schema definitions
│   └── gen/              # Generated code
├── test/                 # Test resources
│   ├── fixtures/         # Test data
│   └── golden/           # Golden test files
└── scripts/              # Build and utility scripts
```

## Development Workflow

### Building

```bash
# Standard build
go build -o skybit cmd/skybit/main.go

# Build with version info
go build -ldflags "-X main.version=1.0.0" -o skybit cmd/skybit/main.go

# Cross-platform builds
GOOS=linux GOARCH=amd64 go build -o skybit-linux cmd/skybit/main.go
GOOS=darwin GOARCH=arm64 go build -o skybit-mac cmd/skybit/main.go
GOOS=windows GOARCH=amd64 go build -o skybit.exe cmd/skybit/main.go

# Development build with race detector
go build -race -o skybit-dev cmd/skybit/main.go
```

### Testing

```bash
# Run all tests
go test ./...

# Run with coverage
go test -cover ./...

# Run specific package tests
go test ./internal/tui/...

# Run with verbose output
go test -v ./...

# Run integration tests
go test -tags=integration ./...

# Generate coverage report
go test -coverprofile=coverage.out ./...
go tool cover -html=coverage.out
```

### Running Locally

```bash
# Run with default settings
go run cmd/skybit/main.go

# Run with debug logging
SKYBIT_DEBUG=true go run cmd/skybit/main.go

# Run with mock backend
SKYBIT_ENV=mock go run cmd/skybit/main.go

# Run with custom config
SKYBIT_CONFIG=./config.yaml go run cmd/skybit/main.go

# Run with custom WebSocket URL
SKYBIT_WS_URL=ws://localhost:8080 go run cmd/skybit/main.go
```

## Code Guidelines

### Go Best Practices

1. **Error Handling**
```go
// Always check and wrap errors with context
result, err := someFunction()
if err != nil {
    return fmt.Errorf("failed to process: %w", err)
}

// Use custom error types for specific cases
type ValidationError struct {
    Field   string
    Message string
}

func (e ValidationError) Error() string {
    return fmt.Sprintf("validation failed for %s: %s", e.Field, e.Message)
}
```

2. **Context Usage**
```go
// Pass context for cancellation and timeouts
func ProcessOrder(ctx context.Context, order Order) error {
    ctx, cancel := context.WithTimeout(ctx, 5*time.Second)
    defer cancel()
    
    select {
    case <-ctx.Done():
        return ctx.Err()
    case result := <-processAsync(order):
        return result
    }
}
```

3. **Concurrency Patterns**
```go
// Use channels for communication
func worker(jobs <-chan Job, results chan<- Result) {
    for job := range jobs {
        results <- process(job)
    }
}

// Use sync.Once for initialization
var (
    instance *Client
    once     sync.Once
)

func GetClient() *Client {
    once.Do(func() {
        instance = newClient()
    })
    return instance
}
```

### TUI Development with Bubble Tea

1. **Model Structure**
```go
type Model struct {
    // State fields
    orders    []Order
    positions []Position
    
    // UI fields
    viewport viewport.Model
    table    table.Model
    input    textinput.Model
    
    // Control fields
    ready bool
    err   error
}
```

2. **Update Pattern**
```go
func (m Model) Update(msg tea.Msg) (tea.Model, tea.Cmd) {
    switch msg := msg.(type) {
    case tea.KeyMsg:
        switch msg.String() {
        case "q", "ctrl+c":
            return m, tea.Quit
        case "tab":
            m.switchPane()
        }
    case tea.WindowSizeMsg:
        m.viewport.Width = msg.Width
        m.viewport.Height = msg.Height - 4
    case OrderUpdateMsg:
        m.orders = msg.Orders
        m.updateTable()
    }
    return m, nil
}
```

3. **View Rendering**
```go
func (m Model) View() string {
    if !m.ready {
        return "Loading..."
    }
    
    return lipgloss.JoinVertical(
        lipgloss.Left,
        m.headerView(),
        m.contentView(),
        m.footerView(),
    )
}
```

### WebSocket Communication

1. **Message Handling**
```go
func (c *Client) handleMessage(msg []byte) error {
    var env Envelope
    if err := json.Unmarshal(msg, &env); err != nil {
        return fmt.Errorf("unmarshal failed: %w", err)
    }
    
    handler, ok := c.handlers[env.Type]
    if !ok {
        return fmt.Errorf("unknown message type: %s", env.Type)
    }
    
    return handler(env)
}
```

2. **Reconnection Logic**
```go
func (c *Client) connectWithRetry(ctx context.Context) error {
    backoff := time.Second
    maxBackoff := time.Minute
    
    for {
        if err := c.connect(); err == nil {
            return nil
        }
        
        select {
        case <-ctx.Done():
            return ctx.Err()
        case <-time.After(backoff):
            backoff = min(backoff*2, maxBackoff)
        }
    }
}
```

## Adding Features

### Adding a New Command

1. **Define the command** in `internal/app/commands.go`:
```go
type Command struct {
    Name        string
    Description string
    Handler     CommandHandler
    Aliases     []string
}

var commands = []Command{
    {
        Name:        "chart",
        Description: "Display price chart for symbol",
        Handler:     handleChart,
        Aliases:     []string{"c"},
    },
}
```

2. **Implement the handler**:
```go
func handleChart(ctx context.Context, args []string) error {
    if len(args) < 1 {
        return errors.New("symbol required")
    }
    
    symbol := args[0]
    data, err := fetchChartData(ctx, symbol)
    if err != nil {
        return fmt.Errorf("fetch chart data: %w", err)
    }
    
    return displayChart(data)
}
```

3. **Add tests**:
```go
func TestHandleChart(t *testing.T) {
    tests := []struct {
        name    string
        args    []string
        wantErr bool
    }{
        {"valid symbol", []string{"AAPL"}, false},
        {"missing symbol", []string{}, true},
    }
    
    for _, tt := range tests {
        t.Run(tt.name, func(t *testing.T) {
            err := handleChart(context.Background(), tt.args)
            if (err != nil) != tt.wantErr {
                t.Errorf("handleChart() error = %v, wantErr %v", err, tt.wantErr)
            }
        })
    }
}
```

### Adding a Market Data Vendor

1. **Implement the interface** in `internal/market/vendors/`:
```go
type NewVendor struct {
    apiKey string
    client *http.Client
}

func (v *NewVendor) GetQuote(ctx context.Context, symbol string) (*Quote, error) {
    // Implementation
}

func (v *NewVendor) GetOrderBook(ctx context.Context, symbol string) (*OrderBook, error) {
    // Implementation
}

func (v *NewVendor) Subscribe(symbols []string) error {
    // Implementation
}
```

2. **Register the vendor**:
```go
func init() {
    RegisterVendor("newvendor", func(config VendorConfig) (MarketDataProvider, error) {
        return &NewVendor{
            apiKey: config.APIKey,
            client: &http.Client{Timeout: 10 * time.Second},
        }, nil
    })
}
```

## Testing

### Unit Testing

```go
func TestOrderValidator(t *testing.T) {
    validator := NewOrderValidator()
    
    order := Order{
        Symbol:   "AAPL",
        Quantity: 100,
        Side:     "buy",
        Type:     "limit",
        Price:    150.00,
    }
    
    err := validator.Validate(order)
    assert.NoError(t, err)
    
    order.Quantity = -1
    err = validator.Validate(order)
    assert.Error(t, err)
}
```

### Integration Testing

```go
// +build integration

func TestSnapTradeIntegration(t *testing.T) {
    if testing.Short() {
        t.Skip("skipping integration test")
    }
    
    client := NewSnapTradeClient(getTestCredentials())
    
    positions, err := client.GetPositions(context.Background())
    require.NoError(t, err)
    assert.NotNil(t, positions)
}
```

### Golden Tests

```go
func TestPreviewGeneration(t *testing.T) {
    testCases := []string{
        "buy_limit_order",
        "sell_market_order",
        "complex_conditional",
    }
    
    for _, tc := range testCases {
        t.Run(tc, func(t *testing.T) {
            input := loadFixture(t, tc+".input.json")
            expected := loadGolden(t, tc+".golden.json")
            
            actual := generatePreview(input)
            
            if *update {
                saveGolden(t, tc+".golden.json", actual)
            } else {
                assert.JSONEq(t, expected, actual)
            }
        })
    }
}
```

## Debugging

### Debug Logging

```go
import "github.com/skybit/cli/pkg/logger"

func debugFunction() {
    logger.Debug("entering function", "param", value)
    defer logger.Debug("exiting function")
    
    // Log structured data
    logger.Info("processing order",
        "symbol", order.Symbol,
        "quantity", order.Quantity,
        "price", order.Price,
    )
}
```

### Performance Profiling

```go
// Enable profiling endpoint
import _ "net/http/pprof"

func init() {
    if os.Getenv("SKYBIT_PROFILE") == "true" {
        go func() {
            log.Println(http.ListenAndServe("localhost:6060", nil))
        }()
    }
}

// Profile CPU usage
go tool pprof http://localhost:6060/debug/pprof/profile

// Profile memory
go tool pprof http://localhost:6060/debug/pprof/heap
```

### TUI Debugging

```go
// Log to file instead of stdout
func setupDebugLog() {
    if debug := os.Getenv("SKYBIT_DEBUG"); debug != "" {
        logFile, _ := os.Create("/tmp/skybit-debug.log")
        log.SetOutput(logFile)
    }
}

// Debug TUI state
func (m Model) debugState() {
    log.Printf("State: orders=%d, positions=%d, view=%s",
        len(m.orders), len(m.positions), m.currentView)
}
```

## Performance Optimization

### Profiling Commands

```bash
# CPU profiling
go test -cpuprofile=cpu.prof -bench=.
go tool pprof cpu.prof

# Memory profiling
go test -memprofile=mem.prof -bench=.
go tool pprof mem.prof

# Trace execution
go test -trace=trace.out
go tool trace trace.out
```

### Optimization Tips

1. **Reduce Allocations**
```go
// Bad: creates new slice each time
func processItems(items []Item) []Result {
    var results []Result
    for _, item := range items {
        results = append(results, process(item))
    }
    return results
}

// Good: preallocate slice
func processItems(items []Item) []Result {
    results := make([]Result, 0, len(items))
    for _, item := range items {
        results = append(results, process(item))
    }
    return results
}
```

2. **Use Buffers for String Building**
```go
// Bad: string concatenation
func buildMessage(parts []string) string {
    result := ""
    for _, part := range parts {
        result += part + " "
    }
    return result
}

// Good: use strings.Builder
func buildMessage(parts []string) string {
    var builder strings.Builder
    for _, part := range parts {
        builder.WriteString(part)
        builder.WriteByte(' ')
    }
    return builder.String()
}
```

## Release Process

### Version Management

```bash
# Tag a new version
git tag -a v1.0.0 -m "Release version 1.0.0"
git push origin v1.0.0

# Build release binaries
make release VERSION=1.0.0
```

### Release Checklist

- [ ] Update version in `version.go`
- [ ] Update CHANGELOG.md
- [ ] Run full test suite
- [ ] Build for all platforms
- [ ] Test binaries on each platform
- [ ] Create GitHub release
- [ ] Update Homebrew formula
- [ ] Announce on Discord/Twitter

## Troubleshooting

### Common Issues

1. **WebSocket connection fails**
   - Check firewall settings
   - Verify API endpoint is reachable
   - Check for proxy configuration
   - Ensure valid authentication token

2. **TUI rendering issues**
   - Verify terminal supports UTF-8
   - Check terminal size (minimum 80x24)
   - Try different terminal emulator
   - Update terminal fonts

3. **Build failures**
   - Ensure Go version 1.21+
   - Run `go mod tidy`
   - Clear module cache: `go clean -modcache`
   - Check for conflicting dependencies

4. **Test failures**
   - Update test fixtures
   - Check for timing-dependent tests
   - Verify mock implementations
   - Run tests individually to isolate issues