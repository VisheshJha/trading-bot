package market

import "trading-bot/internal/models"

type DataHandler struct {
    historicalData []models.MarketData
}

func (dh *DataHandler) Add(data models.MarketData) {
    dh.historicalData = append(dh.historicalData, data)
}

func (dh *DataHandler) GetHistorical(window int) []models.MarketData {
    if window > len(dh.historicalData) {
        return dh.historicalData
    }
    return dh.historicalData[len(dh.historicalData)-window:]
}