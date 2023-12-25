package models

type MongoDBConfig struct {
	URI  string `json:"uri"`
	Name string `json:"name"`
}

type Config struct {
	Port     string        `json:"port"`
	Database MongoDBConfig `json:"database"`
}

type Price struct {
	Amount   float64 `json:"amount"`
	Currency string  `json:"currency"`
}

type Position struct {
	Lat float64 `json:"lat"`
	Lng float64 `json:"lng"`
}

type Offer struct {
	From      Position
	To        Position
	Client_id string
	Price     Price
}
