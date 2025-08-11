# Skybit CLI Architecture

## System Overview

The Skybit CLI is a terminal-based trading interface that connects to a cloud-based trading engine via WebSocket. The architecture follows a clean separation between the presentation layer (TUI), business logic, and external integrations.

```
┌─────────────────────────────────────────────────────────┐
│                     User Terminal                        │
├─────────────────────────────────────────────────────────┤
│                    Skybit CLI (Go)                       │
│  ┌─────────────┐  ┌──────────────┐  ┌──────────────┐  │
│  │     TUI     │  │   Business   │  │   Network    │  │
│  │  (Bubble    │──│    Logic     │──│   (WebSocket │  │
│  │   Tea)      │  │              │  │    Client)   │  │
│  └─────────────┘  └──────────────┘  └──────────────┘  │
└────────────────────────────┬────────────────────────────┘
                             │ WSS
                             ▼
┌─────────────────────────────────────────────────────────┐
│                  Skybit Backend (Cloud)                  │
│  ┌─────────────┐  ┌──────────────┐  ┌──────────────┐  │
│  │  WebSocket  │  │     LLM      │  │     MCP      │  │
│  │   Gateway   │──│   Processor  │──│    Tools     │  │
│  └─────────────┘  └──────────────┘  └──────────────┘  │
│         │                                    │          │
│         ▼                                    ▼          │
│  ┌──────────────┐                   ┌──────────────┐   │
│  │   SnapTrade  │                   │ Market Data  │   │
│  │   (Broker)   │                   │   Vendors    │   │
│  └──────────────┘                   └──────────────┘   │
└─────────────────────────────────────────────────────────┘
```

## Component Architecture

### 1. Terminal UI Layer

**Technology**: Charmbracelet Bubble Tea framework

**Components**:
- **Main Model**: Central state management and event routing
- **Views**: Chat, Preview, Positions, Orders, Market Data
- **Components**: Reusable widgets (tables, inputs, buttons)
- **Theme**: Consistent styling using Lip Gloss

**Responsibilities**:
- Render interactive terminal interface
- Handle keyboard/mouse input
- Display real-time updates
- Manage view transitions

### 2. Business Logic Layer

**Core Modules**:
- **Command Processor**: Parse and validate user commands
- **State Manager**: Maintain application state
- **Session Manager**: Handle authentication and session lifecycle
- **Policy Engine**: Apply trading rules and risk checks locally
- **Cache Manager**: Store frequently accessed data

**Responsibilities**:
- Coordinate between UI and network layers
- Maintain consistent application state
- Apply business rules
- Handle offline capabilities

### 3. Network Layer

**WebSocket Client**:
- Persistent connection management
- Automatic reconnection with exponential backoff
- Message queuing during disconnection
- Heartbeat monitoring
- Compression support

**Protocol Handler**:
- Message serialization/deserialization
- Request/response correlation
- Event stream processing
- Error handling and recovery

### 4. Backend Integration

**API Gateway**:
- WebSocket endpoint: `wss://api.skybit.ai/artifact/v1`
- REST endpoints for non-streaming operations
- Authentication via device-code flow
- Rate limiting and throttling

**MCP Tool Suite**:
- Trading tools (`trade.*`)
- Market data tools (`data.*`)
- Risk management tools (`risk.*`)
- Portfolio tools (`portfolio.*`)
- Notification tools (`notify.*`)

## Data Flow

### Trading Flow

```
User Input → Command Parser → LLM Processing → Plan Generation
    ↓            ↓                ↓                ↓
Terminal    Validation      Backend API      MCP Tools
    ↓            ↓                ↓                ↓
Chat View   Local Policy    Preview Gen    Risk Checks
    ↓            ↓                ↓                ↓
Display     User Review     Approval      Order Submit
    ↓            ↓                ↓                ↓
Status      Confirmation    Execution     Fill Updates
```

### Real-time Data Flow

```
Market Data Vendors → Backend Aggregator → WebSocket Stream
                           ↓                    ↓
                    Normalization         Client Cache
                           ↓                    ↓
                      Filtering          TUI Updates
                           ↓                    ↓
                    Rate Limiting      User Display
```

## Message Protocol

### Envelope Structure

```json
{
  "type": "MessageType",
  "id": "uuid-v4",
  "ts": "2024-01-15T10:30:00Z",
  "session": "session-uuid",
  "data": {}
}
```

### Message Types

| Type | Direction | Purpose |
|------|-----------|---------|
| Auth | Client→Server | Authenticate session |
| Chat | Client→Server | Send user command |
| Plan | Server→Client | Show execution plan |
| Preview | Server→Client | Display order preview |
| Approve | Client→Server | Confirm order execution |
| Cancel | Client→Server | Cancel order/preview |
| Status | Server→Client | Order status update |
| Heartbeat | Bidirectional | Connection health check |
| Error | Server→Client | Error notification |

## State Management

### Client State

```go
type AppState struct {
    Session      SessionState
    Portfolio    PortfolioState
    MarketData   MarketDataState
    Orders       OrderState
    UI           UIState
    Config       ConfigState
}
```

**State Updates**:
- Immutable updates (no direct mutations)
- Event-driven state changes
- Centralized state store
- Subscription-based UI updates

### Session State

- Authentication status
- Connection health
- User preferences
- Active subscriptions

### Portfolio State

- Current positions
- Account balances
- P&L calculations
- Historical trades

## Security Architecture

### Authentication

1. **Device Code Flow**:
   - CLI requests device code
   - User authorizes in browser
   - CLI polls for token
   - Token stored in memory (optional keychain)

2. **Session Management**:
   - Short-lived tokens (15 minutes)
   - Automatic refresh before expiry
   - Secure token storage
   - Session isolation

### Data Security

- **Transport**: TLS 1.3 for all connections
- **Storage**: Encrypted local cache
- **Secrets**: Never logged or displayed
- **Audit**: Tamper-evident logging

### Risk Controls

- **Client-side validation** before submission
- **Server-side policy engine** for enforcement
- **Rate limiting** on all operations
- **Circuit breakers** for system protection

## Performance Considerations

### Optimization Strategies

1. **Rendering**:
   - Differential updates only
   - Virtual viewport for large lists
   - Debounced UI updates
   - Lazy loading of data

2. **Network**:
   - Message compression
   - Binary protocol option
   - Request batching
   - Connection pooling

3. **Memory**:
   - Bounded caches
   - Automatic cleanup
   - Memory pooling
   - Garbage collection tuning

### Latency Targets

| Operation | Target | Maximum |
|-----------|--------|---------|
| Keystroke→Display | 16ms | 50ms |
| Command→Response | 100ms | 500ms |
| Quote Update | 50ms | 200ms |
| Order Submit | 200ms | 1000ms |

## Scalability

### Horizontal Scaling

- Stateless CLI design
- Backend load balancing
- Geographic distribution
- CDN for static assets

### Vertical Scaling

- Configurable resource limits
- Adaptive quality settings
- Progressive enhancement
- Graceful degradation

## Monitoring & Observability

### Metrics

- Connection health
- Message latency
- Error rates
- User actions
- Performance metrics

### Logging

- Structured logging (JSON)
- Log levels (DEBUG, INFO, WARN, ERROR)
- Correlation IDs
- Audit trail

### Tracing

- Distributed tracing support
- Request correlation
- Performance profiling
- Error tracking

## Deployment Architecture

### Client Distribution

- **Platforms**: macOS, Linux, Windows, BSD
- **Package Managers**: Homebrew, APT, YUM
- **Direct Download**: GitHub Releases
- **Auto-update**: Built-in updater

### Backend Deployment

- **Infrastructure**: AWS/GCP/Azure
- **Orchestration**: Kubernetes
- **Service Mesh**: Istio/Linkerd
- **Observability**: Prometheus/Grafana

## Disaster Recovery

### Failure Modes

1. **Network Failure**:
   - Automatic reconnection
   - Message queuing
   - Offline mode
   - State recovery

2. **Backend Failure**:
   - Fallback endpoints
   - Degraded mode
   - Read-only access
   - Status page integration

3. **Data Corruption**:
   - Checksums validation
   - State snapshots
   - Rollback capability
   - Manual recovery tools

## Future Architecture Considerations

### Planned Enhancements

1. **Plugin System**:
   - Custom indicators
   - Strategy modules
   - Third-party integrations

2. **Multi-Account**:
   - Account switching
   - Aggregate views
   - Cross-account strategies

3. **Advanced Features**:
   - Options trading
   - Crypto integration
   - International markets
   - Social features

### Technology Roadmap

- WebAssembly for performance-critical paths
- gRPC for binary protocol option
- GraphQL for flexible data queries
- WebRTC for ultra-low latency