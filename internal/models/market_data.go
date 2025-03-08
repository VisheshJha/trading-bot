package models

type MarketData struct {
    Exchange        string  `json:"exchange"`
    Symbol          string  `json:"symbol"`
    InstrumentToken string  `json:"instrument_token"`
    LTP             float64 `json:"ltp"`
    Open            float64 `json:"open_price"`
    High            float64 `json:"high_price"`
    Low             float64 `json:"low_price"`
    Volume          int64   `json:"volume"`
    Timestamp       int64   `json:"timestamp"`
}