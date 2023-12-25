package models

type Config struct {
	// todo: db
}

// todo: db_struct

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

type Event struct {
	Id              string            `json:"id"`
	Source          string            `json:"source"`
	Type            string            `json:"type"`
	DataContentType string            `json:"datacontenttype"`
	Time            string            `json:"time"`
	Data            map[string]string `json:"data"`
}
