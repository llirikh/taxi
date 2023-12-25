package mongodb

import (
	"context"
	"fmt"
	"go.mongodb.org/mongo-driver/bson"
	"go.mongodb.org/mongo-driver/bson/primitive"
	"go.mongodb.org/mongo-driver/mongo"
	"go.mongodb.org/mongo-driver/mongo/options"
	"time"
)

type Database struct {
	Client *mongo.Client
	Name   string
}

func NewDatabase(uri string, name string) (*Database, error) {
	ctx, cancel := context.WithTimeout(context.Background(), 10*time.Second)
	defer cancel()
	opts := options.Client()
	opts.ApplyURI(uri)

	client, err := mongo.Connect(ctx, opts)
	if err != nil {
		return nil, err
	}

	return &Database{Client: client, Name: name}, nil
}

func (db *Database) Close() {
	if db.Client != nil {
		if err := db.Client.Disconnect(context.Background()); err != nil {
			fmt.Println(err)
		}
	}
}

func (db *Database) GetTripsByUserID(ctx context.Context, userID string) ([]Trip, error) {
	coll := db.Client.Database(db.Name).Collection("trips")

	cursor, err := coll.Find(ctx, bson.M{"client_id": userID})
	if err != nil {
		return nil, err
	}

	var trips []Trip
	err = cursor.All(ctx, &trips)
	if err != nil {
		return nil, err
	}

	return trips, nil
}

func (db *Database) GetTripByID(ctx context.Context, tripID string) (*Trip, error) {
	coll := db.Client.Database(db.Name).Collection("trips")
	currID, err := primitive.ObjectIDFromHex(tripID)
	if err != nil {
		return nil, err
	}

	var trip Trip
	if err := coll.FindOne(ctx, bson.M{"_id": currID}).Decode(&trip); err != nil {
		return nil, err
	}

	return &trip, nil
}

func (db *Database) CancelTripByID(ctx context.Context, tripID string) error {
	coll := db.Client.Database(db.Name).Collection("trips")
	currID, err := primitive.ObjectIDFromHex(tripID)
	if err != nil {
		return err
	}

	if _, err := coll.DeleteOne(ctx, bson.M{"_id": currID}); err != nil {
		return err
	}

	return nil
}

func (db *Database) CreateTrip(ctx context.Context, trip *Trip) error {
	coll := db.Client.Database(db.Name).Collection("trips")

	if _, err := coll.InsertOne(ctx, trip); err != nil {
		return err
	}

	return nil
}
