package main

import (
	"context"
	"log"
	"offering_service/internal/app"
)

func main() {
	ctx := context.Background()
	App := app.NewApp()
	err := App.Start(ctx)
	if err != nil {
		log.Fatal(err)
		return
	}
}
