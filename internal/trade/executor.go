package trade

import (
    "context"
    "trading-bot/internal/api"
    "trading-bot/internal/models"
)

type Executor struct {
    client *api.Client
}

func NewExecutor(client *api.Client) *Executor {
    return &Executor{client: client}
}

func (e *Executor) PlaceOrder(ctx context.Context, order models.Order) (*models.OrderResponse, error) {
    return e.client.PlaceOrder(ctx, order)
}