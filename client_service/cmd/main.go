package main

import (
	"client_service/internal/app"
	"client_service/internal/models"
	"client_service/internal/mongodb"
	"context"
	"fmt"
	"log"
)

func main() {
	ctx := context.Background()
	App := app.NewApp()
	err := App.Start(ctx)
	if err != nil {
		log.Fatal(err)
		return
	}
	offer := models.Offer{Price: models.Price{Amount: 345345, Currency: "RUB"}, From: models.Position{Lat: 15, Lng: 16}, To: models.Position{Lat: 1, Lng: 2}, Client_id: "kirill"}
	trip := mongodb.Trip{Offer_id: "request.Offer_id", Client_id: "userID",
		From:   mongodb.Position{Lat: offer.From.Lat, Lng: offer.From.Lng},
		To:     mongodb.Position{Lat: offer.To.Lat, Lng: offer.To.Lng},
		Price:  mongodb.Price{Amount: offer.Price.Amount, Currency: offer.Price.Currency},
		Status: "DRIVER_SEARCH"}
	err = App.Handler.Database.CreateTrip(&trip)
	if err != nil {
		fmt.Println(err)
	}
	trpis, err := App.Handler.Database.GetTripsByUserID("userID")
	if err != nil {
		fmt.Println(err)
	}
	fmt.Println(trpis)
}
