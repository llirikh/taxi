package main

import (
	"context"
	"fmt"
	"github.com/segmentio/kafka-go"
	"os"
	"strings"
	"time"
)

func newKafkaWriter(kafkaURL, topic string) *kafka.Writer {
	return &kafka.Writer{
		Addr:     kafka.TCP(kafkaURL),
		Topic:    topic,
		Balancer: &kafka.LeastBytes{},
	}
}

func newKafkaReader(kafkaURL, topic, groupID string) *kafka.Reader {
	brokers := strings.Split(kafkaURL, ",")
	return kafka.NewReader(kafka.ReaderConfig{
		Brokers:  brokers,
		GroupID:  groupID,
		Topic:    topic,
		MinBytes: 10e3, // 10KB
		MaxBytes: 10e6, // 10MB
	})
}

func SendMessage(writer *kafka.Writer, i int) {
	key := fmt.Sprintf("Key-%d", i)
	msg := kafka.Message{
		Key:   []byte(key),
		Value: []byte("FROM CLIENT " + string(rune(i))),
	}
	var err error
	err = writer.WriteMessages(context.Background(), msg)

	if err != nil {
		fmt.Println("!!! CLIENT WRITING ERROR : " + err.Error() + " !!!")
	} else {
		fmt.Println("CLIENT PRODUCED -> ", key)
	}
}

func GetMessage(reader *kafka.Reader) {
	m, err := reader.ReadMessage(context.Background())
	if err != nil {
		fmt.Println("!!! CLIENT READING ERROR: " + err.Error() + " !!!")
		time.Sleep(1 * time.Second)
	}
	fmt.Printf("CLIENT CONSUMED -> TOPIC:%v PARTITION:%v OFFSET:%v	%s = %s\n", m.Topic, m.Partition, m.Offset, string(m.Key), string(m.Value))
}

func main() {
	// Init kafka params
	kafkaURL := os.Getenv("kafkaURL")
	topicTo := os.Getenv("topicTO")
	topicFrom := os.Getenv("topicFROM")
	groupId := os.Getenv("GroupID")

	// Init writer
	writer := newKafkaWriter(kafkaURL, topicTo)
	defer writer.Close()

	// Init reader
	reader := newKafkaReader(kafkaURL, topicFrom, groupId)
	defer reader.Close()

	fmt.Println("*** START CLIENT ***")
	for i := 0; ; i++ {
		SendMessage(writer, i)
		time.Sleep(4 * time.Second)
		GetMessage(reader)
	}
}
