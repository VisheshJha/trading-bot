package models

type TradeAction string

const (
    Buy  TradeAction = "BUY"
    Sell TradeAction = "SELL"
    Hold TradeAction = "HOLD"
)

type TradeSignal struct {
    Action          TradeAction
    Symbol          string
    InstrumentToken string
    Quantity        int
}