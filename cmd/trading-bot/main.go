package main

import (
	"context"
	"flag"
	"fmt"
	"log/slog"
	"net/http"
	"os"
	"os/signal"
	"syscall"
	"time"

	"trading-bot/internal/api"
	"trading-bot/internal/auth"
	"trading-bot/internal/config"
	"trading-bot/internal/market"
	"trading-bot/internal/models"
	"trading-bot/internal/strategies"
	"trading-bot/internal/trade"
)

var (
	logLevel = flag.String("log-level", "info", "Log level (debug, info, warn, error)")
)

func main() {
	flag.Parse()
	initLogger(*logLevel)

	slog.Info("Starting trading bot...")
	defer slog.Info("Application shutdown complete")

	// Load configuration
	cfg, err := config.LoadConfig("configs/config.yaml")
	if err != nil {
		slog.Error("Failed to load configuration", "error", err)
		os.Exit(1)
	}
	slog.Debug("Configuration loaded", "config", cfg)

	// Initialize auth components
	tokenStorage := auth.NewFileTokenStorage("token.json")
	authClient := auth.NewClient(
		cfg.Upstox.ClientID,
		cfg.Upstox.ClientSecret,
		cfg.Upstox.RedirectURI,
		tokenStorage,
	)

	// Try to load existing token
    token, err := tokenStorage.Load()
    if err != nil {
        if os.IsNotExist(err) {
            slog.Info("No existing token found, starting OAuth flow")
        } else {
            slog.Warn("Token loading failed", "error", err)
        }
        startOAuthServer(authClient)
    } else {
        authClient.SetToken(token)
        slog.Info("Loaded existing authentication token")
        
        // Verify token validity
        if token.Expiry.Before(time.Now()) {
            slog.Warn("Loaded token expired, starting new OAuth flow")
            startOAuthServer(authClient)
        }
    }

	// Setup token auto-refresh
	tokenManager := auth.NewTokenManager(authClient)
	tokenManager.StartAutoRefresh()
	defer tokenManager.Stop()

    
	// Initialize API clients
	apiClient := api.NewClient(authClient)
	streamClient := market.NewStreamClient(authClient)

	// Connect to market data
	if err := streamClient.Connect(); err != nil {
		slog.Error("Failed to connect to market data stream", "error", err)
		os.Exit(1)
	}
	defer streamClient.Close()

	if err := streamClient.Subscribe([]string{cfg.Trading.InstrumentToken}); err != nil {
		slog.Error("Failed to subscribe to instrument", "instrument", cfg.Trading.InstrumentToken, "error", err)
		os.Exit(1)
	}

	// Initialize trading components
	strategy := strategies.NewEMAStrategy(cfg.Strategy.EMAShort, cfg.Strategy.EMALong)
	executor := trade.NewExecutor(apiClient)

	// Setup signal handling
	sigChan := make(chan os.Signal, 1)
	signal.Notify(sigChan, syscall.SIGINT, syscall.SIGTERM)

	// Start processing loop
	processMarketData(streamClient.ReadMessages(), strategy, executor, sigChan)
}

func startOAuthServer(authClient *auth.Client) {
	server := &http.Server{Addr: ":8080"}
	http.HandleFunc("/callback", func(w http.ResponseWriter, r *http.Request) {
		code := r.URL.Query().Get("code")
		if err := authClient.ExchangeCode(code); err != nil {
			slog.Error("Authentication failed", "error", err)
			fmt.Fprint(w, "Authentication failed! Check logs for details.")
			return
		}
		fmt.Fprint(w, "Authentication successful! Starting trading system...")
		go func() {
			if err := server.Shutdown(context.Background()); err != nil {
				slog.Error("Failed to shutdown auth server", "error", err)
			}
		}()
	})

	go func() {
		slog.Info("Starting OAuth server", "url", authClient.GetAuthorizationURL())
		if err := server.ListenAndServe(); err != http.ErrServerClosed {
			slog.Error("OAuth server failed", "error", err)
		}
	}()
}

func processMarketData(
	dataChan <-chan models.MarketData,
	strategy strategies.Strategy,
	executor *trade.Executor,
	sigChan <-chan os.Signal,
) {
	slog.Info("Bot initialized. Waiting for market data...")

loop:
	for {
		select {
		case data, ok := <-dataChan:
			if !ok {
				slog.Warn("Market data channel closed")
				break loop
			}
			handleMarketData(data, strategy, executor)

		case <-sigChan:
			slog.Info("Shutdown signal received")
			break loop
		}
	}
}

func handleMarketData(
	data models.MarketData,
	strategy strategies.Strategy,
	executor *trade.Executor,
) {
	slog.Debug("Received market data", "data", data)

	signal := strategy.Analyze(data)
	if signal.Action == models.Hold {
		return
	}

	slog.Info("Generating trade signal",
		"action", signal.Action,
		"instrument", data.InstrumentToken,
		"quantity", signal.Quantity,
	)

	_, err := executor.PlaceOrder(context.Background(), models.Order{
		InstrumentToken: data.InstrumentToken,
		Quantity:        signal.Quantity,
		OrderType:       "MARKET",
		Product:         "I",
		Direction:       string(signal.Action),
	})

	if err != nil {
		slog.Error("Failed to execute order",
			"error", err,
			"instrument", data.InstrumentToken,
			"action", signal.Action,
		)
	}
}

func initLogger(level string) {
	var logLevel slog.Level
	switch level {
	case "debug":
		logLevel = slog.LevelDebug
	case "warn":
		logLevel = slog.LevelWarn
	case "error":
		logLevel = slog.LevelError
	default:
		logLevel = slog.LevelInfo
	}

	handler := slog.NewTextHandler(os.Stderr, &slog.HandlerOptions{
		Level: logLevel,
		ReplaceAttr: func(groups []string, a slog.Attr) slog.Attr {
			if a.Key == slog.TimeKey {
				return slog.Attr{}
			}
			return a
		},
	})
	slog.SetDefault(slog.New(handler))
}
