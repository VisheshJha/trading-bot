package strategies

import (
    "trading-bot/internal/models"
)

type EMAStrategy struct {
    shortWindow int
    longWindow  int
    shortEMA    float64
    longEMA     float64
}

func NewEMAStrategy(short, long int) *EMAStrategy {
    return &EMAStrategy{
        shortWindow: short,
        longWindow:  long,
    }
}

func (s *EMAStrategy) Analyze(data models.MarketData) models.TradeSignal {
    s.updateEMA(data.LTP)
    
    if s.shortEMA > s.longEMA {
        return models.TradeSignal{
            Action:          models.Buy,
            Symbol:          data.Symbol,
            InstrumentToken: data.InstrumentToken,
            Quantity:        1,
        }
    }
    
    if s.shortEMA < s.longEMA {
        return models.TradeSignal{
            Action:          models.Sell,
            Symbol:          data.Symbol,
            InstrumentToken: data.InstrumentToken,
            Quantity:        1,
        }
    }
    
    return models.TradeSignal{Action: models.Hold}
}

func (s *EMAStrategy) updateEMA(price float64) {
    if s.shortEMA == 0 {
        s.shortEMA = price
    } else {
        s.shortEMA = (price * 2/float64(s.shortWindow+1)) + (s.shortEMA * (1 - 2/float64(s.shortWindow+1)))
    }

    if s.longEMA == 0 {
        s.longEMA = price
    } else {
        s.longEMA = (price * 2/float64(s.longWindow+1)) + (s.longEMA * (1 - 2/float64(s.longWindow+1)))
    }
}