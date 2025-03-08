package models

type Position struct {
    Symbol         string
    InstrumentToken string
    Quantity       int
    EntryPrice     float64
    Direction      string // LONG/SHORT
}