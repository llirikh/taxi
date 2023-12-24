package app

import (
	"context"
	"errors"
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
		if err := a.Handler.Server.ListenAndServe(); err != nil && !errors.Is(err, http.ErrServerClosed) {
			doneWithErr <- err
		}
	}()

	err := <-doneWithErr
	if err != nil {
		return err
	}

	return nil
}
