package main

import (
	"context"
	"fmt"
	"log"
	"trip/internal/app"
)

func main() {
	ctx := context.Background()
	newApp := app.NewApp()
	err := newApp.Start(ctx)
	for i := 0; i < 50; i++ {
		fmt.Println("START")
	}
	if err != nil {
		log.Fatalln(err)
		return
	}
}
