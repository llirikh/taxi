package app

import (
	"context"
	"errors"
	"fmt"
	"net/http"
	"offering_service/internal/api/handlers"
	"os/signal"
	"syscall"
)

type App struct {
	Handler *handlers.OfferingHandler
}

func NewApp() *App {
	offeringHandler := handlers.NewHandler()
	app := App{Handler: offeringHandler}
	return &app
}

func (a *App) Start(ctx context.Context) error {
	ctx, stop := signal.NotifyContext(context.Background(), syscall.SIGINT, syscall.SIGTERM)
	defer stop()

	doneWithErr := make(chan error)

	go func() {
		fmt.Println("stst")
		if err := a.Handler.Server.ListenAndServe(); err != nil && !errors.Is(err, http.ErrServerClosed) {
			doneWithErr <- err
		}
	}()

	fmt.Println("aaaaaaaaaaa")
	err := <-doneWithErr
	if err != nil {
		return err
	}

	return nil
}
