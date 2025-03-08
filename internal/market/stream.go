package market

import (
	"encoding/json"
	"fmt"
	"net/http"
	"trading-bot/internal/auth"
	"trading-bot/internal/models"

	"github.com/gorilla/websocket"
)

type StreamClient struct {
    conn        *websocket.Conn
    authClient  *auth.Client 
}

func NewStreamClient(authClient *auth.Client) *StreamClient {
    return &StreamClient{
        authClient: authClient,
    }
}

func (s *StreamClient) Connect() error {
    token, err := s.authClient.Token()
    if err != nil {
        return fmt.Errorf("token retrieval failed: %w", err)
    }

    header := http.Header{}
    header.Add("Authorization", "Bearer "+token.AccessToken)
    header.Add("x-api-key", s.authClient.ClientID())

    conn, _, err := websocket.DefaultDialer.Dial(
        "wss://api.upstox.com/v2/feed/market-data-feed",
        header,
    )
    if err != nil {
        return fmt.Errorf("websocket connection failed: %w", err)
    }
    s.conn = conn
    return nil
}

func (s *StreamClient) Close() error {
    if s.conn != nil {
        return s.conn.Close()
    }
    return nil
}

func (s *StreamClient) ReadMessages() <-chan models.MarketData {
    ch := make(chan models.MarketData, 100)
    go func() {
        defer close(ch)
        for {
            _, message, err := s.conn.ReadMessage()
            if err != nil {
                return
            }
            
            var data models.MarketData
            if err := json.Unmarshal(message, &data); err == nil {
                ch <- data
            }
        }
    }()
    return ch
}

func (s *StreamClient) Subscribe(instruments []string) error {
    subscription := map[string]interface{}{
        "action": "subscribe",
        "params": map[string]interface{}{
            "mode":           "full",
            "instrumentKeys": instruments,
        },
    }
    return s.conn.WriteJSON(subscription)
}

