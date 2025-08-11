package tools

import (
    "context"
    "fmt"
    "os"

    "github.com/kirincap/skybit-cli/internal/broker/snaptrade"
)

// Optional account selection from env for MVP
func getDefaultAccountID(ctx context.Context, c *snaptrade.Client) (string, error) {
    if v := os.Getenv("SNAPTRADE_ACCOUNT_ID"); v != "" {
        return v, nil
    }
    accts, err := c.ListAccounts(ctx)
    if err != nil { return "", err }
    if len(accts) == 0 { return "", fmt.Errorf("no accounts linked; run brokers connect") }
    return accts[0].ID, nil
}

func SnapTradeAccounts() Handler {
    return func(ctx context.Context, args map[string]any) (any, error) {
        cfg, err := snaptrade.LoadConfig()
        if err != nil { return nil, err }
        cli := snaptrade.New(cfg)
        accts, err := cli.ListAccounts(ctx)
        if err != nil { return nil, err }
        return map[string]any{"accounts": accts}, nil
    }
}

func SnapTradePositions() Handler {
    return func(ctx context.Context, args map[string]any) (any, error) {
        cfg, err := snaptrade.LoadConfig()
        if err != nil { return nil, err }
        cli := snaptrade.New(cfg)
        acctID, _ := args["account_id"].(string)
        if acctID == "" { acctID, err = getDefaultAccountID(ctx, cli); if err != nil { return nil, err } }
        pos, err := cli.GetPositions(ctx, acctID)
        if err != nil { return nil, err }
        return map[string]any{"positions": pos}, nil
    }
}

func SnapTradePlaceOrder() Handler {
    return func(ctx context.Context, args map[string]any) (any, error) {
        cfg, err := snaptrade.LoadConfig()
        if err != nil { return nil, err }
        cli := snaptrade.New(cfg)
        acctID, _ := args["account_id"].(string)
        if acctID == "" { acctID, err = getDefaultAccountID(ctx, cli); if err != nil { return nil, err } }
        order, _ := args["order"].(map[string]any)
        req := snaptrade.PlaceOrderRequest{
            Symbol:     fmt.Sprintf("%v", order["symbol"]),
            Side:       fmt.Sprintf("%v", order["side"]),
            TIF:        fmt.Sprintf("%v", order["tif"]),
            ClientID:   fmt.Sprintf("%v", order["client_id"]),
        }
        if v, ok := order["qty"].(float64); ok { req.Quantity = v }
        if v, ok := order["limit_price"].(float64); ok { req.LimitPrice = v }
        if v, ok := order["type"].(string); ok { req.Type = v }
        out, err := cli.PlaceOrder(ctx, acctID, req)
        if err != nil { return nil, err }
        return map[string]any{"broker_order_id": out.BrokerOrderID, "status": out.Status}, nil
    }
}

func SnapTradeCancel() Handler {
    return func(ctx context.Context, args map[string]any) (any, error) {
        cfg, err := snaptrade.LoadConfig()
        if err != nil { return nil, err }
        cli := snaptrade.New(cfg)
        acctID, _ := args["account_id"].(string)
        if acctID == "" { acctID, err = getDefaultAccountID(ctx, cli); if err != nil { return nil, err } }
        orderID, _ := args["broker_order_id"].(string)
        if orderID == "" { return nil, fmt.Errorf("broker_order_id required") }
        if err := cli.CancelOrder(ctx, acctID, orderID); err != nil { return nil, err }
        return map[string]any{"status": "CANCELED"}, nil
    }
}

func SnapTradeCancelAll() Handler {
    return func(ctx context.Context, args map[string]any) (any, error) {
        cfg, err := snaptrade.LoadConfig()
        if err != nil { return nil, err }
        cli := snaptrade.New(cfg)
        acctID, _ := args["account_id"].(string)
        if acctID == "" { acctID, err = getDefaultAccountID(ctx, cli); if err != nil { return nil, err } }
        if err := cli.CancelAll(ctx, acctID); err != nil { return nil, err }
        return map[string]any{"status": "CANCELED_ALL"}, nil
    }
}


