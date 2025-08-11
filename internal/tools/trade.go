package tools

import (
    "context"
    "fmt"
    "time"
)

// trade.preview: compute simple totals and fake fees
func TradePreview() Handler {
    return func(ctx context.Context, args map[string]any) (any, error) {
        ordersAny, _ := args["orders"].([]any)
        var total float64
        outOrders := make([]map[string]any, 0, len(ordersAny))
        for _, oi := range ordersAny {
            o, _ := oi.(map[string]any)
            qty, _ := o["qty"].(float64)
            limit, _ := o["limit_price"].(float64)
            val := qty * limit
            total += val
            oo := map[string]any{}
            for k, v := range o { oo[k] = v }
            oo["estimated_total"] = val
            outOrders = append(outOrders, oo)
        }
        impact := map[string]any{
            "total_value": total,
            "commission": 1.00,
            "fees": []any{
                map[string]any{"type": "SEC", "amount": 0.46},
                map[string]any{"type": "TAF", "amount": 0.01},
            },
            "total_cost": total + 1.47,
            "slippage_bps": 0,
            "pnl_impact": 0,
        }
        return map[string]any{
            "orders": outOrders,
            "impact": impact,
        }, nil
    }
}

// trade.place_order: stub accepted
func TradePlaceOrder() Handler {
    return func(ctx context.Context, args map[string]any) (any, error) {
        order, _ := args["order"].(map[string]any)
        return map[string]any{
            "broker_order_id": fmt.Sprintf("SNAP-%d", time.Now().UnixNano()),
            "status":          "ACCEPTED",
            "submitted_at":    time.Now().UTC().Format(time.RFC3339),
            "echo":            order,
        }, nil
    }
}

func TradeCancel() Handler {
    return func(ctx context.Context, args map[string]any) (any, error) {
        return map[string]any{"status": "CANCELED"}, nil
    }
}

func TradeCancelAll() Handler {
    return func(ctx context.Context, args map[string]any) (any, error) {
        return map[string]any{"status": "CANCELED_ALL"}, nil
    }
}


