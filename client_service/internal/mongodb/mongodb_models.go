package mongodb

import "go.mongodb.org/mongo-driver/bson/primitive"

type Position struct {
	Lat float64 `bson:"lat"`
	Lng float64 `bson:"lng"`
}

type Price struct {
	Amount   float64 `bson:"amount"`
	Currency string  `bson:"currency"`
}

type Trip struct {
	ID        primitive.ObjectID `bson:"_id,omitempty"`
	Offer_id  string             `bson:"offer_id"`
	Client_id string             `bson:"client_id"`
	From      Position           `bson:"from"`
	To        Position           `bson:"to"`
	Price     Price              `bson:"price"`
	Status    string             `bson:"status"`
}
