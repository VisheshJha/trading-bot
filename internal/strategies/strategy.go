package strategies

import "trading-bot/internal/models"

type Strategy interface {
    Analyze(data models.MarketData) models.TradeSignal
}