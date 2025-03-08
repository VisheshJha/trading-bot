package auth

import (
	"context"
	"fmt"
	"time"

	"golang.org/x/oauth2"
)

type Client struct {
    config  *oauth2.Config
    token   *oauth2.Token
    storage TokenStorage
}

func NewClient(clientID, clientSecret, redirectURL string, storage TokenStorage) *Client {
    return &Client{
        config: &oauth2.Config{
            ClientID:     clientID,
            ClientSecret: clientSecret,
            RedirectURL:  redirectURL,
            Endpoint: oauth2.Endpoint{
                AuthURL:  "https://api.upstox.com/v2/login/authorization/dialog",
                TokenURL: "https://api.upstox.com/v2/login/authorization/token",
            },
        },
        storage: storage,
    }
}

func (c *Client) GetAuthorizationURL() string {
    return c.config.AuthCodeURL("state-token", 
        oauth2.AccessTypeOffline,
        oauth2.SetAuthURLParam("prompt", "consent"),
    )
}

func (c *Client) ExchangeCode(code string) error {
    token, err := c.config.Exchange(context.Background(), code)
    if err != nil {
        return fmt.Errorf("token exchange failed: %w", err)
    }
    
    if token.RefreshToken == "" {
        return fmt.Errorf("no refresh token received")
    }
    
    c.token = token
    return c.storage.Save(token)
}

func (c *Client) TokenSource() oauth2.TokenSource {
    return c.config.TokenSource(context.Background(), c.token)
}

func (c *Client) RefreshIfNeeded() error {
    if c.token == nil || c.token.RefreshToken == "" {
        return fmt.Errorf("no refresh token available")
    }
    
    if !c.token.Expiry.Before(time.Now().Add(5 * time.Minute)) {
        return nil
    }

    newToken, err := c.config.TokenSource(context.Background(), c.token).Token()
    if err != nil {
        return fmt.Errorf("token refresh failed: %w", err)
    }
    
    c.token = newToken
    return c.storage.Save(newToken)
}

func (c *Client) Token() (*oauth2.Token, error) {
    if c.token == nil {
        return nil, fmt.Errorf("no token available")
    }
    return c.token, nil
}

func (c *Client) ClientID() string {
    return c.config.ClientID
}

func (c *Client) SetToken(token *oauth2.Token) {
    c.token = token
}