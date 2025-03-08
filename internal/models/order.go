package models

type Order struct {
    InstrumentToken string  `json:"instrument_token"`
    Quantity        int     `json:"quantity"`
    OrderType       string  `json:"order_type"`
    Product         string  `json:"product"`
    Direction       string  `json:"transaction_type"`
    Price           float64 `json:"price,omitempty"`
}

type OrderResponse struct {
    OrderID  string `json:"order_id"`
    Message  string `json:"message"`
    Status   string `json:"status"`
}