package models

// CollectionPrice is the mongo collection name of prices
const (
	CollectionPrice = "prices"
)

// Price is a record for cryptocurrency price
type Price struct {
	Symbol    string  `json:"symbol"`
	Open      float64 `json:"open"`
	High      float64 `json:"high"`
	Low       float64 `json:"low"`
	Close     float64 `json:"close"`
	Volume    float64 `json:"volume"`
	Time      string  `json:"time"`
	Timestamp int64   `json:"timestamp"`
}

// NewPrice instatiates a new price
func NewPrice() Price {
	return Price{}
}
