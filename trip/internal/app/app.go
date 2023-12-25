package app

import (
	"context"
	"fmt"
	"github.com/segmentio/kafka-go"
	"os"
	"os/signal"
	"syscall"
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
