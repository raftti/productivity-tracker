package producer

import "github.com/confluentinc/confluent-kafka-go/kafka"

func CreateKafkaProducer() (*kafka.Producer, error) {
	producer, err := kafka.NewProducer(&kafka.ConfigMap{
		"bootstrap.servers": "kafka:9092",
	})

	if err != nil {
		return nil, err
	}
	
	return producer, nil
}