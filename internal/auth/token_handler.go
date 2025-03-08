package auth

import (
	"encoding/json"
	"fmt"
	"os"
	"time"

	"golang.org/x/oauth2"
)

type TokenStorage interface {
    Save(token *oauth2.Token) error
    Load() (*oauth2.Token, error)
}

type FileTokenStorage struct {
    Path string
}

func (s *FileTokenStorage) Save(token *oauth2.Token) error {
    data, err := json.Marshal(token)
    if err != nil {
        return fmt.Errorf("token marshal error: %w", err)
    }
    return os.WriteFile(s.Path, data, 0600)
}

func (s *FileTokenStorage) Load() (*oauth2.Token, error) {
    data, err := os.ReadFile(s.Path)
    if err != nil {
        return nil, fmt.Errorf("token read error: %w", err)
    }
    
    var token oauth2.Token
    if err := json.Unmarshal(data, &token); err != nil {
        return nil, fmt.Errorf("token unmarshal error: %w", err)
    }
    
    if token.Expiry.Before(time.Now()) {
        return nil, fmt.Errorf("token expired")
    }
    
    return &token, nil
}

func NewFileTokenStorage(path string) *FileTokenStorage {
    return &FileTokenStorage{Path: path}
}