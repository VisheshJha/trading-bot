package auth

import (
	"log"
	"time"
)

type TokenManager struct {
    client  *Client
    ticker  *time.Ticker
}

func NewTokenManager(client *Client) *TokenManager {
    return &TokenManager{
        client: client,
        ticker: time.NewTicker(30 * time.Minute),
    }
}

func (m *TokenManager) StartAutoRefresh() {
    go func() {
        for range m.ticker.C {
            if err := m.client.RefreshIfNeeded(); err != nil {
                log.Printf("Auto token refresh failed: %v", err)
            }
        }
    }()
}

func (m *TokenManager) Stop() {
    m.ticker.Stop()
}

func (m *TokenManager) LoadInitialToken() error {
    token, err := m.client.storage.Load()
    if err != nil {
        return err
    }
    m.client.token = token
    return nil
}