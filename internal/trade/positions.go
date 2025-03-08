package trade

import "trading-bot/internal/models"

type PositionTracker struct {
    positions map[string]*models.Position
}

func NewPositionTracker() *PositionTracker {
    return &PositionTracker{
        positions: make(map[string]*models.Position),
    }
}

func (pt *PositionTracker) Update(position *models.Position) {
    pt.positions[position.Symbol] = position
}

func (pt *PositionTracker) GetOpenPositions() []*models.Position {
    positions := make([]*models.Position, 0, len(pt.positions))
    for _, pos := range pt.positions {
        positions = append(positions, pos)
    }
    return positions
}