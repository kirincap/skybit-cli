package snaptrade

import (
    "bytes"
    "context"
    "encoding/json"
    "fmt"
    "net/http"
    "time"
)

type Client struct {
    http  *http.Client
    cfg   Config
    base  string
}

func New(cfg Config) *Client {
    base := "https://api.sandbox.snaptrade.com/api/v1"
    if cfg.Env == "production" { base = "https://api.snaptrade.com/api/v1" }
    return &Client{http: &http.Client{Timeout: 15 * time.Second}, cfg: cfg, base: base}
}

func (c *Client) do(ctx context.Context, method, path string, body any, out any) error {
    var buf *bytes.Reader
    if body != nil {
        b, err := json.Marshal(body)
        if err != nil { return err }
        buf = bytes.NewReader(b)
    } else {
        buf = bytes.NewReader(nil)
    }
    req, err := http.NewRequestWithContext(ctx, method, c.base+path, buf)
    if err != nil { return err }
    req.Header.Set("Content-Type", "application/json")
    req.Header.Set("X-Snaptrade-Client-Id", c.cfg.ClientID)
    req.Header.Set("X-Snaptrade-Client-Secret", c.cfg.ClientSecret)
    resp, err := c.http.Do(req)
    if err != nil { return err }
    defer resp.Body.Close()
    if resp.StatusCode/100 != 2 {
        return fmt.Errorf("snaptrade http %d", resp.StatusCode)
    }
    if out != nil {
        return json.NewDecoder(resp.Body).Decode(out)
    }
    return nil
}

// Minimal shapes for MVP
type Account struct {
    ID   string `json:"id"`
    Name string `json:"name"`
}

type Position struct {
    Symbol   string  `json:"symbol"`
    Quantity float64 `json:"quantity"`
    AvgPrice float64 `json:"avg_price"`
}

func (c *Client) ListAccounts(ctx context.Context) ([]Account, error) {
    var out []Account
    err := c.do(ctx, http.MethodGet, "/accounts", nil, &out)
    return out, err
}

func (c *Client) GetPositions(ctx context.Context, accountID string) ([]Position, error) {
    var out []Position
    err := c.do(ctx, http.MethodGet, fmt.Sprintf("/accounts/%s/positions", accountID), nil, &out)
    return out, err
}

type PlaceOrderRequest struct {
    Symbol     string  `json:"symbol"`
    Side       string  `json:"side"`
    Quantity   float64 `json:"quantity"`
    Type       string  `json:"type"`
    LimitPrice float64 `json:"limit_price,omitempty"`
    TIF        string  `json:"tif"`
    ClientID   string  `json:"client_id"`
}

type PlaceOrderResponse struct {
    BrokerOrderID string `json:"broker_order_id"`
    Status        string `json:"status"`
}

func (c *Client) PlaceOrder(ctx context.Context, accountID string, req PlaceOrderRequest) (PlaceOrderResponse, error) {
    var out PlaceOrderResponse
    err := c.do(ctx, http.MethodPost, fmt.Sprintf("/accounts/%s/orders", accountID), req, &out)
    return out, err
}

func (c *Client) CancelOrder(ctx context.Context, accountID, brokerOrderID string) error {
    return c.do(ctx, http.MethodPost, fmt.Sprintf("/accounts/%s/orders/%s/cancel", accountID, brokerOrderID), nil, nil)
}

func (c *Client) CancelAll(ctx context.Context, accountID string) error {
    return c.do(ctx, http.MethodPost, fmt.Sprintf("/accounts/%s/orders/cancel_all", accountID), nil, nil)
}


