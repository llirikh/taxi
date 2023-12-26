package main

import (
	"context"
	"fmt"
	"os"
	"strings"
	"time"

	"github.com/segmentio/kafka-go"
)

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

func newKafkaWriter(kafkaURL, topic string) *kafka.Writer {
	return &kafka.Writer{
		Addr:     kafka.TCP(kafkaURL),
		Topic:    topic,
		Balancer: &kafka.LeastBytes{},
	}
}

func SendMessage(writer *kafka.Writer, i int) {
	key := fmt.Sprintf("Key-%d", i)
	msg := kafka.Message{
		Key:   []byte(key),
		Value: []byte("FROM TRIP " + string(rune(i))),
	}
	var err error
	err = writer.WriteMessages(context.Background(), msg)

	if err != nil {
		fmt.Println("!!! TRIP WRITING ERROR : " + err.Error() + " !!!")
		fmt.Printf("!!! TOPIC:%v PARTITION:%v OFFSET:%v	%s = %s\n", msg.Topic, msg.Partition, msg.Offset, string(msg.Key), string(msg.Value))
	} else {
		fmt.Println("TRIP PRODUCED -> ", key)
	}
}

func GetMessage(reader *kafka.Reader) {
	m, err := reader.ReadMessage(context.Background())
	if err != nil {
		fmt.Println("!!! TRIP ERROR READING: " + err.Error() + " !!!")
		time.Sleep(1 * time.Second)
	}
	fmt.Printf("TRIP CONSUMED -> TOPIC:%v PARTITION:%v OFFSET:%v	%s = %s\n", m.Topic, m.Partition, m.Offset, string(m.Key), string(m.Value))
}

func main() {
	// Init kafka params
	kafkaURL := os.Getenv("kafkaURL")
	topicFrom := os.Getenv("topicFROM")
	topicClient := os.Getenv("topicCLIENT")
	topicDriver := os.Getenv("topicDRIVER")
	groupID := os.Getenv("GroupID")

	// Init reader
	reader := newKafkaReader(kafkaURL, topicFrom, groupID)
	defer reader.Close()

	// Init writers
	writerClient := newKafkaWriter(kafkaURL, topicClient)
	defer writerClient.Close()

	writerDriver := newKafkaWriter(kafkaURL, topicDriver)
	defer writerDriver.Close()

	fmt.Println("*** START TRIP ***")
	for i := 0; ; i++ {
		SendMessage(writerClient, i)
		fmt.Println("!!! " + topicClient + " !!!")
		SendMessage(writerDriver, i)
		time.Sleep(3 * time.Second)
		GetMessage(reader)
	}
}
