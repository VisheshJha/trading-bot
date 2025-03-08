package trade

import "trading-bot/internal/models"

type PositionManager struct {
    positions map[string]*models.Position
}

func NewPositionManager() *PositionManager {
    return &PositionManager{
        positions: make(map[string]*models.Position),
    }
}

func (pm *PositionManager) Add(position *models.Position) {
    pm.positions[position.Symbol] = position
}

func (pm *PositionManager) Get(symbol string) *models.Position {
    return pm.positions[symbol]
}

func (pm *PositionManager) Remove(symbol string) {
    delete(pm.positions, symbol)
}