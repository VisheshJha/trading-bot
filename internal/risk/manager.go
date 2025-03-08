package risk

import (
    "sync"
)

type Manager struct {
    mu           sync.Mutex
    maxDailyLoss float64
    stopLoss     float64
    dailyPnL     float64
}

func NewManager(maxDailyLoss, stopLoss float64) *Manager {
    return &Manager{
        maxDailyLoss: maxDailyLoss,
        stopLoss:     stopLoss,
    }
}

func (rm *Manager) AllowTrade(symbol string, quantity int, price float64) bool {
    rm.mu.Lock()
    defer rm.mu.Unlock()
    
    potentialLoss := float64(quantity) * price * rm.stopLoss
    return (rm.dailyPnL - potentialLoss) > -rm.maxDailyLoss
}

func (rm *Manager) UpdatePnL(change float64) {
    rm.mu.Lock()
    defer rm.mu.Unlock()
    rm.dailyPnL += change
}