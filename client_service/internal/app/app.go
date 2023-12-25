package app

import (
	"client_service/internal/api/handlers"
	"client_service/internal/config"
	"client_service/internal/mongodb"
	"context"
	"errors"
	"fmt"
	"log"
	"net/http"
	"os/signal"
	"syscall"
)

type App struct {
	Handler *handlers.ClientHandler
}

func NewApp() *App {
	cfg, err := config.InitConfig()
	if err != nil {
		//loging
	}

	db, err := mongodb.NewDatabase(cfg.Database.URI, cfg.Database.Name) // Pass the database name from config
	if err != nil {
		log.Fatal("Error initializing database:", err)
	}

	clientHandler := handlers.NewHandler(db, cfg)
	app := App{Handler: clientHandler}
	return &app
}

func (a *App) Start(ctx context.Context) error {
	ctx, stop := signal.NotifyContext(context.Background(), syscall.SIGINT, syscall.SIGTERM)
	defer stop()

	doneWithErr := make(chan error)

	go func() {
		fmt.Println("server")
		if err := a.Handler.Server.ListenAndServe(); err != nil && !errors.Is(err, http.ErrServerClosed) {
			doneWithErr <- err
		}
	}()

	a.Handler.Database.Close()

	err := <-doneWithErr
	if err != nil {
		return err
	}

	return nil
}
