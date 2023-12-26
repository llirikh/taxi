package app

import (
	"context"
	"fmt"
	"github.com/jmoiron/sqlx"
	"github.com/pressly/goose/v3"
	"github.com/segmentio/kafka-go"
	"os"
	"os/signal"
	"syscall"
	"time"
	"trip/internal/config"
	"trip/internal/models"
	"trip/pkg/kafka-go"
)

type App struct {
	Config *models.Config

	Reader       *kafka.Reader
	WriterClient *kafka.Writer
	WriterDriver *kafka.Writer
}

func NewApp() *App {
	// todo: init config

	// KAFKA CONNECTION
	kafkaUrl := os.Getenv("kafkaURL")
	topicFrom := os.Getenv("topicFROM")
	topicClient := os.Getenv("topicCLIENT")
	topicDriver := os.Getenv("topicDRIVER")
	groupId := os.Getenv("groupID")

	fmt.Println(kafkaUrl, topicFrom, topicClient, topicClient)

	reader := kafka_go.NewReader(kafkaUrl, topicFrom, groupId)
	writerClient := kafka_go.NewWriter(kafkaUrl, topicClient)
	writerDriver := kafka_go.NewWriter(kafkaUrl, topicDriver)

	app := App{
		Config: &models.Config{},

		Reader:       reader,
		WriterClient: writerClient,
		WriterDriver: writerDriver,
	}

	return &app
}

func initializeDatabase(ctx context.Context, cfg *config.DbConfig) (*sqlx.DB, error) {
	// Ensure a 2-second delay before initializing the database (for demonstration purposes)
	time.Sleep(2 * time.Second)

	// Open a connection to the database
	db, err := sqlx.Open(drme, cfg.DSN)
	if err != nil {
		return nil, fmt.Errorf("failed to establish a database connection: %w", err)
	}

	// Set maximum open and idle connections along with connection lifetime
	db.DB.SetMaxOpenConns(100)
	db.DB.SetMaxIdleConns(10)
	db.DB.SetConnMaxLifetime(0)

	// Check if the connection is successful
	if err = db.PingContext(ctx); err != nil {
		return nil, fmt.Errorf("unable to ping the database: %w", err)
	}

	// Database migrations
	fs := os.DirFS(cfg.MigrationsDir)
	goose.SetBaseFS(fs)

	if err = goose.SetDialect(drme); err != nil {
		return nil, fmt.Errorf("failed to set database dialect: %w", err)
	}

	// Apply pending migrations
	if err = goose.UpContext(ctx, db.DB, "."); err != nil {
		return nil, fmt.Errorf("failed to apply database migrations: %w", err)
	}

	return db, nil
}

func (a *App) Start(ctx context.Context) error {
	ctx, stop := signal.NotifyContext(context.Background(), syscall.SIGINT, syscall.SIGTERM)
	defer stop()

	doneWithErr := make(chan error)

	err := <-doneWithErr
	go func() {
		for {
			_, err := kafka_go.ReadBytes(a.Reader)
			if err != nil {
				fmt.Println("ERROR READ: " + err.Error())
				doneWithErr <- err
			}
		}
	}()

	if err != nil {
		return err
	}

	return nil
}
