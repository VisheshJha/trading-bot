package market

type SymbolManager struct {
    allowedSymbols map[string]bool
}

func NewSymbolManager(symbols []string) *SymbolManager {
    sm := &SymbolManager{
        allowedSymbols: make(map[string]bool),
    }
    for _, sym := range symbols {
        sm.allowedSymbols[sym] = true
    }
    return sm
}

func (sm *SymbolManager) IsAllowed(symbol string) bool {
    return sm.allowedSymbols[symbol]
}