package main

import (
	"context"
	"fmt"
	"github.com/segmentio/kafka-go"
	"os"
	"time"
)

func newKafkaWriter(kafkaURL, topic string) *kafka.Writer {
	return &kafka.Writer{
		Addr:     kafka.TCP(kafkaURL),
		Topic:    topic,
		Balancer: &kafka.LeastBytes{},
	}
}

func main() {
	// get kafka writer using environment variables.
	kafkaURL := os.Getenv("kafkaURL")
	topic := os.Getenv("topic")
	writer := newKafkaWriter(kafkaURL, topic)
	defer writer.Close()

	fmt.Println("*** START CLIENT PRODUCE ***")
	for i := 0; ; i++ {
		key := fmt.Sprintf("Key-%d", i)
		msg := kafka.Message{
			Key:   []byte(key),
			Value: []byte("CLIENT MESSAGE " + string(rune(i))),
		}
		var err error
		err = writer.WriteMessages(context.Background(), msg)

		if err != nil {
			fmt.Println("!!! CLIENT ERROR WHILE MESSAGING: " + err.Error() + " !!!")
		} else {
			fmt.Println("CLIENT PRODUCED: ", key)
		}
		time.Sleep(2 * time.Second)
	}
}
