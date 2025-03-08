package api

import (
    "bytes"
    "context"
    "encoding/json"
    "fmt"
    "net/http"
    "time"
    
    "trading-bot/internal/auth"
    "trading-bot/internal/models"
)

type Client struct {
    httpClient  *http.Client
    baseURL     string
    authClient  *auth.Client
}

func NewClient(authClient *auth.Client) *Client {
    return &Client{
        httpClient: &http.Client{Timeout: 10 * time.Second},
        baseURL:    "https://api.upstox.com/v2",
        authClient: authClient,
    }
}

func (c *Client) doRequest(ctx context.Context, method, path string, body interface{}) (*http.Response, error) {
    // Check and refresh token if needed
    if err := c.authClient.RefreshIfNeeded(); err != nil {
        return nil, fmt.Errorf("token refresh failed: %w", err)
    }

    token, err := c.authClient.Token()
    if err != nil {
        return nil, err
    }

    var reqBody []byte
    if body != nil {
        reqBody, err = json.Marshal(body)
        if err != nil {
            return nil, fmt.Errorf("request marshal error: %w", err)
        }
    }

    req, err := http.NewRequestWithContext(ctx, method, c.baseURL+path, bytes.NewBuffer(reqBody))
    if err != nil {
        return nil, fmt.Errorf("request creation failed: %w", err)
    }

    req.Header.Set("Authorization", "Bearer "+token.AccessToken)
    req.Header.Set("x-api-key", c.authClient.ClientID())
    req.Header.Set("Content-Type", "application/json")

    return c.httpClient.Do(req)
}

func (c *Client) PlaceOrder(ctx context.Context, order models.Order) (*models.OrderResponse, error) {
    orderReq := map[string]interface{}{
        "instrument_token": order.InstrumentToken,
        "quantity":        order.Quantity,
        "order_type":      order.OrderType,
        "product":         order.Product,
        "transaction_type": order.Direction,
    }

    resp, err := c.doRequest(ctx, "POST", "/order/place", orderReq)
    if err != nil {
        return nil, err
    }
    defer resp.Body.Close()

    if resp.StatusCode >= 400 {
        return nil, fmt.Errorf("api error: %s", resp.Status)
    }

    var orderResp models.OrderResponse
    if err := json.NewDecoder(resp.Body).Decode(&orderResp); err != nil {
        return nil, fmt.Errorf("response decode failed: %w", err)
    }

    return &orderResp, nil
}