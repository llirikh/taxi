package main

import (
	"client_service/internal/app"
	"context"
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
}
