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

func (db *Database) GetTripsByUserID(userID string) ([]Trip, error) {
	coll := db.Client.Database(db.Name).Collection("trips")
	ctx, cancel := context.WithTimeout(context.Background(), 5*time.Second)
	defer cancel()
	cursor, err := coll.Find(ctx, bson.M{"client_id": userID})
	if err != nil {
		return nil, err
	}

	var trips []Trip
	err = cursor.All(ctx, &trips)
	if err != nil {
		return nil, err
	}
	fmt.Println("got trips from bd")
	fmt.Println(trips)

	return trips, nil
}

func (db *Database) GetTripByID(tripID string) (*Trip, error) {
	coll := db.Client.Database(db.Name).Collection("trips")
	ctx, cancel := context.WithTimeout(context.Background(), 5*time.Second)
	defer cancel()
	currID, err := primitive.ObjectIDFromHex(tripID)
	if err != nil {
		return nil, err
	}

	var trip Trip
	if err := coll.FindOne(ctx, bson.M{"_id": currID}).Decode(&trip); err != nil {
		return nil, err
	}

	fmt.Println("got trip from bd")
	fmt.Println(trip)

	return &trip, nil
}

func (db *Database) CancelTripByID(tripID string) error {
	coll := db.Client.Database(db.Name).Collection("trips")
	ctx, cancel := context.WithTimeout(context.Background(), 5*time.Second)
	defer cancel()
	currID, err := primitive.ObjectIDFromHex(tripID)
	if err != nil {
		return err
	}

	if _, err := coll.DeleteOne(ctx, bson.M{"_id": currID}); err != nil {
		return err
	}
	fmt.Println("cancelled trip")

	return nil
}

func (db *Database) CreateTrip(trip *Trip) error {
	coll := db.Client.Database(db.Name).Collection("trips")
	ctx, cancel := context.WithTimeout(context.Background(), 5*time.Second)
	defer cancel()

	if _, err := coll.InsertOne(ctx, trip); err != nil {
		return err
	}

	fmt.Println("trip created")

	return nil
}
