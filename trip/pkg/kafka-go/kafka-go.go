package kafka_go

import (
	"context"
	"fmt"
	"github.com/segmentio/kafka-go"
	"strings"
)

func NewWriter(kafkaUrl, topic string) *kafka.Writer {
	return &kafka.Writer{
		Addr:     kafka.TCP(kafkaUrl),
		Topic:    topic,
		Balancer: &kafka.LeastBytes{},
	}
}

func NewReader(kafkaURL, topic, groupID string) *kafka.Reader {
	brokers := strings.Split(kafkaURL, ",")
	return kafka.NewReader(kafka.ReaderConfig{
		Brokers:  brokers,
		GroupID:  groupID,
		Topic:    topic,
		MinBytes: 10e2, // 1KB
		MaxBytes: 10e4, // 100KB
	})
}

func WriteBytes(writer *kafka.Writer, key, bytes []byte) error {
	message := kafka.Message{
		Key:   key,
		Value: bytes,
	}

	err := writer.WriteMessages(context.Background(), message)
	if err != nil {
		return err
	}

	WriteLog(&message)

	return nil
}

func ReadBytes(reader *kafka.Reader) (*kafka.Message, error) {
	message, err := reader.ReadMessage(context.Background())
	if err != nil {
		return nil, err
	}

	ReadLog(&message)

	return &message, nil
}

// DEBUG
func WriteLog(m *kafka.Message) {
	fmt.Printf("PRODUSED -> TOPIC:%v PARTITION:%v OFFSET:%v	%s = %s\n", m.Topic, m.Partition, m.Offset, string(m.Key), string(m.Value))
}

func ReadLog(m *kafka.Message) {
	fmt.Printf("CONSUMED -> TOPIC:%v PARTITION:%v OFFSET:%v	%s = %s\n", m.Topic, m.Partition, m.Offset, string(m.Key), string(m.Value))
}
