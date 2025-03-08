package config

import (
	"fmt"
	"os"

	"gopkg.in/yaml.v2"
)

type Config struct {
    Upstox struct {
        ClientID     string `yaml:"client_id"`
        ClientSecret string `yaml:"client_secret"`
        RedirectURI  string `yaml:"redirect_uri"`
    } `yaml:"upstox"`
    
    Trading struct {
        InstrumentToken string `yaml:"instrument_token"`
    } `yaml:"trading"`
    
    Strategy struct {
        EMAShort int `yaml:"ema_short"`
        EMALong  int `yaml:"ema_long"`
    } `yaml:"strategy"`
}

func LoadConfig(path string) (*Config, error) {
    file, err := os.Open(path)
    if err != nil {
        return nil, fmt.Errorf("failed to open config file: %w", err)
    }
    defer file.Close()

    var cfg Config
    if err := yaml.NewDecoder(file).Decode(&cfg); err != nil {
        return nil, fmt.Errorf("failed to decode config: %w", err)
    }
    return &cfg, nil
}